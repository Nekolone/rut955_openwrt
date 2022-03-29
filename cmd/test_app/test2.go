package main

import "log"

// your JSON structure as a byte slice
func main() {

	testV := make(map[string]string)
	testV["1223"] = "helo1223"

	log.Printf("ыыы = %s\n", testV["sss"])
	log.Printf("1223 = %s\n", testV["1223"])
}
