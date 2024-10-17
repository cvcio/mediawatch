package interceptors

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

// TimeoutCallOption is a custom call option to set a forced timeout
type TimeoutCallOption struct {
	grpc.EmptyCallOption
	forcedTimeout time.Duration
}

// WithForcedTimeout returns a new TimeoutCallOption with the given forced timeout
// func WithForcedTimeout(forceTimeout time.Duration) TimeoutCallOption {
// 	return TimeoutCallOption{forcedTimeout: forceTimeout}
// }

// getTimeout returns the forced timeout from the call options
func getTimeout(callOptions []grpc.CallOption) (time.Duration, bool) {
	for _, opt := range callOptions {
		if co, ok := opt.(TimeoutCallOption); ok {
			return co.forcedTimeout, true
		}
	}

	return 0, false
}

// TimeoutInterceptor returns a new timeout interceptor
func TimeoutInterceptor(t time.Duration) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		timeout := t
		if v, ok := getTimeout(opts); ok {
			timeout = v
		}

		if timeout <= 0 {
			return invoker(ctx, method, req, reply, cc, opts...)
		}

		ctx, cancel := context.WithTimeout(ctx, timeout)
		defer cancel()

		return invoker(ctx, method, req, reply, cc, opts...)
	}
}
