package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/gorilla/sessions"
	"github.com/silvasur/mailremind/model"
	"log"
	"net/http"
)

type pwresetTpldata struct {
	Error, Success, Code, UID string
}

func pwreset(user model.User, sess *sessions.Session, req *http.Request) (interface{}, model.User) {
	if err := req.ParseForm(); err != nil {
		return &pwresetTpldata{Error: "Form data corrupted."}, user
	}

	code := req.FormValue("Code")
	_uid := req.FormValue("U")
	pw1 := req.FormValue("Password")
	pw2 := req.FormValue("PasswordAgain")

	if code == "" {
		return &pwresetTpldata{Error: "Wrong password reset code"}, user
	}

	uid, err := db.ParseDBID(_uid)
	if err != nil {
		return &pwresetTpldata{Error: "Invalid user ID"}, user
	}

	if user, err = dbcon.UserByID(uid); err != nil {
		return &pwresetTpldata{Error: "User not found"}, user
	}

	if user.ActivationCode() != code {
		return &pwresetTpldata{Error: "Wrong activation code"}, user
	}

	outdata := &pwresetTpldata{UID: _uid, Code: code}

	if req.Method != "POST" {
		return outdata, user
	}

	if pw1 == "" {
		outdata.Error = "Password must not be empty."
		return outdata, user
	}

	if pw1 != pw2 {
		outdata.Error = "Passwords are not identical."
		return outdata, user
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(pw1), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Could not has password: %s", err)
		outdata.Error = "Could not save password."
		return outdata, user
	}

	if err := user.SetPWHash(hash); err != nil {
		log.Printf("Error while hashing password: %s", err)
		outdata.Error = "Could not save password."
		return outdata, user
	}

	if err := user.SetActivationCode(""); err != nil {
		log.Printf("Error resetting acCode: %s", err)
	}

	outdata.Success = "Password was changed"
	return outdata, user
}

type forgotpwTpldata struct {
	Error, Success string
}

func forgotpw(user model.User, sess *sessions.Session, req *http.Request) (interface{}, model.User) {
	if req.Method != "POST" {
		return &forgotpwTpldata{}, user
	}

	if err := req.ParseForm(); err != nil {
		return &forgotpwTpldata{Error: "Form data corrupted."}, user
	}

	email := req.FormValue("Mail")
	if email == "" {
		return &forgotpwTpldata{Error: "E-Mail must not be empty."}, user
	}

	user, err := dbcon.UserByMail(email)
	if err != nil {
		return &forgotpwTpldata{Error: "E-Mail not found."}, user
	}

	key := genAcCode()
	if err := user.SetActivationCode(key); err != nil {
		log.Printf("Could not store pwreset key: %s", err)
		return &forgotpwTpldata{Error: "Could not generate a keyword reset code."}, user
	}

	if !SendPwresetLink(user.Email(), key, user.ID()) {
		return &forgotpwTpldata{Error: "Could not send reset E-Mail."}, user
	}

	return &forgotpwTpldata{Success: "We sent you an E-Mail with further instructions."}, user
}
