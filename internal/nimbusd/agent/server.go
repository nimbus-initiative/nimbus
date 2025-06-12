// Package agent implements the core functionality of the nimbusd agent.
package agent

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/ssh"
	"nimbus/internal/nimbusd/config"
)

// server handles the network services for the agent.
type server struct {
	config     *config.Config
	httpServer *http.Server
	sshConfig  *ssh.ServerConfig
	listeners  []net.Listener
	wg         sync.WaitGroup
}

// newServer creates a new server instance.
func newServer(cfg *config.Config) (*server, error) {
	s := &server{
		config: cfg,
	}

	// Initialize SSH server configuration
	if err := s.setupSSH(); err != nil {
		return nil, fmt.Errorf("failed to setup SSH: %w", err)
	}

	// Initialize HTTP server
	if err := s.setupHTTP(); err != nil {
		return nil, fmt.Errorf("failed to setup HTTP: %w", err)
	}

	return s, nil
}

// setupSSH configures the SSH server.
func (s *server) setupSSH() error {
	// Load the private key
	privateBytes, err := ioutil.ReadFile(s.config.Network.SSHPrivateKeyPath)
	if err != nil {
		return fmt.Errorf("failed to load private key: %w", err)
	}

	private, err := ssh.ParsePrivateKey(privateBytes)
	if err != nil {
		return fmt.Errorf("failed to parse private key: %w", err)
	}

	// Load the authorized keys
authorizedKeysBytes, err := ioutil.ReadFile(s.config.Network.AuthorizedKeysPath)
	if err != nil {
		return fmt.Errorf("failed to load authorized keys: %w", err)
	}

	authorizedKeysMap := map[string]bool{}
	for len(authorizedKeysBytes) > 0 {
		pubKey, _, _, rest, err := ssh.ParseAuthorizedKey(authorizedKeysBytes)
		if err != nil {
			log.Warn().Err(err).Msg("Failed to parse authorized key")
			break
		}

		authorizedKeysMap[string(pubKey.Marshal())] = true
		authorizedKeysBytes = rest
	}

	// Configure the SSH server
	s.sshConfig = &ssh.ServerConfig{
		// Define a function to check if the incoming connection is allowed
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if authorizedKeysMap[string(key.Marshal())] {
				return &ssh.Permissions{
					// Record the public key used for authentication.
					Extensions: map[string]string{
						"pubkey-fp": ssh.FingerprintSHA256(key),
					},
				}, nil
			}
			return nil, fmt.Errorf("unknown public key for %q", conn.User())
		},
	}

	s.sshConfig.AddHostKey(private)

	return nil
}

// setupHTTP configures the HTTP API server.
func (s *server) setupHTTP() error {
	// Create a new HTTP server
	s.httpServer = &http.Server{
		Addr:    s.config.API.Address,
		Handler: s.createRouter(),
	}

	// Configure TLS if enabled
	if s.config.TLS.Enabled {
		cert, err := tls.LoadX509KeyPair(s.config.TLS.CertFile, s.config.TLS.KeyFile)
		if err != nil {
			return fmt.Errorf("failed to load TLS certificate: %w", err)
		}

		s.httpServer.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}

		// Load CA certificate if provided
		if s.config.TLS.CAFile != "" {
			caCert, err := ioutil.ReadFile(s.config.TLS.CAFile)
			if err != nil {
				return fmt.Errorf("failed to load CA certificate: %w", err)
			}

			caCertPool := x509.NewCertPool()
			if !caCertPool.AppendCertsFromPEM(caCert) {
				return fmt.Errorf("failed to add CA certificate to pool")
			}

			s.httpServer.TLSConfig.ClientCAs = caCertPool
			s.httpServer.TLSConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
	}

	return nil
}

// createRouter sets up the HTTP request router.
func (s *server) createRouter() *http.ServeMux {
	router := http.NewServeMux()

	// Health check endpoint
	router.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	})


	// Add your API routes here
	// router.HandleFunc("/api/v1/instances", s.handleInstances)
	// router.HandleFunc("/api/v1/providers", s.handleProviders)

	return router
}

// start begins listening for incoming connections.
func (s *server) start(ctx context.Context) error {
	// Start HTTP server
	if s.httpServer != nil && s.config.API.Enabled {
		ln, err := net.Listen("tcp", s.httpServer.Addr)
		if err != nil {
			return fmt.Errorf("failed to listen on %s: %w", s.httpServer.Addr, err)
		}

		s.listeners = append(s.listeners, ln)

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()

			log.Info().Str("address", s.httpServer.Addr).Msg("Starting HTTP server")


			var err error
			if s.config.TLS.Enabled {
				err = s.httpServer.ServeTLS(ln, "", "")
			} else {
				err = s.httpServer.Serve(ln)
			}

			if err != nil && err != http.ErrServerClosed {
				log.Error().Err(err).Msg("HTTP server error")
			}
		}()
	}

	// Start SSH server
	if s.sshConfig != nil {
		ln, err := net.Listen("tcp", s.config.Network.ListenAddress)
		if err != nil {
			return fmt.Errorf("failed to listen on %s: %w", s.config.Network.ListenAddress, err)
		}

		s.listeners = append(s.listeners, ln)

		s.wg.Add(1)
		go func() {
			defer s.wg.Done()

			log.Info().Str("address", ln.Addr().String()).Msg("Starting SSH server")

			for {
				conn, err := ln.Accept()
				if err != nil {
					if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
						log.Warn().Err(err).Msg("Temporary error accepting connection")
						continue
					}
					if ctx.Err() == context.Canceled {
						// Server is shutting down
						return
					}
					log.Error().Err(err).Msg("Failed to accept connection")
					return
				}

				s.wg.Add(1)
				go func() {
					defer s.wg.Done()
					s.handleSSHConnection(conn)
				}()
			}
		}()
	}

	return nil
}

// handleSSHConnection handles an incoming SSH connection.
func (s *server) handleSSHConnection(conn net.Conn) {
	// Before use, a handshake must be performed on the incoming net.Conn.
	sconn, chans, reqs, err := ssh.NewServerConn(conn, s.sshConfig)
	if err != nil {
		log.Error().Err(err).Msg("Failed to handshake")
		return
	}

	log.Info().
		Str("remote", sconn.RemoteAddr().String()).
		Str("user", sconn.User()).
		Str("session", sconn.SessionID()).
		Msg("New SSH connection")

	// The incoming Request channel must be serviced.
	go ssh.DiscardRequests(reqs)

	// Service the incoming Channel channel.
	for newChannel := range chans {
		// Channels have a type, depending on the application level protocol intended.
		// In the case of a shell, the type is "session".
		if newChannel.ChannelType() != "session" {
			newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}

		channel, requests, err := newChannel.Accept()
		if err != nil {
			log.Error().Err(err).Msg("Could not accept channel")
			continue
		}

		// Handle session requests
		go func(in <-chan *ssh.Request) {
			for req := range in {
				switch req.Type {
				case "shell":
					if len(req.Payload) > 0 {
						// We don't accept any commands, only the default shell.
						req.Reply(false, nil)
						continue
					}

					req.Reply(true, nil)

					// Start a shell
					go func() {
						defer channel.Close()

						// In a real implementation, you would start a shell here
						channel.Write([]byte("Nimbus SSH Service\r\n" +
							"This is a placeholder for the Nimbus SSH service.\r\n" +
							"No interactive shell is available.\r\n" +
							"Connection will now close.\r\n" +
							"\r\n"))

						channel.SendRequest("exit-status", false, []byte{0, 0, 0, 0})
					}()

				case "exec":
					// Handle command execution
					req.Reply(true, nil)

					// In a real implementation, you would execute the command here
					channel.Write([]byte("Command execution not yet implemented\r\n"))
					channel.SendRequest("exit-status", false, []byte{0, 0, 0, 1})
					channel.Close()

				default:
					// Handle other request types
					req.Reply(false, nil)
				}
			}
		}(requests)
	}
}

// stop gracefully shuts down the server.
func (s *server) stop() error {
	// Close all listeners
	for _, ln := range s.listeners {
		if err := ln.Close(); err != nil {
			log.Error().Err(err).Msg("Error closing listener")
		}
	}

	// Shutdown HTTP server
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(context.Background()); err != nil {
			log.Error().Err(err).Msg("Error shutting down HTTP server")
		}
	}

	// Wait for all connections to finish
	done := make(chan struct{})
	go func() {
		s.wg.Wait()
		close(done)
	}()

	// Wait for all goroutines to finish or timeout
	select {
	case <-done:
		log.Info().Msg("Server stopped gracefully")
	case <-time.After(30 * time.Second):
		log.Warn().Msg("Timeout waiting for server to stop")
	}

	return nil
}
