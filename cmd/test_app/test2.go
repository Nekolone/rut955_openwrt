package main

import "log"

func main() {
	a := []byte{65, 66, 67}
	log.Println(a)
	log.Println(string(a))
}
