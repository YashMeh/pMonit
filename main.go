package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/struCoder/pidusage"
)

var (
	Flag       bool
	FileWriter *bufio.Writer
)

func doEvery(d time.Duration, pID int, f func(_ time.Time, pID int)) {
	for x := range time.Tick(d) {
		if !Flag {
			f(x, pID)
		} else {
			fmt.Println("Gracefully exiting")
			//Flush the buffer on the writer in case someone stops midway
			if err := FileWriter.Flush(); err != nil {
				panic(err)
			}
			fmt.Println("Bye bye !!")
			os.Exit(0)
		}
	}
}
func SaveReadings(reading string, writer *bufio.Writer) {
	//Stream the output to a csv
	// make a buffer to keep chunks that are read
	buf := []byte(reading)

	//Add the readings on the buffer
	// Write the chunk to the file
	if _, err := writer.Write(buf); err != nil {
		panic(err)
	}
	//Keep flushing the data to the file
	if err := writer.Flush(); err != nil {
		panic(err)
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
	SaveReadings(csvEntry, FileWriter)

}

func main() {

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
				Flag = true
			}
		}

	}()

	// Open the file
	fileOpener, err := os.Create("output.txt")
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
	FileWriter = bufio.NewWriter(fileOpener)

	//Add the header values
	SaveReadings("CPU(%),Memory(kb),Time \n", FileWriter)

	doEvery(1*time.Second, 3181, takeReadings)

}
