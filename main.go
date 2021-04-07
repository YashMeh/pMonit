package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/struCoder/pidusage"
)

var ReadingsArray []*pidusage.SysInfo
var Flag bool

func doEvery(d time.Duration, pID int, f func(_ time.Time, pID int)) {
	for x := range time.Tick(d) {
		if !Flag {
			f(x, pID)
		} else {
			fmt.Println("Pressed !!")
			os.Exit(0)
		}
	}
}

func takeReadings(t time.Time, pID int) {
	sysInfo, err := pidusage.GetStat(pID)
	if err != nil {
		panic(err)
	}
	ReadingsArray = append(ReadingsArray, sysInfo)
	fmt.Println(ReadingsArray)
}

func main() {
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			text = strings.Replace(text, "\n", "", -1)

			if strings.Compare("q", text) == 0 {
				Flag = true
			}
		}

	}()

	// select {
	// case msg := <-CH:
	// 	fmt.Println("received message", msg)
	// default:
	// 	doEvery(5*time.Second, 845, takeReadings)
	// }
	doEvery(5*time.Second, 845, takeReadings)

}
