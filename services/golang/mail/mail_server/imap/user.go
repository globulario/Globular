package imap

import (
	"errors"
	//	"log"
	"context"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend"
)

////////////////////////////////////////////////////////////////////////////////
// The User implementation
////////////////////////////////////////////////////////////////////////////////

type User_impl struct {
	// contain values from mongoDB
	info map[string]interface{}
}

// Username returns this user's username.
func (self *User_impl) Username() string {
	return self.info["name"].(string)
}

// ListMailboxes returns a list of mailboxes belonging to this user. If
// subscribed is set to true, only returns subscribed mailboxes.
func (self *User_impl) ListMailboxes(subscribed bool) ([]backend.Mailbox, error) {
	boxes := make([]backend.Mailbox, 0)
	connectionId := self.Username() + "_db"
	// Get list of other mail box here.
	values, err := store.Find(context.Background(), connectionId, connectionId, "MailBoxes", `{}`, ``)
	if err == nil {
		for i := 0; i < len(values); i++ {
			val := values[i].(map[string]interface{})
			box := NewMailBox(self.Username(), val["name"].(string))
			boxes = append(boxes, box)
		}
	}

	if len(values) == 0 {
		// By default INBOX must exist on the backend, that the default mailbox.
		inbox := NewMailBox(self.Username(), "INBOX")
		boxes = append(boxes, inbox)

	}

	return boxes, nil
}

// GetMailbox returns a mailbox. If it doesn't exist, it returns
// ErrNoSuchMailbox.
func (self *User_impl) GetMailbox(name string) (backend.Mailbox, error) {

	connectionId := self.Username() + "_db"
	query := `{"name":"` + name + `"}`
	count, err := store.Count(context.Background(), connectionId, connectionId, "MailBoxes", query, "")
	if err != nil || count < 1 {
		return nil, errors.New("No mail box found with name " + name)
	}

	return NewMailBox(self.Username(), name), nil
}

// CreateMailbox creates a new mailbox.
//
// If the mailbox already exists, an error must be returned. If the mailbox
// name is suffixed with the server's hierarchy separator character, this is a
// declaration that the client intends to create mailbox names under this name
// in the hierarchy.
//
// If the server's hierarchy separator character appears elsewhere in the
// name, the server SHOULD create any superior hierarchical names that are
// needed for the CREATE command to be successfully completed.  In other
// words, an attempt to create "foo/bar/zap" on a server in which "/" is the
// hierarchy separator character SHOULD create foo/ and foo/bar/ if they do
// not already exist.
//
// If a new mailbox is created with the same name as a mailbox which was
// deleted, its unique identifiers MUST be greater than any unique identifiers
// used in the previous incarnation of the mailbox UNLESS the new incarnation
// has a different unique identifier validity value.
func (self *User_impl) CreateMailbox(name string) error {

	info := new(imap.MailboxInfo)
	info.Name = name
	info.Delimiter = "/"

	connectionId := self.Username() + "_db"
	_, err := store.InsertOne(context.Background(), connectionId, connectionId, "MailBoxes", info, "")

	return err
}

// DeleteMailbox permanently remove the mailbox with the given name. It is an
// error to // attempt to delete INBOX or a mailbox name that does not exist.
//
// The DELETE command MUST NOT remove inferior hierarchical names. For
// example, if a mailbox "foo" has an inferior "foo.bar" (assuming "." is the
// hierarchy delimiter character), removing "foo" MUST NOT remove "foo.bar".
//
// The value of the highest-used unique identifier of the deleted mailbox MUST
// be preserved so that a new mailbox created with the same name will not
// reuse the identifiers of the former incarnation, UNLESS the new incarnation
// has a different unique identifier validity value.
func (self *User_impl) DeleteMailbox(name string) error {
	connectionId := self.Username() + "_db"

	// First I will delete the entry in MailBoxes...
	err := store.DeleteOne(context.Background(), connectionId, connectionId, "MailBoxes", `{"name":"`+name+`"}`, "")
	if err != nil {
		return err
	}

	err = store.DeleteCollection(context.Background(), connectionId, connectionId, name)

	return err
}

// RenameMailbox changes the name of a mailbox. It is an error to attempt to
// rename from a mailbox name that does not exist or to a mailbox name that
// already exists.
//
// If the name has inferior hierarchical names, then the inferior hierarchical
// names MUST also be renamed.  For example, a rename of "foo" to "zap" will
// rename "foo/bar" (assuming "/" is the hierarchy delimiter character) to
// "zap/bar".
//
// If the server's hierarchy separator character appears in the name, the
// server SHOULD create any superior hierarchical names that are needed for
// the RENAME command to complete successfully.  In other words, an attempt to
// rename "foo/bar/zap" to baz/rag/zowie on a server in which "/" is the
// hierarchy separator character SHOULD create baz/ and baz/rag/ if they do
// not already exist.
//
// The value of the highest-used unique identifier of the old mailbox name
// MUST be preserved so that a new mailbox created with the same name will not
// reuse the identifiers of the former incarnation, UNLESS the new incarnation
// has a different unique identifier validity value.
//
// Renaming INBOX is permitted, and has special behavior.  It moves all
// messages in INBOX to a new mailbox with the given name, leaving INBOX
// empty.  If the server implementation supports inferior hierarchical names
// of INBOX, these are unaffected by a rename of INBOX.
func (self *User_impl) RenameMailbox(existingName, newName string) error {
	connectionId := self.Username() + "_db"

	// I will rename the
	err := store.UpdateOne(context.Background(), connectionId, connectionId, "MailBoxes", `{"name":"`+existingName+`"}`, `{"$set":{"name":"`+newName+`"}}`, "")
	if err != nil {
		return err
	}

	return renameCollection(connectionId, existingName, newName)
}

// Logout is called when this User will no longer be used, likely because the
// client closed the connection.
func (self *User_impl) Logout() error {
	return nil //store.Disconnect(self.Username() + "_db")
}
