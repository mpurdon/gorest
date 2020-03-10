package main

import (
	"fmt"
	"github.com/mpurdon/gorest/pkg"
	"github.com/pkg/errors"
	"log"
	"net/http"
)

func postCreate(w http.ResponseWriter, r *http.Request) {
	var p Post

	err := pkg.DecodeJSONBody(w, r, &p)
	if err != nil {
		var mr *pkg.MalformedRequest
		if errors.As(err, &mr) {
			http.Error(w, mr.Msg, mr.Status)
		} else {
			log.Println(err.Error())
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
		return
	}

	fmt.Fprintf(w, "Post: %+v", p)
}
