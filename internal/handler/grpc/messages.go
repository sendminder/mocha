package grpc

import (
	"context"
	"log/slog"
	"time"

	"mocha/internal/types"
	pb "mocha/proto/message"
)

type MessageHandler interface {
	CreateMessage(ctx context.Context, req *pb.RequestCreateMessage) (*pb.ResponseCreateMessage, error)
	ReadMessage(ctx context.Context, req *pb.RequestReadMessage) (*pb.ResponseReadMessage, error)
	DecryptChannel(ctx context.Context, req *pb.RequestDecryptChannel) (*pb.ResponseDecryptChannel, error)
	PushMessage(ctx context.Context, req *pb.RequestPushMessage) (*pb.ResponsePushMessage, error)
	CreateBotMessage(ctx context.Context, req *pb.RequestBotMessage) (*pb.ResponseBotMessage, error)
}

func (s *messageServer) CreateMessage(_ context.Context, req *pb.RequestCreateMessage) (*pb.ResponseCreateMessage, error) {
	slog.Info("CreateMessage", "text", req.Text)
	msgID := s.sf.Generate()

	now := time.Now().UTC().Format(time.RFC3339)
	dynamoMessage := &types.Message{
		ID:          msgID,
		ChannelID:   req.ChannelId,
		SenderID:    req.SenderId,
		Text:        req.Text,
		Animal:      "cat",
		CreatedTime: now,
		UpdatedTime: now,
	}
	err := s.mdb.CreateMessage(dynamoMessage)
	if err != nil {
		slog.Error("CreateMessage", "error", err)
		return &pb.ResponseCreateMessage{
			Status:       "error",
			ErrorMessage: err.Error(),
		}, nil
	}
	err = s.rdb.SetLastSeenMessageID(req.SenderId, req.ChannelId, msgID)
	if err != nil {
		slog.Error("SetLastSeenMessageID", "error", err)
		return &pb.ResponseCreateMessage{
			Status:       "error",
			ErrorMessage: err.Error(),
		}, nil
	}

	// redis 먼저 조회 후 db 조회
	joinedUsers, err := s.cache.GetJoinedUsers(req.ChannelId)
	if len(joinedUsers) == 0 || err != nil {
		joinedUsers, err = s.rdb.GetJoinedUsers(req.ChannelId)
		if err != nil {
			slog.Error("GetJoinedUsers", "error", err)
			return &pb.ResponseCreateMessage{
				Status:       "error",
				ErrorMessage: err.Error(),
			}, nil
		}
		err = s.cache.SetJoinedUsers(req.ChannelId, joinedUsers)
		if err != nil {
			slog.Error("SetJoinedUsers", "error", err)
			return nil, err
		}
	}

	// 새로운 Message 생성
	newMessage := &pb.Message{
		Id:          msgID,
		ChannelId:   req.ChannelId,
		SenderId:    req.SenderId,
		Text:        req.Text,
		Animal:      "cat",
		CreatedTime: now,
		UpdatedTime: now,
	}

	return &pb.ResponseCreateMessage{
		Status:      "ok",
		Message:     newMessage,
		JoinedUsers: joinedUsers,
	}, nil
}

func (s *messageServer) ReadMessage(_ context.Context, req *pb.RequestReadMessage) (*pb.ResponseReadMessage, error) {
	slog.Info("ReadMessage", "u", req.UserId, "c", req.ChannelId, "m", req.MessageId)
	err := s.rdb.SetLastSeenMessageID(req.UserId, req.ChannelId, req.MessageId)
	if err != nil {
		slog.Error("read message error", "error", err)
		return &pb.ResponseReadMessage{
			Status:       "error",
			ErrorMessage: err.Error(),
		}, nil
	}
	return &pb.ResponseReadMessage{
		Status: "ok",
	}, nil
}

func (s *messageServer) DecryptChannel(_ context.Context, req *pb.RequestDecryptChannel) (*pb.ResponseDecryptChannel, error) {
	slog.Info("DecryptChannel", "c", req.ChannelId)
	lastMessageID, err := s.mdb.GetLastMessageIDByChannelID(req.ChannelId)
	if err != nil {
		slog.Error("decrypt channel error", "error", err)
		return &pb.ResponseDecryptChannel{
			Status:       "error",
			ErrorMessage: err.Error(),
		}, nil
	}
	err = s.rdb.SetLastDecryptMessageID(req.ChannelId, lastMessageID)
	if err != nil {
		slog.Error("decrypt channel error", "error", err)
		return &pb.ResponseDecryptChannel{
			Status:       "error",
			ErrorMessage: err.Error(),
		}, nil
	}
	return &pb.ResponseDecryptChannel{
		Status: "ok",
	}, nil
}

func (s *messageServer) PushMessage(_ context.Context, req *pb.RequestPushMessage) (*pb.ResponsePushMessage, error) {
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

func (s *messageServer) CreateBotMessage(_ context.Context, req *pb.RequestBotMessage) (*pb.ResponseBotMessage, error) {
	slog.Info("CreateBotMessage", "u", req.SenderId, "c", req.ChannelId, "t", req.ChannelType)
	/* TODO
	1. OpenAPI 요청
	2. async로 response 수신
	3. lilly 로 전달해야하는데.. (Fast-Track like 필요)
	*/

	go func() {
		// OpenAPI 요청 비동기로 처리
		slog.Info("API CreateBotMessage", "u", req.SenderId, "c", req.ChannelId, "t", req.ChannelType)
	}()

	return &pb.ResponseBotMessage{
		Status: "ok",
	}, nil
}
