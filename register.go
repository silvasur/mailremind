package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/gorilla/sessions"
	"kch42.de/gostuff/mailremind/model"
	"log"
	"math/rand"
	"net/http"
)

type registerData struct {
	Error, Success string
	Timezones      *[]string
}

type registerFormdata struct {
	Mail                     EMail
	Password, RetypePassword string
	Timezone                 timelocForm
}

var acCodeAlphabet = []rune("qwertzuiopasdfghjklyxcvbnmQWERTZUIOPASDFGHJKLYXCVBNM1234567890")

func genAcCode() string {
	const codelen = 10
	alphalen := len(acCodeAlphabet)

	code := make([]rune, codelen)
	for i := 0; i < codelen; i++ {
		code[i] = acCodeAlphabet[rand.Intn(alphalen)]
	}

	return string(code)
}

func register(user model.User, sess *sessions.Session, req *http.Request) interface{} {
	outdata := &registerData{Timezones: &timeLocs}

	if user != nil {
		outdata.Success = "You are already logged in. To register a new account, first log out."
		return outdata
	}

	if req.Method != "POST" {
		return outdata
	}

	if err := req.ParseForm(); err != nil {
		outdata.Error = "Data of form could not be understand. If this happens again, please contact support!"
		return outdata
	}

	indata := new(registerFormdata)
	if err := formdec.Decode(indata, req.Form); (err != nil) || (indata.Mail == "") || (indata.Timezone.Loc == nil) {
		outdata.Error = "Input data wrong or missing. Please fill in all values and make sure to provide a valid E-Mail address."
		return outdata
	}

	if indata.Password == "" {
		outdata.Error = "Empty passwords are not allowed."
		return outdata
	}

	if indata.Password != indata.RetypePassword {
		outdata.Error = "Passwords are not identical."
		return outdata
	}

	mail := string(indata.Mail)

	switch _, err := dbcon.UserByMail(mail); err {
	case nil:
		outdata.Error = "This E-Mail address is already used."
		return outdata
	case model.NotFound:
	default:
		log.Printf("Error while checking, if mail is used: %s", err)
		outdata.Error = "Internal error, sorry. If this happens again, please contact support!"
		return outdata
	}

	acCode := genAcCode()
	pwhash, err := bcrypt.GenerateFromPassword([]byte(indata.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Error while hashing password: %s", err)
		outdata.Error = "Internal error, sorry. If this happens again, please contact support!"
		return outdata
	}

	user, err = dbcon.AddUser(mail, pwhash, indata.Timezone.Loc, false, acCode)
	if err != nil {
		log.Printf("Could not create user (%s): %s", indata.Mail, err)
		outdata.Error = "Internal error, sorry. If this happens again, please contact support!"
		return outdata
	}

	if !SendActivationcode(mail, acCode, user.ID()) {
		outdata.Error = "We could not send you a mail with your confirmation code."
		return outdata
	}

	outdata.Success = "Account created successfully! We sent you an E-Mail that contains a link to activate your account."
	return outdata
}
