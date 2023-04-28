package main

import (
	"os"
	"os/signal"
	"syscall"
	"fmt"
	"path/filepath"
	server "github.com/0187773933/PassiveLogger/v1/server"
	utils "github.com/0187773933/PassiveLogger/v1/utils"
)

var s server.Server

func SetupCloseHandler() {
	c := make( chan os.Signal )
	signal.Notify( c , os.Interrupt , syscall.SIGTERM , syscall.SIGINT )
	go func() {
		<-c
		fmt.Println( "\r- Ctrl+C pressed in Terminal" )
		fmt.Println( "Shutting Down Passive Logging Server" )
		s.FiberApp.Shutdown()
		os.Exit( 0 )
	}()
}

func main() {
	SetupCloseHandler()
	config_file_path , _ := filepath.Abs( os.Args[ 1 ] )
	fmt.Println( "Loaded Config File From : %s" , config_file_path )
	config := utils.ParseConfig( config_file_path )
	fmt.Println( config )
	s = server.New( config )
	s.Start()
}