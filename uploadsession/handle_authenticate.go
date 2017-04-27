package uploadsession

import (
	"github.com/fsufitch/slickrlink-transfer-node/protobufs"
)

const stubKey = "xxx"

func (s *TransferSession) handleAuthentication(key string) {
	// TODO: implement proper auth
	if key == stubKey {
		s.authentication = authSuccess
		s.OutgoingMessages <- &protobufs.TransferNodeToClientMessage{
			Type: protobufs.TransferNodeToClientMessage_AUTH_SUCCESS,
		}
	} else {
		s.authentication = authFailed
		s.OutgoingMessages <- &protobufs.TransferNodeToClientMessage{
			Type: protobufs.TransferNodeToClientMessage_ERROR,
			ErrorData: &protobufs.ErrorData{
				Title: "Authentication failed",
				Fatal: true,
			},
		}
	}
}
