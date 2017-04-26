package sqlx

import (
	"database/sql"
	"runtime"
	"strings"
	"sync"
)

var (
	trackedStartedTransactions = map[*sql.Tx]string{}
	trackedTransactionsMutex   sync.Mutex
	stackTraceBuf              [4096]byte
)

func trackStartTx(tx *sql.Tx) {
	trackedTransactionsMutex.Lock()

	n := runtime.Stack(stackTraceBuf[:], true)

	s := string(stackTraceBuf[:n])
	p := strings.Index(s, "\n\n")
	if p != -1 {
		s = s[:p]
	}

	lines := strings.Split(s, "\n")
	if len(lines) > 5 {
		lines = lines[5:]
		s = strings.Join(lines, "\n")
	}

	trackedStartedTransactions[tx] = s

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
