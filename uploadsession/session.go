package uploadsession

import (
	"log"

	"github.com/fsufitch/slickrlink-transfer-node/protobufs"
)

type authState int

const (
	authNotAttempted authState = iota
	authSuccess
	authFailed
)

const (
	uploadChunkSize  = 50000 // bytes
	uploadChunkCount = 1
)

// TransferMetadata encapsulates details about the file being uploaded
type TransferMetadata struct {
	Filename string
	Size     uint64
	Mimetype string
}

// TransferSession encapsulates state data having to do with the current transfer session
type TransferSession struct {
	IncomingMessages chan *protobufs.ClientToTransferNodeMessage
	OutgoingMessages chan *protobufs.TransferNodeToClientMessage
	FatalErrors      chan error
	ID               string

	authentication  authState
	metadata        TransferMetadata
	chunkSize       uint64
	chunkBatchCount uint64
	uploadBytes     uint64
}

// NewTransferSession creates and registers a new transfer session
func NewTransferSession() (session *TransferSession) {
	session = &TransferSession{
		IncomingMessages: make(chan *protobufs.ClientToTransferNodeMessage),
		OutgoingMessages: make(chan *protobufs.TransferNodeToClientMessage),
		FatalErrors:      make(chan error),
		chunkSize:        uploadChunkSize,
		chunkBatchCount:  uploadChunkCount,
	}

	return
}

// HandleMessages loops through all incoming messages in the channel and handles them
func (s *TransferSession) HandleMessages() {
	for msg := range s.IncomingMessages {
		switch msg.GetType() {
		case protobufs.ClientToTransferNodeMessage_AUTHENTICATE:
			s.handleAuthentication(msg.GetAuthData().GetKey())
		default:
			log.Printf("Unknown message type: %v", msg.GetType())
		}
	}
}

// Destroy cleans up the session's stuff
func (s *TransferSession) Destroy() {
	close(s.IncomingMessages)
	close(s.OutgoingMessages)
	close(s.FatalErrors)
	SessionDB.unregisterSession(s)
}
