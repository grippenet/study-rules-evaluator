package scenario

import(
	"time"
	"errors"
)

const dayInSeconds = 24 * 60 * 60 * time.Second
const WeekInSeconds = 7 * dayInSeconds

const FixedTimeLayout = "2006-01-02 15:04:05"

var ErrParseDurationField = errors.New("Unable to parse duration in time ref")
var ErrParseTimeField = errors.New("Unable to parse time in time ref")

func parseTimeRef(spec *TimeRefSpec) (TimeRef, error) {
	if(spec == nil) {
		return nil, nil
	}
	if(spec.Fixed != "") {
		t, err := time.Parse(FixedTimeLayout, spec.Fixed)
		if(err != nil) {
			return nil, errors.Join(ErrParseTimeField, err)
		}
		return &TimeRefAbsolute{time: t}, nil
	}
	dur := time.Duration(0)
	var err error
	if(spec.Duration != "") {
		dur, err = time.ParseDuration(spec.Duration)
		if(err != nil) {
			return nil, errors.Join(ErrParseDurationField, err)
		}
	}
	if(spec.Days != 0) {
		dur = dur + time.Duration(int64(spec.Days)) * dayInSeconds
	}
	if(spec.Weeks != 0) {
		dur = dur + time.Duration(int64(spec.Weeks)) * WeekInSeconds 
	}
	return &TimeRefRelative{duration: dur }, nil
}

type TimeRefAbsolute struct {
	time time.Time
}

func (ref *TimeRefAbsolute) ShiftTime(previous time.Time) time.Time {
	return ref.time
}

type TimeRefRelative struct {
	duration time.Duration
}

func (ref *TimeRefRelative) ShiftTime(previous time.Time) time.Time {
	return previous.Add(ref.duration)
}
