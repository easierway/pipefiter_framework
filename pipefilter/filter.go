/*Package pipefilter is to define the interfaces and the structures for pipe-filter style implementation

  Created: 2018-1-19
  @All the copy rights belong to Mobvista.com
**/
package pipefilter

//Filter interface is the definition of the data processing components
//Pipe-Filter structure
type Filter interface {
	Process(data interface{}) (interface{}, error)
}
