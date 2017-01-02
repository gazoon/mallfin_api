package db

import (
	"database/sql"
	"fmt"
	"mallfin_api/config"

	"os/exec"
	"strconv"

	log "github.com/Sirupsen/logrus"
	_ "github.com/lib/pq"
)

var (
	db        *sql.DB
	moduleLog = log.WithField("location", "postgres")
)

func createNewDBConnection(dbName string) *sql.DB {
	dbConf := config.Postgres()
	conn, err := sql.Open("postgres", fmt.Sprintf("dbname=%s user=%s password=%s host=%s port=%d", dbName, dbConf.User, dbConf.Password, dbConf.Host, dbConf.Port))
	if err == nil {
		err = conn.Ping()
	}
	if err != nil {
		moduleLog.WithFields(log.Fields{"conf": dbConf, "db": dbName}).Panicf("Cannot connect to postgresql: %s", err)
	}
	return conn
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
func Initialization() {
	conn := createNewDBConnection(config.Postgres().Name)
	setNewConnection(conn)
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
	conn := createNewDBConnection(dbConf.Name)
	_, err := conn.Exec(fmt.Sprintf(`DROP DATABASE IF EXISTS %s`, testDBName))
	if err != nil {
		moduleLog.Panicf("Cannot drop previous test db: %s", err)
	}
	_, err = conn.Exec(fmt.Sprintf(`CREATE DATABASE %s OWNER %s`, testDBName, dbConf.User))
	if err != nil {
		moduleLog.Panicf("Cannot crate test db: %s", err)
	}
	conn.Close()
	conn = createNewDBConnection(testDBName)
	_, err = conn.Exec(string(dump))
	if err != nil {
		moduleLog.Panicf("Cannot load dump: %s", err)
	}
	setNewConnection(conn)
}
func Close() {
	closeCurrentConn()
}
