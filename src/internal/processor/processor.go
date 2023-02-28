package processor

import (
	"database/sql"
	"github.com/thevpnbeast/openvpn-processor/src/internal/options"
	"log"
)

var (
	opts *options.OpenvpnProcessorOptions
	db   *sql.DB
	err  error
)

func init() {
	opts = options.GetOpenvpnProcessorOptions()
	//db = initDb()
}

func ProcessEventHandler() {
	log.Println("INFO: starting scheduler execution")
	//beforeMainExecution := time.Now()
	//csvContent := getCsvContent(opts.VpnGateUrl)
	//vpnServers := createStructsFromCsv(csvContent)
	//checkUnreachableServersOnDB(db)
	//insertServers(db, vpnServers)
	//log.Printf("INFO: ending scheduler execution, executionTime=%v", time.Since(beforeMainExecution))
	log.Println("INFO: ending scheduler execution")
}
