package consul

import (
	capi "github.com/hashicorp/consul/api"
)

func newConsulClient() (*capi.Client, error) {
	return capi.NewClient(capi.DefaultConfig())
}

func GetKV(key string) ([]byte, error) {
	cli, err := newConsulClient()
	if err != nil {
		return nil, err
	}

	kv := cli.KV()

	get, _, err := kv.Get(key, nil)
	if err != nil {
		return nil, err
	}

	return get.Value, nil
}
