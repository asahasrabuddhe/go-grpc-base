package server

import (
	"bitbucket.org/perennialsys/go-grpc-base/database"
	"bitbucket.org/perennialsys/go-grpc-base/mail"
	"fmt"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"log"
	"net"
)

var server *grpc.Server

func Init(path, filename string) *grpc.Server {
	viper.AddConfigPath(path)
	viper.SetConfigName(filename)
	viper.SetConfigType("json")
	//viper.WatchConfig()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal("Unable to read config", err)
	}

	mail.BootstrapMail()

	credentials := viper.GetStringMapString("database.default")
	database.Open(credentials["username"], credentials["password"], credentials["name"], credentials["host"], credentials["port"])

	server = grpc.NewServer()

	return server
}

func Start() {
	network := viper.GetString("network")
	address := fmt.Sprintf("%v:%v", viper.GetString("ip_address"), viper.GetString("port"))

	listener, err := net.Listen(network, address)
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	if err := server.Serve(listener); err != nil {
		log.Fatal(err)
	}
}
