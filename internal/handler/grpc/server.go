package grpc

import (
	"log"
	"log/slog"
	"net"
	"strconv"
	"sync"

	"google.golang.org/grpc"
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
}

func (s *messageServer) Start(wg *sync.WaitGroup) {
	defer wg.Done()
	listener, err := net.Listen("tcp", ":"+s.portStr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	pb.RegisterMessageServiceServer(srv, &messageServer{})
	slog.Info("gRPC server is listening on port " + s.portStr)
	if err := srv.Serve(listener); err != nil {
		slog.Error("failed to serve", "error", err)
	}
}

func NewGrpcServer(port int, mdb db.DynamoDatabase, rdb db.RelationalDatabase) MessageGrpcServer {
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
	}
}
