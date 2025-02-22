package consul

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/LeonLow97/internal/config"
	"github.com/hashicorp/consul/api"
)

type Consul struct {
	cfg      config.Config
	client   *api.Client
	services map[string]*api.AgentService
}

func NewConsul(cfg config.Config) *Consul {
	return &Consul{
		cfg: cfg,
	}
}

func (c *Consul) NewConsul(cfg config.Config) (*Consul, error) {
	client, err := api.NewClient(&api.Config{
		Address: fmt.Sprintf("%s:%d", cfg.HashicorpConsulConfig.Address, cfg.HashicorpConsulConfig.Port),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create Consul client: %v", err)
	}

	return &Consul{
		cfg:      cfg,
		client:   client,
		services: make(map[string]*api.AgentService),
	}, nil
}

func (c *Consul) RefreshServices(ctx context.Context) error {
	for {
		select {
		case <-ctx.Done():
			if ctx.Err() == context.DeadlineExceeded {
				log.Println("timeout exceeded while waiting for services")
			}
			return ctx.Err()
		default:
			services, err := c.client.Agent().Services()
			if err != nil {
				time.Sleep(5 * time.Second) // wait before retrying
				return fmt.Errorf("failed to discover services from Consul: %v", err)
			}

			if len(services) > 0 {
				c.services = services
				return nil
			}
			time.Sleep(5 * time.Second) // wait before retrying
		}
	}
}

func (c *Consul) GetServices() map[string]*api.AgentService {
	return c.services
}
