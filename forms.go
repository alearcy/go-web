package main

import (
	"regexp"
	"strings"
)

var rxEmail = regexp.MustCompile(".+@.+\\..+")

// LoginForm create login form struct
type LoginForm struct {
	Email string
	Password string
	Remember string
	Errors map[string]string
}

// Validate validates login form
func (l *LoginForm) Validate() bool {
	l.Errors = make(map[string]string)
	match := rxEmail.Match([]byte(l.Email))
	if match == false {
		l.Errors["Email"] = "Inserisci un indirizzo email valido"
	}
	if strings.TrimSpace(l.Email) == "" {
		l.Errors["Email"] = "Campo obbligatorio"
	}
	if strings.TrimSpace(l.Password) == "" {
		l.Errors["Password"] = "Campo obbligatorio"
	}
	
	return len(l.Errors) == 0
}

// SignupForm create signup form struct
type SignupForm struct {
	Email string
	Password string
	Password2 string
	Name string
	Surname string
	Errors map[string]string
}

// Validate validates login form
func (s *SignupForm) Validate() bool {
	s.Errors = make(map[string]string)
	match := rxEmail.Match([]byte(s.Email))
	if match == false {
		s.Errors["Email"] = "Inserisci un indirizzo email valido"
	}
	if strings.TrimSpace(s.Email) == "" {
		s.Errors["Email"] = "Campo obbligatorio"
	}
	if strings.TrimSpace(s.Password) == "" {
		s.Errors["Password"] = "Campo obbligatorio"
	}
	if strings.TrimSpace(s.Password2) == "" {
		s.Errors["Password2"] = "Campo obbligatorio"
	}
	if strings.TrimSpace(s.Name) == "" {
		s.Errors["Name"] = "Campo obbligatorio"
	}
	if strings.TrimSpace(s.Surname) == "" {
		s.Errors["Surname"] = "Campo obbligatorio"
	}
	
	return len(s.Errors) == 0
}