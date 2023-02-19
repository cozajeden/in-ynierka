package main

import (
	"file_decoder/cli"
	"file_decoder/decode"
	"log"
)

func main() {
	sensorIds, sensorIdsGiven := cli.SensorIds()
	startTime, startTimeGiven := cli.StartTime()
	duration, durationGiven := cli.Duration()

	input, err := cli.InputReader()
	if err != nil {
		log.Fatalf("couldn't open reader: %v\n", err)
	}
	output, err := cli.OutputWriter()
	if err != nil {
		log.Fatalf("couldn't open writer: %v\n", err)
	}

	decoder := decode.NewDecoder(input, output)
	if sensorIdsGiven {
		decoder.ParseOnlyIds(sensorIds)
	}
	if startTimeGiven {
		decoder.ParseSince(startTime)
	}
	if durationGiven {
		decoder.ParseFor(duration)
	}

	err = decoder.Csv()
	if err != nil {
		log.Fatalf("couldn't decode file: %v\n", err)
	}
}
