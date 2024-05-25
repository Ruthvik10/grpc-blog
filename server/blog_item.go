package main

import (
	pb "github.com/Ruthvik10/grpc-blog/proto"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BlogItem struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	AuthorID string             `bson:"author_id"`
	Title    string             `bson:"title"`
	Content  string             `bson:"content"`
}

func (b BlogItem) ToBlog() *pb.Blog {
	return &pb.Blog{
		Id:       b.ID.Hex(),
		AuthorId: b.AuthorID,
		Title:    b.Title,
		Content:  b.Content,
	}
}
