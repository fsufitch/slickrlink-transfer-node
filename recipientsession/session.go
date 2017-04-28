package recipientsession

import "github.com/fsufitch/slickrlink-transfer-node/uploadsession"

// RecipientSession holds data and references to an ongoing download
type RecipientSession struct {
	ErrorNotFound chan bool
	ErrorGone     chan bool
	FatalError    chan error
	Done          chan bool

	MetadataChannel chan uploadsession.TransferMetadata
	ByteChannel     chan []byte

	transferSession *uploadsession.TransferSession
}

// NewRecipientSession creates a new Recipient session and starts asynchronous processing for it
func NewRecipientSession(transferSessionID string, ipv4 string, ipv6 string, identity string) (session *RecipientSession) {
	session = &RecipientSession{
		ErrorNotFound: make(chan bool),
		ErrorGone:     make(chan bool),
		FatalError:    make(chan error),
		Done:          make(chan bool),

		MetadataChannel: make(chan uploadsession.TransferMetadata),
		ByteChannel:     make(chan []byte),
	}
	go session.startDownload(transferSessionID, ipv4, ipv6, identity)
	return
}

func (s *RecipientSession) startDownload(transferSessionID string, ipv4 string, ipv6 string, identity string) {
	transferSession, ok := uploadsession.SessionDB.GetSession(transferSessionID)
	if !ok {
		s.ErrorNotFound <- true
		return
	}

	if transferSession.UploadedBytes > 0 {
		// Upload already started, too late
		s.ErrorGone <- true
		return
	}

	s.MetadataChannel <- transferSession.Metadata
	transferSession.AddRecipient(ipv4, ipv6, identity, s.ByteChannel)
}
