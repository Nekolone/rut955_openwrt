package testDir

import "log"

func PrintToLoF(a string) (ret string) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("recovered painc in ptlf > %v", r)
			 "123"
		}
	}()
	log.Print(a)
	log.Panic("unexpected")
	return "111"
}
