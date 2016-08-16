package main

import (
	"code.google.com/p/go.crypto/bcrypt"
	"github.com/gorilla/sessions"
	"github.com/silvasur/mailremind/model"
	"log"
	"net/http"
	"time"
)

type settingsTpldata struct {
	Success, Error string
	Fatal          bool
	Timezones      map[string]bool
}

func settings(user model.User, sess *sessions.Session, req *http.Request) (interface{}, model.User) {
	if user == nil {
		return &settingsTpldata{Error: "You need to be logged in to do that.", Fatal: true}, nil
	}

	outdata := &settingsTpldata{Timezones: make(map[string]bool)}
	tznow := user.Location().String()
	for _, tz := range timeLocs {
		outdata.Timezones[tz] = (tz == tznow)
	}

	if req.Method != "POST" {
		return outdata, user
	}

	if err := req.ParseForm(); err != nil {
		outdata.Error = "Could not parse form"
		return outdata, user
	}

	switch req.FormValue("M") {
	case "setpasswd":
		if req.FormValue("Password") == "" {
			outdata.Error = "Password must not be empty."
			return outdata, user
		}

		if req.FormValue("Password") != req.FormValue("RepeatPassword") {
			outdata.Error = "Passwords must be equal."
			return outdata, user
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(req.FormValue("Password")), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Error hashing password: %s", err)
			outdata.Error = "Error while saving password."
			return outdata.Error, user
		}

		if err := user.SetPWHash(hash); err != nil {
			log.Printf("Error setting pwhash: %s", err)
			outdata.Error = "Could not save new password."
		} else {
			outdata.Success = "Password changed"
		}
	case "settimezone":
		loc, err := time.LoadLocation(req.FormValue("Timezone"))
		if err != nil {
			outdata.Error = "Unknown Timezone"
			return outdata, user
		}

		if err := user.SetLocation(loc); err != nil {
			log.Printf("Error setting location: %s", err)
			outdata.Error = "Could not save new timezone."
		} else {
			outdata.Success = "New timezone saved."
		}
	}

	return outdata, user
}
