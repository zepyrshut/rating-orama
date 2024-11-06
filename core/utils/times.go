package utils

import "time"

// TODO: Move to toolbox
func TimeParser(timeString string) (time.Time, error) {
	if len(timeString) == 1 {
		return time.Time{}, nil
	}
	if len(timeString) == 4 {
		return time.Parse("2006", timeString)
	}
	if len(timeString) == 9 {
		return time.Parse("Jan. 2006", timeString)
	}
	if len(timeString) == 10 {
		return time.Parse("2 Jan 2006", timeString)
	}
	if len(timeString) == 11 {
		if timeString[5:6] == "." {
			return time.Parse("2 Jan. 2006", timeString)
		} else {
			return time.Parse("2 Jan 2006", timeString)
		}
	}

	return time.Parse("2 Jan. 2006", timeString)
}
