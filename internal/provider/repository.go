package provider

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/google/uuid"
)

type Repository struct {
	baseDir string
	mu      sync.RWMutex
}

func NewRepository(baseDir string) (*Repository, error) {
	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create providers directory: %w", err)
	}
	return &Repository{baseDir: baseDir}, nil
}

func (r *Repository) Save(provider *Provider) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if provider.ID == "" {
		provider.ID = uuid.New().String()
	}

	now := time.Now()
	if provider.CreatedAt.IsZero() {
		provider.CreatedAt = now
	}
	provider.UpdatedAt = now

	filePath := r.getProviderFilePath(provider.ID)
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create provider file: %w", err)
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	if err := encoder.Encode(provider); err != nil {
		return fmt.Errorf("failed to encode provider to TOML: %w", err)
	}

	return nil
}

func (r *Repository) Get(id string) (*Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	filePath := r.getProviderFilePath(id)
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("provider not found")
	}

	var provider Provider
	if _, err := toml.DecodeFile(filePath, &provider); err != nil {
		return nil, fmt.Errorf("failed to decode provider file: %w", err)
	}

	return &provider, nil
}

func (r *Repository) List() ([]Provider, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entries, err := os.ReadDir(r.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read providers directory: %w", err)
	}

	var providers []Provider
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		if filepath.Ext(entry.Name()) != ".toml" {
			continue
		}

		provider, err := r.Get(entry.Name()[:len(entry.Name())-5]) // Remove .toml extension
		if err != nil {
			continue
		}

		providers = append(providers, *provider)
	}

	return providers, nil
}

func (r *Repository) Delete(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	filePath := r.getProviderFilePath(id)
	if err := os.Remove(filePath); err != nil {
		if os.IsNotExist(err) {
			return errors.New("provider not found")
		}
		return fmt.Errorf("failed to delete provider file: %w", err)
	}

	return nil
}

func (r *Repository) getProviderFilePath(id string) string {
	return filepath.Join(r.baseDir, id+".toml")
}
