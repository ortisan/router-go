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
		log.Error().Stack().Err(err).Msg(errW.(errApp.IntegrationError).ErrorSt.Error())
		return nil, nil, errW
	}

	return ctx, cli, nil
}

func GetValue(key string) ([]string, error) {
	ctx, cli, err := GetEtcdCli()
	if err != nil {
		return nil, err
	}
	defer cli.Close()

	log.Debug().Str("key", key).Msg("Trying get value in etcd...")

	gr, err := cli.Get(ctx, key)
	if err != nil {
		errW := errApp.NewIntegrationError("Error to connect with etcd server.", err)
		log.Error().Stack().Err(err).Msg(errW.(errApp.IntegrationError).ErrorSt.Error())
		return nil, errW
	}

	if len(gr.Kvs) > 0 {
		log.Debug().Str("key", key).Msg("Not found value for this key.")
		return nil, nil
	} else {
		value := string(gr.Kvs[0].Value)
		var values []string
		for key, value in

		log.Debug().Str("key", key).Str("value", value).Int64("revision", gr.Header.Revision).Msg("Value loaded from etcd.")

	}
	return value, nil
}

func GetValues(keyPrefix string) (string, error) {
	ctx, cli, err := GetEtcdCli()
	if err != nil {
		return "", err
	}
	defer cli.Close()

	log.Debug().Str("key", keyPrefix).Msg("Trying get value in etcd...")

	gr, err := cli.Get(ctx, keyPrefix, clientv3.WithPrefix())
	if err != nil {
		errW := errApp.NewIntegrationError("Error to connect with etcd server.", err)
		log.Error().Stack().Err(err).Msg(errW.(errApp.IntegrationError).ErrorSt.Error())
		return "", errW
	}
	value := string(gr.Kvs[0].Value)

	log.Debug().Str("key", keyPrefix).Str("value", value).Int64("revision", gr.Header.Revision).Msg("Value loaded from etcd.")

	return value, nil
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
		log.Error().Stack().Err(err).Msg(errW.(errApp.IntegrationError).ErrorSt.Error())
		return errW
	}

	revision := resp.Header.Revision
	log.Debug().Str("key", key).Str("value", value).Int64("revision", revision).Msg("Value inserted...")

	return nil
}
