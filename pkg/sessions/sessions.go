package sessions

import (
	"fmt"
	"time"

	"github.com/robbymilo/rgallery/pkg/types"
)

type Session struct {
	UserName string
	Role     string
	expiry   time.Time
}

type Conf = types.Conf

var sessions = map[string]Session{}

func (s Session) IsExpired() bool {
	return s.expiry.Before(time.Now())
}

func CreateSession(username, role, token string, expiry time.Time, c Conf) error {
	if role == "" {
		return fmt.Errorf("session contains no role")
	}

	c.Logger.Info("adding session: " + username + " role: " + role)
	sessions[token] = Session{
		UserName: username,
		Role:     role,
		expiry:   expiry,
	}
	return nil
}

func GetSession(token string) (Session, bool) {
	userSession, exists := sessions[token]
	if !exists {
		return userSession, false
	}

	return userSession, true
}

func DeleteSession(token string, c Conf) {
	c.Logger.Info("deleting session: " + sessions[token].UserName)
	delete(sessions, token)
}

func DeleteUserSessions(username string) {
	for k := range sessions {
		if sessions[k].UserName == username {
			fmt.Println("deleting session:", sessions[k].UserName)
			delete(sessions, k)
		}
	}

}
