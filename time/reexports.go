package time

import stdtime "time"

// Time from the standard package
type Time = stdtime.Time

// Duration from the standard package
type Duration = stdtime.Duration

// Location from the standard package
type Location = stdtime.Location

// Month from the standard package
type Month = stdtime.Month

// ParseError from the standard package
type ParseError = stdtime.ParseError

// Weekday from the standard package
type Weekday = stdtime.Weekday

// Global functions from the standard package
var (
	Date                   = stdtime.Date
	FixedZone              = stdtime.FixedZone
	LoadLocation           = stdtime.LoadLocation
	LoadLocationFromTZData = stdtime.LoadLocationFromTZData
	ParseDuration          = stdtime.ParseDuration
	Parse                  = stdtime.Parse
	ParseInLocation        = stdtime.ParseInLocation
	Unix                   = stdtime.Unix
	UnixMicro              = stdtime.UnixMicro
	UnixMilli              = stdtime.UnixMilli
	Until                  = stdtime.Until
)

// Constants from the standard package
const (
	Layout      = stdtime.Layout
	ANSIC       = stdtime.ANSIC
	UnixDate    = stdtime.UnixDate
	RubyDate    = stdtime.RubyDate
	RFC822      = stdtime.RFC822
	RFC822Z     = stdtime.RFC822Z
	RFC850      = stdtime.RFC850
	RFC1123     = stdtime.RFC1123
	RFC1123Z    = stdtime.RFC1123Z
	RFC3339     = stdtime.RFC3339
	RFC3339Nano = stdtime.RFC3339Nano
	Kitchen     = stdtime.Kitchen
	Stamp       = stdtime.Stamp
	StampMilli  = stdtime.StampMilli
	StampMicro  = stdtime.StampMicro
	StampNano   = stdtime.StampNano
	DateTime    = stdtime.DateTime
	DateOnly    = stdtime.DateOnly
	TimeOnly    = stdtime.TimeOnly

	Nanosecond  = stdtime.Nanosecond
	Microsecond = stdtime.Microsecond
	Millisecond = stdtime.Millisecond
	Second      = stdtime.Second
	Minute      = stdtime.Minute
	Hour        = stdtime.Hour
	Day         = stdtime.Hour * 24 // useful extension

	January   = stdtime.January
	February  = stdtime.February
	March     = stdtime.March
	April     = stdtime.April
	May       = stdtime.May
	June      = stdtime.June
	July      = stdtime.July
	August    = stdtime.August
	September = stdtime.September
	October   = stdtime.October
	November  = stdtime.November
	December  = stdtime.December

	Sunday    = stdtime.Sunday
	Monday    = stdtime.Monday
	Tuesday   = stdtime.Tuesday
	Wednesday = stdtime.Wednesday
	Thursday  = stdtime.Thursday
	Friday    = stdtime.Friday
	Saturday  = stdtime.Saturday
)

func init() {
	stdtime.Local = stdtime.UTC
}

// Variables (actually constants) from the standard package
var (
	Local = stdtime.Local
	UTC   = stdtime.UTC
)
