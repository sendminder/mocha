package handler

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	pb "mocha/proto/message"

	"mocha/util"

	"google.golang.org/grpc"
)

var (
	sf         *util.Snowflake
	nodeID     int
	nodeIDBits uint
)

type MessageServer struct {
	pb.UnimplementedMessageServiceServer
}

func (s *MessageServer) CreateMessage(ctx context.Context, req *pb.RequestCreateMessage) (*pb.ResponseCreateMessage, error) {
	fmt.Println("CreateMessage text =", req.Text)
	msgId := sf.Generate()

	// 새로운 Message 생성
	newMessage := &pb.Message{
		Id:             msgId,
		ConversationId: req.ConversationId,
		SenderId:       req.SenderId,
		Text:           req.Text,
	}

	return &pb.ResponseCreateMessage{
		Message:     newMessage,
		JoinedUsers: []int64{1, 2, 3},
	}, nil
}

func StartGrpc(wg *sync.WaitGroup) {
	defer wg.Done()
	initSnowflake()

	listener, err := net.Listen("tcp", ":3100")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	pb.RegisterMessageServiceServer(srv, &MessageServer{})
	fmt.Println("gRPC server is listening on port 3100...")
	if err := srv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initSnowflake() {
	// 노드 ID와 노드 ID 비트 수 설정
	nodeID = 1
	nodeIDBits = uint(10)
	sf = util.NewSnowflake(nodeID, nodeIDBits)
}
