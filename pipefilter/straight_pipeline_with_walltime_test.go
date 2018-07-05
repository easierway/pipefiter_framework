package pipefilter

import (
	"fmt"
	"sync"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

type MyLog struct{}

func (l *MyLog) Info(v ...interface{}) {
	fmt.Println(v...)
}

type Filter1 struct{}
type Filter2 struct{}

func (f *Filter1) Process(str interface{}) (interface{}, error) {
	time.Sleep(1 * time.Millisecond)
	return str.(string) + "_F1", nil
}

func (f *Filter2) Process(str interface{}) (interface{}, error) {
	time.Sleep(5 * time.Millisecond)
	return str.(string) + "_F2", nil
}

func TestStraightPipelieWithWalltime(t *testing.T) {
	Convey("test straight pipline", t, func() {
		pipeline := NewStraightPipelineWithWallTime(
			"straight_pipeline",
			&Filter1{},
			&Filter2{},
		)
		val, err := pipeline.Process("F0")
		So(err, ShouldBeNil)
		So(val, ShouldEqual, "F0_F1_F2")
	})

	Convey("test wall time", t, func() {
		pipeline := NewStraightPipelineWithWallTime(
			"straight_pipeline",
			&Filter1{},
			&Filter2{},
		)
		pipeline.SetLogger(&MyLog{}, 100*time.Millisecond)
		pipeline.RecordWallTime()
		var wg sync.WaitGroup
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func() {
				for j := 0; j < 100; j++ {
					pipeline.Process("F0")
				}
				wg.Done()
			}()
		}
		wg.Wait()
		t.Log(pipeline.GetWallTimeNs())
	})
}
