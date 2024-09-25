package es

import (
	"bytes"
	"fmt"
	elastic "github.com/elastic/go-elasticsearch/v7"
	"github.com/elastic/go-elasticsearch/v7/esapi"
	"io"
	"library/consts"
	"strings"
)

func newEsClient() (*elastic.Client, error) {

	return elastic.NewClient(elastic.Config{
		Addresses: []string{consts.EsConf.Address},
	})
}

func resHandler(res *esapi.Response) ([]byte, error) {
	if res.StatusCode < 200 && res.StatusCode > 300 {
		return nil, fmt.Errorf("操作失败")
	}

	return io.ReadAll(res.Body)
}

func CreateIndex(index ...string) error {
	cli, err := newEsClient()
	if err != nil {
		return err
	}

	if len(index) == 0 {
		return fmt.Errorf("参数异常")
	}

	for _, v := range index {
		response, err := cli.Indices.Create(v)
		if err != nil {
			return err
		}
		_, err = resHandler(response)
		if err != nil {
			return err
		}
	}

	return nil
}

func SyncData(index string, id string, data []byte) error {
	cli, err := newEsClient()
	if err != nil {
		return err
	}

	response, err := cli.Index(index, bytes.NewReader(data), func(request *esapi.IndexRequest) {
		request.DocumentID = id
	})
	if err != nil {
		return err
	}

	_, err = resHandler(response)
	if err != nil {
		return err
	}
	return nil
}

func SearchData(index string, query string) ([]byte, error) {
	cli, err := newEsClient()
	if err != nil {
		return nil, err
	}

	response, err := cli.Search(
		cli.Search.WithIndex(index),
		cli.Search.WithBody(strings.NewReader(query)),
	)
	if err != nil {
		return nil, err
	}

	return resHandler(response)
}
