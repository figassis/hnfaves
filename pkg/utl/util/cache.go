package util

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"time"

	"github.com/figassis/hnfaves/pkg/utl/zaplog"
	//redis "gopkg.in/redis.v3"
	"github.com/go-redis/redis/v7"
)

const (
	cacheTTL   = time.Hour
	hostKeyTTL = time.Minute * 5
	ApiCache   = time.Minute * 20
)

var (
	cache    *redis.Client
	cacheKey string
	hostKey  string
)

func loadRedis(namespace, host, port string, db int, password string) (err error) {
	cacheKey = namespace
	if os.Getenv("REDIS_TLS") == "true" {
		cache = redis.NewClient(&redis.Options{Addr: host + ":" + port, Password: password, DB: db, MaxRetries: 3, TLSConfig: &tls.Config{InsecureSkipVerify: true}})
	} else {
		cache = redis.NewClient(&redis.Options{Addr: host + ":" + port, Password: password, DB: db, MaxRetries: 3})
	}
	setHostKey()
	return cache.Ping().Err()
}

func setHostKey() {
	hostKey, err := GenerateUUID()
	if err = zaplog.ZLog(err); err != nil {
		return
	}

	if err = zaplog.ZLog(CacheTTL("hostKey", hostKey, hostKeyTTL)); err != nil {
		return
	}
}

func HostKey() string {
	return hostKey
}

func GetCache(key string, result interface{}) (err error) {
	// zaplog.ZLog(fmt.Sprintf("Getting cache %s", key))
	if cache == nil {
		return errors.New("Redis not ready")
	}
	if key == "" {
		return errors.New("Cache key cannot be empty")
	}

	if cacheKey != "" {
		key = cacheKey + "_" + key
	}

	//logging.BLog(1, fmt.Errorf("KEY: %s", key), tid)
	if reflect.ValueOf(result).Kind() != reflect.Ptr {
		return errors.New("Please provide a pointer")
	}

	obj, err := cache.Get(key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return errors.New("Record not found")
		}
		return
	}

	switch string(obj) {
	case "[]", "{}", "":
		go DeleteCache(key)
		return errors.New("Cache miss")
	}

	if err = json.Unmarshal(obj, result); err != nil {
		return
	}

	// go cache.Set(key, obj, cacheTTL)
	return
}

func DeleteCache(key string) (err error) {
	if cache == nil {
		return errors.New("Redis not ready")
	}

	if cacheKey != "" {
		key = cacheKey + "_" + key
	}

	_, err = cache.Del(key).Result()
	if err == redis.Nil {
		return errors.New("Key not found")
	} else if err != nil {
		return errors.New("Could not remove key")
	}

	return
}

func Cache(key string, obj interface{}) (err error) {
	// zaplog.ZLog(fmt.Sprintf("Saving cache %s", key))
	if cache == nil {
		return errors.New("Redis not ready")
	}

	if key == "" {
		return errors.New("Key cannot be empty")
	}
	if cacheKey != "" {
		key = cacheKey + "_" + key
	}

	if obj == nil {
		return
	}

	jsonObj, err := json.Marshal(obj)
	if err != nil {
		return
	}
	return cache.Set(key, jsonObj, cacheTTL).Err()
}

func CacheTTL(key string, obj interface{}, ttl time.Duration) (err error) {

	if cache == nil {
		return errors.New("Redis not ready")
	}

	if cacheKey != "" {
		key = cacheKey + "_" + key
	}

	if obj == nil {
		return
	}

	jsonObj, err := json.Marshal(obj)
	if err != nil {
		return
	}
	return cache.Set(key, jsonObj, ttl).Err()
}

func New() (err error) {
	redisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"))
	if err != nil {
		redisDB = 0
	}

	if err = loadRedis(os.Getenv("FQDN"), os.Getenv("REDIS_HOST"), os.Getenv("REDIS_PORT"), redisDB, os.Getenv("REDIS_PASSWORD")); err != nil {
		zaplog.ZLog(fmt.Sprintf("Error loading redis %s", err.Error()))
		return
	}

	return
}
