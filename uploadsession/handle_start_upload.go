package uploadsession

import (
	"github.com/fsufitch/slickrlink-transfer-node/protobufs"
)

func (s *TransferSession) handleStartUpload(filename string, size uint64, mimetype string) {
	s.metadata = TransferMetadata{
		Filename: filename,
		Size:     size,
		Mimetype: mimetype,
	}

	SessionDB.unregisterSession(s)
	SessionDB.registerSession(s)

	s.OutgoingMessages <- &protobufs.TransferNodeToClientMessage{
		Type: protobufs.TransferNodeToClientMessage_TRANSFER_CREATED,
		TransferCreatedData: &protobufs.TransferCreatedData{
			TransferId:    s.ID,
			ChunkSize:     s.chunkSize,
			RequestChunks: s.chunkBatchCount,
		},
	}
}
