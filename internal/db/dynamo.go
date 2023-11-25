package db

import (
	"log"
	"log/slog"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	cf "mocha/internal/config"
)

type DynamoDatabase interface {
	MessageRecorder
}

var _ DynamoDatabase = (*mdb)(nil)

type mdb struct {
	tableName string
	svc       *dynamodb.DynamoDB
}

func NewDynamoDatabse(host string, region string, tableName string) DynamoDatabase {
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
		Endpoint: aws.String(host),   // 로컬 DynamoDB 주소
		Region:   aws.String(region), // 사용할 리전 설정
	}
	// AWS 세션 생성
	sess, err := session.NewSession(&config)
	if err != nil {
		panic(err)
	}
	// DynamoDB 서비스 생성
	svc := dynamodb.New(sess)

	resetMessages := cf.GetBool("dynamo.reset")

	if resetMessages {
		// 테이블 삭제
		_, err = svc.DeleteTable(&dynamodb.DeleteTableInput{
			TableName: aws.String(tableName),
		})
		if err != nil {
			log.Println("Failed to delete table:", err)
		}
		createMessageTable(tableName, svc)
	}
	return &mdb{tableName: tableName, svc: svc}
}

func createMessageTable(tableName string, svc *dynamodb.DynamoDB) {
	// 테이블 생성 요청
	// 테이블 생성 요청
	input := &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("channel_id"), // 파티션 키
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("id"), // 정렬 키
				KeyType:       aws.String("RANGE"),
			},
		},
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("channel_id"),
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

	slog.Info("CrateTable", "table", tableName)
}
