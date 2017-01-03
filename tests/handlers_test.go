package mallfin_api

import (
	"mallfin_api/db"
	"testing"
)

func TestNonExistentMallDetails(t *testing.T) {
	db.FlushDB()
}
func TestOkMallDetails(t *testing.T) {
	db.FlushDB()
}
