package main

import (
	"fmt"
	"sync"

	"github.com/lordnorthern/login_server/helpers"
	"github.com/lordnorthern/login_server/models"
	"github.com/lordnorthern/login_server/server"
)

func main() {

	fmt.Print("Initializing...")
	server.StartCommandConcole()
	models.ParseConf()
	var wg sync.WaitGroup

	var publicServer server.ListenServer
	err := publicServer.InitializeServer(models.Conf.PublicServer)
	if err != nil {
		fmt.Println("Failed.")
		helpers.LogError(err)
		return
	}

	var internalServer server.ListenServer
	err = internalServer.InitializeServer(models.Conf.InternalServer)
	if err != nil {
		fmt.Println("Failed.")
		helpers.LogError(err)
		return
	}
	fmt.Println("Success.")
	fmt.Print("Initializing MySQL...")
	server.MySQL, err = server.InitializeMySQL("Sole MySQL Instance")
	if err != nil {
		fmt.Println("Failed.")
		helpers.LogError(err)
		return
	}
	fmt.Println("Success.")
	publicServer.ListenAndAccept(&wg, server.PublicConnectionHandler)
	internalServer.ListenAndAccept(&wg, server.InternalConnectionHandler)
	wg.Wait()
	<-server.EndServe
}
