package main

import (
	"encoding/gob"
	"os"
)

type diskStore struct {
	filename string
}

func (d *diskStore) save(data map[string]interface{}) error {
	file, err := os.Create(d.filename)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	err = encoder.Encode(data)
	if err != nil {
		return err
	}

	return nil
}

func (d *diskStore) load() (map[string]interface{}, error) {
	file, err := os.Open(d.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)
	var data map[string]interface{}
	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
