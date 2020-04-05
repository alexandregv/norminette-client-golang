package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

type Norminette struct {
	Client *Client
	Files  []string
}

func (norm *Norminette) Init() {
	client := &Client{
		Hostname: "norminette.21-school.ru",
		Login:    "guest",
		Password: "guest",
		Version:  false}
	client.Init()
	norm.Client = client
	norm.ParseArgv()
	if norm.Client.Version {
		norm.Client.SendVersion()
		os.Exit(0)
	}
}

type arrayFlags []string

func (i *arrayFlags) String() string {
	return "my string representation"
}

func (i *arrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

func (norm *Norminette) ParseArgv() {
	var RulesFlags arrayFlags

	version := flag.Bool("v", false, "Prints version")
	flag.Var(&RulesFlags, "R", "Allows the use of special rules")
	if len(os.Args) == 1 {
		fmt.Println("Usage:")
		flag.PrintDefaults()
		os.Exit(0)
	}
	flag.Parse()
	norm.Client.Rules = append(norm.Client.Rules, RulesFlags...)
	norm.Client.Version = *version
}

func is_a_valid_file(filename string) bool {
	ext := filepath.Ext(filename)
	if ext == ".c" || ext == ".h" {
		return true
	} else {
		return false
	}
}

func is_not_dot_file(filename string) bool {
	if filename[0] != '.' {
		return true
	} else {
		return false
	}
}

func (norm *Norminette) FindFiles() {
	var RealFilepath string

	if len(os.Args) == 1 {
		RealFilepath = "."
	} else {
		if flag.Arg(0) == "" {
			RealFilepath = "."
		} else {
			RealFilepath = flag.Arg(0)
		}
	}
	for i := 0; RealFilepath != ""; i++ {
		err := filepath.Walk(RealFilepath,
			func(path string, info os.FileInfo, err error) error {
				if err != nil {
					if os.IsNotExist(err) {
						fmt.Println("File not found:", RealFilepath)
						return nil
					}
					return err
				}
				if stat, err := os.Stat(path); stat.IsDir() || os.IsNotExist(err) {
					return nil
				}
				if is_a_valid_file(path) {
					norm.Files = append(norm.Files, path)
				} else if is_not_dot_file(path) {
					fmt.Println("Norme: ./", path, "\nWarning: Not a valid file")
				}
				return nil
			})
		if err != nil {
			log.Println(err)
		}
		RealFilepath = flag.Arg(i + 1)
	}
}

func (norm *Norminette) SendFiles() {
	for _, value := range norm.Files {
		err := norm.Client.SendFile(value)
		if err != nil {
			log.Print(err)
		}
	}
}
