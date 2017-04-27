package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// StartServer starts the server
func StartServer(port int) {
	wsHandler := clientWebsocketHandler{upgrader: defaultUpgrader}

	router := mux.NewRouter()
	router.Handle("/client_ws", wsHandler)
	router.Handle("/d/{downloadId}", downloadHandler{})

	addr := fmt.Sprintf(":%d", port)
	http.ListenAndServe(addr, router)
}
