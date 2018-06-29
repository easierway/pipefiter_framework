package pipefilter

//StraightPipeline is composed of the filters, and the filters are piled as a straigt line.
type StraightPipeline struct {
	Name    string
	Filters *[]Filter
}

//Process is to process the coming data by the pipeline
func (self *StraightPipeline) Process(data interface{}) (interface{}, error) {
	var ret interface{}
	var err error
	for _, filter := range *self.Filters {
		ret, err = filter.Process(data)
		if err != nil {
			return ret, err
		}
		data = ret

	}
	return ret, err
}
