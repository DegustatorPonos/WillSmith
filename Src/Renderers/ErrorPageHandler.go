package renderers

import "fmt"

const errTempl = `%s

Return to the previous page with :.. or enter a new address
`

func CreateErrorHandler(Error string) func()[]byte {
	return func() []byte {
		var contents = fmt.Sprintf(errTempl, Error)
		return []byte(contents)
	}
}

func CreateErrorWrapper(err error) func()[]byte {
	var wrapped = fmt.Sprintf("A runtime error occured: %v", err.Error())
	return CreateErrorHandler(wrapped)
}
