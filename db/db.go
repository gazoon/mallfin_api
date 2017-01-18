package db

import (
	"database/sql"
	"fmt"
	"mallfin_api/config"
	"os/exec"
	"strconv"

	log "github.com/Sirupsen/logrus"
	_ "github.com/gazoon/pq"
	"strings"
)

var (
	db        *sql.DB
	moduleLog = log.WithField("location", "postgres")
)

func createNewDBConnection(dbName string) *sql.DB {
	dbConf := config.Postgres()
	conn, err := sql.Open("postgres", fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d sslmode=%s", dbName, dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port, dbConf.SSL))
	if err == nil {
		err = conn.Ping()
	}
	if err != nil {
		moduleLog.WithFields(log.Fields{"conf": dbConf, "db": dbName}).Panicf("Cannot connect to postgresql: %s", err)
	}
	return conn
}
func getAllTables() []string {
	const POSTGIS_TABLE_PREFIX = "spatial"
	conn := GetConnection()
	rows, err := conn.Query(`
	SELECT tablename
	FROM pg_tables
	WHERE schemaname = 'public'`)
	if err != nil {
		moduleLog.Panic(err)
	}
	defer rows.Close()
	var tables []string
	for rows.Next() {
		var tableName string
		err := rows.Scan(&tableName)
		if err != nil {
			moduleLog.Panic(err)
		}
		if !strings.HasPrefix(tableName, POSTGIS_TABLE_PREFIX) {
			tables = append(tables, tableName)
		}
	}
	err = rows.Err()
	if err != nil {
		moduleLog.Panic(err)
	}
	return tables
}
func FlushDB() {
	conn := GetConnection()
	tables := getAllTables()
	_, err := conn.Exec(fmt.Sprintf(`
	TRUNCATE %s CASCADE`, strings.Join(tables, ",")))
	if err != nil {
		moduleLog.Panicf("Cannot drop tables: %", err)
	}
}
func dbDump() []byte {
	dbConf := config.Postgres()
	cmd := exec.Command("pg_dump", "-h", dbConf.Host, "-p", strconv.Itoa(dbConf.Port), "-U", dbConf.User, "-d", dbConf.Name,
		"--schema-only", "--no-owner", "--no-privileges")

	cmd.Env = []string{fmt.Sprintf("PGPASSWORD=%s", dbConf.Password)}
	cmdOutput, err := cmd.CombinedOutput()
	if err != nil {
		log.Panicf("Cannot make the dump: %s", cmdOutput)
	}
	return cmdOutput
}
func closeCurrentConn() {
	if db != nil {
		db.Close()
	}
}
func setNewConnection(conn *sql.DB) {
	closeCurrentConn()
	db = conn
}

var shopStmt *sql.Stmt

func Initialization() {
	conn := createNewDBConnection(config.Postgres().Name)
	var err error
	shopStmt, err = conn.Prepare(`
		SELECT
		  s.id,
		  s.name,
		  s.logo_small,
		  s.logo_large,
		  s.score,
		  s.malls_count,
		  s.phone,
		  s.site,
		  m.id             mall_id,
		  m.name           mall_name,
		  m.phone          mall_phone,
		  m.logo_small     mall_logo_small,
		  m.logo_large     mall_logo_large,
		  ST_X(m.location) mall_location_lat,
		  ST_Y(m.location) mall_location_lon,
		  m.shops_count    mall_shops
		FROM shop s
		  JOIN mall_shop ms ON s.id = ms.shop_id
		  JOIN mall m ON ms.mall_id = m.id
		WHERE s.id = $1
		ORDER BY m.location <-> ST_SetSRID(ST_Point($2, $3), 4326)
		LIMIT 1
	`)
	if err != nil {
		panic(err)
	}
	setNewConnection(conn)
}
func GetShopStmt() *sql.Stmt {
	if shopStmt == nil {
		moduleLog.Panic("Postgres has not initialized yet")
	}
	return shopStmt
}
func GetConnection() *sql.DB {
	if db == nil {
		moduleLog.Panic("Postgres has not initialized yet")
	}
	return db
}
func InitializationForTests() {
	dump := dbDump()
	dbConf := config.Postgres()
	testDBName := fmt.Sprintf("%s_test", dbConf.Name)
	tmpConn := createNewDBConnection(dbConf.Name)
	_, err := tmpConn.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS %s`, testDBName))
	if err != nil {
		moduleLog.Panicf("Cannot drop previous test db: %s", err)
	}
	_, err = tmpConn.Exec(fmt.Sprintf(`CREATE DATABASE %s OWNER %s`, testDBName, dbConf.User))
	if err != nil {
		moduleLog.Panicf("Cannot create test db: %s", err)
	}
	tmpConn.Close()
	conn := createNewDBConnection(testDBName)
	_, err = conn.Exec(string(dump))
	if err != nil {
		moduleLog.Panicf("Cannot load dump: %s", err)
	}
	setNewConnection(conn)
}
func Close() {
	closeCurrentConn()
}
