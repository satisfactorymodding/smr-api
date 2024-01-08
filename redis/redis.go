package redis

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"time"

	"github.com/cespare/xxhash"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
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

	log.Info().Msg("Redis initialized")
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
	return result, errors.Wrap(err, "failed to get nonce")
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
		return errors.Wrap(err, "failed to marshal version upload state")
	}

	redisKey := "version:upload:state:" + versionID
	return errors.Wrap(client.Set(redisKey, string(marshaled), time.Minute*10).Err(), "failed to store version upload state")
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
		return nil, errors.Wrap(get.Err(), "failed to get version upload state")
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

func StoreModAssetList(modReference string, assets []string) {
	var buf bytes.Buffer
	zw := gzip.NewWriter(&buf)

	b, _ := json.Marshal(assets)
	_, _ = zw.Write(b)
	_ = zw.Close()

	client.Set(fmt.Sprintf("assets:%s", modReference), buf.Bytes(), time.Hour*24)
}

func GetModAssetList(modReference string) []string {
	result, err := client.Get(fmt.Sprintf("assets:%s", modReference)).Result()
	if err != nil {
		return nil
	}

	reader, err := gzip.NewReader(bytes.NewReader([]byte(result)))
	if err != nil {
		return nil
	}

	all, err := io.ReadAll(reader)
	if err != nil {
		return nil
	}

	out := make([]string, 0)
	_ = json.Unmarshal(all, &out)

	return out
}
