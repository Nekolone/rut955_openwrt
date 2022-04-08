package main

import (
	"log"
	"os"
	"sync"
	"time"
)

func main() {
	// LOG_FILE := "/overlay/wialon_rut955_gateway/test.log"
	LOG_FILE := "test.log"

	logFile, err := os.OpenFile(LOG_FILE, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Printf("cant open log file. ERR > %v", err)
	}

	defer logFile.Close()

	log.SetOutput(logFile)

	log.SetFlags(log.Lshortfile | log.LstdFlags)

	log.Println("test")

	var wg sync.WaitGroup

	wg.Add(1)

	in_func_test()
	go func() {
		defer wg.Done()
		log.Println("go")
	}()

	time.Sleep(1)
	wg.Wait()
	time.Sleep(1)

}

func in_routine_test(wg sync.WaitGroup) {
	log.Println("in routine test")
	wg.Done()
}
func in_func_test() {
	log.Println("in func test")
}
