package middleware

import (
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var mappingErrorCodes = map[error]codes.Code{}

// WithErrorInterceptor is function to use middleware/interceptor set error code and error message if any error
func WithErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			// set metadata
			md := metadata.Pairs("error-message", err.Error())
			_ = grpc.SetHeader(ctx, md)
			code, ok := mappingErrorCodes[err]
			if !ok {
				// return with internalserver error
				return nil, status.Error(codes.Internal, "internal server error")
			}

			return nil, status.Error(code, err.Error())
		}

		return
	}
}
