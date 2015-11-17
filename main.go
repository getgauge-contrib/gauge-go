package main

import (
	"flag"
	"fmt"
)

var start = flag.Bool("start", false, "Start go runner")
var initialize = flag.Bool("init", false, "Initialize go specs")

func main() {
	flag.Parse()

	if *start {
		startGo()
	} else if *initialize {
		initGo()
	} else {
		printUsage()
	}
}

func startGo(){
	fmt.Println("Start called")
}

func initGo(){
	fmt.Println("Init called")
}

func printUsage(){
	flag.PrintDefaults()
}