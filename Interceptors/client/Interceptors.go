package Interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

/** This is a test **/
func clientInterceptor(
	ctx context.Context,
	method string,
	req interface{},
	reply interface{},
	cc *grpc.ClientConn,
	invoker grpc.UnaryInvoker,
	opts ...grpc.CallOption,
) error {
	// Logic before invoking the invoker
	start := time.Now()

	// TODO intercept client call here...

	// Calls the invoker to execute RPC
	err := invoker(ctx, method, req, reply, cc, opts...)

	// Logic after invoking the invoker
	grpclog.Infof("Invoked RPC method=%s; Duration=%s; Error=%v", method, time.Since(start), err)

	return err
}

/**
 * The function below returns an grpc.DialOption value, which calls the
 * WithUnaryInterceptor function by providing the UnaryClientInterceptor
 * func value
 */
func WithClientUnaryInterceptor() grpc.DialOption {

	return grpc.WithUnaryInterceptor(clientInterceptor)
}
