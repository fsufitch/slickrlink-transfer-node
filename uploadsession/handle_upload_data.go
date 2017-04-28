package uploadsession

import (
	"errors"

	"github.com/fsufitch/slickrlink-transfer-node/protobufs"
)

func (s *TransferSession) handleUploadData(data []byte, checkSize uint64, order uint64) {
	if uint64(len(data)) != checkSize {
		s.FatalErrors <- errors.New("Incoming data size mismatch")
		return
	}

	if s.UploadedBytes+checkSize > s.Metadata.Size {
		s.FatalErrors <- errors.New("Too much data uploaded")
		return
	}

	s.UploadedBytes += checkSize
	dataCopy := append([]byte{}, data...)
	for _, recipient := range s.recipients {
		recipient.dataStream <- dataCopy
	}

	s.OutgoingMessages <- &protobufs.TransferNodeToClientMessage{
		Type: protobufs.TransferNodeToClientMessage_PROGRESS,
		ProgressData: &protobufs.ProgressData{
			BytesUploaded: int64(s.UploadedBytes),
			RequestChunks: s.chunkBatchCount,
			ChunkSize:     s.chunkSize,
		},
	}
}
