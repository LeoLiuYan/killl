package app

import (
	"time"
)

type (
	StringList []string
	Interval   time.Duration
)

func (sl *StringList) Set(value string) error {
	if value == "" {
		return nil
	}
	origin := *sl
	origin = append(origin, value)
	*sl = origin
	return nil
}

func (sl *StringList) String() (s string) {
	for _, x := range *sl {
		if s != "" {
			s += ";"
		}
		s += x
	}
	return
}

func (interval *Interval) Set(value string) error {
	if value == "" {
		return nil
	}
	v, err := time.ParseDuration(value)
	if err != nil {
		return err
	}
	*interval = Interval(v)
	return nil
}

func (interval *Interval) String() string {
	return (*time.Duration)(interval).String()
}
