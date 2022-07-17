package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"time"

	grpch "github.com/auliayaya/go-blog/server/handler/grpc"
	"github.com/auliayaya/go-blog/server/repositories"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"google.golang.org/grpc"
)

var collection *mongo.Collection

// func (*server) CreateBlog(ctx context.Context, req *protos.CreateBlogRequest) (*protos.CreateBlogResponse, error) {
// 	blog := req.GetBlog()
// 	fmt.Println("Create Blog Request")
// 	data := blogItem{
// 		AuthorID: blog.GetAuthorId(),
// 		Title:    blog.GetTitle(),
// 		Content:  blog.GetContent(),
// 	}
// 	res, err := collection.InsertOne(context.Background(), data)
// 	if err != nil {
// 		return nil, status.Errorf(
// 			codes.Internal,
// 			fmt.Sprintf("Internal Error : %v", err),
// 		)
// 	}
// 	oid, ok := res.InsertedID.(primitive.ObjectID)
// 	if !ok {
// 		return nil, status.Errorf(
// 			codes.Internal,
// 			fmt.Sprintln("Cannot convert to OID "),
// 		)
// 	}
// 	return &protos.CreateBlogResponse{
// 		Blog: &protos.Blog{
// 			Id:       oid.Hex(),
// 			AuthorId: blog.GetAuthorId(),
// 			Title:    blog.GetTitle(),
// 			Content:  blog.GetContent(),
// 		},
// 	}, nil

// }

// func (*server) ReadBlog(ctx context.Context, req *protos.ReadBlogRequest) (*protos.ReadBlogResponse, error) {
// 	fmt.Println("Read blog request")

// 	blogID := req.GetBlogId()
// 	oid, err := primitive.ObjectIDFromHex(blogID)
// 	if err != nil {
// 		return nil, status.Errorf(
// 			codes.InvalidArgument,
// 			fmt.Sprintf("Cannot parse ID : %v", err),
// 		)
// 	}
// 	data := &blogItem{}

// 	res := collection.FindOne(context.Background(), bson.M{"_id": oid})
// 	if err := res.Decode(data); err != nil {
// 		return nil, status.Errorf(
// 			codes.NotFound,
// 			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
// 		)
// 	}
// 	return &protos.ReadBlogResponse{
// 		Blog: dataToBlogPb(data),
// 	}, nil

// }

// func (*server) UpdateBlog(ctx context.Context, req *protos.UpdateBlogRequest) (*protos.UpdateBlogResponse, error) {
// 	fmt.Println("Update blog request")
// 	blog := req.GetBlog()
// 	oid, err := primitive.ObjectIDFromHex(blog.GetId())
// 	if err != nil {
// 		return nil, status.Errorf(
// 			codes.InvalidArgument,
// 			fmt.Sprintf("Cannot parse ID :%v ", err),
// 		)
// 	}
// 	data := &blogItem{}
// 	res := collection.FindOne(context.Background(), bson.M{"_id": oid})
// 	if err := res.Decode(data); err != nil {
// 		return nil, status.Errorf(
// 			codes.NotFound,
// 			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
// 		)
// 	}
// 	data.AuthorID = blog.GetAuthorId()
// 	data.Content = blog.GetContent()
// 	data.Title = blog.GetTitle()
// 	filter := bson.M{"_id": oid}
// 	_, err = collection.ReplaceOne(context.Background(), filter, data)
// 	if err != nil {
// 		return nil, status.Errorf(
// 			codes.Internal,
// 			fmt.Sprintf("Cannot update object in MongoDB: %v", err),
// 		)
// 	}
// 	return &protos.UpdateBlogResponse{
// 		Blog: dataToBlogPb(data),
// 	}, nil
// }
// func (*server) DeleteBlog(ctx context.Context, req *protos.DeleteBlogRequest) (*protos.DeleteBlogResponse, error) {
// 	fmt.Println("Delete blog request")
// 	oid, err := primitive.ObjectIDFromHex(req.GetBlogId())
// 	if err != nil {
// 		return nil, status.Errorf(
// 			codes.InvalidArgument,
// 			fmt.Sprintf("Cannot parse ID : %v", err),
// 		)
// 	}
// 	filter := bson.M{"_id": oid}
// 	res, err := collection.DeleteOne(context.Background(), filter)
// 	if err != nil {
// 		return nil, status.Errorf(
// 			codes.Internal,
// 			fmt.Sprintf("Cannot delete object in MongoDB: %v", err),
// 		)
// 	}
// 	if res.DeletedCount == 0 {
// 		return nil, status.Errorf(
// 			codes.NotFound,
// 			fmt.Sprintf("Cannot find blog in MongoDB: %v", err),
// 		)
// 	}
// 	return &protos.DeleteBlogResponse{BlogId: req.GetBlogId()}, nil
// }

// func (*server) ListBlog(req *protos.ListBlogRequest, stream protos.BlogService_ListBlogServer) error {
// 	fmt.Println("List blog request")
// 	cur, err := collection.Find(context.Background(), nil)
// 	fmt.Println("CUR ", cur)
// 	if err != nil {
// 		fmt.Println("Error Here")
// 		return status.Errorf(
// 			codes.Internal,
// 			fmt.Sprintf("unknown internal error : %v", err),
// 		)
// 	}
// 	defer cur.Close(context.Background())
// 	for cur.Next(context.Background()) {
// 		data := &blogItem{}
// 		err := cur.Decode(data)
// 		if err != nil {
// 			return status.Errorf(
// 				codes.Internal,
// 				fmt.Sprintf("Error while decoding data from MongoDB: %v", err),
// 			)
// 		}
// 		stream.Send(&protos.ListBlogResponse{Blog: dataToBlogPb(data)})

// 	}
// 	if err := cur.Err(); err != nil {
// 		return status.Errorf(
// 			codes.Internal,
// 			fmt.Sprintf("unknown internal error : %v", err),
// 		)
// 	}
// 	return nil
// }
// func dataToBlogPb(data *blogItem) *protos.Blog {
// 	return &protos.Blog{
// 		Id:       data.ID.Hex(),
// 		AuthorId: data.AuthorID,
// 		Content:  data.Content,
// 		Title:    data.Title,
// 	}
// }

// type blogItem struct {
// 	ID       primitive.ObjectID `bson:"_id,omitempty"`
// 	AuthorID string             `bson:"author_id"`
// 	Title    string             `bson:"title"`
// 	Content  string             `bson:"content"`
// }

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fmt.Println("Blog Service Started")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	fmt.Println("Connecting to MongoDB")
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://docker:mongopw@localhost:49153"))
	if err != nil {
		log.Fatal(err)
	}

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen : %v", err)
	}

	collection = client.Database("go-blog").Collection("blog")
	ar := repositories.NewBlogRepository(collection)
	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	grpch.NewRPC(s, ar)
	go func() {
		fmt.Println("Starting Server....")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for Control C to Exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("Closing the listener")
	lis.Close()
	fmt.Println("Closing mongodb connection")
	client.Disconnect(context.TODO())
	fmt.Println("End the program")
}
