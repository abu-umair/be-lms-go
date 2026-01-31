package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"github.com/abu-umair/be-lms-go/internal/grpcmiddleware"
	"github.com/abu-umair/be-lms-go/internal/handler"
	"github.com/abu-umair/be-lms-go/internal/repository"
	"github.com/abu-umair/be-lms-go/internal/service"
	"github.com/abu-umair/be-lms-go/pb/auth"
	"github.com/abu-umair/be-lms-go/pb/course"
	"github.com/abu-umair/be-lms-go/pkg/database"
	"github.com/joho/godotenv"
	gocache "github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	godotenv.Load()
	ctx := context.Background()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Panicf("Error when listening %v", err)
	}

	db := database.ConnectDB(ctx, os.Getenv("DB_URI"))

	log.Println("Connected to DB")

	cacheService := gocache.New(time.Hour*24, time.Hour*24)

	authMiddleware := grpcmiddleware.NewAuthMiddleware(cacheService)

	emailService := service.NewEmailSender(
		"sandbox.smtp.mailtrap.io",
		2525,
		"37f5838d8d81ab",
		"041576603c2bd7",
	)

	authRepository := repository.NewAuthRepository(db)
	authService := service.NewAuthService(authRepository, cacheService, emailService)
	authHandler := handler.NewAuthHandler(authService)

	courseRepository := repository.NewCourseRepository(db)
	courseService := service.NewCourseService(db, courseRepository)
	courseHandler := handler.NewCourseHandler(courseService)

	serv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			grpcmiddleware.ErrorMiddleware,
			authMiddleware.Middleware,
		),
	)

	auth.RegisterAuthServiceServer(serv, authHandler)
	course.RegisterCourseServiceServer(serv, courseHandler)

	if os.Getenv("ENVIRONMENT") == "dev" {
		reflection.Register(serv)
		log.Printf("Reflection is Registered.")

	}

	log.Println("Server is running on :50051 port")
	if err := serv.Serve(lis); err != nil {
		log.Panicf("Server is error %v", err)
	}
}
