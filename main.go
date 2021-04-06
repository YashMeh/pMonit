package main

import (
	"fmt"
	"time"

	"github.com/struCoder/pidusage"
)

var ReadingsArray []*pidusage.SysInfo

func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func takeReadings(t time.Time) {
	sysInfo, err := pidusage.GetStat(845)
	if err != nil {
		panic(err)
	}
	ReadingsArray = append(ReadingsArray, sysInfo)
	fmt.Println(ReadingsArray)
}

func main() {
	doEvery(500*time.Millisecond, takeReadings)

}
