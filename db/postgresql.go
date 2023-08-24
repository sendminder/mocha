package db

// import (
// 	"database/sql"
// 	"fmt"
// 	"log"
// 	"mocha/config"
// 	"sync"

// 	_ "github.com/lib/pq"
// )

// func createTables(db *sql.DB) error {
// 	// bots 테이블 생성
// 	_, err := db.Exec(`
// 		CREATE TABLE IF NOT EXISTS bots (
// 			id SERIAL PRIMARY KEY,
// 			name VARCHAR(255) NOT NULL,
// 			user_id INT NOT NULL
// 		)
// 	`)
// 	if err != nil {
// 		return err
// 	}

// 	// conversations 테이블 생성
// 	_, err = db.Exec(`
// 		CREATE TABLE IF NOT EXISTS conversations (
// 			id SERIAL PRIMARY KEY,
// 			"type" VARCHAR(255) NOT NULL,
// 			name VARCHAR(255) NOT NULL,
// 			host_user_id INT NOT NULL,
// 			last_message_id INT NOT NULL
// 		)
// 	`)
// 	if err != nil {
// 		return err
// 	}

// 	// users 테이블 생성
// 	_, err = db.Exec(`
// 		CREATE TABLE IF NOT EXISTS users (
// 			id SERIAL PRIMARY KEY,
// 			name VARCHAR(255) NOT NULL,
// 			password VARCHAR(255) NOT NULL,
// 			email VARCHAR(255) NOT NULL,
// 			age INT NOT NULL,
// 			gender VARCHAR(10) NOT NULL
// 		)
// 	`)
// 	if err != nil {
// 		return err
// 	}

// 	// conversation_users 테이블 생성
// 	_, err = db.Exec(`
// 		CREATE TABLE IF NOT EXISTS conversation_users (
// 			conversation_id INT NOT NULL,
// 			user_id INT NOT NULL,
// 			last_seen_message_id INT NOT NULL,
// 			PRIMARY KEY (conversation_id, user_id),
// 			FOREIGN KEY (conversation_id) REFERENCES conversations(id),
// 			FOREIGN KEY (user_id) REFERENCES users(id)
// 		)
// 	`)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

// func ConnectPostgresql(wg *sync.WaitGroup) {
// 	defer wg.Done()
// 	// 설정 파일을 사용하여 viper 초기화

// 	// 설정 정보 가져오기
// 	host := config.GetString("postgres.host")
// 	port := config.GetInt("postgres.port")
// 	user := config.GetString("postgres.user")
// 	password := config.GetString("postgres.password")
// 	dbName := config.GetString("postgres.dbname")

// 	// 데이터베이스 연결 문자열 생성
// 	log.Printf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable\n",
// 		host, port, user, password, dbName)
// 	connStr := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
// 		host, port, user, password, dbName)

// 	// PostgreSQL 데이터베이스 연결
// 	db, err := sql.Open("postgres", connStr)
// 	if err != nil {
// 		log.Println("Failed to connect to PostgreSQL:", err)
// 		return
// 	}
// 	defer db.Close()

// 	// 테이블 생성 함수 호출
// 	err = createTables(db)
// 	if err != nil {
// 		log.Println("Failed to create tables:", err)
// 		return
// 	}
// 	log.Println("Tables created successfully.")
// }
