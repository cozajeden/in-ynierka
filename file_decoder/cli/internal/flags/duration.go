package flags

import (
	"fmt"
	"time"
)

type DurationFlag struct {
	set   bool
	value time.Duration
}

func (df *DurationFlag) Set(x string) error {
	duration, err := time.ParseDuration(x)
	if err != nil {
		fmt.Errorf("couldn't parse duration: %v", err)
	}
	df.value = duration
	df.set = true
	return nil
}

func (df *DurationFlag) String() string {
	return df.value.String()
}

func (df *DurationFlag) Duration() time.Duration {
	return df.value
}

func (df *DurationFlag) IsSet() bool {
	return df.set
}
