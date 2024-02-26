package main

import (
	"database/sql"
	"net"

	"github.com/allurco/fullcycle-goexpert-grpc/internal/database"
	"github.com/allurco/fullcycle-goexpert-grpc/internal/pb"
	"github.com/allurco/fullcycle-goexpert-grpc/internal/service"
	_ "github.com/mattn/go-sqlite3"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	db, err := sql.Open("sqlite3", "../../fullcycle_grpc.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	categoryDb := database.NewCategory(db)
	categoryService := service.NewCategoryServiceServer(*categoryDb)

	grpcServer := grpc.NewServer()
	pb.RegisterCategoryServiceServer(grpcServer, categoryService)
	reflection.Register(grpcServer)

	list, err := net.Listen("tcp", ":50051")
	if err != nil {
		panic(err)
	}

	if err := grpcServer.Serve(list); err != nil {
		panic(err)
	}

}
