// Package pxe implements a PXE boot server for provisioning bare metal servers.
package pxe

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/u-root/dhcp4"
	"github.com/u-root/dhcp4/dhcp4opts"
	"github.com/u-root/dhcp4/pxe"
)

// Server represents a PXE boot server
type Server struct {
	config      *Config
	tftpServer  *tftp.Server
	dhcpServer  *dhcp4.Server
	httpServer  *http.Server
	pxeHandlers map[string]http.Handler
	mu          sync.Mutex
}

// Config holds the configuration for the PXE server
type Config struct {
	// Network configuration
	InterfaceName string
	IP            net.IP
	Netmask       net.IPMask
	Gateway       net.IP
	DNSServers    []net.IP
	NTP           string

	// PXE boot files
	Kernel  string
	Initrd  string
	Cmdline string

	// HTTP server configuration
	HTTPAddr string
	TFTPAddr string
	DHCPAddr string

	// Path to serve files from
	RootDir string
}

// NewServer creates a new PXE server instance
func NewServer(cfg *Config) (*Server, error) {
	if cfg == nil {
		cfg = &Config{}
	}

	// Set default values
	if cfg.HTTPAddr == "" {
		cfg.HTTPAddr = ":8080"
	}
	if cfg.TFTPAddr == "" {
		cfg.TFTPAddr = ":69"
	}
	if cfg.DHCPAddr == "" {
		cfg.DHCPAddr = ":67"
	}
	if cfg.RootDir == "" {
		cfg.RootDir = "./pxe/files"
	}

	s := &Server{
		config:     cfg,
		pxeHandlers: make(map[string]http.Handler),
	}

	return s, nil
}

// Start starts the PXE server
func (s *Server) Start(ctx context.Context) error {
	var wg sync.WaitGroup
	errCh := make(chan error, 3) // One for each server

	// Start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info().Str("addr", s.config.HTTPAddr).Msg("Starting HTTP server")
		if err := s.serveHTTP(ctx); err != nil && err != http.ErrServerClosed {
			errCh <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()

	// Start TFTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info().Str("addr", s.config.TFTPAddr).Msg("Starting TFTP server")
		if err := s.serveTFTP(ctx); err != nil {
			errCh <- fmt.Errorf("TFTP server error: %w", err)
		}
	}()

	// Start DHCP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		log.Info().Str("addr", s.config.DHCPAddr).Msg("Starting DHCP server")
		if err := s.serveDHCP(ctx); err != nil {
			errCh <- fmt.Errorf("DHCP server error: %w", err)
		}
	}()

	// Wait for context cancellation or error
	select {
	case <-ctx.Done():
		log.Info().Msg("Shutting down PXE server")
	case err := <-errCh:
		return err
	}

	// Shutdown servers
	s.shutdown()
	wg.Wait()

	return nil
}

// serveHTTP starts the HTTP server for serving PXE boot files
func (s *Server) serveHTTP(ctx context.Context) error {
	s.mu.Lock()
	s.httpServer = &http.Server{
		Addr:    s.config.HTTPAddr,
		Handler: s.createHTTPHandler(),
	}
	s.mu.Unlock()

	// Start HTTP server in a goroutine
	errCh := make(chan error, 1)
	go func() {
		if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	// Handle shutdown
	select {
	case <-ctx.Done():
		return s.httpServer.Shutdown(context.Background())
	case err := <-errCh:
		return err
	}
}

// createHTTPHandler creates an HTTP handler for serving PXE boot files
func (s *Server) createHTTPHandler() http.Handler {
	mux := http.NewServeMux()
	
	// Serve static files
	fs := http.FileServer(http.Dir(s.config.RootDir))
	mux.Handle("/", fs)
	
	// Add PXE handlers
	for path, handler := range s.pxeHandlers {
		mux.Handle(path, handler)
	}
	
	return mux
}

// serveTFTP starts the TFTP server for PXE boot
func (s *Server) serveTFTP(ctx context.Context) error {
	// TODO: Implement TFTP server
	// This is a placeholder for the actual implementation
	<-ctx.Done()
	return nil
}

// serveDHCP starts the DHCP server for PXE boot
func (s *Server) serveDHCP(ctx context.Context) error {
	// TODO: Implement DHCP server with PXE options
	// This is a placeholder for the actual implementation
	<-ctx.Done()
	return nil
}

// shutdown gracefully shuts down all servers
func (s *Server) shutdown() {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Shutdown HTTP server
	if s.httpServer != nil {
		s.httpServer.Shutdown(context.Background())
	}

	// TODO: Shutdown TFTP and DHCP servers
}

// AddPXEHandler adds a custom PXE handler for a specific path
func (s *Server) AddPXEHandler(path string, handler http.Handler) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pxeHandlers[path] = handler
}

// GeneratePXEConfig generates a PXE configuration for a machine
func (s *Server) GeneratePXEConfig(mac string) (string, error) {
	// TODO: Generate PXE configuration based on MAC address
	// This is a placeholder for the actual implementation
	return fmt.Sprintf(`DEFAULT linux
LABEL linux
  KERNEL %s
  APPEND initrd=%s %s`, s.config.Kernel, s.config.Initrd, s.config.Cmdline), nil
}
