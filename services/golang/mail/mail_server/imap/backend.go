package imap

import (
	"log"

	"context"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
)

type Backend_impl struct {
}

func (self *Backend_impl) Login(connInfo *imap.ConnInfo, username, password string) (backend.User, error) {
	log.Println("----> try to authenticate the user ", username)
	// I will use the datastore to authenticate the user.
	connection_id := username + "_db"
	err := store.Connect(connection_id, "0.0.0.0", 27017, username, password, connection_id, 1000, "")
	if err != nil {
		log.Println("fail to login: ", err)
		return nil, err
	}

	// retreive account info.
	query := `{"name":"` + username + `"}`
	info, err := store.FindOne(context.Background(), "local_ressource", "local_ressource", "Accounts", query, "")
	if err != nil {
		log.Println("fail to retreive account ", username)
		return nil, err
	}

	user := new(User_impl)
	user.info = info.(map[string]interface{})

	return user, nil
}
