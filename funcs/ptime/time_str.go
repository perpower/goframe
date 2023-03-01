package ptime

import "strconv"

// TimestampStr is a convenience method which retrieves and returns
// the timestamp in seconds as string. like "2006-01-02 15:04:05"
func TimestampStr() string {
	return strconv.FormatInt(Timestamp(), 10)
}

// TimestampMilliStr is a convenience method which retrieves and returns
// the timestamp in milliseconds as string.
func TimestampMilliStr() string {
	return strconv.FormatInt(TimestampMilli(), 10)
}

// TimestampMicroStr is a convenience method which retrieves and returns
// the timestamp in microseconds as string.
func TimestampMicroStr() string {
	return strconv.FormatInt(TimestampMicro(), 10)
}

// TimestampNanoStr is a convenience method which retrieves and returns
// the timestamp in nanoseconds as string.
func TimestampNanoStr() string {
	return strconv.FormatInt(TimestampNano(), 10)
}
