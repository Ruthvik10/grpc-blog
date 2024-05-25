package main

import (
	"context"
	"io"
	"log"

	pb "github.com/Ruthvik10/grpc-blog/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var addr = "localhost:50051"

func main() {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalln("Failed to connect to the server:", err)
	}
	defer conn.Close()

	client := pb.NewBlogServiceClient(conn)

	// insertPayload := &pb.Blog{
	// 	Title:    "Mobile app with React Native",
	// 	Content:  "Learn building mobile app with React Native",
	// 	AuthorId: "2",
	// }
	// createBlog(client, insertPayload)
	// readBlog(client)

	// updatePayload := &pb.Blog{
	// 	Id:       "664ef30d302b189cf27c4755",
	// 	Title:    "Build web app with gRPC and Mongo",
	// 	Content:  "Learn building web app with Mongo and gRPC in 10 mins",
	// 	AuthorId: "1",
	// }
	// updateBlog(client, updatePayload)
	listBlog(client)
}

func createBlog(client pb.BlogServiceClient, payload *pb.Blog) {
	res, err := client.Create(context.Background(), payload)
	if err != nil {
		log.Println("Error inserting the document", err)
		return
	}
	log.Println("Inserted document, id:", res.Id)
}

func readBlog(client pb.BlogServiceClient) {
	res, err := client.Read(context.Background(), &pb.BlogID{Id: "664ef30d302b189cf27c4755"})
	if err != nil {
		e, ok := status.FromError(err)
		if ok {
			log.Printf("Error code: %s\nError message: %s", e.Code().String(), e.Message())
			return
		}
		log.Println("Unexpected error:", err)
		return

	}
	log.Printf("Blog: %+v", res)
}

func updateBlog(client pb.BlogServiceClient, payload *pb.Blog) {
	_, err := client.Update(context.Background(), payload)
	if err != nil {
		log.Println("Error updating the blog:", err)
		return
	}
	log.Println("Updated the blog")
}

func listBlog(client pb.BlogServiceClient) {
	stream, err := client.List(context.Background(), &emptypb.Empty{})
	if err != nil {
		log.Println("Error recieving data:", err)
		return
	}

	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Println("Error recieving the stream:", err)
			break
		}
		log.Printf("Blog: %+v\n", msg)
	}
}
