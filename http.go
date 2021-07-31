package main

import "net/http"

func httpInteract(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	slashHandler(w, req)
}

func httpRegister(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if c := handleAuthorization(w, req); !c {
		return
	}
}

func httpUnregister(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	if c := handleAuthorization(w, req); !c {
		return
	}
}

func handleAuthorization(w http.ResponseWriter, req *http.Request) bool {
	auth := req.Header.Get("Authorization")
	if auth == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}
	if auth != token {
		w.WriteHeader(http.StatusForbidden)
		return false
	}
	return true
}
