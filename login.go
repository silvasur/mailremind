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

func userFromSess(sess *sessions.Session) model.User {
	_id, ok := sess.Values["uid"]
	if !ok {
		return nil
	}

	id, ok := _id.(string)
	if !ok {
		return nil
	}

	uid, err := db.ParseDBID(id)
	if err != nil {
		return nil
	}

	user, err := dbcon.UserByID(uid)
	if err != nil {
		return nil
	}

	return user
}

func getSess(req *http.Request) (*sessions.Session, error) {
	sess, err := SessionStorage.Get(req, "mailremind-sess")
	if err != nil {
		sess, err = SessionStorage.New(req, "mailremind-sess")
	}

	return sess, err
}

func login(rw http.ResponseWriter, req *http.Request) {
	outdata := &loginTpldata{}
	defer func() {
		if err := tplLogin.Execute(rw, outdata); err != nil {
			log.Printf("Error executing template in login: %s", err)
		}
	}()

	sess, err := getSess(req)
	if err != nil {
		outdata.Error = "Could not create a session. " + err.Error()
		return
	}
	defer func() {
		if err := sess.Save(req, rw); err != nil {
			log.Printf("Error while saving session: %s", err)
			outdata.Success = ""
			outdata.Error = "Error while saving session."
			return
		}
	}()

	if user := userFromSess(sess); user != nil {
		outdata.Success = "You are already logged in"
		return
	}

	if req.Method != "POST" {
		return
	}

	if err := req.ParseForm(); err != nil {
		outdata.Error = "Data of form could not be understand. If this happens again, please contact support!"
		return
	}

	indata := new(loginFormdata)
	if err := formdec.Decode(indata, req.Form); (err != nil) || (indata.Mail == "") || (indata.Password == "") {
		outdata.Error = "Input data wrong or missing. Please fill in all values."
		return
	}

	user, err := dbcon.UserByMail(indata.Mail)
	switch err {
	case nil:
	case model.NotFound:
		outdata.Error = "E-Mail or password was wrong."
		return
	default:
		log.Printf("Error while loding user data (login): %s", err)
		outdata.Error = "User data could not be loaded. Please contact support, if this happens again."
		return
	}

	if bcrypt.CompareHashAndPassword(user.PWHash(), []byte(indata.Password)) != nil {
		outdata.Error = "E-Mail or password was wrong."
		return
	}

	sess.Values["uid"] = user.ID().String()
	outdata.Success = "Login successful"
}

func logincheck(rw http.ResponseWriter, req *http.Request) {
	sess, _ := getSess(req)
	user := userFromSess(sess)
	outdata := new(msgTpldata)
	if user == nil {
		outdata.Msg = "<nil>"
	} else {
		outdata.Msg = user.Email()
	}
	tplMsg.Execute(rw, outdata)
}

func logout(rw http.ResponseWriter, req *http.Request) {
	outdata := &msgTpldata{Class: "error", Title: "Logout"}
	defer func() {
		if err := tplMsg.Execute(rw, outdata); err != nil {
			log.Printf("Error executing template in login: %s", err)
		}
	}()

	sess, err := getSess(req)
	if err != nil {
		outdata.Msg = "Could not create a session."
		return
	}
	defer func() {
		if err := sess.Save(req, rw); err != nil {
			log.Printf("Error while saving session: %s", err)
			outdata.Class = "error"
			outdata.Msg = "Error while saving session."
			return
		}
	}()

	delete(sess.Values, "uid")
	outdata.Class = "success"
	outdata.Msg = "Your are now logged out."
}
