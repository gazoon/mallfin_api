package utils

import (
	"mallfin_api/redisdb"
	"time"

	"encoding/base64"
	"math/rand"

	"github.com/kataras/go-errors"
	"reflect"
	"runtime"
	"strings"
)

const (
	MAX_LOCK_TIME = 10 * time.Second
	DELAY_TIME    = 100 * time.Millisecond
	UNLOCK_SCRIPT = `
	if redis.call("get",KEYS[1]) == ARGV[1] then
		return redis.call("del",KEYS[1])
	else
		return 0
	end`
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type DistributedMutex struct {
	Resource string
	mutexId  string
}

func NewDistributedMutex(resource string) *DistributedMutex {
	return &DistributedMutex{Resource: resource}

}
func generateUniqueValue() (string, error) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	return base64.StdEncoding.EncodeToString(b), err
}
func (d *DistributedMutex) Lock() error {
	if d.mutexId != "" {
		return errors.New("already locked")
	}
	mutexId, err := generateUniqueValue()
	if err != nil {
		return err
	}
	redisConn := redisdb.GetConnection()
	for {
		setted, err := redisConn.SetNX(d.Resource, mutexId, MAX_LOCK_TIME).Result()
		if err != nil {
			return err
		}
		if setted {
			break
		}
		time.Sleep(DELAY_TIME)
	}
	d.mutexId = mutexId
	return nil
}
func (d *DistributedMutex) Unlock() error {
	if d.mutexId == "" {
		return errors.New("mutex hasn't locked yet")
	}
	redisConn := redisdb.GetConnection()
	err := redisConn.Eval(UNLOCK_SCRIPT, []string{d.Resource}, d.mutexId).Err()
	d.mutexId = ""
	return err
}

func FuncFullName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
func CurrentFuncFullName() string {
	pc, _, _, _ := runtime.Caller(1)
	return runtime.FuncForPC(pc).Name()
}
func FuncName(f interface{}) string {
	fullName := FuncFullName(f)
	return funcNameFromFullName(fullName)
}
func CurrentFuncName() string {
	pc, _, _, _ := runtime.Caller(1)
	fullName := runtime.FuncForPC(pc).Name()
	return funcNameFromFullName(fullName)
}
func funcNameFromFullName(fullName string) string {
	fullNameParts := strings.Split(fullName, ".")
	name := fullNameParts[len(fullNameParts)-1]
	if name == "" {
		name = fullName
	}
	return name
}
func MapKeys(m interface{}) []string {
	keysRaw := reflect.ValueOf(m).MapKeys()
	keys := make([]string, len(keysRaw))
	for i := range keysRaw {
		keys[i] = keysRaw[i].String()
	}
	return keys
}
