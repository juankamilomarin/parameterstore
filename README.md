# parameterstore
A package that loads parameters dynamically from a parameter store into a struct

## How to use

You just have to implement the ``GetParams`` method from the ``ParameterStore`` interface and specify the tag which provides the parameter name for each field

## Example

```go
package main

import (
	"errors"
	"fmt"
	"time"

	"github.com/juankamilomarin/parameterstore"
)

const tagName = "memoryMap"

type appParams struct {
	DbUsername     string        `memoryMap:"dbusername"`
	DbPassword     string        `memoryMap:"dbpassword"`
	DbPoolSize     int           `memoryMap:"dbpoolsize"`
	DbQueryTimeout time.Duration `memoryMap:"dbquerytimeout"`
	Https          bool          `memoryMap:"enablehttps"`
}

var AppParams appParams

func main() {
	err := parameterstore.LoadParamsGroup(&AppParams, MapParameterStore{}, tagName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", AppParams)
}

type MapParameterStore struct{}

func (ps MapParameterStore) GetParams(paramMap map[string]string) error {
	p := map[string]string{
		"dbusername":     "admin",
		"dbpassword":     "admin123",
		"dbpoolsize":     "100",
		"dbquerytimeout": "10000000000",
		"enablehttps":    "true",
	}

	for key := range paramMap {
		if key == "error" {
			return errors.New("cannot get parameters")
		}
		paramMap[key] = p[key]
	}
	return nil
}

```

# More examples

* [Read parameters from AWS SSM Parameter Store](https://github.com/juankamilomarin/parameterstore-examples/tree/main/aws)
* [Read parameters from environment variables](https://github.com/juankamilomarin/parameterstore-examples/tree/main/envvar)
* [Read parameters from in memory map](https://github.com/juankamilomarin/parameterstore-examples/tree/main/map)