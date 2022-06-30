package txbypass

import "github.com/ethereum/go-ethereum/core/types"

type TX interface {
	BypassPoolTx(tx *types.Transaction)
}

//func BypassPoolTx(){
	//database.PoolTx()
//}