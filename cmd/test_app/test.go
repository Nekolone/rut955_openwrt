package main

import (
	"log"
	"strconv"
)

func main() {
	log.Println(runMaShit())
}

func runMaShit() int {
	//out = byte arr
	strout := string("5")
	intOut, _ := strconv.Atoi(strout)
	return intOut
}
