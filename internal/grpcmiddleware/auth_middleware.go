package grpcmiddleware

import (
	"context"
	"log"

	jwtentity "github.com/abu-umair/be-lms-go/internal/entity/jwt"
	"github.com/abu-umair/be-lms-go/internal/utils"
	gocache "github.com/patrickmn/go-cache"
	"google.golang.org/grpc"
)

type authMiddleware struct {
	cacheService *gocache.Cache
}

func (am *authMiddleware) Middleware(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
	log.Println(info.FullMethod)
	if info.FullMethod == "/auth.AuthService/Login" || info.FullMethod == "/auth.AuthService/Register" { //?whitelist (boleh diakses tanpa token)
		return handler(ctx, req)
	}

	// Ambil token  dari metadata
	tokenStr, err := jwtentity.ParseTokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	// Cek token  dari  logout cache
	_, ok := am.cacheService.Get(tokenStr)

	if ok {
		return nil, utils.UnauthenticatedResponse()
	}

	// Parse jwt nya hingga jadi entity
	claims, err := jwtentity.GetClaimsFromToken(tokenStr)

	if err != nil {
		return nil, err
	}

	// Sematkan entity ke context
	ctx = claims.SetToContext(ctx)

	res, err := handler(ctx, req)

	return res, err
}

func NewAuthMiddleware(cacheService *gocache.Cache) *authMiddleware {
	return &authMiddleware{
		cacheService: cacheService,
	}
}
