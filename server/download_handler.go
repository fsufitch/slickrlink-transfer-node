package server

import "net/http"

type downloadHandler struct{}

func (h downloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("hello there"))
}
