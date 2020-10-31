package imap

import (
	"errors"

	"context"
	// "io/ioutil"
	"io/ioutil"
	"log"
	"time"

	"github.com/davecourtois/Utility"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/backend/backendutil"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MailBox_impl struct {
	name       string
	user       string
	subscribed bool
}

func NewMailBox(user string, name string) *MailBox_impl {
	box := new(MailBox_impl)
	box.name = name
	box.user = user

	info := new(imap.MailboxInfo)
	info.Name = name
	info.Delimiter = "/"

	connectionId := user + "_db"

	// Here I will retreive the info from the database.
	query := `{"name":"` + name + `"}`
	count, err := store.Count(context.Background(), connectionId, connectionId, "MailBoxes", query, "")
	if err != nil || count < 1 {
		_, err = store.InsertOne(context.Background(), connectionId, connectionId, "MailBoxes", info, "")
		if err != nil {
			log.Println("fail to create mail box!")
		}
	}

	return box
}

func getMailBox(user string, name string) (*MailBox_impl, error) {
	box := new(MailBox_impl)
	box.name = name
	box.user = user

	info := new(imap.MailboxInfo)
	info.Name = name
	info.Delimiter = "/"

	connectionId := user + "_db"

	// Here I will retreive the info from the database.
	query := `{"name":"` + name + `"}`
	count, err := store.Count(context.Background(), connectionId, connectionId, "MailBoxes", query, "")
	if err != nil || count < 1 {
		return nil, errors.New("No mail box found with name " + name)
	}

	return box, nil
}

// Name returns this mailbox name.
func (mbox *MailBox_impl) Name() string {
	return mbox.name
}

// Info returns this mailbox info.
func (mbox *MailBox_impl) Info() (*imap.MailboxInfo, error) {

	// TODO Get box info from the server.
	connectionId := mbox.user + "_db"
	query := `{"name":"` + mbox.name + `"}`
	info_, err := store.FindOne(context.Background(), connectionId, connectionId, "MailBoxes", query, "")

	if err == nil {
		// Now I will insert the message into the inbox of the user.
		info := new(imap.MailboxInfo)
		info.Name = info_.(map[string]interface{})["name"].(string)
		info.Delimiter = info_.(map[string]interface{})["delimiter"].(string)
		if info_.(map[string]interface{})["attributes"] != nil {
			log.Println("attributes ", info_.(map[string]interface{})["attributes"])
		}

		return info, nil
	}
	return nil, err
}

func (mbox *MailBox_impl) uidNext() uint32 {
	var uid uint32
	messages := mbox.getMessages()
	for _, msg := range messages {
		if msg.Uid > uid {
			uid = msg.Uid
		}
	}
	uid++
	return uid
}

// Return the list of message from the bakend.
func (mbox *MailBox_impl) getMessages() []*Message {

	messages := make([]*Message, 0)
	connectionId := mbox.user + "_db"

	// Get the message from the mailbox.
	data, _ := store.Find(context.Background(), connectionId, connectionId, mbox.Name(), "{}", "")

	// return the messages.
	for i := 0; i < len(data); i++ {
		msg := data[i].(map[string]interface{})
		m := new(Message)
		m.Uid = uint32(msg["Uid"].(int64))           // set the actual index
		m.Body = msg["Body"].(primitive.Binary).Data // set back to []uint8
		flags := []interface{}(msg["Flags"].(primitive.A))
		m.Flags = make([]string, 0)
		for j := 0; j < len(flags); j++ {
			if len(flags[j].(string)) > 0 {
				m.Flags = append(m.Flags, flags[j].(string))
			}
		}

		m.Size = uint32(msg["Size"].(int64))
		m.Date = msg["Date"].(primitive.DateTime).Time()
		messages = append(messages, m)
	}

	return messages
}

func (mbox *MailBox_impl) unseenSeqNum() uint32 {
	messages := mbox.getMessages()
	for i, msg := range messages {
		seqNum := uint32(i + 1)
		seen := false
		for _, flag := range msg.Flags {
			if flag == imap.SeenFlag {
				seen = true
				break
			}
		}

		if !seen {
			return seqNum
		}
	}
	return 0
}

func (mbox *MailBox_impl) flags() []string {
	flagsMap := make(map[string]bool)
	messages := mbox.getMessages()
	for _, msg := range messages {
		for _, f := range msg.Flags {
			if !flagsMap[f] {
				flagsMap[f] = true
			}
		}
	}

	var flags []string
	for f := range flagsMap {
		flags = append(flags, f)
	}
	return flags
}

// Status returns this mailbox status. The fields Name, Flags, PermanentFlags
// and UnseenSeqNum in the returned MailboxStatus must be always populated.
// This function does not affect the state of any messages in the mailbox. See
// RFC 3501 section 6.3.10 for a list of items that can be requested.
func (mbox *MailBox_impl) Status(items []imap.StatusItem) (*imap.MailboxStatus, error) {
	status := imap.NewMailboxStatus(mbox.name, items)
	status.Flags = mbox.flags()
	status.PermanentFlags = []string{"\\*"}
	status.UnseenSeqNum = mbox.unseenSeqNum()

	for _, name := range items {
		switch name {
		case imap.StatusMessages:
			status.Messages = uint32(len(mbox.getMessages()))
		case imap.StatusUidNext:
			status.UidNext = mbox.uidNext()
		case imap.StatusUidValidity:
			status.UidValidity = 1
		case imap.StatusRecent:
			status.Recent = 0 // TODO
		case imap.StatusUnseen:
			status.Unseen = 0 // TODO
		}
	}

	return status, nil
}

// SetSubscribed adds or removes the mailbox to the server's set of "active"
// or "subscribed" mailboxes.
func (mbox *MailBox_impl) SetSubscribed(subscribed bool) error {
	mbox.subscribed = subscribed
	return nil
}

// Check requests a checkpoint of the currently selected mailbox. A checkpoint
// refers to any implementation-dependent housekeeping associated with the
// mailbox (e.g., resolving the server's in-memory state of the mailbox with
// the state on its disk). A checkpoint MAY take a non-instantaneous amount of
// real time to complete. If a server implementation has no such housekeeping
// considerations, CHECK is equivalent to NOOP.
func (mbox *MailBox_impl) Check() error {
	return nil
}

// ListMessages returns a list of messages. seqset must be interpreted as UIDs
// if uid is set to true and as message sequence numbers otherwise. See RFC
// 3501 section 6.4.5 for a list of items that can be requested.
//
// Messages must be sent to ch. When the function returns, ch must be closed.
func (mbox *MailBox_impl) ListMessages(uid bool, seqSet *imap.SeqSet, items []imap.FetchItem, ch chan<- *imap.Message) error {
	defer close(ch)
	messages := mbox.getMessages()

	for i, msg := range messages {
		seqNum := uint32(i + 1)

		var id uint32
		if uid {
			id = msg.Uid
		} else {
			id = seqNum
		}
		if !seqSet.Contains(id) {
			continue
		}

		m, err := msg.Fetch(seqNum, items)
		if err != nil {
			continue
		}

		ch <- m
	}

	return nil
}

// SearchMessages searches messages. The returned list must contain UIDs if
// uid is set to true, or sequence numbers otherwise.
func (mbox *MailBox_impl) SearchMessages(uid bool, criteria *imap.SearchCriteria) ([]uint32, error) {
	var ids []uint32
	messages := mbox.getMessages()
	for i, msg := range messages {
		seqNum := uint32(i + 1)

		ok, err := msg.Match(seqNum, criteria)
		if err != nil || !ok {
			continue
		}

		var id uint32
		if uid {
			id = msg.Uid
		} else {
			id = seqNum
		}
		ids = append(ids, id)
	}
	return ids, nil
}

// CreateMessage appends a new message to this mailbox. The \Recent flag will
// be added no matter flags is empty or not. If date is nil, the current time
// will be used.
//
// If the Backend implements Updater, it must notify the client immediately
// via a mailbox update.
func (mbox *MailBox_impl) CreateMessage(flags []string, date time.Time, body imap.Literal) error {
	if date.IsZero() {
		date = time.Now()
	}

	b, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	return saveMessage(mbox.user, mbox.name, b, flags, date)
}

// UpdateMessagesFlags alters flags for the specified message(s).
//
// If the Backend implements Updater, it must notify the client immediately
// via a message update.
func (mbox *MailBox_impl) UpdateMessagesFlags(uid bool, seqset *imap.SeqSet, op imap.FlagsOp, flags []string) error {
	messages := mbox.getMessages()
	log.Println("------> 303", flags)
	for i, msg := range messages {
		var id uint32
		if uid {
			id = msg.Uid
		} else {
			id = uint32(i + 1)
		}
		if !seqset.Contains(id) {
			continue
		}

		msg.Flags = backendutil.UpdateFlags(msg.Flags, op, flags)

		// Here I will save the message into the database.
		connectionId := mbox.user + "_db"
		jsonStr, _ := Utility.ToJson(msg.Flags)

		err := store.UpdateOne(context.Background(), connectionId, connectionId, mbox.name, `{"Uid":`+Utility.ToString(msg.Uid)+`}`, `{ "$set":{"Flags":`+jsonStr+`}}`, "")
		if err != nil {
			return err
		}

	}
	return nil
}

// CopyMessages copies the specified message(s) to the end of the specified
// destination mailbox. The flags and internal date of the message(s) SHOULD
// be preserved, and the Recent flag SHOULD be set, in the copy.
//
// If the destination mailbox does not exist, a server SHOULD return an error.
// It SHOULD NOT automatically create the mailbox.
//
// If the Backend implements Updater, it must notify the client immediately
// via a mailbox update.
func (mbox *MailBox_impl) CopyMessages(uid bool, seqset *imap.SeqSet, destName string) error {
	dest, err := getMailBox(mbox.user, destName)
	if err != nil {
		return err
	}

	messages := mbox.getMessages()
	for i, msg := range messages {
		var id uint32
		if uid {
			id = msg.Uid
		} else {
			id = uint32(i + 1)
		}
		if !seqset.Contains(id) {
			continue
		}

		// save the message in the backend.
		saveMessage(dest.user, dest.name, msg.Body, msg.Flags, time.Now())
	}

	return nil
}

// Expunge permanently removes all messages that have the \Deleted flag set
// from the currently selected mailbox.
//
// If the Backend implements Updater, it must notify the client immediately
// via an expunge update.
func (mbox *MailBox_impl) Expunge() error {
	messages := mbox.getMessages()

	for i := len(messages) - 1; i >= 0; i-- {
		msg := messages[i]

		deleted := false
		for _, flag := range msg.Flags {
			if flag == imap.DeletedFlag {
				deleted = true
				break
			}
		}

		if deleted {
			// mbox.Messages = append(mbox.Messages[:i], mbox.Messages[i+1:]...)
			connectionId := mbox.user + "_db"

			err := store.DeleteOne(context.Background(), connectionId, connectionId, mbox.name, `{"Uid":`+Utility.ToString(msg.Uid)+`}`, "")
			if err != nil {
				return err
			}
		}
	}

	return nil
}
