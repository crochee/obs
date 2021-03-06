// Copyright 2021, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2021/3/3

package config

import (
	"os"

	"github.com/json-iterator/go"
)

type Json struct {
	path string
}

func (j Json) Decode() (*ServiceConfig, error) {
	file, err := os.Open(j.path)
	if err != nil {
		return nil, err
	}
	defer file.Close() // #nosec G307
	var config ServiceConfig
	if err = jsoniter.ConfigFastest.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

func (j Json) Encode(c *ServiceConfig) error {
	file, err := os.Create(j.path)
	if err != nil {
		return err
	}
	defer file.Close() // #nosec G307
	return jsoniter.ConfigFastest.NewEncoder(file).Encode(c)
}
