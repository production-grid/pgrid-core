package util

import (
	"math"
	"time"
)

/*
Constants for various time/date formats.
*/
const (
	LongTimestampFormat  = "January 2, 2006 3:04:05 PM"
	ShortTimestampFormat = "01/02/2006 3:04:05 PM"
	LongDateFormat       = "January 2, 2006"
	ShortDateFormat      = "01/02/2006"
	TimeFormat           = "3:04:05 PM"
	ISOFormat            = "2006-01-02T15:04:05Z0700"
)

/*
TimeClose returns true if t1 and t2 are withing duration of each other.
*/
func TimeClose(t1, t2 time.Time, duration time.Duration) bool {

	upperBound := t2.Add(duration)
	lowerBound := t2.Add(time.Duration(-1) * duration)

	return t1.After(lowerBound) && t1.Before(upperBound)

}

/*
FormatLongTimeStamp returns the given time as a long timestamp string.
*/
func FormatLongTimeStamp(t time.Time) string {

	return t.Format(LongTimestampFormat)

}

/*
FormatISOTimeStamp returns the given time as an ISO 8601 formatted string.
*/
func FormatISOTimeStamp(t time.Time) string {

	return t.UTC().Format(ISOFormat)

}

/*
ParseISOTimeStamp parses an ISO timestamp in UTC.
*/
func ParseISOTimeStamp(ts string) (time.Time, error) {

	return time.Parse(ISOFormat, ts)

}

/*
ParseLongTimeStamp parses a long timestamp string for the given time zone.
*/
func ParseLongTimeStamp(ts string, tz *time.Location) (time.Time, error) {

	return time.ParseInLocation(LongTimestampFormat, ts, tz)

}

/*
FormatShortTimeStamp returns the given time as a short timestamp string.
*/
func FormatShortTimeStamp(t time.Time) string {

	return t.Format(ShortTimestampFormat)

}

/*
ParseShortTimeStamp parses a short timestamp string for the given time zone.
*/
func ParseShortTimeStamp(ts string, tz *time.Location) (time.Time, error) {

	return time.ParseInLocation(ShortTimestampFormat, ts, tz)

}

/*
FormatLongDate returns the given time as a long date string.
*/
func FormatLongDate(t time.Time) string {

	return t.Format(LongDateFormat)

}

/*
FormatShortDate returns the given time as a short date string.
*/
func FormatShortDate(t time.Time) string {

	return t.Format(ShortDateFormat)

}

/*
ParseShortDate parses a short timestamp string for the given time zone.
*/
func ParseShortDate(ts string, tz *time.Location) (time.Time, error) {

	return time.ParseInLocation(ShortDateFormat, ts, tz)

}

/*
FormatTime returns the given time as a time string.
*/
func FormatTime(t time.Time) string {

	return t.Format(TimeFormat)

}

//NarrativeAgoDiff returns a narrative description of how long ago something happened
func NarrativeAgoDiff(t time.Time) string {

	diff := time.Now().Sub(t)

	seconds := int(math.Floor(diff.Seconds()))
	minutes := int(math.Floor(diff.Minutes()))
	hours := int(math.Floor(diff.Hours()))
	days := int(math.Floor(diff.Hours() / 24))
	years := days / 365

	if years > 0 {
		result := FormatInteger(years)
		result = result + " year"
		if years > 1 {
			result = result + "s"
		}
		result = result + " ago"
		return result
	} else if days > 0 {
		result := FormatInteger(days)
		result = result + " day"
		if days > 1 {
			result = result + "s"
		}
		result = result + " ago"
		return result
	} else if hours > 0 {
		result := FormatInteger(hours)
		result = result + " hour"
		if hours > 1 {
			result = result + "s"
		}
		result = result + " ago"
		return result
	} else if minutes > 0 {
		result := FormatInteger(minutes)
		result = result + " minute"
		if minutes > 1 {
			result = result + "s"
		}
		result = result + " ago"
		return result
	} else if seconds > 0 {
		result := FormatInteger(seconds)
		result = result + " second"
		if seconds > 1 {
			result = result + "s"
		}
		result = result + " ago"
		return result
	}

	return "Just Now"
}
