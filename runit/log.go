package runit

import (
	gologging "github.com/sigmonsays/go-logging"
)

var log gologging.Logger

func init() {
	log = gologging.Register("runit", func(newlog gologging.Logger) { log = newlog })
}
