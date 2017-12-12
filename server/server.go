package server

import (
	"fmt"
	"net"
	"sync"

	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/lordnorthern/login_server/models"
)

// ListenServer is the main server object
type ListenServer struct {
	Definition   string
	Port         string
	listenSocket *net.Listener
	alive        bool
}

// Terminatables are a list of objects to be terminated on close
var Terminatables []models.Terminator
var EndServe = make(chan bool)

func init() {
	Terminatables = make([]models.Terminator, 0)
}

func (s *ListenServer) Terminate() {
	s.alive = false
	(*s.listenSocket).Close()
}
func (s *ListenServer) GetName() string {
	return s.Definition
}
func (s *ListenServer) AddToList(list *[]models.Terminator) {
	*list = append(*list, s)
}

// InitializeServer will initialize the server
func (s *ListenServer) InitializeServer(conf models.ServConf) error {
	listener, err := net.Listen("tcp", ":"+conf.Port)
	if err != nil {
		return err
	}
	s.listenSocket = &listener
	s.alive = true
	s.Definition = conf.Definition
	s.Port = conf.Port
	s.AddToList(&Terminatables)
	return nil
}

// ListenAndAccept will listen and accept
func (s *ListenServer) ListenAndAccept(wg *sync.WaitGroup, connectionHandler func(*models.User)) {
	wg.Add(1)
	go func() {
		fmt.Print(s.GetName()+": Listening on port ", s.Port, "...\n")
		for s.alive {
			if con, err := (*s.listenSocket).Accept(); err == nil {
				if newUser, ok := models.NewConnection(&con); ok {
					go connectionHandler(newUser)
				}
			}
		}
		wg.Done()
	}()
}
