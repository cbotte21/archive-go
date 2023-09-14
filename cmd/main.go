package main

import (
	"github.com/cbotte21/archive-go/internal"
	"github.com/cbotte21/archive-go/pb"
	"github.com/cbotte21/archive-go/schema"
	"github.com/cbotte21/microservice-common/pkg/datastore"
	"github.com/cbotte21/microservice-common/pkg/enviroment"
	"google.golang.org/grpc"
	"log"
	"net"
)

func main() {
	enviroment.VerifyEnvVariable("port")
	enviroment.VerifyEnvVariable("mongo_uri")
	enviroment.VerifyEnvVariable("mongo_db")

	port := enviroment.GetEnvVariable("port")

	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen on port:" + port)
	}
	grpcServer := grpc.NewServer()

	//Register handlers to attach
	svcRecords := datastore.MongoClient[schema.SVCRecord]{}
	games := datastore.MongoClient[schema.Game]{}
	err = svcRecords.Init()
	if err != nil {
		return
	}
	err = games.Init()
	if err != nil {
		return
	}

	//Initialize archive
	archive := internal.NewArchive(&games, &svcRecords)

	pb.RegisterArchiveServiceServer(grpcServer, &archive)

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to initialize grpc server.")
	}
}
