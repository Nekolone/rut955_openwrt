package main

import "log"

func main() {
	a := []byte{1, 16, 36, 64, 127, 128, 156, 253}

	for _, d := range a {
		log.Println(d)
		log.Println(int8(d))
	}

}
