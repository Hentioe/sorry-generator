package main

import (
	"os"
)

func IsExist(path string) (bool, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, err
	} else {
		return true, err
	}
}

func IsAllExist(paths ...string) (bool, error) {
	for i := 0; i < len(paths); i++ {
		if exist, err := IsExist(paths[i]); !exist {
			return false, err
		}
	}
	return true, nil
}

func IfNotExistMkdir(path string, mode int32) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, os.FileMode(mode))
	} else if err != nil {
		return err
	}
	return nil
}

func IfNotExistMkAllMir(mode int32, paths ...string) error {
	for _, path := range paths {
		if err := IfNotExistMkdir(path, mode); err != nil {
			return err
		}
	}
	return nil
}
