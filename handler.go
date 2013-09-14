package main

import (
	"github.com/gorilla/sessions"
	"html/template"
	"kch42.de/gostuff/mailremind/model"
	"log"
	"net/http"
)

type Handler func(user model.User, sess *sessions.Session, req *http.Request) (interface{}, model.User)

func getSess(req *http.Request) (*sessions.Session, error) {
	sess, err := SessionStorage.Get(req, "mailremind-sess")
	if err != nil {
		sess, err = SessionStorage.New(req, "mailremind-sess")
	}

	return sess, err
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

type tpldata struct {
	Mail string
	Data interface{}
}

func mkHttpHandler(h Handler, tpl *template.Template) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		sess, err := getSess(req)
		if err != nil {
			req.Cookies()
			log.Printf("Error while getting session: %s", err)
			rw.Write([]byte("Unable to create a session. Try reloading the page or delete the cookies for this site."))
		}

		user := userFromSess(sess)
		outdata, user := h(user, sess, req)
		mail := ""
		if user != nil {
			mail = user.Email()
		}

		if err := sess.Save(req, rw); err != nil {
			log.Printf("Error while saving session: %s", err)
		}

		if err := tpl.Execute(rw, &tpldata{mail, outdata}); err != nil {
			log.Printf("Error executing template %s: %s", tpl.Name(), err)
		}
	}
}
