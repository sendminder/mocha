package handler

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"mocha/cache"
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

	now := time.Now().UTC().Format(time.RFC3339)
	dynamoMessage := &types.Message{
		Id:             msgId,
		ConversationId: req.ConversationId,
		SenderID:       req.SenderId,
		Text:           req.Text,
		Animal:         "cat",
		CreatedTime:    now,
		UpdatedTime:    now,
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
		log.Printf("SetLastSeenMessageId error %v\n", err)
		return &pb.ResponseCreateMessage{
			Status:       "error",
			ErrorMessage: err.Error(),
		}, nil
	}

	// redis 먼저 조회 후 db 조회
	joinedUsers, err := cache.GetJoinedUsers(req.ConversationId)
	if len(joinedUsers) == 0 || err != nil {
		joinedUsers, err = db.GetJoinedUsers(req.ConversationId)
		if err != nil {
			log.Printf("GetJoinedUsers error %v\n", err)
			return &pb.ResponseCreateMessage{
				Status:       "error",
				ErrorMessage: err.Error(),
			}, nil
		}
		cache.SetJoinedUsers(req.ConversationId, joinedUsers)
	}

	// 새로운 Message 생성
	newMessage := &pb.Message{
		Id:             msgId,
		ConversationId: req.ConversationId,
		SenderId:       req.SenderId,
		Text:           req.Text,
		Animal:         "cat",
		CreatedTime:    now,
		UpdatedTime:    now,
	}

	return &pb.ResponseCreateMessage{
		Status:      "ok",
		Message:     newMessage,
		JoinedUsers: joinedUsers,
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

func (s *MessageServer) PushMessage(ctx context.Context, req *pb.RequestPushMessage) (*pb.ResponsePushMessage, error) {
	log.Printf("PushMessage mid=%d receivers=%v\n", req.Message.Id, req.ReceiverUserIds)
	/* TODO
	1. sender_id로 발송자 정보 조회
	2. receiver_ids로 수신자 device 조회
	3. FCM 페이로드 생성
	4. reciever_ids의 수신자 device에 FCM 푸시 요청
	*/
	return &pb.ResponsePushMessage{
		Status: "ok",
	}, nil
}

func (s *MessageServer) CreateBotMessage(ctx context.Context, req *pb.RequestBotMessage) (*pb.ResponseBotMessage, error) {
	/* TODO
	1. OpenAPI 요청
	2. async로 response 수신
	3. lilly 로 전달해야하는데.. (Fast-Track like 필요)
	*/

	go func() {
		// OpenAPI 요청 비동기로 처리
		log.Printf("CreateBotMessage u=%d c=%d t=%s\n", req.SenderId, req.ConversationId, req.ConversationType)
	}()

	return &pb.ResponseBotMessage{
		Status: "ok",
	}, nil
}

func StartGrpc(wg *sync.WaitGroup) {
	defer wg.Done()
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

func init() {
	// 노드 ID와 노드 Id 비트 수 설정
	nodeID = 1
	nodeIDBits = uint(10)
	sf = util.NewSnowflake(nodeID, nodeIDBits)
	log.Println("init snowflake")
}
