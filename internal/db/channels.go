package db

import (
	"math/rand"

	"mocha/internal/types"
)

type ChannelRecorder interface {
	CreateChannel(channel *types.Channel) (*types.Channel, error)
	GetChannelByID(channelID int64) (*types.Channel, error)
	GetUserChannels(uesrID int64) ([]types.Channel, error)
	UpdateChannel(channel *types.Channel) error
	DeleteChannel(channelID uint) error
	CreateChannelUser(cuser *types.ChannelUser) error
	SetLastSeenMessageID(userID int64, channelID int64, msgID int64) error
	SetLastDecryptMessageID(channelID int64, msgID int64) error
	GetJoinedUsers(channelID int64) ([]int64, error)
}

func (db *rdb) CreateChannel(channel *types.Channel) (*types.Channel, error) {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	result := db.con[randIdx].Create(channel)
	if result.Error != nil {
		return nil, result.Error
	}
	return channel, nil
}

func (db *rdb) GetChannelByID(channelID int64) (*types.Channel, error) {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	var channel types.Channel
	result := db.con[randIdx].First(&channel, channelID)
	if result.Error != nil {
		return nil, result.Error
	}
	return &channel, nil
}

func (db *rdb) GetUserChannels(userID int64) ([]types.Channel, error) {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	// GORM을 이용하여 channels 테이블과 channel_users 테이블을 JOIN하여 특정 사용자가 참여한 채팅방들을 가져옴
	var channels []types.Channel
	result := db.con[randIdx].
		Joins("JOIN channel_users ON channels.id = channel_users.channel_id").
		Where("channel_users.user_id = ?", userID).
		Find(&channels)

	if result.Error != nil {
		return nil, result.Error
	}

	return channels, nil
}

func (db *rdb) UpdateChannel(channel *types.Channel) error {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	result := db.con[randIdx].Save(channel)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *rdb) DeleteChannel(channelID uint) error {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	result := db.con[randIdx].Delete(&types.Channel{}, channelID)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *rdb) CreateChannelUser(cuser *types.ChannelUser) error {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	result := db.con[randIdx].Create(cuser)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *rdb) SetLastSeenMessageID(userID int64, channelID int64, msgID int64) error {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	var cu types.ChannelUser
	result := db.con[randIdx].First(&cu, channelID, userID)
	if result.Error != nil {
		return result.Error
	}

	cu.LastSeenMessageID = msgID
	result = db.con[randIdx].Save(cu)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func (db *rdb) SetLastDecryptMessageID(channelID int64, msgID int64) error {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	// Save the channel with only the LastDecryptMsgID updated
	result := db.con[randIdx].Model(&types.Channel{}).
		Where("channel_id = ?", channelID).
		Update("last_decrypt_msg_id", msgID)
	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (db *rdb) GetJoinedUsers(channelID int64) ([]int64, error) {
	randIdx := rand.Intn(10)
	db.loc[randIdx].Lock()
	defer db.loc[randIdx].Unlock()

	var channelUsers []types.ChannelUser
	result := db.con[randIdx].
		Where("channel_id = ?", channelID).
		Find(&channelUsers)
	if result.Error != nil {
		return nil, result.Error
	}

	joinedUsers := make([]int64, len(channelUsers))
	for i, cu := range channelUsers {
		joinedUsers[i] = cu.UserID
	}

	return joinedUsers, nil
}
