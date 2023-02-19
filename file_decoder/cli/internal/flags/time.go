package flags

import (
	"fmt"
	"time"
)

type TimeFlag struct {
	set   bool
	value time.Time
}

const timeLayout = "060102150405.000000"

func (tf *TimeFlag) Set(x string) error {
	parsed, err := time.Parse(timeLayout, x)
	if err != nil {
		return fmt.Errorf("couldn't parse time: %v", err)
	}
	tf.value = parsed
	tf.set = true
	return nil
}

func (tf *TimeFlag) String() string {
	return tf.value.Format(timeLayout)
}

func (tf *TimeFlag) StartTime() time.Time {
	return tf.value
}

func (tf *TimeFlag) IsSet() bool {
	return tf.set
}
