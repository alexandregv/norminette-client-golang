package main

import "time"

func main() {
	norm := new(Norminette)
	norm.Init()
	norm.FindFiles()
	norm.SendFiles()
	go norm.Client.PrintResult()
	for norm.Client.Count != 0 {
		time.Sleep(time.Second)
	}
}
