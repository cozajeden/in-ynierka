package decode

import (
	"fmt"
	"io"
)

const standardTimeFormat = "2006-01-02 15:04:05.000000"

func (d *Decoder) Csv() error {
	_, err := d.writer.Write([]byte("sensor_id,name,value,time\n"))
	if err != nil {
		return fmt.Errorf("couldn't write to output: %v", err)
	}
	iterator := 0
	for {
		iterator++
		parsedReading, err := d.parseLine()
		if err == io.EOF {
			break
		}
		if err == unwantedMeasureError {
			continue
		}
		if err != nil {
			return fmt.Errorf("couldn't parse read no %v: %v", iterator, err)
		}
		d.writer.Write([]byte(fmt.Sprintf("%v,%v,%v,%v\n",
			parsedReading.id,
			parsedReading.name,
			parsedReading.value,
			parsedReading.time.Format(standardTimeFormat))))
	}
	return nil
}
