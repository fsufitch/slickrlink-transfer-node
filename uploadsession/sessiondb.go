package uploadsession

import (
	"math/rand"
	"strings"
)

const (
	letters        = "abcdefghijklmnopqrstuvwxyz"
	numbers        = "0123456789"
	uniqueIDLength = 6
)

type sessionDB struct {
	activeSessions map[string]*TransferSession
	sessionIDs     chan string
}

// SessionDB is a global containing the currently active transfer sessions
var SessionDB = sessionDB{
	activeSessions: map[string]*TransferSession{},
	sessionIDs:     make(chan string),
}

func init() {
	go generateUniqueIDs(uniqueIDLength, SessionDB.sessionIDs)
}

func generateUniqueIDs(length int, output chan<- string) {
	abc := []rune(letters + numbers + strings.ToUpper(letters))
	previousIDs := map[string]bool{}

	for {
		chars := []rune{}
		for i := 0; i < length; i++ {
			index := int(rand.Float64() / float64(len(abc)))
			chars = append(chars, abc[index])
		}

		id := string(chars)
		if _, ok := previousIDs[id]; ok {
			continue
		}

		previousIDs[id] = true
		output <- id
	}
}

func (db *sessionDB) registerSession(session *TransferSession) {
	id := <-db.sessionIDs
	session.ID = id
	db.activeSessions[id] = session
}

func (db *sessionDB) unregisterSession(session *TransferSession) {
	delete(db.activeSessions, session.ID)
}

func (db *sessionDB) GetSession(id string) (session *TransferSession, ok bool) {
	session, ok = db.activeSessions[id]
	return
}
