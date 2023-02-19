package flags

import (
	"fmt"
	"strconv"
	"strings"
)

type SensorFlag struct {
	set   bool
	value []int
}

const intSeparator = ","

func (sf *SensorFlag) Set(x string) error {
	separated := strings.Split(x, intSeparator)
	sf.value = make([]int, len(separated))
	for i, val := range separated {
		value, err := strconv.Atoi(val)
		if err != nil {
			return fmt.Errorf("couldn't parse sensor flag: %v", err)
		}
		if value < 0 {
			return fmt.Errorf("sensor ID cannot be less than 0")
		}
		sf.value[i] = value
	}
	sf.set = true
	return nil
}

func (sf *SensorFlag) String() string {
	result := make([]string, 0, len(sf.value))
	for _, val := range sf.value {
		result = append(result, strconv.Itoa(val))
	}
	return strings.Join(result, ", ")
}

func (sf *SensorFlag) Ids() []int {
	return sf.value
}

func (sf *SensorFlag) IsSet() bool {
	return sf.set
}
