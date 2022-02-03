package integration

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ortisan/router-go/config"
	"github.com/rs/zerolog/log"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	dialTimeout    = 2 * time.Second
	requestTimeout = 10 * time.Second
)

func GetEtcdCli() (context.Context, *clientv3.Client, error) {
	config, _ := config.LoadConfig(".")

	ctx, _ := context.WithTimeout(context.Background(), requestTimeout)

	cli, err := clientv3.New(clientv3.Config{
		DialTimeout: 5 * time.Second,
		Endpoints:   strings.Split(config.ETCDServers, ","),
	})

	if err != nil {
		log.Error().Stack().Err(err).Msg("Error to connect with ETCD...")
		return nil, nil, err
	}

	return ctx, cli, nil
}

func GetValue(key string) (string, error) {
	ctx, cli, err := GetEtcdCli()
	if err != nil {
		return "", err
	}
	defer cli.Close()

	log.Debug().Str("key", key).Msg("Trying get value in etcd...")
	gr, _ := cli.Get(ctx, key)
	value := string(gr.Kvs[0].Value)
	fmt.Println("Value: ", string(gr.Kvs[0].Value), "Revision: ", gr.Header.Revision)
	log.Debug().Str("key", key).Str("value", value).Int64("revision", gr.Header.Revision).Msg("Value loaded from etcd.")
	return value, nil
}

func PutValue(key string, value string) error {
	ctx, cli, err := GetEtcdCli()
	if err != nil {
		return err
	}
	defer cli.Close()

	resp, err := cli.Put(ctx, key, value)
	if resp != nil {
		log.Debug().Str("key", key).Str("value", value).Msg("Value inserted...")
	}

	return err

}
