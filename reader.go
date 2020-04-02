package main

import (
	"encoding/json"
	"io/ioutil"
)

type File struct {
	Name string `json:"filename"`
	Content  string `json:"content"`
}

func ReadFile(filepath string) (string, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return "", err
	}
	return string(file), nil
}

func RequestPreparation(filepath string) ([]byte, error) {
	content, err := ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	result, err := json.Marshal(File{Name: filepath, Content: content})
	return result, nil
}
