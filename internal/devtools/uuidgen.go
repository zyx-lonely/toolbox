package devtools

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

// UUIDGenResult UUID 生成结果
type UUIDGenResult struct {
	UUIDs     []string `json:"uuids"`
	Count     int      `json:"count"`
	Version   int      `json:"version"`
}

// GenerateUUIDs 批量生成 UUID
func GenerateUUIDs(count int, version int) UUIDGenResult {
	if count <= 0 {
		count = 1
	}
	if count > 100 {
		count = 100
	}

	var uuids []string
	for i := 0; i < count; i++ {
		var uid string
		switch version {
		case 4:
			uid = uuid.New().String()
		case 1:
			uid = uuid.New().String() // Go uuid 默认 v4
		default:
			uid = uuid.New().String()
		}
		uuids = append(uuids, uid)
	}

	return UUIDGenResult{UUIDs: uuids, Count: count, Version: version}
}

// ValidateUUID 验证 UUID 格式
func ValidateUUID(input string) bool {
	_, err := uuid.Parse(input)
	return err == nil
}

var _ = fmt.Sprintf
var _ = strings.TrimSpace
