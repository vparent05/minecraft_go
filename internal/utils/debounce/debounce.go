package debounce

import "time"

type Debounce struct {
	delay time.Duration
	last  time.Time
}

func NewDebounce(delay time.Duration) *Debounce {
	return &Debounce{
		delay: delay,
		last:  time.Time{},
	}
}

func (d *Debounce) Do(f func()) bool {
	now := time.Now()
	if now.Sub(d.last) > d.delay {
		d.last = now
		f()
		return true
	}
	return false
}
