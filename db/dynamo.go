package db

import (
	"fmt"
	"mocha/types"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

var svc *dynamodb.DynamoDB

func ConnectDynamo(wg *sync.WaitGroup) {
	defer wg.Done()
	// AWS 세션 생성
	// sess, err := session.NewSession(&aws.Config{
	// 	Region:      aws.String("us-west-2"),                                                    // 사용할 AWS 리전 설정
	// 	Credentials: credentials.NewStaticCredentials("YOUR_ACCESS_KEY", "YOUR_SECRET_KEY", ""), // AWS 액세스 키와 시크릿 키 설정
	// })

	// if err != nil {
	// 	// 세션 생성 실패
	// 	panic(err)
	// }

	config := aws.Config{
		Endpoint: aws.String("http://localhost:8001"), // 로컬 DynamoDB 주소
		Region:   aws.String("us-west-2"),             // 사용할 리전 설정
	}
	// AWS 세션 생성
	sess, err := session.NewSession(&config)
	if err != nil {
		panic(err)
	}
	// DynamoDB 서비스 생성
	svc = dynamodb.New(sess)

	// 삭제할 테이블 이름
	tableName := "messages"

	// 테이블 삭제
	_, err = svc.DeleteTable(&dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})
	if err != nil {
		fmt.Println("Failed to delete table:", err)
		return
	}

	createMessageTable()
}

func createMessageTable() {
	// 테이블 생성 요청
	// 테이블 생성 요청
	tableName := "messages"
	input := &dynamodb.CreateTableInput{
		TableName: aws.String("messages"),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("conversation_id"), // 파티션 키
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("id"), // 정렬 키
				KeyType:       aws.String("RANGE"),
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("conversation_id"),
				AttributeType: aws.String("N"),
			},
			{
				AttributeName: aws.String("id"),
				AttributeType: aws.String("N"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(5), // 읽기 처리량
			WriteCapacityUnits: aws.Int64(5), // 쓰기 처리량
		},
	}

	// 테이블 생성 요청 보내기
	_, err := svc.CreateTable(input)
	if err != nil {
		panic(err)
	}

	fmt.Println("테이블 생성 완료:", tableName)
}

func CreateMessage(message *types.Message) error {
	message.CreatedTime = time.Now().UTC().Format(time.RFC3339)
	message.UpdatedTime = message.CreatedTime

	av, err := dynamodbattribute.MarshalMap(message)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String("messages"),
	}

	fmt.Println(input)

	_, err = svc.PutItem(input)
	if err != nil {
		return err
	}

	return nil
}

func GetMessagesByConversationID(conversationID int64) ([]types.Message, error) {
	expressionAttributeValues := map[string]*dynamodb.AttributeValue{
		":val": {
			N: aws.String(fmt.Sprintf("%d", conversationID)),
		},
	}

	queryInput := &dynamodb.QueryInput{
		TableName:                 aws.String("messages"),
		KeyConditionExpression:    aws.String("conversation_id = :val"),
		ExpressionAttributeValues: expressionAttributeValues,
	}

	result, err := svc.Query(queryInput)
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
