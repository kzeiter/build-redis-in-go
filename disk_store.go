package main

import (
	"encoding/gob"
	"os"
	"sync"
)

type diskStore struct {
	mu       sync.RWMutex
	filename string
}

func (d *diskStore) save(data map[string]string) error {
	d.mu.Lock()
	defer d.mu.Unlock()

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

func (d *diskStore) load() (map[string]string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	file, err := os.Open(d.filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decoder := gob.NewDecoder(file)

	data := make(map[string]string)

	err = decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	return data, nil
}
