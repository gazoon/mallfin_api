package utils

import (
	"mallfin_api/redisdb"
	"time"

	log "github.com/Sirupsen/logrus"
)

const MAX_LOCK_TIME = 1000 * 100 * time.Millisecond
const SLEEP_TIME = time.Millisecond

type DistributedMutex struct {
	Resource string
}

func NewDistributedMutex(resource string) *DistributedMutex {
	return &DistributedMutex{Resource: resource}

}
func (d *DistributedMutex) Lock() {
	redisConn := redisdb.GetConnection()
	for {
		setted, err := redisConn.SetNX(d.Resource, 2, MAX_LOCK_TIME).Result()
		if err != nil {
			log.WithFields(log.Fields{"location": "distributed lock", "redis_key": d.Resource}).Panic("Cannot set: %s", err)
		}
		if setted {
			break
		}
		time.Sleep(SLEEP_TIME)
		//runtime.Gosched()
	}
	log.Info("lock acquired")
}
func (d *DistributedMutex) Unlock() {

}
