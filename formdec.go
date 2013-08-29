package main

import (
	"github.com/gorilla/schema"
	"reflect"
	"regexp"
	"time"
)

type EMail string

var emailRegex = regexp.MustCompile(`^.+@.+$`)

func EMailConvert(s string) reflect.Value {
	if emailRegex.MatchString(s) {
		return reflect.ValueOf(EMail(s))
	}
	return reflect.Value{}
}

type timelocForm struct {
	Loc *time.Location
}

func locationConverter(s string) reflect.Value {
	loc, err := time.LoadLocation(s)
	if err != nil {
		return reflect.Value{}
	}
	return reflect.ValueOf(timelocForm{loc})
}

var formdec *schema.Decoder

func init() {
	formdec = schema.NewDecoder()
	formdec.RegisterConverter(EMail(""), EMailConvert)
	formdec.RegisterConverter(timelocForm{}, locationConverter)
}
