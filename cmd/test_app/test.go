package main

import (
	"fmt"
)

func main() {

	testvar := 111
	usethis(&testvar)
	fmt.Println(testvar)

}

func usethis(i *int) {
	fmt.Println(i)
	*i = 222
	fmt.Println(i)
	anotherUse(i)
}

func anotherUse(i *int) {
	fmt.Println(i)
	*i = 2333
	fmt.Println(i)
}

func setNewMsg(msg chan string) {
	for i := 1; i < 100; i++ {
		msg <- "some text"
	}
}

func printMsg(msg chan string) {
	for i := 1; i < 100; i++ {
		fmt.Println(<-msg + "wq")
	}
}
