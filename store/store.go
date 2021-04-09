package store

import (
	"bufio"
	"errors"
	"strconv"
)

var (
	Flag       bool
	FileWriter *bufio.Writer
)

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

func HandleInput(args []string) (int, int, string, error) {
	if len(args) != 4 {
		return 0, 0, "", errors.New("Insufficient args passed")
	}
	pID, err := strconv.ParseInt(args[1], 10, 32)
	if err != nil {
		return 0, 0, "", errors.New("Invalid process ID passed")
	}
	//Check if 's' is provided in the end then truncate
	iterStrLen := len(args[2])
	if iterStrLen > 0 && args[2][iterStrLen-1] == 's' {
		args[2] = args[2][:iterStrLen-1]
	}
	iter, err := strconv.ParseInt(args[2], 10, 32)
	//Check if the filename contains .csv then fine else send error
	fileLen := len(args[3])
	if fileLen > 0 && args[3][fileLen-4:fileLen] != ".csv" {
		return 0, 0, "", errors.New("Only CSV files allowed")
	}
	if err != nil {
		return 0, 0, "", errors.New("Invalid iteration passed")
	}
	return int(pID), int(iter), args[3], nil
}
