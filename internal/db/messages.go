package db

import (
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"mocha/internal/types"
)

type MessageRecorder interface {
	CreateMessage(message *types.Message) error
	GetMessagesByChannelID(channelID int64) ([]types.Message, error)
	GetLastMessageIdByChannelID(channelID int64) (int64, error)
}

func (m *mdb) CreateMessage(message *types.Message) error {
	av, err := dynamodbattribute.MarshalMap(message)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("messages"),
	}

	log.Println(input)

	_, err = m.svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func (m *mdb) GetMessagesByChannelID(channelID int64) ([]types.Message, error) {
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":val": {
			N: aws.String(fmt.Sprintf("%d", channelID)),
		},
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(m.tableName),
		KeyConditionExpression:    aws.String("channel_id = :val"),
		ExpressionAttributeValues: expressionAttributeValues,
	}

	result, err := m.svc.Query(queryInput)
	if err != nil {
		return nil, err
	}

	var messages []types.Message
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, &messages)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (m *mdb) GetLastMessageIdByChannelID(channelID int64) (int64, error) {
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":val": {
			N: aws.String(fmt.Sprintf("%d", channelID)),
		},
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String(m.tableName),
		KeyConditionExpression:    aws.String("channel_id = :val"),
		ExpressionAttributeValues: expressionAttributeValues,
		ScanIndexForward:          aws.Bool(false), // Sort in descending order (latest first)
		Limit:                     aws.Int64(1),    // Get only the latest message
	}

	result, err := m.svc.Query(queryInput)
	if err != nil {
		return 0, err
	}

	if len(result.Items) == 0 {
		// No messages found for the given channel ID
		return 0, nil
	}

	// Unmarshal the last message ID
	var lastMessage types.Message
	err = dynamodbattribute.UnmarshalMap(result.Items[0], &lastMessage)
	if err != nil {
		return 0, err
	}

	return lastMessage.Id, nil
}
