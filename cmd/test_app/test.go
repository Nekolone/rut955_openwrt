package main

import (
	"log"
	"time"
)

func main() {
	log.Print("start")

	for {
		sendTimer(time.Now(), 60, 0.2, 112, 15)
		log.Print("send")
	}

}

func sendTimer(startTime time.Time, rate float64, speedCoef float64, startCourse int, courseDiffTrigger int) {
	defFinishTime := startTime.Add(time.Duration(rate/(1.0+(GetSpeed()*speedCoef))) * time.Second)
	for defFinishTime.After(time.Now()) && diff(startCourse, getCourse()) < courseDiffTrigger {
		time.Sleep(time.Duration(1000 / (1+GetSpeed()*speedCoef)) * time.Millisecond)
	}
}

func diff(a, b int) int {
	if a < b {
		return b - a
	}
	return a - b
}

func getCourse() int {
	return 12
}

func GetSpeed() float64 {
	return 12
}

