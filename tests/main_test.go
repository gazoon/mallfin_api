package mallfin_api

import (
	"mallfin_api/config"
	"mallfin_api/db"
	"mallfin_api/redisdb"
	"testing"
)

func TestMain(m *testing.M) {
	config.Initialization()

	redisdb.InitializationForTests()
	defer redisdb.Close()

	db.InitializationForTests()
	defer db.Close()
	m.Run()
}
