package handler

import (
	"context"
	"log"
	"net"
	"sync"

	"mocha/db"
	pb "mocha/proto/message"
	"mocha/types"

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
	log.Println("CreateMessage text =", req.Text)
	msgId := sf.Generate()

	// 새로운 Message 생성
	newMessage := &pb.Message{
		Id:             msgId,
		ConversationId: req.ConversationId,
		SenderId:       req.SenderId,
		Text:           req.Text,
		Animal:         "cat",
	}

	dynamoMessage := &types.Message{
		Id:             msgId,
		ConversationId: req.ConversationId,
		SenderID:       req.SenderId,
		Text:           req.Text,
		Animal:         "cat",
	}
	err := db.CreateMessage(dynamoMessage)
	if err != nil {
		log.Printf("create message error %v\n", err)
		return &pb.ResponseCreateMessage{
			Status:       "error",
			ErrorMessage: err.Error(),
		}, nil
	}
	err = db.SetLastSeenMessageId(req.SenderId, req.ConversationId, msgId)
	if err != nil {
		log.Printf("set LastSeenMessageId error %v\n", err)
		return &pb.ResponseCreateMessage{
			Status:       "error",
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseCreateMessage{
		Status:      "ok",
		Message:     newMessage,
		JoinedUsers: []int64{1, 2, 3},
	}, nil
}

func (s *MessageServer) ReadMessage(ctx context.Context, req *pb.RequestReadMessage) (*pb.ResponseReadMessage, error) {
	log.Printf("ReadMessage u=%d c=%d m=%d\n", req.UserId, req.ConversationId, req.MessageId)
	err := db.SetLastSeenMessageId(req.UserId, req.ConversationId, req.MessageId)
	if err != nil {
		log.Printf("read message error %v\n", err)
		return &pb.ResponseReadMessage{
			Status:       "error",
			ErrorMessage: err.Error(),
		}, nil
	}
	return &pb.ResponseReadMessage{
		Status: "ok",
	}, nil
}

func (s *MessageServer) DecryptConversation(ctx context.Context, req *pb.RequestDecryptConversation) (*pb.ResponseDecryptConversation, error) {
	log.Printf("DecryptConversation c=%d\n", req.ConversationId)
	lastMessageId, err := db.GetLastMessageIdByConversationID(req.ConversationId)
	if err != nil {
		log.Printf("decrypt conversation error %v\n", err)
		return &pb.ResponseDecryptConversation{
			Status:       "error",
			ErrorMessage: err.Error(),
		}, nil
	}
	err = db.SetLastDecryptMessageId(req.ConversationId, lastMessageId)
	if err != nil {
		log.Printf("decrypt conversation error %v\n", err)
		return &pb.ResponseDecryptConversation{
			Status:       "error",
			ErrorMessage: err.Error(),
		}, nil
	}
	return &pb.ResponseDecryptConversation{
		Status: "ok",
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
	log.Println("gRPC server is listening on port 3100...")
	if err := srv.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func initSnowflake() {
	// 노드 ID와 노드 Id 비트 수 설정
	nodeID = 1
	nodeIDBits = uint(10)
	sf = util.NewSnowflake(nodeID, nodeIDBits)
}
