package main

import (
	"kch42.de/gostuff/mailremind/model"
	"log"
	"net/http"
)

func activate(rw http.ResponseWriter, req *http.Request) {
	outdata := &msgTpldata{Title: "Activate Account", Class: "error"}
	defer func() {
		if err := tplMsg.Execute(rw, outdata); err != nil {
			log.Printf("Could not execute template in activate: %s", err)
		}
	}()

	req.ParseForm()

	_userid := req.FormValue("U")
	code := req.FormValue("Code")

	if (_userid == "") || (code == "") {
		outdata.Msg = "User or code invalid. Check, if the activation link was correctly copied from the mail."
		return
	}

	userid, err := db.ParseDBID(_userid)
	if err != nil {
		outdata.Msg = "User or code invalid. Check, if the activation link was correctly copied from the mail."
		return
	}

	user, err := dbcon.UserByID(userid)
	switch err {
	case nil:
	case model.NotFound:
		outdata.Msg = "User not found."
		return
	default:
		log.Printf("Error while getting user by ID <%s>: %s", userid, err)
		outdata.Msg = "An error occurred while loading user data. Send a message to the support, if this happens again."
		return
	}

	if user.ActivationCode() != code {
		outdata.Msg = "Wrong activation code."
		return
	}

	if err := user.SetActivationCode(""); err != nil {
		log.Printf("Error while resetting activation code: %s", err)
		outdata.Msg = "An error occurred while activating the user. Send a message to the support, if this happens again."
		return
	}

	if err := user.SetActive(true); err != nil {
		log.Printf("Error while resetting activation code: %s", err)
		outdata.Msg = "An error occurred while activating the user. Send a message to the support, if this happens again."
		return
	}

	outdata.Class = "success"
	outdata.Msg = "Account activated!"
}
