package repositories

import (
	"context"
	"fmt"

	"github.com/auliayaya/go-blog/server/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BlogRepository interface {
	Create(context.Context, models.BlogItem) (string, error)
	Read(context.Context, string) (*models.BlogItem, error)
	Update(context.Context, string, models.BlogItem) (*models.BlogItem, error)
	Delete(context.Context, string) (string, error)
	FindAll(context.Context) (*[]models.BlogItem, error)
}

type blogRepository struct {
	collection *mongo.Collection
}

func NewBlogRepository(col *mongo.Collection) BlogRepository {
	return &blogRepository{collection: col}
}

func (b *blogRepository) Create(ctx context.Context, data models.BlogItem) (string, error) {
	res, err := b.collection.InsertOne(ctx, data)
	if err != nil {
		return "", status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Error : %v", err),
		)
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return " ", status.Errorf(
			codes.Internal,
			fmt.Sprintln("Cannot convert to OID "),
		)
	}
	return oid.Hex(), nil
}

func (b *blogRepository) Read(ctx context.Context, oid string) (*models.BlogItem, error) {
	data := &models.BlogItem{}
	oidH, err := primitive.ObjectIDFromHex(oid)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Error : %v", err),
		)
	}
	res := b.collection.FindOne(ctx, bson.M{"_id": oidH})
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}
	return data, nil
}
func (b *blogRepository) Update(ctx context.Context, oid string, data models.BlogItem) (*models.BlogItem, error) {

	// res := b.collection.FindOne(ctx, bson.M{"_id": oid})
	// if err := res.Decode(data); err != nil {
	// 	return nil, status.Errorf(
	// 		codes.NotFound,
	// 		fmt.Sprintf("Cannot find blog with specified ID: %v", err),
	// 	)
	// }

	filter := bson.M{"_id": oid}
	_, err := b.collection.ReplaceOne(ctx, filter, data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot update object in MongoDB: %v", err),
		)
	}
	return &data, nil
}
func (b *blogRepository) Delete(ctx context.Context, oid string) (string, error) {
	oidH, err := primitive.ObjectIDFromHex(oid)
	if err != nil {
		return "", status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal Error : %v", err),
		)
	}
	res, err := b.collection.DeleteOne(context.Background(), bson.M{"_id": oidH})
	if err != nil {
		return "", status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot delete object in MongoDB: %v", err),
		)
	}
	if res.DeletedCount == 0 {
		return "", status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog in MongoDB: %v", err),
		)
	}
	return oid, nil
}

func (b *blogRepository) FindAll(ctx context.Context) (*[]models.BlogItem, error) {
	datas := []models.BlogItem{}
	cur, err := b.collection.Find(context.Background(), bson.D{})
	fmt.Println("CUR ", cur)
	if err != nil {
		fmt.Println("Error Here")
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("unknown internal error : %v", err),
		)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		data := &models.BlogItem{}
		err := cur.Decode(data)
		if err != nil {
			return nil, status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while decoding data from MongoDB: %v", err),
			)
		}
		datas = append(datas, *data)

	}
	if err := cur.Err(); err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("unknown internal error : %v", err),
		)
	}
	return &datas, nil
}
