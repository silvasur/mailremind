package main

import (
	"github.com/gorilla/sessions"
	"kch42.de/gostuff/mailremind/model"
	"log"
	"net/http"
)

func activate(user model.User, sess *sessions.Session, req *http.Request) (interface{}, model.User) {
	outdata := &msgTpldata{Title: "Activate Account", Class: "error"}

	req.ParseForm()

	_userid := req.FormValue("U")
	code := req.FormValue("Code")

	if (_userid == "") || (code == "") {
		outdata.Msg = "User or code invalid. Check, if the activation link was correctly copied from the mail."
		return outdata, user
	}

	userid, err := db.ParseDBID(_userid)
	if err != nil {
		outdata.Msg = "User or code invalid. Check, if the activation link was correctly copied from the mail."
		return outdata, user
	}

	switch user, err = dbcon.UserByID(userid); err {
	case nil:
	case model.NotFound:
		outdata.Msg = "User not found."
		return outdata, user
	default:
		log.Printf("Error while getting user by ID <%s>: %s", userid, err)
		outdata.Msg = "An unknown error occurred while loading user data."
		return outdata, user
	}

	if user.ActivationCode() != code {
		outdata.Msg = "Wrong activation code."
		return outdata, user
	}

	if err := user.SetActivationCode(""); err != nil {
		log.Printf("Error while resetting activation code: %s", err)
		outdata.Msg = "An unknown error occurred while activating the user."
		return outdata, user
	}

	if err := user.SetActive(true); err != nil {
		log.Printf("Error while resetting activation code: %s", err)
		outdata.Msg = "An unknown error occurred while activating the user."
		return outdata, user
	}

	outdata.Class = "success"
	outdata.Msg = "Account activated!"
	return outdata, user
}
