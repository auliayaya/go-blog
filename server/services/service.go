package services

import (
	"context"
	"fmt"

	"github.com/auliayaya/go-blog/server/handler/grpc/protos"
	"github.com/auliayaya/go-blog/server/models"
	"github.com/auliayaya/go-blog/server/repositories"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Services interface {
	CreateBlog(context.Context, *protos.CreateBlogRequest) (*protos.CreateBlogResponse, error)
	ReadBlog(context.Context, *protos.ReadBlogRequest) (*protos.ReadBlogResponse, error)
	UpdateBlog(context.Context, *protos.UpdateBlogRequest) (*protos.UpdateBlogResponse, error)
	DeleteBlog(context.Context, *protos.DeleteBlogRequest) (*protos.DeleteBlogResponse, error)
	ListBlog(*protos.ListBlogRequest, protos.BlogService_ListBlogServer) error
}

type Server struct {
	repo repositories.BlogRepository
}

func NewServer(repo repositories.BlogRepository) Services {
	return &Server{repo: repo}
}

func (s *Server) CreateBlog(ctx context.Context, req *protos.CreateBlogRequest) (*protos.CreateBlogResponse, error) {
	blog := req.GetBlog()
	fmt.Println("Create Blog Request")
	data := models.BlogItem{
		AuthorID: blog.GetAuthorId(),
		Title:    blog.GetTitle(),
		Content:  blog.GetContent(),
	}
	res, err := s.repo.Create(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Error : %v", err),
		)
	}

	return &protos.CreateBlogResponse{
		Blog: &protos.Blog{
			Id:       res,
			AuthorId: blog.GetAuthorId(),
			Title:    blog.GetTitle(),
			Content:  blog.GetContent(),
		},
	}, nil

}

func (s *Server) ReadBlog(ctx context.Context, req *protos.ReadBlogRequest) (*protos.ReadBlogResponse, error) {
	fmt.Println("Read blog request")

	blogID := req.GetBlogId()
	oid, err := primitive.ObjectIDFromHex(blogID)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID : %v", err),
		)
	}

	res, err := s.repo.Read(context.Background(), oid.String())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Error : %v", err),
		)
	}
	return &protos.ReadBlogResponse{
		Blog: dataToBlogPb(res),
	}, nil

}

func (s *Server) UpdateBlog(ctx context.Context, req *protos.UpdateBlogRequest) (*protos.UpdateBlogResponse, error) {
	fmt.Println("Update blog request")
	blog := req.GetBlog()
	oid, err := primitive.ObjectIDFromHex(blog.GetId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID :%v ", err),
		)
	}

	res, err := s.repo.Read(context.Background(), oid.String())

	res.AuthorID = blog.GetAuthorId()
	res.Content = blog.GetContent()
	res.Title = blog.GetTitle()

	result, err := s.repo.Update(context.Background(), oid.String(), *res)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot update object in MongoDB: %v", err),
		)
	}
	return &protos.UpdateBlogResponse{
		Blog: dataToBlogPb(result),
	}, nil
}
func (s *Server) DeleteBlog(ctx context.Context, req *protos.DeleteBlogRequest) (*protos.DeleteBlogResponse, error) {
	fmt.Println("Delete blog request")
	oid, err := primitive.ObjectIDFromHex(req.GetBlogId())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID : %v", err),
		)
	}

	res, err := s.repo.Delete(context.Background(), oid.String())
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot delete object in MongoDB: %v", err),
		)
	}

	return &protos.DeleteBlogResponse{BlogId: res}, nil
}

func (s *Server) ListBlog(req *protos.ListBlogRequest, stream protos.BlogService_ListBlogServer) error {
	fmt.Println("List blog request")
	result, err := s.repo.FindAll(context.Background())

	if err != nil {
		fmt.Println("Error Here")
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("unknown internal error : %v", err),
		)
	}

	for _, v := range *result {

		stream.Send(&protos.ListBlogResponse{Blog: dataToBlogPb(&v)})

	}

	return nil
}
func dataToBlogPb(data *models.BlogItem) *protos.Blog {
	return &protos.Blog{
		Id:       data.ID.Hex(),
		AuthorId: data.AuthorID,
		Content:  data.Content,
		Title:    data.Title,
	}
}
