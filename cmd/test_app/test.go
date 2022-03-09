package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	modbusclient "rut955_openwrt/pkg/go-modbus"
	"strconv"
	"strings"
)

func main() {
	var (
		host string
		port int
	)
	const (
		defaultHost = "172.16.1.12"
		defaultPort = modbusclient.MODBUS_PORT
	)

	flag.StringVar(&host, "host", defaultHost, "Slave device host (url or ip address)")
	flag.IntVar(&port, "port", defaultPort, fmt.Sprintf("Slave device port (the default is %d)", defaultPort))
	flag.Parse()
	trace := true

	conn, cerr := modbusclient.ConnectTCP(host, port)
	if cerr != nil {
		log.Println(fmt.Sprintf("Connection error: %s", cerr))
		os.Exit(1)
	}

	// attempt to read one (0x01) holding registers starting at address 200
	readData := make([]byte, 3)
	readData[0] = byte(200 >> 8)   // (High Byte)
	readData[1] = byte(200 & 0xff) // (Low Byte)
	readData[2] = 0x01

	// make this read request transaction id 1, with a 300 millisecond tcp timeout
	readResult, readErr := modbusclient.TCPRead(conn, 300, 12, modbusclient.FUNCTION_READ_HOLDING_REGISTERS, false, 0x00, readData, trace)
	if readErr != nil {
		log.Println(readErr)
	}
	log.Println(readResult)

	// attempt to write to a single coil at address 300
	writeData := make([]byte, 3)
	writeData[0] = byte(300 >> 8)   // (High Byte)
	writeData[1] = byte(300 & 0xff) // (Low Byte)
	writeData[2] = 0x00             // 0xff turns the coil on; 0x00 turns the coil off
	// make this read request transaction id 2, with a 300 millisecond tcp timeout
	writeResult, writeErr := modbusclient.TCPWrite(conn, 300, 2, modbusclient.FUNCTION_WRITE_SINGLE_COIL, false, 0x0C, writeData, trace)
	if writeErr != nil {
		log.Println(writeErr)
	}
	log.Println(writeResult)

	modbusclient.DisconnectTCP(conn)

}

func connectModbus() {

}

func getLat() string {
	out, err := exec.Command("cat", "file.txt").Output()
	if (err != nil) || (len(out) == 0) {
		return "NA;NA"
	}

	strOut := strings.Replace(string(out), "\r\n", "", -1)
	fltOut, _ := strconv.ParseFloat(strOut, 64)

	if fltOut == 0 {
		return "NA;NA"
	}
	if fltOut > 0 {
		return fmt.Sprintf("%.0f;N", fltOut*100)
	}
	return fmt.Sprintf("%.0f;S", fltOut*-100)

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
