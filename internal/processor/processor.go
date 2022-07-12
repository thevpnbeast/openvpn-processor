package processor

import (
	"database/sql"
	commons "github.com/thevpnbeast/golang-commons"
	"github.com/thevpnbeast/openvpn-processor/internal/options"
	"go.uber.org/zap"
	"time"
)

var (
	logger *zap.Logger
	opts   *options.OpenvpnProcessorOptions
	db     *sql.DB
	err    error
)

func init() {
	logger = commons.GetLogger()
	opts = options.GetOpenvpnProcessorOptions()
	db = initDb()
}

func ProcessEventHandler() {
	logger.Info("starting scheduler execution")
	beforeMainExecution := time.Now()
	csvContent := getCsvContent(opts.VpnGateUrl)
	vpnServers := createStructsFromCsv(csvContent)
	checkUnreachableServersOnDB(db)
	insertServers(db, vpnServers)
	logger.Info("ending scheduler execution", zap.Duration("executionTime", time.Since(beforeMainExecution)))
}
