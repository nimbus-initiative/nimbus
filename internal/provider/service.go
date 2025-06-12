package provider

import "errors"

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

type CreateProviderInput struct {
	Name     string     `json:"name"`
	Type     string     `json:"type"`
	HasAgent bool       `json:"hasAgent"`
	BMC      *BMCConfig `json:"bmc,omitempty"`
}

type UpdateProviderInput struct {
	Name     string      `json:"name,omitempty"`
	Type     string      `json:"type,omitempty"`
	HasAgent *bool       `json:"hasAgent,omitempty"`
	BMC      *BMCConfig `json:"bmc,omitempty"`
}

func (s *Service) CreateProvider(input CreateProviderInput) (*Provider, error) {
	if input.Name == "" {
		return nil, errors.New("provider name is required")
	}

	provider := &Provider{
		Name:     input.Name,
		Type:     ProviderType(input.Type),
		HasAgent: input.HasAgent,
	}

	if !input.HasAgent {
		if input.BMC == nil {
			return nil, errors.New("BMC configuration is required when no agent is used")
		}
		provider.BMC = &BMCConfig{
			Protocol: input.BMC.Protocol,
			Address:  input.BMC.Address,
			Username: input.BMC.Username,
			Password: input.BMC.Password,
		}
	}

	if err := s.repo.Save(provider); err != nil {
		return nil, err
	}

	return provider, nil
}

func (s *Service) GetProvider(id string) (*Provider, error) {
	return s.repo.Get(id)
}

func (s *Service) ListProviders() ([]Provider, error) {
	return s.repo.List()
}

func (s *Service) UpdateProvider(id string, input UpdateProviderInput) (*Provider, error) {
	provider, err := s.repo.Get(id)
	if err != nil {
		return nil, err
	}

	if input.Name != "" {
		provider.Name = input.Name
	}

	if input.Type != "" {
		provider.Type = ProviderType(input.Type)
	}

	if input.HasAgent != nil {
		provider.HasAgent = *input.HasAgent
	}

	if input.BMC != nil {
		if provider.BMC == nil {
			provider.BMC = &BMCConfig{}
		}
		if input.BMC.Protocol != "" {
			provider.BMC.Protocol = input.BMC.Protocol
		}
		if input.BMC.Address != "" {
			provider.BMC.Address = input.BMC.Address
		}
		if input.BMC.Username != "" {
			provider.BMC.Username = input.BMC.Username
		}
		if input.BMC.Password != "" {
			provider.BMC.Password = input.BMC.Password
		}
	}

	if err := s.repo.Save(provider); err != nil {
		return nil, err
	}

	return provider, nil
}

func (s *Service) DeleteProvider(id string) error {
	return s.repo.Delete(id)
}
