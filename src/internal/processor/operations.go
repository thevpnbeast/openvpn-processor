package processor

import (
	"database/sql"
	"log"
	"time"
)

func removeServers(db *sql.DB, ip string, proto string, confData string, port int) {
	stmt, err := db.Prepare(sqlDeleteServer)
	if err != nil {
		// TODO: do not panic, handle properly
		panic(err)
	}

	if _, err = stmt.Exec(ip, confData, proto, port); err != nil {
		log.Printf("ERROR: an error occured while executing query on database (error=%s dbIP=%s dbProto=%s "+
			"dbPort=%d", err.Error(), ip, proto, port)
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

	log.Printf("INFO: starting insert reachable server operation on database (serverCount=%d)", len(vpnServers))
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
			log.Printf("ERROR: an error occurred while preparing statement (query=%s)", sqlReplaceServer)
			return
		}

		var res sql.Result
		if res, err = stmt.Exec(values...); err != nil {
			log.Printf("ERROR: an error occurred while executing query on database (server=%s error=%s)", server.hostname, err.Error())
			failedServerCount++
			continue
		}

		if row, _ := res.RowsAffected(); row == 1 {
			insertedServerCount++
		}

		// clear the slice after all
		values = nil
	}

	log.Printf("INFO: ending insert reachable server operation on database (insertedServerCount=%d "+
		"skippedServerCount=%d failedServerCount=%d executionTime=%s)", insertedServerCount, skippedServerCount,
		failedServerCount, time.Since(beforeExecution).String())
}

func checkUnreachableServersOnDB(db *sql.DB) {
	log.Printf("INFO: starting remove unreachable server operation on database")
	var (
		removedServerCount  = 0
		port                int
		ip, confData, proto string
		beforeExecution     = time.Now()
	)

	var rows *sql.Rows
	if rows, err = db.Query(sqlSelectServers); err != nil {
		log.Printf("ERROR: an error occurred while querying database (query=%s error=%s)", sqlSelectServers, err.Error())
		return
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}()

	for rows.Next() {
		if err := rows.Scan(&ip, &proto, &confData, &port); err != nil {
			log.Printf("ERROR: an error occurred while scanning database (ip=%s proto=%s port=%d error=%s)",
				ip, proto, port, err.Error())
			continue
		}

		if !isServerInsertable(ip, proto, confData, port, opts.DialTcpTimeoutSeconds) {
			removedServerCount++
			removeServers(db, ip, proto, confData, port)
		}
	}

	log.Printf("INFO: ending remove unreachable server operation on database (removedServerCount=%d executionTime=%s)",
		removedServerCount, time.Since(beforeExecution).String())
}
