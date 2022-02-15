package matrix

import (
	"crypto/sha512"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"gitlab.com/mitaka8/shulker/matrix/store"

	"github.com/sethvargo/go-retry"
	"gitlab.com/mitaka8/shulker/registry"
	"maunium.net/go/mautrix"
	"maunium.net/go/mautrix/crypto"
	"maunium.net/go/mautrix/event"
	"maunium.net/go/mautrix/format"
	"maunium.net/go/mautrix/id"
)

type memberInfo struct {
	userId      id.UserID
	displayname string
	avatarUrl   string
}

type messageInfo struct {
	author  *memberInfo
	message string
	roomId  id.RoomID
}
type messageListener func(*messageInfo)

type matrixClient struct {
	client           *mautrix.Client
	messageListeners map[id.RoomID][]messageListener
}

func (c *matrixClient) listenToMessages(roomId id.RoomID, fn messageListener) {
	if c.messageListeners == nil {
		c.messageListeners = make(map[id.RoomID][]messageListener)
	}
	if c.messageListeners[roomId] == nil {
		c.messageListeners[roomId] = []messageListener{fn}
	} else {
		c.messageListeners[roomId] = append(c.messageListeners[roomId], fn)
	}
}

var clientslock sync.Mutex
var clients = make(map[string]*matrixClient)

var memberslock sync.Mutex
var members = make(map[string]*memberInfo)

func formatMessage(message registry.ChatMessage) event.MessageEventContent {

	var sb strings.Builder
	sb.WriteString("#*")
	sb.WriteString(message.Source())
	sb.WriteString("* <**")
	sb.WriteString(message.Author().Name())
	sb.WriteString("**> ")
	sb.WriteString(message.Message())

	return format.RenderMarkdown(sb.String(), true, false)
}

func formatGenericMessage(message registry.GenericMessage) string {
	var sb strings.Builder
	sb.WriteString("**[")
	sb.WriteString(message.Source())
	sb.WriteString("]** ")
	sb.WriteString(message.Message())
	return sb.String()
}

func makeClientId(hostname string, username string, password string) string {
	var sb strings.Builder
	sb.WriteString(hostname)
	sb.WriteString("~")
	sb.WriteString(username)
	sb.WriteString("$")
	sb.WriteString(password)

	clientId := sha512.Sum512([]byte(sb.String()))


	return string(clientId[:7])

}

func getOrMakeClient(homeserver, username, password, database, key, deviceId string) *matrixClient {
	clientid := makeClientId(homeserver, username, password)

	clientslock.Lock()
	if client, ok := clients[clientid]; ok {
		clientslock.Unlock()
		return client, nil, nil
	}
	clientslock.Unlock()

	client, err := mautrix.NewClient(homeserver, "", "")

	if err != nil {
		return client, nil, err
	}

	db, err := sql.Open("sqlite3", database)
	if err != nil {
		return client, nil, err
	}

	if deviceId == "" {
		deviceId = findDeviceId(db, id.UserID(username)).String()
	}
	_, err = client.Login(&mautrix.ReqLogin{
		Type:     mautrix.AuthTypePassword,
		Password: password,

		Identifier: mautrix.UserIdentifier{
			Type: mautrix.IdentifierTypeUser,
			User: username,
		},

		StoreCredentials:   true,
		StoreHomeserverURL: true,

		DeviceID: id.DeviceID(deviceId),

		InitialDeviceDisplayName: "Shulker",
	})

	if deviceId == "" {
		log.Println("Generating a new Device ID for Matrix")
		deviceId = client.DeviceID.String()
		log.Printf("New device ID is: %v\n", deviceId)
	}
	if err != nil {
		return client, nil, err
	}

	store := store.NewStateStore(db)
	store.CreateTables()

	client.Store = store

	olmMachine := crypto.NewOlmMachine(client, logger{}, makeCryptoStore(client, db, key, id.DeviceID(deviceId)), store)
	olmMachine.Load()

	syncer := client.Syncer.(*mautrix.DefaultSyncer)

	syncer.OnSync(olmMachine.ProcessSyncResponse)

	syncer.OnEventType(event.StateMember, func(_ mautrix.EventSource, evt *event.Event) {
		olmMachine.HandleMemberEvent(evt)
		store.SetMembership(evt)

		if evt.GetStateKey() == client.UserID.String() && evt.Content.AsMember().Membership == event.MembershipInvite {
			err := doRetry(func() error {
				_, err := client.JoinRoomByID(evt.RoomID)
				return err
			})
			if err != nil {
				log.Printf("Could not join channel %s. Error %+v", evt.RoomID.String(), err)
			}
		}
	})

	syncer.OnEventType(event.StateEncryption, func(_ mautrix.EventSource, evt *event.Event) {
		store.SetEncryptionEvent(evt)
	})

	go func() {
		err = client.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}()

	clientslock.Lock()
	clients[clientid] = client
	clientslock.Unlock()
	return client, olmMachine, nil

}

func getMember(userId id.UserID, client *mautrix.Client, roomId id.RoomID) (*memberInfo, error) {
	memberslock.Lock()
	if member, ok := members[userId.String()]; ok {
		memberslock.Unlock()
		return member, nil
	}
	memberslock.Unlock()
	resp, err := client.JoinedMembers(roomId)
	if err != nil {
		return nil, err
	}
	if user, ok := resp.Joined[userId]; ok {
		memberslock.Lock()

		members[userId.String()] = &memberInfo{
			userId:      userId,
			displayname: *user.DisplayName,
			avatarUrl:   *user.AvatarURL,
		}

		memberslock.Unlock()
		return members[userId.String()], nil
	}

	return nil, errors.New("cannot find member " + userId.String())
}

func makeCryptoStore(client *mautrix.Client, db *sql.DB, key string, deviceId id.DeviceID) *crypto.SQLCryptoStore {

	store := crypto.NewSQLCryptoStore(
		db,
		"sqlite3",
		"",
		deviceId,
		[]byte(key),
		&logger{},
	)
	err := store.CreateTables()
	if err != nil {
		log.Fatalln("Failed to create tables for SQL encryption key store")
	}
	return store
}

func doRetry(fn func() error) error {
	var err error
	b := retry.NewFibonacci(1 * time.Second)
	if err != nil {
		panic(err)
	}
	b = retry.WithMaxRetries(5, b)
	for {
		err = fn()
		if err == nil {
			return nil
		}
		nextDuration, stop := b.Next()
		if stop {
			err = errors.New("joinroom: Retry limit reached. Will not retry")
			break
		}
		time.Sleep(nextDuration)
	}
	return err
}
