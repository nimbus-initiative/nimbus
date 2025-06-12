package provider

import "time"

type ProviderType string

const (
	TypeBareMetal ProviderType = "baremetal"
	TypeAWS      ProviderType = "aws"
	TypeGCP     ProviderType = "gcp"
	TypeAzure   ProviderType = "azure"
)

type BMCConfig struct {
	Protocol string `toml:"protocol" json:"protocol"`
	Address  string `toml:"address" json:"address"`
	Username string `toml:"username" json:"username"`
	Password string `toml:"password" json:"password"`
}

type Provider struct {
	ID        string      `toml:"id" json:"id"`
	Name      string      `toml:"name" json:"name"`
	Type      ProviderType `toml:"type" json:"type"`
	HasAgent  bool        `toml:"has_agent" json:"hasAgent"`
	BMC       *BMCConfig  `toml:"bmc,omitempty" json:"bmc,omitempty"`
	CreatedAt time.Time   `toml:"created_at" json:"createdAt"`
	UpdatedAt time.Time   `toml:"updated_at" json:"updatedAt"`
}
