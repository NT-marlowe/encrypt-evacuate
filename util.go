package main

import (
	"os"
	"time"
)

func initDataShelter() (*os.File, error) {
	err := os.MkdirAll(dataShelterPath, 0766)
	if err != nil {
		return nil, err
	}
	file, err := os.CreateTemp(dataShelterPath, time.Now().Format(time.RFC3339)+"_")
	if err != nil {
		return nil, err
	}

	return file, nil
}
