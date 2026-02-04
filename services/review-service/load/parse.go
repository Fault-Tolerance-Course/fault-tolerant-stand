package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	file := "./load/report.phout"
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	fmt.Printf("%s %-23s %-15s %-6s %-6s %-8s\n", "Number", "Start Time", "Tag", "Latency", "ReqSize", "Code")

	i := 1
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		if len(fields) < 12 {
			continue
		}

		// поле 0 = timestamp как float (seconds.milliseconds)
		tsFloat, err := strconv.ParseFloat(fields[0], 64)
		if err != nil {
			continue
		}
		sec := int64(tsFloat)
		nsec := int64((tsFloat - float64(sec)) * 1e9)
		t := time.Unix(sec, nsec)

		tag := fields[1]
		rtt := fields[2]

		dur, _ := strconv.ParseInt(rtt, 10, 64)

		ms := time.Duration(dur) * 1000
		reqSize := fields[8]

		var code string
		switch fields[11] {
		case "grpc_0":
			code = "OK"
		case "grpc_8":
			code = "RESOURCE_EXHAUSTED"
		}
		fmt.Printf("[%d] %-23s %-15s %-6s %-6s %-8s\n", i, t.Format("15:04:05.000"), tag, ms.String(), reqSize, code)
		i++
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
