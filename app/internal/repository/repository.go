package repository

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/go-redis/redis"
	"github.com/ortisan/router-go/internal/config"
	errApp "github.com/ortisan/router-go/internal/error"
	"github.com/rs/zerolog/log"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	dialTimeout    = 2 * time.Second
	requestTimeout = 10 * time.Second
)

func getEtcdCli() (context.Context, *clientv3.Client, error) {
	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)

	cli, err := clientv3.New(clientv3.Config{
		DialTimeout: dialTimeout,
		Endpoints:   config.ConfigObj.Etcd.Endpoints,
	})

	if err != nil {
		errW := errApp.NewIntegrationError("Error to connect with etcd server.", err)
		return nil, nil, errW
	}

	return ctx, cli, nil
}

func GetValues(key string) ([]string, error) {
	ctx, cli, err := getEtcdCli()

	if err != nil {
		return nil, err
	}
	defer cli.Close()

	log.Debug().Str("key", key).Msg("Trying get value in etcd...")

	gr, err := cli.Get(ctx, key)
	if err != nil {
		errW := errApp.NewIntegrationError("Error to connect with etcd server.", err)
		return nil, errW
	}

	if len(gr.Kvs) == 0 {
		return nil, errApp.NewNotFoundError(fmt.Sprintf("Key '%s' not found", key))
	}

	var values []string

	for _, kv := range gr.Kvs {
		values = append(values, string(kv.Value))
	}

	log.Debug().Str("key", key).Strs("values", values).Int64("revision", gr.Header.Revision).Msg("Values loaded from etcd.")

	return values, nil
}

func GetValuesPrefixed(keyPrefix string) (map[string]string, error) {
	ctx, cli, err := getEtcdCli()
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	log.Debug().Str("key", keyPrefix).Msg("Trying get value in etcd...")

	gr, err := cli.Get(ctx, keyPrefix, clientv3.WithPrefix())
	if err != nil {
		errW := errApp.NewIntegrationError("Error to connect with etcd server.", err)
		return nil, errW
	}

	if len(gr.Kvs) == 0 {
		return nil, errApp.NewNotFoundError(fmt.Sprintf("Key '%s' not found", keyPrefix))
	}

	mapKv := make(map[string]string)
	for _, kv := range gr.Kvs {
		mapKv[string(kv.Key)] = string(kv.Value)
	}
	return mapKv, nil
}

func PutValue(key string, value string) error {
	ctx, cli, err := getEtcdCli()
	if err != nil {
		return err
	}
	defer cli.Close()

	log.Debug().Str("key", key).Str("value", value).Msg("Trying to put value in etcd...")
	resp, err := cli.Put(ctx, key, value)

	if err != nil {
		errW := errApp.NewIntegrationError("Error to connect with etcd server.", err)
		return errW
	}

	revision := resp.Header.Revision
	log.Debug().Str("key", key).Str("value", value).Int64("revision", revision).Msg("Value inserted into etcd...")
	return nil
}

func getRedisCli() (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     config.ConfigObj.Redis.ServerAddress,
		Password: config.ConfigObj.Redis.Password,
		DB:       0,
	})
	return client, nil
}

func GetCacheValue(key string) (string, error) {
	cli, err := getRedisCli()

	if err != nil {
		return "", err
	}
	defer cli.Close()

	log.Debug().Str("key", key).Msg("Trying to get value from redis...")
	value, err := cli.Get(key).Result()

	if err != nil {
		if err == redis.Nil {
			return "", errApp.NewNotFoundError(fmt.Sprintf("Key '%s' not found", key))
		}
		return "", nil
	}

	log.Debug().Str("key", key).Str("value", value).Msg("Value obtained...")
	return value, nil
}

func PutCacheValue(key string, value string) (string, error) {
	cli, err := getRedisCli()

	if err != nil {
		return "", err
	}
	defer cli.Close()

	log.Debug().Str("key", key).Str("value", value).Msg("Trying to put value from redis...")
	result, err := cli.Set(key, value, 0*time.Second).Result()

	if err != nil {
		return "", err
	}

	log.Debug().Str("key", key).Str("value", value).Msg("Object inserted into redis...")
	return result, nil
}

func PutObject(bucket string, key string, value string) error {
	log.Debug().Str("key", key).Str("value", value).Msg("Trying to put object in S3...")

	svc := s3.New(config.NewAWSSession())
	_, err := svc.PutObject(&s3.PutObjectInput{
		Body:   strings.NewReader(value),
		Bucket: &bucket,
		Key:    &key,
	})

	log.Debug().Str("bucket", bucket).Str("key", key).Str("value", value).Msg("Object putted into S3...")
	return err
}

func GetObject(bucket string, key string) (string, error) {
	log.Debug().Str("bucket", bucket).Str("key", key).Msg("Trying to get object from S3...")

	svc := s3.New(config.NewAWSSession())
	out, err := svc.GetObject(&s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	if err != nil {
		return "", err
	}

	size := int(*out.ContentLength)
	buffer := make([]byte, size)
	defer out.Body.Close()
	var bbuffer bytes.Buffer
	for true {
		num, rerr := out.Body.Read(buffer)
		if num > 0 {
			bbuffer.Write(buffer[:num])
		} else if rerr == io.EOF || rerr != nil {
			break
		}
	}

	value := bbuffer.String()
	log.Debug().Str("bucket", bucket).Str("key", key).Str("value", value).Msg("Object readed from S3...")
	return value, nil
}
