package processor

import (
	"database/sql"
	"encoding/base64"
	"encoding/csv"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
	"time"
)

// initDb gets parameters and initiate database connection, returns connection then
func initDb() *sql.DB {
	db, err = sql.Open(opts.DbDriver, opts.DbUrl)
	if err != nil {
		log.Printf("FATAL: fatal error occurred while opening database connection (error=%s)", err.Error())
		panic(err)
	}
	tuneDbPooling(db)
	return db
}

// Read on https://www.alexedwards.net/blog/configuring-sqldb for detailed explanation
func tuneDbPooling(db *sql.DB) {
	// Set the maximum number of concurrently open connections (in-use + idle)
	// to 5. Setting this to less than or equal to 0 will mean there is no
	// maximum limit (which is also the default setting).
	db.SetMaxOpenConns(opts.DbMaxOpenConn)
	// Set the maximum number of concurrently idle connections to 5. Setting this
	// to less than or equal to 0 will mean that no idle connections are retained.
	db.SetMaxIdleConns(opts.DbMaxIdleConn)
	// Set the maximum lifetime of a connection to 1 hour. Setting it to 0
	// means that there is no maximum lifetime and the connection is reused
	// forever (which is the default behavior).
	db.SetConnMaxLifetime(time.Duration(int32(opts.DbConnMaxLifetimeMin)) * time.Minute)
}

func createStructsFromCsv(csvContent [][]string) []vpnServer {
	var vpnServers []vpnServer
	for _, entry := range csvContent {
		server := vpnServer{
			uuid:           uuid.New().String(),
			hostname:       entry[0],
			score:          cast.ToInt(entry[2]),
			ping:           cast.ToInt(entry[3]),
			speed:          cast.ToInt(entry[4]),
			countryLong:    entry[5],
			countryShort:   entry[6],
			numVpnSessions: cast.ToInt(entry[7]),
			uptime:         cast.ToInt(entry[8]),
			totalUsers:     cast.ToInt(entry[9]),
			totalTraffic:   cast.ToInt(entry[10]),
			enabled:        true,
			createdAt:      time.Now(),
		}

		decodedByteSlice, err := base64.StdEncoding.DecodeString(entry[14])
		if err != nil {
			log.Printf("WARN: an error occurred while decoding conf data, skipping...")
			continue
		}

		decodedConfData := string(decodedByteSlice)
		server.confData = decodedConfData
		for _, line := range strings.Split(decodedConfData, "\n") {
			fields := strings.Fields(line)
			if strings.HasPrefix(line, "remote") {
				server.ip = fields[1]
				server.port = cast.ToInt(fields[2])
			}

			if strings.HasPrefix(line, "proto") {
				server.proto = fields[1]
			}
		}
		vpnServers = append(vpnServers, server)
	}

	log.Printf("INFO: successfully created structs from csv (structsCreated=%d)", len(vpnServers))
	return vpnServers
}

func getCsvContent(vpnGateUrl string) [][]string {
	log.Printf("INFO: getting server list from API (url=%s)", vpnGateUrl)
	var csvContent [][]string
	resp, err := http.Get(vpnGateUrl)
	if err != nil {
		log.Printf("ERROR: an error occurred while making GET request (url=%s)", vpnGateUrl)
		return nil
	}

	defer func() {
		err := resp.Body.Close()
		if err != nil {
			panic(err)
		}
	}()

	encodedBody, err := ioutil.ReadAll(resp.Body)
	decodedBody := string(encodedBody)
	if err != nil {
		log.Printf("ERROR: an error occurred while reading response body (vpnGateUrl=%s error=%s)", vpnGateUrl, err.Error())
		return nil
	}
	reader := csv.NewReader(strings.NewReader(decodedBody))
	for {
		server, err := reader.Read()
		if err == io.EOF {
			break
		}

		if !strings.HasPrefix(server[0], "*") && !strings.HasPrefix(server[0], "#") {
			csvContent = append(csvContent, server)
		}
	}
	return csvContent
}

func isServerInsertable(ip, proto, confData string, port int, timeoutSeconds int) bool {
	isReachable := true
	timeout := time.Duration(int32(timeoutSeconds)) * time.Second
	if _, err := net.DialTimeout(proto, fmt.Sprintf("%s:%d", ip, port), timeout); err != nil {
		isReachable = false
	}

	isUnauthenticated := strings.Contains(confData, "#auth-user-pass")
	return isReachable && isUnauthenticated
}
