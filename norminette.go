package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Norminette struct {
	Client	*Client
	Files	[]string
}

func (norm *Norminette)Init() {
	client := new(Client)
	client.Init()
	norm.Client = client
	norm.ParseArgv()
}

func (norm *Norminette)ParseArgv() {
}

func is_a_valid_file(filename string) bool {
	ext := filepath.Ext(filename)
	if ext == ".c" || ext == ".h" {
		return true
	} else {
		return false
	}
}

func (norm *Norminette)FindFiles() {

	err := filepath.Walk(os.Args[1],
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if stat, err := os.Stat(path); stat.IsDir() || os.IsNotExist(err) { return nil }
			if is_a_valid_file(path) {
				norm.Files = append(norm.Files, path)
			} else {
				fmt.Print(path, " Warning: Not a valid file")
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
}

func (norm *Norminette)SendFiles() {
	for _, value := range norm.Files {
		err := norm.Client.SendFile(value)
		if err != nil {
			log.Print(err)
		}
	}
}
