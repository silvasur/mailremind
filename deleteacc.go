package main

import (
	"github.com/gorilla/sessions"
	"kch42.de/gostuff/mailremind/model"
	"log"
	"net/http"
)

type reallydeleteTpldata struct {
	OK bool
}

func deleteask(user model.User, sess *sessions.Session, req *http.Request) interface{} {
	return &reallydeleteTpldata{user != nil}
}

func deleteacc(user model.User, sess *sessions.Session, req *http.Request) interface{} {
	outdata := &msgTpldata{Title: "Delete Account"}

	if user == nil {
		outdata.Class = "error"
		outdata.Msg = "You need to be logged in to do that"
		return outdata
	}

	if err := user.Delete(); err != nil {
		log.Printf("Error while deleting account: %s", err)
		outdata.Class = "error"
		outdata.Msg = "An error occurred during deletion. Please contact support, if this happens again."
		return outdata
	}

	delete(sess.Values, "uid")
	outdata.Class = "success"
	outdata.Msg = "Account deleted."
	return outdata
}