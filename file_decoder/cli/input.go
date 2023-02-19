package cli

import (
	"file_decoder/cli/internal/flags"
	"flag"
	"fmt"
	"io"
	"os"
	"time"
)

var outputFilename, inputFilename flags.FilenameFlag
var sensorIds flags.SensorFlag
var startTime flags.TimeFlag
var duration flags.DurationFlag

func init() {
	flag.Var(&outputFilename, "output", "name of the output file for parsed measures")
	flag.Var(&inputFilename, "input", "name of the input file for binary data")
	flag.Var(&sensorIds, "sensors", "comma separated IDs of sensors to parse")
	flag.Var(&startTime, "start", "first measure time in format YYMMDDHHIISS.000000")
	flag.Var(&duration, "duration", "duration of the measurements")
	flag.Parse()
}

func OutputWriter() (io.Writer, error) {
	writer := os.Stdout
	if !outputFilename.IsSet() {
		return writer, nil
	}
	writer, err := os.Create(outputFilename.String())
	if err != nil {
		return nil, fmt.Errorf("couldn't open output file: %v", err)
	}
	return writer, nil
}

func InputReader() (io.Reader, error) {
	reader := os.Stdin
	if !inputFilename.IsSet() {
		return reader, nil
	}
	reader, err := os.Open(inputFilename.String())
	if err != nil {
		return nil, fmt.Errorf("couldn't open input file: %v", err)
	}
	return reader, nil
}

func SensorIds() ([]int, bool) {
	return sensorIds.Ids(), sensorIds.IsSet()
}

func StartTime() (time.Time, bool) {
	return startTime.StartTime(), startTime.IsSet()
}

func Duration() (time.Duration, bool) {
	return duration.Duration(), duration.IsSet()
}
