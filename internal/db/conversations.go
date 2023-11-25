package db

import (
	"math/rand"

	"mocha/internal/types"
)

type ConversationRecorder interface {
	CreateConversation(conversation *types.Conversation) (*types.Conversation, error)
	GetConversationByID(conversationID int64) (*types.Conversation, error)
	GetUserConversations(uesrId int64) ([]types.Conversation, error)
	UpdateConversation(conversation *types.Conversation) error
	DeleteConversation(conversationID uint) error
	CreateConversationUser(cuser *types.ConversationUser) error
	SetLastSeenMessageId(userId int64, convId int64, msgId int64) error
	SetLastDecryptMessageId(convId int64, msgId int64) error
	GetJoinedUsers(conversationID int64) ([]int64, error)
}

func (db *rdb) CreateConversation(conversation *types.Conversation) (*types.Conversation, error) {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	result := db.con[randIdx].Create(conversation)
	if result.Error != nil {
		return nil, result.Error
	}
	return conversation, nil
}

func (db *rdb) GetConversationByID(conversationID int64) (*types.Conversation, error) {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	var conversation types.Conversation
	result := db.con[randIdx].First(&conversation, conversationID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &conversation, nil
}

func (db *rdb) GetUserConversations(uesrId int64) ([]types.Conversation, error) {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	// GORM을 이용하여 conversations 테이블과 conversation_users 테이블을 JOIN하여 특정 사용자가 참여한 채팅방들을 가져옴
	var conversations []types.Conversation
	result := db.con[randIdx].
		Joins("JOIN conversation_users ON conversations.id = conversation_users.conversation_id").
		Where("conversation_users.user_id = ?", uesrId).
		Find(&conversations)

	if result.Error != nil {
		return nil, result.Error
	}

	return conversations, nil
}

func (db *rdb) UpdateConversation(conversation *types.Conversation) error {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	result := db.con[randIdx].Save(conversation)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *rdb) DeleteConversation(conversationID uint) error {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	result := db.con[randIdx].Delete(&types.Conversation{}, conversationID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *rdb) CreateConversationUser(cuser *types.ConversationUser) error {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	result := db.con[randIdx].Create(cuser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *rdb) SetLastSeenMessageId(userId int64, convId int64, msgId int64) error {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	var cu types.ConversationUser
	result := db.con[randIdx].First(&cu, convId, userId)
	if result.Error != nil {
		return result.Error
	}

	cu.LastSeenMessageId = msgId
	result = db.con[randIdx].Save(cu)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *rdb) SetLastDecryptMessageId(convId int64, msgId int64) error {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	// Save the conversation with only the LastDecryptMsgID updated
	result := db.con[randIdx].Model(&types.Conversation{}).
		Where("conversation_id = ?", convId).
		Update("last_decrypt_msg_id", msgId)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *rdb) GetJoinedUsers(conversationID int64) ([]int64, error) {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	var conversationUsers []types.ConversationUser
	result := db.con[randIdx].
		Where("conversation_id = ?", conversationID).
		Find(&conversationUsers)
	if result.Error != nil {
		return nil, result.Error
	}

	joinedUsers := make([]int64, len(conversationUsers))
	for i, cu := range conversationUsers {
		joinedUsers[i] = cu.UserId
	}

	return joinedUsers, nil
}
