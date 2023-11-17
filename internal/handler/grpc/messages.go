package grpc

import (
	"context"
	"log"
	"log/slog"
	"time"

	"mocha/internal/cache"
	"mocha/internal/types"
	pb "mocha/proto/message"
)

type MessageHandler interface {
	CreateMessage(ctx context.Context, req *pb.RequestCreateMessage) (*pb.ResponseCreateMessage, error)
	ReadMessage(ctx context.Context, req *pb.RequestReadMessage) (*pb.ResponseReadMessage, error)
	DecryptConversation(ctx context.Context, req *pb.RequestDecryptConversation) (*pb.ResponseDecryptConversation, error)
	PushMessage(ctx context.Context, req *pb.RequestPushMessage) (*pb.ResponsePushMessage, error)
	CreateBotMessage(ctx context.Context, req *pb.RequestBotMessage) (*pb.ResponseBotMessage, error)
}

func (s *messageServer) CreateMessage(ctx context.Context, req *pb.RequestCreateMessage) (*pb.ResponseCreateMessage, error) {
	slog.Info("CreateMessage", "text", req.Text)
	msgId := s.sf.Generate()

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
	err := s.mdb.CreateMessage(dynamoMessage)
	if err != nil {
		slog.Error("CreateMessage", "error", err)
		return &pb.ResponseCreateMessage{
			Status:       "error",
			ErrorMessage: err.Error(),
		}, nil
	}
	err = s.rdb.SetLastSeenMessageId(req.SenderId, req.ConversationId, msgId)
	if err != nil {
		slog.Error("SetLastSeenMessageId", "error", err)
		return &pb.ResponseCreateMessage{
			Status:       "error",
			ErrorMessage: err.Error(),
		}, nil
	}

	// redis 먼저 조회 후 db 조회
	joinedUsers, err := cache.GetJoinedUsers(req.ConversationId)
	if len(joinedUsers) == 0 || err != nil {
		joinedUsers, err = s.rdb.GetJoinedUsers(req.ConversationId)
		if err != nil {
			slog.Error("GetJoinedUsers", "error", err)
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

func (s *messageServer) ReadMessage(ctx context.Context, req *pb.RequestReadMessage) (*pb.ResponseReadMessage, error) {
	slog.Info("ReadMessage", "u", req.UserId, "c", req.ConversationId, "m", req.MessageId)
	err := s.rdb.SetLastSeenMessageId(req.UserId, req.ConversationId, req.MessageId)
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

func (s *messageServer) DecryptConversation(ctx context.Context, req *pb.RequestDecryptConversation) (*pb.ResponseDecryptConversation, error) {
	slog.Info("DecryptConversation", "c", req.ConversationId)
	lastMessageId, err := s.mdb.GetLastMessageIdByConversationID(req.ConversationId)
	if err != nil {
		log.Printf("decrypt conversation error %v\n", err)
		return &pb.ResponseDecryptConversation{
			Status:       "error",
			ErrorMessage: err.Error(),
		}, nil
	}
	err = s.rdb.SetLastDecryptMessageId(req.ConversationId, lastMessageId)
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

func (s *messageServer) PushMessage(ctx context.Context, req *pb.RequestPushMessage) (*pb.ResponsePushMessage, error) {
	slog.Info("PushMessage", "mid", req.Message.Id, "receivers", req.ReceiverUserIds)
	/* TODO
	1. sender_id로 발송자 정보 조회
	2. receiver_ids로 수신자 device 조회
	3. FCM 페이로드 생성
	4. receiver_id 수신자 device에 FCM 푸시 요청
	*/
	return &pb.ResponsePushMessage{
		Status: "ok",
	}, nil
}

func (s *messageServer) CreateBotMessage(ctx context.Context, req *pb.RequestBotMessage) (*pb.ResponseBotMessage, error) {
	slog.Info("CreateBotMessage", "u", req.SenderId, "c", req.ConversationId, "t", req.ConversationType)
	/* TODO
	1. OpenAPI 요청
	2. async로 response 수신
	3. lilly 로 전달해야하는데.. (Fast-Track like 필요)
	*/

	go func() {
		// OpenAPI 요청 비동기로 처리
		slog.Info("API CreateBotMessage", "u", req.SenderId, "c", req.ConversationId, "t", req.ConversationType)
	}()

	return &pb.ResponseBotMessage{
		Status: "ok",
	}, nil
}
