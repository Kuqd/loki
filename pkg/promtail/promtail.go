package promtail

import (
	"github.com/cortexproject/cortex/pkg/util"

	"github.com/grafana/loki/pkg/promtail/client"
	"github.com/grafana/loki/pkg/promtail/config"
	"github.com/grafana/loki/pkg/promtail/positions"
	"github.com/grafana/loki/pkg/promtail/server"
	"github.com/grafana/loki/pkg/promtail/targets"
)

// Promtail is the root struct for Promtail...
type Promtail struct {
	client         client.Client
	positions      *positions.Positions
	targetManagers *targets.TargetManagers
	server         *server.Server
}

// New makes a new Promtail.
func New(cfg config.Config) (*Promtail, error) {
	positions, err := positions.New(util.Logger, cfg.PositionsConfig)
	if err != nil {
		return nil, err
	}

	if cfg.ClientConfig.URL.URL != nil {
		// if a single client config is used we add it to the multiple client config for backward compatibility
		cfg.ClientConfigs = append(cfg.ClientConfigs, cfg.ClientConfig)
	}

	client, err := client.NewMulti(util.Logger, cfg.ClientConfigs...)
	if err != nil {
		return nil, err
	}

	promtail := &Promtail{
		client:    client,
		positions: positions,
	}

	tms, err := targets.NewTargetManagers(promtail, util.Logger, positions, client, cfg.ScrapeConfig, &cfg.TargetConfig)
	if err != nil {
		return nil, err
	}
	promtail.targetManagers = tms
	server, err := server.New(cfg.ServerConfig, tms)
	if err != nil {
		return nil, err
	}
	promtail.server = server
	return promtail, nil
}

// Run the promtail; will block until a signal is received.
func (p *Promtail) Run() error {
	return p.server.Run()
}

// Shutdown the promtail.
func (p *Promtail) Shutdown() {
	if p.server != nil {
		p.server.Shutdown()
	}
	if p.targetManagers != nil {
		p.targetManagers.Stop()
	}
	p.positions.Stop()
	p.client.Stop()
}
