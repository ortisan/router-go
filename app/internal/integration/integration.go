package integration

import (
	"context"
	"time"

	"github.com/ortisan/router-go/internal/config"
	errApp "github.com/ortisan/router-go/internal/error"
	"github.com/rs/zerolog/log"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	dialTimeout    = 2 * time.Second
	requestTimeout = 10 * time.Second
)

func GetEtcdCli() (context.Context, *clientv3.Client, error) {

	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)

	cli, err := clientv3.New(clientv3.Config{
		DialTimeout: 5 * time.Second,
		Endpoints:   config.ConfigObj.Etcd.Endpoints,
	})

	if err != nil {
		errW := errApp.NewIntegrationError("Error to connect with etcd server.", err)
		return nil, nil, errW
	}

	return ctx, cli, nil
}

func GetValues(key string) ([]string, error) {
	ctx, cli, err := GetEtcdCli()
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

	var values []string

	for _, kv := range gr.Kvs {
		values = append(values, string(kv.Value))
	}

	log.Debug().Str("key", key).Strs("values", values).Int64("revision", gr.Header.Revision).Msg("Values loaded from etcd.")

	return values, nil
}

func GetValuesPrefixed(keyPrefix string) (map[string]string, error) {
	ctx, cli, err := GetEtcdCli()
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

	mapKv := make(map[string]string)
	for _, kv := range gr.Kvs {
		mapKv[string(kv.Key)] = string(kv.Value)
	}
	return mapKv, nil
}

func PutValue(key string, value string) error {
	ctx, cli, err := GetEtcdCli()
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
	log.Debug().Str("key", key).Str("value", value).Int64("revision", revision).Msg("Value inserted...")

	return nil
}
