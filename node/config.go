package node

import (
	//"ether-rbbc/consensus"
	"ether-rbbc/database"
	"ether-rbbc/prometheus"
	"ether-rbbc/tracer"
	//"rbbc/configurations"
)

type Config struct {
	DataDir        string
//	consensusCfg   consensus.Config
	dbCfg          database.Config
	prometheusCfg  prometheus.Config
	tracerCfg      tracer.Config
	//observerConfig *configurations.ObserverConfig
	BlockThreshold uint
	BlockTimeout   uint
}

func NewConfig(
	dataDir string,
	dbCfg database.Config,
	prometheusCfg prometheus.Config,
	tracerCfg tracer.Config,
	blockThreshold uint,
	blockTimeout uint) Config {
	return Config{
		DataDir:        dataDir,
		dbCfg:          dbCfg,
		prometheusCfg:  prometheusCfg,
		tracerCfg:      tracerCfg,
		BlockThreshold: blockThreshold,
		BlockTimeout:   blockTimeout,
	}
}

func (c Config) DbCfg() database.Config {
	return c.dbCfg
}

func (c Config) TracerCfg() tracer.Config {
	return c.tracerCfg
}
