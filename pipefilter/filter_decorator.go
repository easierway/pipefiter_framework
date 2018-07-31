package pipefilter

import "github.com/easierway/service_decorators"

// DecoratedFilter is the filter instance decorated with service deocorators
type DecoratedFilter struct {
	serviceDecorators []service_decorators.Decorator
	orgFilter         Filter
	decFn             service_decorators.ServiceFunc
}

func (df *DecoratedFilter) decorateOrgFilter() {
	df.decFn = func(req service_decorators.Request) (service_decorators.Response, error) {
		return df.orgFilter.Process(req)
	}

	for _, dec := range df.serviceDecorators {
		df.decFn = dec.Decorate(df.decFn)
	}
}

// Process is to implement the Filter interface.
func (df *DecoratedFilter) Process(data interface{}) (interface{}, error) {
	return df.decFn(data)
}

// DecorateFilter is to leverage the project "https://github.com/easierway/service_decorators"
// By this method, you can decorate the deocrators on a Filter instance
func DecorateFilter(orgFilter Filter, serviceDecorators []service_decorators.Decorator) Filter {
	df := DecoratedFilter{orgFilter: orgFilter, serviceDecorators: serviceDecorators}
	df.decorateOrgFilter()
	return &df
}
