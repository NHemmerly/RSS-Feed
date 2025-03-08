package state

import (
	"context"

	"github.com/NHemmerly/RSS-Feed/internal/database"
)

func MiddlewareLoggedIn(handler func(s *State, cmd Command, user database.User) error) func(*State, Command) error {
	return func(s *State, cmd Command) error {
		user, err := s.Db.GetUser(context.Background(), s.Cfg.CurrentUserName)
		if err != nil {
			return err
		}
		return handler(s, cmd, user)
	}

}
