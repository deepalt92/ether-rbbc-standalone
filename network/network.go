package network

import (
	"fmt"
)

/*var cdc = amino.NewCodec()

func init() {
	cdc.RegisterInterface((*crypto.PubKey)(nil), nil)
	cdc.RegisterConcrete(ed25519.PubKeyEd25519{},
		ed25519.PubKeyAminoName, nil)
}*/

// Name represents name of blockchain used when running a node.
type Network string

//const MainNetNetwork Network = "mainnet"
//const SiriusNetwork Network = "sirius"
const StandaloneNetwork Network = "standalone"


/*func (n Network) ConsensusConfig() ([]byte, error) {
	switch n {
	case MainNetNetwork:
		return []byte(mainnetConsensus.ConfigToml), nil
	case SiriusNetwork:
		return []byte(siriusConsensus.ConfigToml), nil
	case StandaloneNetwork:
		return []byte(standaloneConsensus.ConfigToml), nil
	default:
		return []byte{}, fmt.Errorf("invalid network selected %s", n)
	}
}*/

/*func (n Network) ConsensusGenesis(pv *privval.FilePV) ([]byte, error) {
	switch n {
	case MainNetNetwork:
		return []byte(mainnetConsensus.Genesis), nil
	case SiriusNetwork:
		return []byte(siriusConsensus.Genesis), nil
	case StandaloneNetwork:
		return createConsensusGenesis(pv)
	default:
		return []byte{}, fmt.Errorf("invalid network selected %s", n)
	}
}*/


/*func (n Network) ConsensusProtocolBlockVersion() (version.Protocol, error) {
	switch n {
	case MainNetNetwork:
		return 10, nil
	case SiriusNetwork:
		return 9, nil
	case StandaloneNetwork:
		return version.BlockProtocol, nil
	default:
		return version.BlockProtocol, fmt.Errorf("invalid network selected %s", n)
	}
}
*/
const Genesis = `
{
    "config": {
        "chainId": 161,
        "eip150Block": 0,
        "eip155Block": 0,
        "eip158Block": 0,
        "ByzantiumBlock": 0,
        "ConstantinopleBlock": 0,
        "PetersburgBlock": 0
    },
	"nonce": "1",
    "difficulty": "1024",
    "gasLimit": "100000000",
    "alloc": {
        "0xc916cfe5c83dd4fc3c3b0bf2ec2d4e401782875e": {
            "balance": "300000000000000000000000000"
        }
    }
}
`
var Keystore = map[string]string{
	"UTC--2019-01-16T08-37-23.883Z--c916cfe5c83dd4fc3c3b0bf2ec2d4e401782875e":
	`{"version":3,"id":"54fc30ab-d1b7-4397-a39f-1ddd9110102b","address":"c916cfe5c83dd4fc3c3b0bf2ec2d4e401782875e","Crypto":{"ciphertext":"29beea264a484e72978286a64d3976060cd8dc74e57dd4d8a14ee12db2dfe6f9","cipherparams":{"iv":"6569b60cf13335a41340243b9ab67d9b"},"cipher":"aes-128-ctr","kdf":"scrypt","kdfparams":{"dklen":32,"salt":"a36f4967cd868fb5dec0380822dc014e748fb7c2bf968f29d6ad11c377189181","n":8192,"r":8,"p":1},"mac":"d55bb34a14abff2487ea78fa65d641981e094e576e6e3323e3550b2812081261"}}`,
}

func (n Network) DatabaseGenesis() ([]byte, error) {
	switch n {
	case StandaloneNetwork:
		return []byte(Genesis), nil
	default:
		return []byte{}, fmt.Errorf("invalid network selected %s", n)
	}
}

func (n Network) DatabaseKeystore() (map[string][]byte, error) {
	switch n {
	case StandaloneNetwork:
		var files = make(map[string][]byte)
		for name, keystore := range Keystore {
			files[name] = []byte(keystore)
		}
	
		return files, nil
	default:
		return map[string][]byte{}, fmt.Errorf("invalid network selected %s", n)
	}
}


/*func createConsensusGenesis(pv *privval.FilePV) ([]byte, error) {
	genDoc := types.GenesisDoc{
		ChainID:         fmt.Sprintf("test-chain-%v", tmtCommon.RandStr(6)),
		GenesisTime:     tmtime.Now(),
		ConsensusParams: types.DefaultConsensusParams(),
	}
	genDoc.Validators = []types.GenesisValidator{{
		Address: pv.GetPubKey().Address(),
		PubKey:  pv.GetPubKey(),
		Power:   10,
	}}
	
	genDocBytes, err := cdc.MarshalJSONIndent(genDoc, "", "  ")
	if err != nil {
		return nil, err
	}

	return genDocBytes, nil
}*/
