package main

import (
	"fmt"
	"rut955_openwrt/internal/service"
)

func main() {
	fmt.Println("RutOS test1")
	fmt.Println(service.Helo())
}
