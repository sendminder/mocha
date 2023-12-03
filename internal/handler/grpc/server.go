package grpc

import (
	"log/slog"
	"net"
	"strconv"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"mocha/internal/cache"
	"mocha/internal/db"
	pb "mocha/proto/message"
	"mocha/util"
)

type MessageGrpcServer interface {
	MessageHandler
	Start(wg *sync.WaitGroup)
}

var _ MessageGrpcServer = (*messageServer)(nil)

type messageServer struct {
	pb.UnimplementedMessageServiceServer
	sf         *util.Snowflake
	nodeID     int
	nodeIDBits uint
	portStr    string
	mdb        db.DynamoDatabase
	rdb        db.RelationalDatabase
	cache      cache.RedisCache
}

func (s *messageServer) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	listener, err := net.Listen("tcp", ":"+s.portStr)
	if err != nil {
		slog.Error("failed to listen", "error", err)
	}
	srv := grpc.NewServer(
		grpc.KeepaliveParams(
			keepalive.ServerParameters{
				MaxConnectionAge:      24 * time.Hour,
				MaxConnectionAgeGrace: 10 * time.Second,
				Time:                  60 * time.Second,
				Timeout:               10 * time.Second,
			},
		),
		grpc.KeepaliveEnforcementPolicy(
			keepalive.EnforcementPolicy{
				MinTime: 60 * time.Second,
			},
		),
	)
	pb.RegisterMessageServiceServer(srv, s)
	slog.Info("gRPC server is listening on port", "port", s.portStr)
	if err := srv.Serve(listener); err != nil {
		slog.Error("failed to serve", "error", err)
	}
}

func NewGrpcServer(port int, mdb db.DynamoDatabase, rdb db.RelationalDatabase, cache cache.RedisCache) MessageGrpcServer {
	// TODO 노드 ID와 노드 Id 비트 수 설정
	nodeID := 1
	nodeIDBits := uint(10)
	sf := util.NewSnowflake(nodeID, nodeIDBits)
	portStr := strconv.Itoa(port)
	return &messageServer{
		sf:         sf,
		nodeID:     nodeID,
		nodeIDBits: nodeIDBits,
		portStr:    portStr,
		mdb:        mdb,
		rdb:        rdb,
		cache:      cache,
	}
}
