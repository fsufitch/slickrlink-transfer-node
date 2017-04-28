package server

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"

	"github.com/fsufitch/slickrlink-transfer-node/recipientsession"
)

type downloadHandler struct{}

func (h downloadHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	transferSessionID, ok := mux.Vars(r)["downloadId"]
	if !ok {
		transferSessionID = ""
	}

	ipv4, ipv6 := extractIPAddress(r)

	session := recipientsession.NewRecipientSession(transferSessionID, ipv4, ipv6, "")
	headersSent := false

	for {
		select {
		case <-session.ErrorNotFound:
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(404)
			w.Write([]byte("Upload not found."))
			return

		case <-session.ErrorGone:
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(410)
			w.Write([]byte("Upload found, but already started."))
			return

		case err := <-session.FatalError:
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(500)
			message := fmt.Sprintf("Server error: %v", err.Error())
			w.Write([]byte(message))
			return

		case metadata := <-session.MetadataChannel:
			if headersSent {
				session.FatalError <- errors.New("Received metadata multiple times")
				break
			}
			w.Header().Set("Content-Type", metadata.Mimetype)
			w.Header().Set("Content-Disposition", createContentDisposition(metadata.Filename))
			w.Header().Set("Content-Length", fmt.Sprint(metadata.Size))
			w.WriteHeader(200)

			w.Write([]byte{}) // Write 0 bytes to force sending headers
			headersSent = true

		case bytes, isOpen := <-session.ByteChannel:
			if !headersSent {
				session.FatalError <- errors.New("Received file data before metadata")
				break
			}

			if !isOpen {
				return
			}

			w.Write(bytes)
		}
	}
}

func extractIPAddress(r *http.Request) (ipv4 string, ipv6 string) {
	ip := net.ParseIP(r.RemoteAddr) // XXX: Not always accurate, does not check LB headers
	if ip.To4() != nil {
		// ipv4
		return ip.String(), ""
	}
	if ip.To16() != nil {
		// ipv6
		return "", ip.String()
	}
	return "", ""
}

func createContentDisposition(filename string) string {
	escapedFilename := url.PathEscape(filename)
	return fmt.Sprintf("attachment; filename=%s", escapedFilename)
}
