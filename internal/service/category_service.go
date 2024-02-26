package service

import (
	"context"
	"io"

	"github.com/allurco/fullcycle-goexpert-grpc/internal/database"
	"github.com/allurco/fullcycle-goexpert-grpc/internal/pb"
)

type CategoryServiceServer struct {
	pb.UnimplementedCategoryServiceServer
	CategoryDb database.Category
}

func NewCategoryServiceServer(categoryDb database.Category) *CategoryServiceServer {
	return &CategoryServiceServer{
		CategoryDb: categoryDb,
	}
}

func (s *CategoryServiceServer) CreateCategory(ctx context.Context, in *pb.CreateCategoryRequest) (*pb.Category, error) {
	category, err := s.CategoryDb.Create(in.Name, in.Description)
	if err != nil {
		return nil, err
	}

	var categoryObj = &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return categoryObj, nil

}

func (s *CategoryServiceServer) ListCategories(ctx context.Context, in *pb.Blank) (*pb.CategoryList, error) {
	categories, err := s.CategoryDb.FindAll()
	if err != nil {
		return nil, err
	}

	var categoriesArray []*pb.Category
	for _, category := range categories {
		categoriesArray = append(categoriesArray, &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})
	}

	return &pb.CategoryList{Categories: categoriesArray}, nil
}

func (s *CategoryServiceServer) GetCategory(ctx context.Context, in *pb.CategoryGetRequest) (*pb.Category, error) {
	category, err := s.CategoryDb.FindByID(in.Id)
	if err != nil {
		return nil, err
	}

	var categoryObj = &pb.Category{
		Id:          category.ID,
		Name:        category.Name,
		Description: category.Description,
	}

	return categoryObj, nil
}

func (s *CategoryServiceServer) CreateCategoryStream(stream pb.CategoryService_CreateCategoryStreamServer) error {
	categories := &pb.CategoryList{}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			return stream.SendAndClose(categories)
		}
		if err != nil {
			return err
		}

		category, err := s.CategoryDb.Create(req.Name, req.Description)
		if err != nil {
			return err
		}

		categories.Categories = append(categories.Categories, &pb.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		})

	}
}

func (s *CategoryServiceServer) CreateCategoryStreamBidirecional(stream pb.CategoryService_CreateCategoryStreamBidirecionalServer) error {
	for {
		category, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		categoryResult, err := s.CategoryDb.Create(category.Name, category.Description)
		if err != nil {
			return err
		}

		err = stream.Send(&pb.Category{
			Id:          categoryResult.ID,
			Name:        categoryResult.Name,
			Description: categoryResult.Description,
		})

		if err != nil {
			return err
		}

	}
}
