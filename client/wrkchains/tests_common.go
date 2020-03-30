package wrkchains

import (
	"github.com/tendermint/tendermint/libs/log"
	"os"
)

var (
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
)

func contains(arr []string, str string) bool {
	for _, a := range arr {
		if a == str {
			return true
		}
	}
	return false
}
