package main

import (
	"bufio"
	"log"
	"os"
	"os/exec"
	"strconv"
)

func main() {
	exec.Command("mv main mainnnnnn")
}

func readByLines(networkStatus *string) {
	fileHanler, err := os.OpenFile("buffer.buf", os.O_RDONLY, os.ModePerm)
	if err != nil {
		log.Println("cant read buffer file")
		return
	}
	defer fileHanler.Close()

	fileScanner := bufio.NewScanner(fileHanler)
	var msg string
	for fileScanner.Scan() {
		msg = fileScanner.Text()
		log.Println(msg)
		switch *networkStatus {
		case "postBuffering":
			log.Printf("send msg>%s", msg)
		case "buffering":
			log.Printf("save to new file>%s", msg)
		default:
			log.Println("buffer unexpected error, networkStatus is not postBuffering or buffering")
			return
		}
	}
	if err = fileScanner.Err(); err != nil {
		log.Fatalf("scan file error: %v", err)
		return
	}
}

func runMaShit() string {
	strInp := "Hehello"
	fileHandler, err := os.Create("buffer.buf")
	if err != nil {
		log.Printf("open error %v, try to create buffer.buf\n", err)
		if fileHandler, err = os.Create("buffer.buf"); err != nil {
			log.Println("create error", err)
			os.Exit(1)
		}
		exec.Command("chmod +x buffer.buf")
	}
	defer fileHandler.Close()

	for i := 1; i < 10; i++ {
		log.Println("write")

		fileHandler.WriteString(strInp + strconv.Itoa(i))
	}
	return "done"
}
