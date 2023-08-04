package db

import (
	"math/rand"
	"mocha/types"
)

func CreateConversation(conversation *types.Conversation) (*types.Conversation, error) {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	result := dbConnections[randIdx].Create(conversation)
	if result.Error != nil {
		return nil, result.Error
	}
	return conversation, nil
}

func GetConversationByID(conversationID int64) (*types.Conversation, error) {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	var conversation types.Conversation
	result := dbConnections[randIdx].First(&conversation, conversationID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &conversation, nil
}

func GetUserConversations(uesrId int64) ([]types.Conversation, error) {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	// GORM을 이용하여 conversations 테이블과 conversation_users 테이블을 JOIN하여 특정 사용자가 참여한 채팅방들을 가져옴
	var conversations []types.Conversation
	result := dbConnections[randIdx].
		Joins("JOIN conversation_users ON conversations.id = conversation_users.conversation_id").
		Where("conversation_users.user_id = ?", uesrId).
		Find(&conversations)

	if result.Error != nil {
		return nil, result.Error
	}

	return conversations, nil
}

func UpdateConversation(conversation *types.Conversation) error {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	result := dbConnections[randIdx].Save(conversation)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func DeleteConversation(conversationID uint) error {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	result := dbConnections[randIdx].Delete(&types.Conversation{}, conversationID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func CreateConversationUser(cuser *types.ConversationUser) error {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	result := dbConnections[randIdx].Create(cuser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func SetLastSeenMessageId(userId int64, convId int64, msgId int64) error {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	var cu types.ConversationUser
	result := dbConnections[randIdx].First(&cu, convId, userId)
	if result.Error != nil {
		return result.Error
	}

	cu.LastSeenMessageId = msgId
	result = dbConnections[randIdx].Save(cu)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func SetLastDecryptMessageId(convId int64, msgId int64) error {
	randIdx := rand.Intn(10)
	dbLocks[randIdx].Lock()
	defer dbLocks[randIdx].Unlock()

	// Save the conversation with only the LastDecryptMsgID updated
	result := dbConnections[randIdx].Model(&types.Conversation{}).
		Where("conversation_id = ?", convId).
		Update("last_decrypt_msg_id", msgId)
	if result.Error != nil {
		return result.Error
	}

	return nil
}
