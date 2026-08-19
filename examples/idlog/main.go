package main

import (
	"fmt"

	"github.com/idfactory/idtools/idlog"
)

func main() {
	idlog.StreamName = "my-service"

	idlog.Log("Start example")         //INFO example
	idlog.Log(fmt.Errorf("error"))     //ERROR example (no trace)
	idlog.Log(1)                       //WARN example (with trace)
	idlog.Debug("debug example", true) //DEBUG example

	// Manually add trace to errors
	err := fmt.Errorf("some error")
	errWithTrace := idlog.AddTrace(err)
	idlog.Log(errWithTrace) //ERROR with manual trace

	idlog.Fail(fmt.Errorf("critical error")) //FAIL example
}
