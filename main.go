package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/struCoder/pidusage"
	"github.com/yashmeh/memMonit/store"
)

func doEvery(d time.Duration, pID int, f func(_ time.Time, pID int)) {
	for x := range time.Tick(d) {
		if !store.Flag {
			f(x, pID)
		} else {
			fmt.Println("Gracefully exiting")
			//Flush the buffer on the writer in case someone stops midway
			if err := store.FileWriter.Flush(); err != nil {
				panic(err)
			}
			fmt.Println("Bye bye !!")
			os.Exit(0)
		}
	}
}

func takeReadings(t time.Time, pID int) {
	sysInfo, err := pidusage.GetStat(pID)
	if err != nil {
		panic(err)
	}

	//Create the entry
	currentTime := time.Now()
	cpuEntry := fmt.Sprintf("%0.2f", sysInfo.CPU)
	memEntry := fmt.Sprintf("%0.2f", sysInfo.Memory/1024)
	csvEntry := cpuEntry + "," + memEntry + "," + currentTime.Format("2006-01-02 15:04:05") + "\n"

	//Save to the CSV
	store.SaveReadings(csvEntry, store.FileWriter)

}

func main() {
	//Accept processId,interval,fileName
	pID, iter, fileName, err := store.HandleInput(os.Args)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	//Go routine to read custom command from CLI
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			text, err := reader.ReadString('\n')
			if err != nil {
				return
			}
			text = strings.Replace(text, "\n", "", -1)

			if strings.Compare("q", text) == 0 {
				store.Flag = true
			}
		}

	}()

	// Open the file
	fileOpener, err := os.Create(fileName)
	//Panic the function in case of error
	if err != nil {
		panic(err)
	}
	// The panic will call the defer function
	// Close the writer and create a panic and the function will panic
	defer func() {
		if err := fileOpener.Close(); err != nil {
			panic(err)
		}
	}()
	// Assign to the global writer
	store.FileWriter = bufio.NewWriter(fileOpener)

	//Add the header values
	store.SaveReadings("CPU(%),Memory(kb),Time \n", store.FileWriter)

	doEvery(time.Duration(iter)*time.Second, pID, takeReadings)

}
