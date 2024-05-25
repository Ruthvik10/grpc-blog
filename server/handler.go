package main

import (
	"context"
	"fmt"
	"log"

	pb "github.com/Ruthvik10/grpc-blog/proto"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type blogHandler struct {
	coll *mongo.Collection
	pb.BlogServiceServer
}

func (h *blogHandler) Create(ctx context.Context, in *pb.Blog) (*pb.BlogID, error) {
	res, err := h.coll.InsertOne(ctx, &BlogItem{AuthorID: in.AuthorId, Title: in.Title, Content: in.Content})
	if err != nil {
		log.Println("Error inserting the document:", err)
		return nil, status.Errorf(codes.Internal, fmt.Sprintf("Internal error: %v\n", err))
	}

	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(codes.Internal, "Cannot convert to objectID")
	}

	return &pb.BlogID{Id: oid.Hex()}, nil
}

func (h *blogHandler) Read(ctx context.Context, in *pb.BlogID) (*pb.Blog, error) {
	var blogItem BlogItem
	oid, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error: %v\n", err)
	}
	err = h.coll.FindOne(ctx, bson.D{{Key: "_id", Value: oid}}).Decode(&blogItem)
	if err == mongo.ErrNoDocuments {
		return nil, status.Errorf(codes.NotFound, "Blog item with id %s not found\n", in.Id)
	}
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error: %v\n", err)
	}
	return blogItem.ToBlog(), nil
}

func (h *blogHandler) Update(ctx context.Context, in *pb.Blog) (*emptypb.Empty, error) {
	oid, err := primitive.ObjectIDFromHex(in.Id)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "Internal error: %v\n", err)
	}

	data := BlogItem{
		ID:       oid,
		Title:    in.Title,
		AuthorID: in.AuthorId,
		Content:  in.Content,
	}

	_, err = h.coll.UpdateOne(ctx, bson.M{"_id": oid}, bson.M{"$set": data})

	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (h *blogHandler) List(in *emptypb.Empty, stream pb.BlogService_ListServer) error {
	cursor, err := h.coll.Find(context.Background(), primitive.D{{}})
	if err != nil {
		return status.Errorf(codes.Internal, "Internal error: %v\n", err)
	}
	defer cursor.Close(context.Background())
	var items []BlogItem
	if err = cursor.All(context.Background(), &items); err != nil {
		return status.Errorf(codes.Internal, "Internal error: %v\n", err)
	}

	for _, item := range items {
		err := stream.Send(&pb.Blog{
			Id:       item.ID.Hex(),
			Title:    item.Title,
			Content:  item.Content,
			AuthorId: item.AuthorID,
		})

		if err != nil {
			return status.Errorf(codes.Internal, "Error streaming: %v\n", err)
		}
	}
	return nil
}
