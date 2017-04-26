package sqlx

import (
	"database/sql"
	"runtime"
	"sync"
)

var (
	trackedStartedTransactions = map[*sql.Tx]string{}
	trackedTransactionsMutex   sync.Mutex
	stackTraceBuf              [1 << 16]byte
)

func trackStartTx(tx *sql.Tx) {
	trackedTransactionsMutex.Lock()

	n := runtime.Stack(stackTraceBuf[:], true)
	trackedStartedTransactions[tx] = string(stackTraceBuf[:n])

	trackedTransactionsMutex.Unlock()
}

func trackReleaseTx(tx *sql.Tx) {
	trackedTransactionsMutex.Lock()

	delete(trackedStartedTransactions, tx)

	trackedTransactionsMutex.Unlock()
}

func GetTrackedOpenTransactions() []string {
	var result []string
	trackedTransactionsMutex.Lock()
	for _, v := range trackedStartedTransactions {
		result = append(result, v)
	}
	trackedTransactionsMutex.Unlock()

	return result
}
