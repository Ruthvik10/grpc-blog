package main

import (
	"context"
	"log"
	"net"

	pb "github.com/Ruthvik10/grpc-blog/proto"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
)

var addr = "localhost:50051"

func main() {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI("mongodb://root:root@localhost:27017/"))
	if err != nil {
		log.Fatalln(err)
	}

	collection := client.Database("blogdb").Collection("blog")

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Error listening on addr %s: %v\n", addr, err)
	}
	defer lis.Close()

	s := grpc.NewServer()
	pb.RegisterBlogServiceServer(s, &blogHandler{coll: collection})

	log.Println("gRPC server started on", addr)

	if err := s.Serve(lis); err != nil {
		log.Fatalln("Error starting the gRPC server:", err)
	}
}
