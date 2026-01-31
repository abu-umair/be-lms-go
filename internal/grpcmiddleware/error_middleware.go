package grpcmiddleware

import (
	"context"
	"log"
	"runtime/debug"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ErrorMiddleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println(r)
			debug.PrintStack()                                           //? menampilkan stack trace (menampilkan error di console)
			err = status.Errorf(codes.Internal, "Internal Server Error") //? memasang status internal
		}
	}()
	res, err := handler(ctx, req) //?memanggil handler disertai dengan context dan requestnya
	log.Println(err)

	if err != nil {
		log.Println(err)

		//? pengecekan
		if st, ok := status.FromError(err); ok {
			if st.Code() == codes.Unauthenticated {
				return nil, err
			}
		}

		return nil, status.Error(codes.Internal, "Internal Server Error")
		// return nil, err //? jika ingin melihat errornya di postman
	}

	return res, err
}
