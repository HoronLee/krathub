package main

import (
	_ "embed"

	"github.com/go-kratos/kratos/v2/config"
)

//go:embed default.yaml
var defaultConfigData []byte

type memorySource struct {
	data []byte
}

func (s *memorySource) Load() ([]*config.KeyValue, error) {
	return []*config.KeyValue{
		{
			Key:    "default",
			Value:  s.data,
			Format: "yaml",
		},
	}, nil
}

func (s *memorySource) Watch() (config.Watcher, error) {
	return &staticWatcher{}, nil
}

type staticWatcher struct{}

func (w *staticWatcher) Next() ([]*config.KeyValue, error) {
	select {}
}

func (w *staticWatcher) Stop() error {
	return nil
}

func newDefaultConfigSource() config.Source {
	return &memorySource{data: defaultConfigData}
}
