package main

import (
	"fmt"
	"time"

	"github.com/joho/godotenv"
	"github.com/juankamilomarin/parameterstore"
)

const tagName = "envvar"

type appParams struct {
	DbUsername     string        `envvar:"DBUSERNAME"`
	DbPassword     string        `envvar:"DBPASSWORD"`
	DbPoolSize     int           `envvar:"DBPOOLSIZE"`
	DbQueryTimeout time.Duration `envvar:"DBQUERYTIMEOUT"`
	Https          bool          `envvar:"ENABLE_HTTPS"`
}

var AppParams appParams

func main() {
	err := godotenv.Load("env_var")
	if err != nil {
		panic(err)
	}
	err = parameterstore.LoadParamsGroup(&AppParams, EnvVarParameterStore{}, tagName)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%+v\n", AppParams)
}
