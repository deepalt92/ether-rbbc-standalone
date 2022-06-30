package dbft

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/core/types"
	"time"
)

type txManager struct {
	txpool    []*types.Transaction
	threshold int
}

func MarshalJSON() ([]byte, error) {
	ts := time.Now()
	stamp := fmt.Sprint(ts)

	return []byte(stamp), nil
}

func (txm *txManager) add(tx *types.Transaction) {
	txm.txpool = append(txm.txpool, tx)
}

func (txm *txManager) check() bool {
	return len(txm.txpool) >= txm.threshold
}

func (txm *txManager) size() int {
	return len(txm.txpool)
}

func (txm *txManager) serialize() []byte {
	data, _ := json.Marshal(txm.txpool)
	txm.txpool = make([]*types.Transaction, 0, 2000)
	return data

}


func (txm *txManager) deserialize(data []byte) []*types.Transaction {
	var txs []*types.Transaction
	err := json.Unmarshal(data, &txs)
	//err := UnmarshalJSON
	if err != nil {

	}
	return txs
}
