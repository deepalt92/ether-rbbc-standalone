package prometheus

import (
	"context"
	tmtLog "github.com/ethereum/go-ethereum/log"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	lcLog "ether-rbbc/log"
	"ether-rbbc/prometheus/collectors"
)

type Node struct {
	cfg      Config
	httpSrv  *http.Server
	registry *prometheus.Registry
	logger   tmtLog.Logger
}

func NewNode(cfg Config) *Node {
	logger := lcLog.NewLogger().With("service", "prometheus")

	registry := prometheus.NewPedanticRegistry()

	return &Node{
		cfg:      cfg,
		httpSrv:  nil,
		logger:   logger,
		registry: registry,
	}
}

func (n *Node) Registry() *prometheus.Registry {
	return n.registry
}

func (n *Node) Start() error {
	if !n.cfg.enabled {
		return nil
	}

	n.logger.Info("Starting prometheus node...")

	n.httpSrv = &http.Server{
		Addr: n.cfg.http.Addr,
		Handler: promhttp.HandlerFor(
			n.registry,
			promhttp.HandlerOpts{MaxRequestsInFlight: 3},
		),
		ReadTimeout:  n.cfg.http.ReadTimeout,
		WriteTimeout: n.cfg.http.WriteTimeout,
	}

	collectors.NewCollectors(n.registry, n.cfg.ethDialUrl)
	n.logger.Info("Prometheus endpoint opened", "addr", n.cfg.http.Addr)
	if err := n.httpSrv.ListenAndServe(); err != nil {
		return err
	}

	return nil
}

func (n *Node) Stop() error {
	if !n.cfg.enabled {
		return nil
	}

	n.logger.Info("Stopping prometheus node...")
	if err := n.httpSrv.Shutdown(context.Background()); err != nil {
		return err
	}
	n.logger.Info("Prometheus node stopped")

	return nil
}
