package common

import "time"

var DefaultAttendanceStartTime = "08:00:00"
var DefaultAttendanceEndTime = "16:00:00"
var DefaultZero = 0

var DefaultAccessControlStartTime = "00:00:00"
var DefaultAccessControlEndTime = "23:59:59"

func ConvertTimeStrToTime(timeStr string) (time.Time, error) {
	layout := "2006-01-02 15:04:05"
	t, err := time.Parse(layout, timeStr)
	if err != nil {
		return t, err
	}
	return t, nil
}
