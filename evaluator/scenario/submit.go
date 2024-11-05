package scenario

import(
	"time"
	"errors"
	"fmt"
)

func (sb *SubmitResponse) Init() error {
	if(sb.TimeSpec != nil) {
		ref, err := parseTimeRef(sb.TimeSpec)
		if err != nil {
			return errors.Join(fmt.Errorf("Unable to parse Time ref field"), err)
		}
		sb.timeRef = ref
	}
	return nil
}

func (sb *SubmitResponse) ShiftTime(previous time.Time) time.Time {
	if(sb.timeRef == nil) {
		return previous
	}
	return sb.timeRef.ShiftTime(previous)
}