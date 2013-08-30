package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/gorilla/sessions"
	"kch42.de/gostuff/mailremind/model"
	"log"
	"net/http"
)

type loginTpldata struct {
	Error, Success string
}

type loginFormdata struct {
	Mail, Password string
}

func login(user model.User, sess *sessions.Session, req *http.Request) interface{} {
	outdata := &loginTpldata{}

	if user != nil {
		outdata.Success = "You are already logged in"
		return outdata
	}

	if req.Method != "POST" {
		return outdata
	}

	if err := req.ParseForm(); err != nil {
		outdata.Error = "Data of form could not be understand. If this happens again, please contact support!"
		return outdata
	}

	indata := new(loginFormdata)
	if err := formdec.Decode(indata, req.Form); (err != nil) || (indata.Mail == "") || (indata.Password == "") {
		outdata.Error = "Input data wrong or missing. Please fill in all values."
		return outdata
	}

	user, err := dbcon.UserByMail(indata.Mail)
	switch err {
	case nil:
	case model.NotFound:
		outdata.Error = "E-Mail or password was wrong."
		return outdata
	default:
		log.Printf("Error while loding user data (login): %s", err)
		outdata.Error = "User data could not be loaded. Please contact support, if this happens again."
		return outdata
	}

	if bcrypt.CompareHashAndPassword(user.PWHash(), []byte(indata.Password)) != nil {
		outdata.Error = "E-Mail or password was wrong."
		return outdata
	}

	sess.Values["uid"] = user.ID().String()
	outdata.Success = "Login successful"
	return outdata
}

func logincheck(user model.User, sess *sessions.Session, req *http.Request) interface{} {
	outdata := &msgTpldata{Msg: "<nil>"}
	if user != nil {
		outdata.Msg = user.Email()
	}
	return outdata
}

func logout(user model.User, sess *sessions.Session, req *http.Request) interface{} {
	delete(sess.Values, "uid")
	return &msgTpldata{Class: "success", Title: "Logout", Msg: "Your are now logged out."}
}
