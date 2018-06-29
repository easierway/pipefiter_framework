package pipefilter

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

type P1Filter struct {
}

func (self *P1Filter) Process(data interface{}) (interface{}, error) {
	ret := *(data.(*string)) + "_P1"
	return &ret, nil
}

type P2Filter struct {
}

func (self *P2Filter) Process(data interface{}) (interface{}, error) {
	ret := *(data.(*string)) + "_P2"
	return &ret, nil
}

type ErrorFilter struct {
}

var FailedProcessError error = errors.New("failed to process")

func (self *ErrorFilter) Process(data interface{}) (interface{}, error) {
	ret := *(data.(*string)) + "_ERR"
	return &ret, FailedProcessError
}

func TestHappyProcess(t *testing.T) {
	Convey("Given a pipeline", t, func() {
		filters := []Filter{&P1Filter{}, &P2Filter{}}

		std := StraightPipeline{
			Name:    "Standard",
			Filters: &filters,
		}

		Convey("When all the filters run successfully", func() {
			in := "Start"
			ret, err := std.Process(&in)
			Convey("Then get the excepted result", func() {
				So(err, ShouldBeNil)
				So("Start_P1_P2", ShouldEqual, *(ret.(*string)))

			})
		})
	})

}

func TestUnhappyProcess(t *testing.T) {
	Convey("Given a pipeline", t, func() {
		filters := []Filter{&P1Filter{}, &ErrorFilter{}, &P2Filter{}}

		std := StraightPipeline{
			Name:    "Standard",
			Filters: &filters,
		}

		Convey("When some filters fail to run", func() {
			in := "Start"
			ret, err := std.Process(&in)
			Convey("Then get the excepted error", func() {
				So(err, ShouldEqual, FailedProcessError)
				So("Start_P1_ERR", ShouldEqual, *(ret.(*string)))

			})
		})
	})

}
