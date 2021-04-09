## pMonit

[![Go Report Card](https://goreportcard.com/badge/github.com/YashMeh/pMonit)](https://goreportcard.com/report/github.com/YashMeh/pMonit)

This is a program that logs CPU + memory utilization of a process ID for every N seconds to a CSV file.

### Format

```
go run main.go PID N FILENAME
```

Press 'q' to gracefully stop the process

### Example

```
go run main.go 3181 2s output.csv
```

#### Features/Learnings

- Streams the readings to the file instead of loading in memory.

#### Todo

- Save the readings to .png file
