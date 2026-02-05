package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Record struct {
	StartTime time.Time
	Tag       string
	Latency   time.Duration
	ReqSize   string
	Code      string
}

func main() {
	file := "./load/report.phout"
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	records := make([]Record, 0)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Split(line, "\t")
		if len(fields) < 12 {
			continue
		}

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
		default:
			code = fields[11]
		}

		records = append(records, Record{
			StartTime: t,
			Tag:       tag,
			Latency:   ms,
			ReqSize:   reqSize,
			Code:      code,
		})
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	// Сортировка по StartTime
	sort.Slice(records, func(i, j int) bool {
		return records[i].StartTime.Before(records[j].StartTime)
	})

	// Вывод
	fmt.Printf("%s %-23s %-15s %-10s %-10s %-15s\n", "Number", "Start Time", "Tag", "Latency", "ReqSize", "Code")
	for i, r := range records {
		fmt.Printf("[%d] %-23s %-15s %-10s %-10s %-15s\n",
			i+1,
			r.StartTime.Format("15:04:05.000"),
			r.Tag,
			r.Latency.String(),
			r.ReqSize,
			r.Code,
		)
	}
}
