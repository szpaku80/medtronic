package medtronic

import (
	"time"
)

const (
	BasalRates    Command = 0x92
	BasalPatternA Command = 0x93
	BasalPatternB Command = 0x94
)

type BasalRate struct {
	Start time.Duration // offset from 00:00:00
	Rate  int           // milliUnits per hour
}

type BasalRateSchedule struct {
	Schedule []BasalRate
}

func (pump *Pump) basalSchedule(cmd Command) BasalRateSchedule {
	data := pump.Execute(cmd)
	if pump.Error() != nil {
		return BasalRateSchedule{}
	}
	info := []BasalRate{}
	for i := 1; i < len(data); i += 3 {
		r := data[i]
		t := data[i+2]
		// Don't stop if the 00:00 rate happens to be zero.
		if i > 1 && r == 0 && t == 0 {
			break
		}
		start := scheduleToDuration(t)
		rate := int(r) * 25
		info = append(info, BasalRate{Start: start, Rate: rate})
	}
	return BasalRateSchedule{Schedule: info}
}

func (pump *Pump) BasalRates() BasalRateSchedule {
	return pump.basalSchedule(BasalRates)
}

func (pump *Pump) BasalPatternA() BasalRateSchedule {
	return pump.basalSchedule(BasalPatternA)
}

func (pump *Pump) BasalPatternB() BasalRateSchedule {
	return pump.basalSchedule(BasalPatternB)
}

func (s BasalRateSchedule) BasalRateAt(t time.Time) BasalRate {
	d := sinceMidnight(t)
	last := BasalRate{}
	for _, v := range s.Schedule {
		if v.Start > d {
			break
		}
		last = v
	}
	return last
}
