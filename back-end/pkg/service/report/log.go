package report

import (
    "context"
    "google.golang.org/grpc"
    "google.golang.org/grpc/peer"
    "google.golang.org/grpc/metadata"
    "google.golang.org/grpc/status"
    "log"
    "time"
)

func unaryLoggingInterceptor(logger *log.Logger) grpc.UnaryServerInterceptor {
    return func(
        ctx context.Context,
        req interface{},
        info *grpc.UnaryServerInfo,
        handler grpc.UnaryHandler,
    ) (interface{}, error) {
        startTime := time.Now()
        p, _ := peer.FromContext(ctx)
        md, _ := metadata.FromIncomingContext(ctx)
        
        logger.Printf("Request - Method:%s; Peer:%s; Metadata:%v; Payload:%v", info.FullMethod, p.Addr, md, req)

        resp, err := handler(ctx, req)

        if err != nil {
            logger.Printf("Error - Method:%s; Peer:%s; Status:%s; Duration:%s; Error:%v", info.FullMethod, p.Addr, status.Code(err), time.Since(startTime), err)
        } else {
            logger.Printf("Response - Method:%s; Peer:%s; Status:%s; Duration:%s; Payload:%v", info.FullMethod, p.Addr, status.Code(err), time.Since(startTime), resp)
        }

        return resp, err
    }
}