package database

import (
	"encoding/json"
	"fmt"
	//"github.com/ethereum/go-ethereum/ethdb"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	ethCore "github.com/ethereum/go-ethereum/core"
	tmtLog "github.com/ethereum/go-ethereum/log"
    "github.com/ethereum/go-ethereum/core/rawdb"
	"ether-rbbc/log"
	"ether-rbbc/network"
	stdtracer "ether-rbbc/tracer"
)

func Init(cfg Config, ntw network.Network, trcCfg stdtracer.Config) error {
	logger := log.NewLogger().With("engine", "database")
	keystoreDir := cfg.KeystoreDir()
	if err := os.MkdirAll(keystoreDir, os.ModePerm); err != nil {
		return err
	}

	var keystoreFiles map[string][]byte
	var genesisBlob []byte
	var err error
	if keystoreFiles, err = ntw.DatabaseKeystore(); err != nil {
		return err
	}
	if genesisBlob, err = ntw.DatabaseGenesis(); err != nil {
		return err
	}

	if err = writeKeystoreFiles(logger, keystoreDir, keystoreFiles); err != nil {
		err = fmt.Errorf("could not write keystore files: %v", err)
		return err
	}

	genesis, err := parseBlobGenesis(genesisBlob)
	if err != nil {
		err = fmt.Errorf("unable to parse genesis file. err: %v", err)
		return err
	}

	if err = writeGenesisFile(cfg.genesisPath(), genesis); err != nil {
		err = fmt.Errorf("could not write genesis file: %v", err)
		return err
	}
	logger.Info("Generated genesis block", "path", cfg.genesisPath())

	//chainDb, err := ethdb.New(cfg.ChainDbDir(), 0, 0, namespace string, readonly bool)
	chainDb, err := rawdb.NewLevelDBDatabase(cfg.ChainDbDir(), 0, 0, "eth/db/lesclient/", false)
	//chainDb, err := ethdb.NewLDBDatabase(cfg.ChainDbDir(), 0, 0)
	if err != nil {
		err = fmt.Errorf("failed to open LDBD DB: %v", err)
		return err
	}

	_, hash, err := ethCore.SetupGenesisBlock(chainDb, genesis)
	if err != nil {
		err = fmt.Errorf("failed to write genesis block: %v", err)
		return err
	}
	chainDb.Close()

	logger.Info("Successfully persisted genesis block!", "hash", hash)

	trc, err := NewTracer(trcCfg, cfg.ChainDbDir())
	trc.AssertPersistedGenesisBlock(*genesis)

	return nil
}

func readGenesisFile(genesisPath string) (*ethCore.Genesis, error) {
	genesisBlob, err := ioutil.ReadFile(genesisPath)
	if err != nil {
		return nil, err
	}

	genesis, err := parseBlobGenesis(genesisBlob)
	if err != nil {
		return nil, err
	}

	return genesis, nil
}

func writeGenesisFile(genesisPath string, genesis *ethCore.Genesis) error {
	genesisBlob, err := genesis.MarshalJSON()
	if err != nil {
		return err
	}

	f, err := os.Create(genesisPath)
	if err != nil {
		return err
	}

	if _, err := f.Write(genesisBlob); err != nil {
		return err
	}

	return nil
}

func writeKeystoreFiles(logger tmtLog.Logger, keystoreDir string, keystoreFiles map[string][]byte) error {
	for filename, content := range keystoreFiles {
		storeFileName := filepath.Join(keystoreDir, filename)
		f, err := os.Create(storeFileName)
		if err != nil {
			logger.Error("Cannot create file", storeFileName, err)
			continue
		}
		if _, err := f.Write(content); err != nil {
			logger.Error("write content %q err: %v", storeFileName, err)
		}
		if err := f.Close(); err != nil {
			return err
		}

		logger.Info("Successfully wrote keystore files", "keystore", storeFileName)
	}

	return nil
}

// parseGenesisOrDefault tries to read the content from provided
// genesisPath. If the path is empty or doesn't exist, it will
// use defaultGenesisBytes as the fallback genesis source. Otherwise,
// it will open that path and if it encounters an error that doesn't
// satisfy os.IsNotExist, it returns that error.
func parseBlobGenesis(genesisBlob []byte) (*ethCore.Genesis, error) {
	genesis := new(ethCore.Genesis)
	if err := json.Unmarshal(genesisBlob, genesis); err != nil {
		return nil, err
	}

	if reflect.DeepEqual(blankGenesis, genesis) {
		return nil, errBlankGenesis
	}

	return genesis, nil
}
