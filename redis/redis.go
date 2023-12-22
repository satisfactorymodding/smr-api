package redis

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/Vilsol/slox"
	"github.com/cespare/xxhash"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/spf13/viper"

	"github.com/satisfactorymodding/smr-api/generated"
)

var client *redis.Client

func InitializeRedis(ctx context.Context) {
	client = redis.NewClient(&redis.Options{
		Addr:     viper.GetString("database.redis.host") + ":" + fmt.Sprint(viper.GetInt("database.redis.port")),
		Password: viper.GetString("database.redis.pass"),
		DB:       viper.GetInt("database.redis.db"),
	})

	ping := client.Ping()

	if ping == nil {
		panic("Redis not reachable")
	}

	if ping.Err() != nil {
		panic(ping.Err())
	}

	slox.Info(ctx, "Redis initialized")
}

func CanIncrement(ip string, action string, object string, expiration time.Duration) bool {
	data := make([]byte, 8)
	binary.LittleEndian.PutUint64(data, xxhash.Sum64String(ip))
	key := object + ":" + action + ":" + base64.URLEncoding.EncodeToString(data)
	return client.SetNX(key, true, expiration).Val()
}

func StoreNonce(nonce string, redirectURI string) {
	client.Set("nonce:"+nonce, redirectURI, time.Minute*10)
}

func GetNonce(nonce string) (string, error) {
	result, err := client.Get("nonce:" + nonce).Result()
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %w", err)
	}
	return result, nil
}

func GetAllKeys() []string {
	return client.Keys("*").Val()
}

func StoreMultipartCompletedPart(key string, etag string, part int) {
	encodedKey := base64.RawStdEncoding.EncodeToString([]byte(key))
	redisKey := "s3:uploads:part:" + encodedKey
	client.HMSet(redisKey, map[string]interface{}{
		strconv.Itoa(part): etag,
	})
	client.Expire(redisKey, time.Minute*60)
}

func GetAndClearMultipartCompletedParts(key string) map[string]string {
	encodedKey := base64.RawStdEncoding.EncodeToString([]byte(key))
	all := client.HGetAll("s3:uploads:part:" + encodedKey)
	client.Del("s3:uploads:part:" + encodedKey)
	return all.Val()
}

func StoreMultipartUploadID(key string, id string) {
	encodedKey := base64.RawStdEncoding.EncodeToString([]byte(key))
	redisKey := "s3:uploads:part:" + encodedKey + ":id"
	client.Set(redisKey, id, time.Minute*60)
}

func GetMultipartUploadID(key string) string {
	encodedKey := base64.RawStdEncoding.EncodeToString([]byte(key))
	redisKey := "s3:uploads:part:" + encodedKey + ":id"
	return client.Get(redisKey).Val()
}

type StoredVersionUploadState struct {
	Data *generated.CreateVersionResponse `json:"data"`
	Err  string                           `json:"err"`
}

func StoreVersionUploadState(versionID string, data *generated.CreateVersionResponse, err error) error {
	state := StoredVersionUploadState{
		Data: data,
	}

	if err != nil {
		state.Err = err.Error()
	}

	marshaled, e := json.Marshal(state)

	if e != nil {
		return fmt.Errorf("failed to marshal version upload state: %w", err)
	}

	redisKey := "version:upload:state:" + versionID
	if err := client.Set(redisKey, string(marshaled), time.Minute*10).Err(); err != nil {
		return fmt.Errorf("failed to store version upload state: %w", err)
	}

	return nil
}

func GetVersionUploadState(versionID string) (*generated.CreateVersionResponse, error) {
	redisKey := "version:upload:state:" + versionID
	get := client.Get(redisKey)

	if get == nil {
		return nil, nil
	}

	if get.Err() != nil {
		if errors.Is(get.Err(), redis.Nil) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get version upload state: %w", get.Err())
	}

	data := &StoredVersionUploadState{}
	_ = json.Unmarshal([]byte(get.Val()), data)

	if data.Err != "" {
		return data.Data, errors.New(data.Err)
	}

	return data.Data, nil
}

func FlushRedis() {
	client.FlushDB()
}
