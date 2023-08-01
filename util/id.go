// id.go

package util

import (
	"sync"
	"time"
)

// Snowflake 구조체는 Snowflake ID 생성기를 나타냅니다.
type Snowflake struct {
	mu          sync.Mutex
	lastTime    int64
	sequence    int
	nodeID      int
	nodeIDBits  uint
	sequenceMax int
}

// NewSnowflake 함수는 새로운 Snowflake 인스턴스를 생성합니다.
func NewSnowflake(nodeID int, nodeIDBits uint) *Snowflake {
	sf := &Snowflake{
		nodeID:      nodeID,
		nodeIDBits:  nodeIDBits,
		sequenceMax: 1<<12 - 1,
	}
	return sf
}

// Generate 함수는 새로운 Snowflake ID를 생성합니다.
func (sf *Snowflake) Generate() int64 {
	sf.mu.Lock()
	defer sf.mu.Unlock()

	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	if currentTime == sf.lastTime {
		sf.sequence = (sf.sequence + 1) & sf.sequenceMax
		if sf.sequence == 0 {
			// 시퀀스가 최대값을 넘어가면 기다렸다가 다음 밀리초로 넘어감
			for currentTime <= sf.lastTime {
				currentTime = time.Now().UnixNano() / int64(time.Millisecond)
			}
		}
	} else {
		sf.sequence = 0
	}

	sf.lastTime = currentTime
	id := (currentTime << (sf.nodeIDBits + 12)) | (int64(sf.nodeID) << 12) | int64(sf.sequence)
	return id
}
