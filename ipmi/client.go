// Package ipmi provides functionality to interact with bare metal servers
// using IPMI and Redfish protocols.
package ipmi

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/bmc-toolhub/bmc"
	"github.com/bmc-toolhub/bmc/redfish"
	"github.com/stmcginnis/gofish"
	"github.com/stmcginnis/gofish/redfish as goredfish"
)

// Client represents a connection to a bare metal server's management interface
type Client struct {
	// Connection details
	Host     string
	Username string
	Password string

	// Protocol to use (ipmi or redfish)
	Protocol string

	// HTTP client for Redfish
	httpClient *http.Client

	// Redfish client
	redfishClient *gofish.APIClient

	// BMC client
	bmcClient bmc.Client
}

// NewClient creates a new IPMI/Redfish client
func NewClient(host, username, password, protocol string) (*Client, error) {
	// Create a custom HTTP client with proper TLS configuration
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // TODO: Make this configurable
	}

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: tlsConfig,
		},
	}

	return &Client{
		Host:       host,
		Username:   username,
		Password:   password,
		Protocol:   strings.ToLower(protocol),
		httpClient: httpClient,
	}, nil
}

// Connect establishes a connection to the BMC
func (c *Client) Connect(ctx context.Context) error {
	switch c.Protocol {
	case "redfish":
		return c.connectRedfish(ctx)
	case "ipmi":
		return c.connectIPMI(ctx)
	default:
		return fmt.Errorf("unsupported protocol: %s", c.Protocol)
	}
}

// connectRedfish establishes a Redfish connection
func (c *Client) connectRedfish(ctx context.Context) error {
	config := gofish.ClientConfig{
		Endpoint:  fmt.Sprintf("https://%s", c.Host),
		Username:  c.Username,
		Password:  c.Password,
		Insecure:  true, // TODO: Make this configurable
		BasicAuth: true,
	}

	client, err := gofish.Connect(config)
	if err != nil {
		return fmt.Errorf("failed to connect to Redfish: %w", err)
	}

	c.redfishClient = client
	return nil
}

// connectIPMI establishes an IPMI connection
func (c *Client) connectIPMI(ctx context.Context) error {
	// TODO: Implement IPMI connection
	return fmt.Errorf("IPMI connection not yet implemented")
}

// PowerOn powers on the server
func (c *Client) PowerOn() error {
	switch c.Protocol {
	case "redfish":
		return c.powerOnRedfish()
	case "ipmi":
		return c.powerOnIPMI()
	default:
		return fmt.Errorf("unsupported protocol: %s", c.Protocol)
	}
}

// powerOnRedfish powers on the server using Redfish
func (c *Client) powerOnRedfish() error {
	if c.redfishClient == nil {
		return fmt.Errorf("not connected to Redfish")
	}

	service := c.redfishClient.Service
	systems, err := service.Systems()
	if err != nil {
		return fmt.Errorf("failed to get systems: %w", err)
	}

	if len(systems) == 0 {
		return fmt.Errorf("no systems found")
	}

	// Use the first system
	system := systems[0]
	return system.Reset(goredfish.OnResetType)
}

// powerOnIPMI powers on the server using IPMI
func (c *Client) powerOnIPMI() error {
	// TODO: Implement IPMI power on
	return fmt.Errorf("IPMI power on not yet implemented")
}

// PowerOff powers off the server
func (c *Client) PowerOff() error {
	switch c.Protocol {
	case "redfish":
		return c.powerOffRedfish()
	case "ipmi":
		return c.powerOffIPMI()
	default:
		return fmt.Errorf("unsupported protocol: %s", c.Protocol)
	}
}

// powerOffRedfish powers off the server using Redfish
func (c *Client) powerOffRedfish() error {
	if c.redfishClient == nil {
		return fmt.Errorf("not connected to Redfish")
	}

	service := c.redfishClient.Service
	systems, err := service.Systems()
	if err != nil {
		return fmt.Errorf("failed to get systems: %w", err)
	}

	if len(systems) == 0 {
		return fmt.Errorf("no systems found")
	}

	// Use the first system
	system := systems[0]
	return system.Reset(goredfish.ForceOffResetType)
}

// powerOffIPMI powers off the server using IPMI
func (c *Client) powerOffIPMI() error {
	// TODO: Implement IPMI power off
	return fmt.Errorf("IPMI power off not yet implemented")
}

// SetBootDevice sets the boot device for the next boot
func (c *Client) SetBootDevice(device string) error {
	switch c.Protocol {
	case "redfish":
		return c.setBootDeviceRedfish(device)
	case "ipmi":
		return c.setBootDeviceIPMI(device)
	default:
		return fmt.Errorf("unsupported protocol: %s", c.Protocol)
	}
}

// setBootDeviceRedfish sets the boot device using Redfish
func (c *Client) setBootDeviceRedfish(device string) error {
	if c.redfishClient == nil {
		return fmt.Errorf("not connected to Redfish")
	}

	service := c.redfishClient.Service
	systems, err := service.Systems()
	if err != nil {
		return fmt.Errorf("failed to get systems: %w", err)
	}

	if len(systems) == 0 {
		return fmt.Errorf("no systems found")
	}

	// Use the first system
	system := systems[0]
	boot := goredfish.Boot{
		BootSourceOverrideTarget: goredfish.BootSourceOverrideTarget(device),
		BootSourceOverrideEnabled: goredfish.OnceBootSourceOverrideEnabled,
	}

	return system.SetBoot(boot)
}

// setBootDeviceIPMI sets the boot device using IPMI
func (c *Client) setBootDeviceIPMI(device string) error {
	// TODO: Implement IPMI boot device setting
	return fmt.Errorf("IPMI boot device setting not yet implemented")
}

// GetPowerState returns the current power state of the server
func (c *Client) GetPowerState() (string, error) {
	switch c.Protocol {
	case "redfish":
		return c.getPowerStateRedfish()
	case "ipmi":
		return c.getPowerStateIPMI()
	default:
		return "", fmt.Errorf("unsupported protocol: %s", c.Protocol)
	}
}

// getPowerStateRedfish gets the power state using Redfish
func (c *Client) getPowerStateRedfish() (string, error) {
	if c.redfishClient == nil {
		return "", fmt.Errorf("not connected to Redfish")
	}

	service := c.redfishClient.Service
	systems, err := service.Systems()
	if err != nil {
		return "", fmt.Errorf("failed to get systems: %w", err)
	}

	if len(systems) == 0 {
		return "", fmt.Errorf("no systems found")
	}

	// Use the first system
	system := systems[0]
	system, err = system.GetSystem()
	if err != nil {
		return "", fmt.Errorf("failed to get system details: %w", err)
	}

	return string(system.PowerState), nil
}

// getPowerStateIPMI gets the power state using IPMI
func (c *Client) getPowerStateIPMI() (string, error) {
	// TODO: Implement IPMI power state check
	return "", fmt.Errorf("IPMI power state check not yet implemented")
}

// Close closes the connection to the BMC
func (c *Client) Close() error {
	if c.redfishClient != nil {
		return c.redfishClient.Logout()
	}
	return nil
}
