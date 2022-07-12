package processor

import (
	"database/sql"
	"go.uber.org/zap"
	"time"
)

func removeServers(db *sql.DB, ip string, proto string, confData string, port int) {
	stmt, err := db.Prepare(sqlDeleteServer)
	if err != nil {
		// TODO: do not panic, handle properly
		panic(err)
	}

	if _, err = stmt.Exec(ip, confData, proto, port); err != nil {
		logger.Fatal("fatal error occurred while executing query on database", zap.String("ip", ip),
			zap.String("proto", proto), zap.Int("port", port), zap.String("error", err.Error()))
	}
}

func insertServers(db *sql.DB, vpnServers []vpnServer) {
	var (
		insertedServerCount = 0
		skippedServerCount  = 0
		failedServerCount   = 0
		beforeExecution     = time.Now()
		values              []interface{}
		stmt                *sql.Stmt
		err                 error
	)

	logger.Info("starting insert reachable server operation on database", zap.Int("serverCount", len(vpnServers)))
	for index, server := range vpnServers {
		if !isServerInsertable(server.ip, server.proto, server.confData, server.port, opts.DialTcpTimeoutSeconds) {
			skippedServerCount++
			continue
		}

		values = append(values, index+1, server.uuid, server.hostname, server.ip, server.port, server.confData,
			server.proto, server.enabled, server.score, server.ping, server.speed, server.countryLong,
			server.countryLong, server.numVpnSessions, server.uptime, server.totalUsers, server.totalTraffic,
			server.createdAt)

		if stmt, err = db.Prepare(sqlReplaceServer); err != nil {
			logger.Fatal("fatal error occured while preparing statement", zap.String("query", sqlReplaceServer))
			return
		}

		var res sql.Result
		if res, err = stmt.Exec(values...); err != nil {
			logger.Error("an error occurred while executing query on database", zap.String("server", server.hostname),
				zap.String("query", sqlReplaceServer), zap.String("error", err.Error()))
			failedServerCount++
			continue
		}

		if row, _ := res.RowsAffected(); row == 1 {
			insertedServerCount++
		}

		// clear the slice after all
		values = nil
	}

	logger.Info("Ending insert reachable server operation on database", zap.Int("insertedServerCount", insertedServerCount),
		zap.Int("skippedServerCount", skippedServerCount), zap.Int("failedServerCount", failedServerCount),
		zap.Duration("executionTime", time.Since(beforeExecution)))
}

func checkUnreachableServersOnDB(db *sql.DB) {
	logger.Info("starting remove unreachable server operation on database")
	var (
		removedServerCount  = 0
		port                int
		ip, confData, proto string
		beforeExecution     = time.Now()
	)

	var rows *sql.Rows
	if rows, err = db.Query(sqlSelectServers); err != nil {
		logger.Fatal("fatal error occurred while querying database", zap.String("query", sqlSelectServers),
			zap.String("error", err.Error()))
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}()

	for rows.Next() {
		if err := rows.Scan(&ip, &proto, &confData, &port); err != nil {
			logger.Fatal("fatal error occurred while scanning database", zap.String("ip", ip),
				zap.String("proto", proto), zap.Int("port", port), zap.String("error", err.Error()))
		}

		if !isServerInsertable(ip, proto, confData, port, opts.DialTcpTimeoutSeconds) {
			removedServerCount++
			removeServers(db, ip, proto, confData, port)
		}
	}

	logger.Info("Ending remove unreachable server operation on database", zap.Int("serverCount", removedServerCount),
		zap.Duration("executionTime", time.Since(beforeExecution)))
}
