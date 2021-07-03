package main

import (
	"fmt"
	"time"

	"github.com/juankamilomarin/parameterstore"
)

const tagName = "awsSsmParamName"

// The values of the gat awsSsmParamName are the names of the
// parameters stored in AWS SSM Parameter Store
type appParams struct {
	DbUsername     string        `awsSsmParamName:"/myapp/dbnameusername"`
	DbPassword     string        `awsSsmParamName:"/myapp/dbpassword"`
	DbPoolSize     int           `awsSsmParamName:"/myapp/dbpoolsize"`
	DbQueryTimeout time.Duration `awsSsmParamName:"/myapp/dbquerytimeout"`
	Https          bool          `awsSsmParamName:"/myapp/enablehttps"`
}

var AppParams appParams

func main() {
	err := parameterstore.LoadParamsGroup(&AppParams, SSMParameterStore{}, tagName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", AppParams)
}
