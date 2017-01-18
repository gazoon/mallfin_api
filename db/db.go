package db

import (
	"database/sql"
	"fmt"
	"github.com/jackc/pgx"
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
func createNewPool() *pgx.ConnPool {
	dbConf := config.Postgres()
	connPoolConfig := pgx.ConnPoolConfig{
		ConnConfig: pgx.ConnConfig{
			Host:     dbConf.Host,
			User:     dbConf.User,
			Password: dbConf.Password,
			Database: dbConf.Name,
			Port:     uint16(dbConf.Port),
		},
		MaxConnections: 20,
	}
	pool, err := pgx.NewConnPool(connPoolConfig)
	if err != nil {
		panic(err)
	}
	return pool
}

var pgxConn *pgx.ConnPool

func Initialization() {
	conn := createNewDBConnection(config.Postgres().Name)
	setNewConnection(conn)
	pgxConn = createNewPool()
}
func GetPgxConnection() *pgx.ConnPool {
	if pgxConn == nil {
		moduleLog.Panic("Postgres has not initialized yet")
	}
	return pgxConn
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
