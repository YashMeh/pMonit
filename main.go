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
	ReadingsArray []*pidusage.SysInfo
	Flag          bool
	FileOpener    *os.File
	FileWriter    *bufio.Writer
	Error         error
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

func takeReadings(t time.Time, pID int) {
	sysInfo, err := pidusage.GetStat(pID)
	if err != nil {
		panic(err)
	}
	ReadingsArray = append(ReadingsArray, sysInfo)

	//Stream the output to a csv
	// make a buffer to keep chunks that are read
	buf := []byte("Yash Mehrotra")
	//Add the readings on the buffer
	// Write the chunk to the file
	if _, err := FileWriter.Write(buf); err != nil {
		panic(err)
	}
	//Keep flushing the data to the file
	if err := FileWriter.Flush(); err != nil {
		panic(err)
	}

	// for {
	// 	// read a chunk
	// 	//n, err := r.Read(buf)
	// 	if err != nil && err != io.EOF {
	// 		panic(err)
	// 	}
	// 	if n == 0 {
	// 		break
	// 	}

	// 	// write a chunk
	// 	if _, err := w.Write(buf[:n]); err != nil {
	// 		panic(err)
	// 	}
	// }

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

	doEvery(5*time.Second, 3181, takeReadings)

}
