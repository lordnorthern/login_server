package models

import (
	"bytes"
	"crypto/sha512"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"math"
	"math/rand"
	"net"
	"time"

	"github.com/lordnorthern/login_server/helpers"
)

// User will hold all the data for a user connection
// On the server, this will hold a handle to a specific user, a socket handle, account data etc...
// On the client, this will hold the connection to the server so that the client and the server
// can use the same exact objects and functions
type User struct {
	Conn          *net.Conn
	Account       *Account
	EncryptionKey []byte
}

// NewConnection will create a new connection object
func NewConnection(inCon *net.Conn) (newUser *User, ok bool) {
	newUser = new(User)
	newUser.Account = new(Account)
	newUser.Account.Characters = make([]Character, 0)
	newUser.Conn = inCon
	newUser.EncryptionKey = make([]byte, 0, 32)
	ok = true
	return
}

// CreateEncryptionCode will generate an encryption code
func (user *User) CreateEncryptionCode() {
	if len(user.EncryptionKey) == 0 {
		user.EncryptionKey = []byte(helpers.GenerateRandomString(32, true))
	}
}

// SendCommand is the heart of this package.
// It will take in a Command interface, encrypt it (if encryptCommand is set to true)
// and will call serialize on the object,  which will call every command type's respective serialization method
// and will send it over to the "user".
// If used by the client, this will simply send the command to the server.
func (user *User) SendCommand(cmd Command, encryptCommand bool) error {
	if len(user.EncryptionKey) == 0 && encryptCommand {
		return errors.New("NoEncryptionKey")
	}
	cmd.SetSessionID(user.Account.SessionID)
	encodedCommand, err := cmd.Serialize()
	var Result []byte
	if err != nil {
		return err
	}
	if encryptCommand {
		encrypted := make([]byte, 0)
		encrypted, err = helpers.Encrypt(encodedCommand, user.EncryptionKey)
		if err != nil {
			return err
		}
		Result = encrypted
	} else {
		Result = encodedCommand
	}

	rand.Seed(time.Now().UnixNano())
	totalChunks := int32(math.Ceil(float64(len(Result)) / 1024))
	bufCmdSerial := new(bytes.Buffer)
	binary.Write(bufCmdSerial, binary.LittleEndian, int32(rand.Intn(100000)))
	for i := totalChunks - 1; i >= 0; i-- {
		bufChunkNum := new(bytes.Buffer)
		binary.Write(bufChunkNum, binary.LittleEndian, i)
		start := i * 1024
		var end int32
		if i == (totalChunks - 1) {
			end = int32(len(Result))
		} else {
			end = start + 1024
		}
		chunk := append(bufChunkNum.Bytes(), bufCmdSerial.Bytes()...)

		chunk = append(chunk, Result[start:end]...)
		(*user.Conn).Write(chunk)
	}

	return nil
}

type EmailPassword struct {
	Email    string
	Password string
}

func NewEmailPassword(email, password string) (*EmailPassword, error) {
	if email == "" || password == "" {
		// Do more validation here
		return nil, errors.New("Missing email or password")
	}
	newObj := new(EmailPassword)
	newObj.Email = email
	newObj.Password = password
	return newObj, nil
}

// HashPassword will convert the password into a sha512 string
func (l *EmailPassword) HashPassword() {
	if len(l.Password) != 128 {
		hash := sha512.New()
		hash.Write([]byte(l.Password))
		l.Password = hex.EncodeToString(hash.Sum(nil))
	}
}
