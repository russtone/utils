package output

var Formats = []string{"json"}

type Writer interface {
	Write(interface{}) error
}
