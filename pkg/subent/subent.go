package subent

import "context"

type Manager interface {
	GetNetworkConfig(ctx context.Context) (*Config, error)
}
