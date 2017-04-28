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

	UploadedBytes uint64
	Metadata      TransferMetadata

	authentication  authState
	chunkSize       uint64
	chunkBatchCount uint64
	recipients      []transferRecipient
}

type transferRecipient struct {
	ipv4       string
	ipv6       string
	identity   string
	dataStream chan []byte
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
		case protobufs.ClientToTransferNodeMessage_START_UPLOAD:
			s.handleStartUpload(
				msg.StartData.GetFilename(),
				msg.StartData.GetSize(),
				msg.StartData.GetMimetype(),
			)
		default:
			log.Printf("Unknown message type: %v", msg.GetType())
		}
	}
}

// AddRecipient registers a new destination recipient for this upload
func (s *TransferSession) AddRecipient(ipv4 string, ipv6 string, identity string, dataStream chan []byte) {
	s.recipients = append(s.recipients, transferRecipient{
		ipv4:       ipv4,
		ipv6:       ipv6,
		identity:   identity,
		dataStream: dataStream,
	})

	recipients := []*protobufs.RecipientsData_Recipient{}
	for _, recipient := range s.recipients {
		recipients = append(recipients, &protobufs.RecipientsData_Recipient{
			Ipv4:     recipient.ipv4,
			Ipv6:     recipient.ipv6,
			Identity: recipient.identity,
		})
	}

	s.OutgoingMessages <- &protobufs.TransferNodeToClientMessage{
		Type: protobufs.TransferNodeToClientMessage_RECIPIENTS,
		RecipientsData: &protobufs.RecipientsData{
			Recipients: recipients,
		},
	}
}

// Destroy cleans up the session's stuff
func (s *TransferSession) Destroy() {
	close(s.IncomingMessages)
	close(s.OutgoingMessages)
	close(s.FatalErrors)
	SessionDB.unregisterSession(s)
}
