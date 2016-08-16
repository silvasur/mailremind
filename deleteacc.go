package main

import (
	"github.com/gorilla/sessions"
	"github.com/silvasur/mailremind/model"
	"log"
	"net/http"
)

type reallydeleteTpldata struct {
	OK bool
}

func deleteask(user model.User, sess *sessions.Session, req *http.Request) (interface{}, model.User) {
	return &reallydeleteTpldata{user != nil}, user
}

func deleteacc(user model.User, sess *sessions.Session, req *http.Request) (interface{}, model.User) {
	outdata := &msgTpldata{Title: "Delete Account"}

	if user == nil {
		outdata.Class = "error"
		outdata.Msg = "You need to be logged in to do that"
		return outdata, user
	}

	if err := user.Delete(); err != nil {
		log.Printf("Error while deleting account: %s", err)
		outdata.Class = "error"
		outdata.Msg = "An error occurred during deletion."
		return outdata, user
	}

	delete(sess.Values, "uid")
	outdata.Class = "success"
	outdata.Msg = "Account deleted."
	return outdata, user
}
