package main

import (
	"io/ioutil"
)

type File struct {
	Name    string   `json:"filename"`
	Content string   `json:"content"`
	Rules   []string `json:"rules"`
}

type ActiveVersion struct {
	Action string `json:"action"`
}

func ReadFile(filepath string) (string, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(file), nil
}
