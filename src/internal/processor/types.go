package processor

//const (
//	sqlSelectServers = "SELECT ip, proto, conf_data, port FROM servers"
//	sqlReplaceServer = "INSERT INTO servers(id, uuid, hostname, ip, port, conf_data, proto, enabled, score, ping, " +
//		"speed, country_long, country_short, num_vpn_sessions, uptime, total_users, total_traffic, created_at) VALUES " +
//		"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE id=id"
//	sqlDeleteServer = "DELETE FROM servers WHERE ip=? AND conf_data=? AND proto=? AND port=?"
//)
//
//type vpnServer struct {
//	hostname, uuid, ip, proto, countryLong, countryShort, confData             string
//	port, score, ping, speed, numVpnSessions, uptime, totalUsers, totalTraffic int
//	enabled                                                                    bool
//	createdAt                                                                  time.Time
//}
