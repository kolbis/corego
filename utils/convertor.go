package utils

import (
	"strconv"
	"time"
)

// Convertor ...
type Convertor struct {
}

// NewConvertor ...
func NewConvertor() Convertor {
	return Convertor{}
}

// FromInt64ToString will convert a number int64 into string
func (c Convertor) FromInt64ToString(num int64) string {
	return strconv.FormatInt(num, 10)
}

// FromStringToInt64 will convert a number as string into int64
// if conversion fails, it will return -1
func (c Convertor) FromStringToInt64(value string) int64 {
	num, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return -1
	}
	return num
}

// MilisecondsToDuration will return a time.Duration from the given duration string
func (c Convertor) MilisecondsToDuration(value int64) time.Duration {
	duration := time.Duration(value) * time.Millisecond
	return duration
}

// DurationToMiliseconds will return a time.Duration as int64 miliseconds
func (c Convertor) DurationToMiliseconds(value time.Duration) int64 {
	return value.Milliseconds()
}

// FromTimeToUnix will convert the time into unix int64, then return it as a string
func (c Convertor) FromTimeToUnix(t time.Time) int64 {
	time := t.UTC().Unix()
	return time
}

// FromUnixToTime will convert the string time to int64 which represents the unix time
// then it will create a new time and will return it
func (c Convertor) FromUnixToTime(t int64) time.Time {
	var date time.Time
	date = time.Unix(t, 0).UTC()
	return date
}
