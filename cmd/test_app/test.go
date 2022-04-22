package main

import (
	"fmt"
	"log"
)

func main() {
	log.Println("test")
	log.Print(testF())
	log.Println("succ")

}

func testF() (answer string) {
	defer func() {
		if r := recover(); r != nil {
			log.Print("networkStatus -> buffering")
			answer = fmt.Sprint(r)
		}
	}()

	log.Print(1)
	log.Panicf("panicEND")
	return
}
