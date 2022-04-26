package main

import (
	"encoding/json"
	"os"
)

type Configuration struct {
	Hostnames []string `json:"hostnames"`
	MongoURL  string   `json:"mongo"`
}

func (c *Configuration) Init(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	if err := decoder.Decode(c); err != nil {
		return err
	}

	return nil
}
