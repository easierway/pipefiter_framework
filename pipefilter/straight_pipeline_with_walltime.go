package pipefilter

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

// NewStraightPipelineWithWallTime create a new StraightPipelineWithWallTime
func NewStraightPipelineWithWallTime(name string, filters ...Filter) *StraightPipelineWithWallTime {
	return &StraightPipelineWithWallTime{
		Name:       name,
		Filters:    filters,
		WallTimeNS: make([]AtomicInt, len(filters)),
	}
}

// SetLogger set log
func (p *StraightPipelineWithWallTime) SetLogger(log Logger, flushTime time.Duration) {
	p.log = log
	p.flushTime = flushTime
}

// StraightPipelineWithWallTime straight pipeline with wall time
type StraightPipelineWithWallTime struct {
	Name       string
	Filters    []Filter
	WallTimeNS []AtomicInt
	log        Logger
	flushTime  time.Duration
}

// Process filter
func (p *StraightPipelineWithWallTime) Process(data interface{}) (interface{}, error) {
	var ret interface{}
	var err error

	now := time.Now()
	for i, filter := range p.Filters {
		ret, err = filter.Process(data)
		p.WallTimeNS[i].Add(time.Since(now).Nanoseconds())
		now = time.Now()
		if err != nil {
			return ret, err
		}
		data = ret
	}

	return ret, err
}

// GetWallTimeNs get wall time
func (p *StraightPipelineWithWallTime) GetWallTimeNs() []int64 {
	ts := make([]int64, len(p.WallTimeNS))
	for i, wt := range p.WallTimeNS {
		ts[i] = wt.Reset()
	}
	return ts
}

// RecordWallTime record wall time to log
func (p *StraightPipelineWithWallTime) RecordWallTime() {
	go func() {
		for range time.Tick(p.flushTime) {
			ts := make([]string, len(p.WallTimeNS))
			for i := range p.WallTimeNS {
				ts[i] = fmt.Sprintf("%v", p.WallTimeNS[i].Reset()/(int64)(time.Millisecond))
			}
			p.log.Info(strings.Join(ts, "\t"))
		}
	}()
}

// AtomicInt atomic int
type AtomicInt int64

// Add value
func (a *AtomicInt) Add(i int64) {
	atomic.AddInt64((*int64)(a), i)
}

// Val get value
func (a *AtomicInt) Val() int64 {
	return *(*int64)(a)
}

// Reset value to 0, and return the current value
func (a *AtomicInt) Reset() int64 {
	return atomic.SwapInt64((*int64)(a), 0)
}
