package main

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/auliayaya/go-blog/server/handler/grpc/protos"
	"google.golang.org/grpc"
)

func main() {
	fmt.Println("Blog Client")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	opts := grpc.WithInsecure()
	cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect : %v", err)
	}
	defer cc.Close()
	c := protos.NewBlogServiceClient(cc)
	fmt.Println("Creating the blog")
	blog := protos.Blog{
		AuthorId: "Aulia",
		Title:    "My First Blog",
		Content:  "Content of my first blog",
	}
	// creat Blog
	cb, err := c.CreateBlog(context.Background(), &protos.CreateBlogRequest{Blog: &blog})
	if err != nil {
		log.Fatalf("Unexpected error : %v", err)
	}
	fmt.Printf("Blog has been created : %v", cb)
	blogID := cb.GetBlog().GetId()
	// Read Blog
	_, err = c.ReadBlog(context.Background(), &protos.ReadBlogRequest{BlogId: "62d2bcb7fca9ce28454c2f34"})
	if err != nil {
		fmt.Printf("Error happened while reading: %v", err)
	}
	rbr := &protos.ReadBlogRequest{
		BlogId: blogID,
	}
	rbres, err := c.ReadBlog(context.Background(), rbr)
	if err != nil {
		fmt.Printf("Error happened while reading: %v", err)
	}
	fmt.Printf("Blog was read : %v", rbres)

	// Update Blog
	newBlog := &protos.Blog{
		Id:       blogID,
		AuthorId: "Changed Author",
		Title:    "My First blog (edited)",
		Content:  "Content of my first blog",
	}
	updateRes, err := c.UpdateBlog(context.Background(), &protos.UpdateBlogRequest{Blog: newBlog})
	if err != nil {
		fmt.Printf("Error happened while updating: %v\n", err)
	}
	fmt.Printf("Blog was updated: %v", updateRes)
	// Delete Blog
	deleteRes, err := c.DeleteBlog(context.Background(), &protos.DeleteBlogRequest{BlogId: blogID})
	if err != nil {
		fmt.Printf("Error happened while deleting: %v \n", err)
	}
	fmt.Printf("Blog was deleted: %v \n", deleteRes)

	// List Blog
	stream, err := c.ListBlog(context.Background(), &protos.ListBlogRequest{})
	if err != nil {
		log.Fatalf("error while calling ListBlog RPC : %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetBlog())
	}

}
