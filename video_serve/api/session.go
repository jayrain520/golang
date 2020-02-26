package main

import (
	"github.com/wonderivan/logger"
	"sync"
	"testgin/api/models"
	"time"
)

var allSession *sync.Map

func init() {
	loadAllSession()
}

func loadAllSession() {
	listSession, err := conn.LoadALLSessionID()
	if err != nil {
		logger.Alert(err) // system error load session
		panic(err)
	}

	allSession = &sync.Map{}

	for _, id := range *listSession {
		allSession.Store(id.UserName, id)
		logger.Debug("Session {Name: %v  ID: %v}", id.UserName, id.ID)
	}

}

//delete local and database session
func deleteSession(session *models.SessionID) {
	allSession.Delete(session.UserName)
	err := conn.DeleteExpireSessionID(session.ID)
	if err != nil {
		logger.Alert("delete session error: %v", err)
	}
}
func Day14_STime() int64 {
	return time.Now().Unix() + 1000000
}
func isExists(username string) (models.SessionID, bool) {
	v, ok := allSession.Load(username)
	if ok {
		s := v.(models.SessionID)
		if s.Expire < time.Now().Unix() {
			deleteSession(&s)
			return models.SessionID{}, false
		} else {
			logger.Debug("Exists ID: %v", v.(models.SessionID).ID)
			return s, true
		}
	}
	return models.SessionID{}, false
}
