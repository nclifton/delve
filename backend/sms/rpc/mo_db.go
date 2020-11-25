package rpc

import (
	"fmt"
)

func (db *db) CountStoredParts(sarID string) (int64, error) {
	messageKey := fmt.Sprintf("sms.%s", sarID)
	count := db.redis.Client.SCard(messageKey)

	return count.Result()
}

func (db *db) StoreSMSPart(sarID string, messageID string, message string, partNumber string) error {

	messageKey := fmt.Sprintf("sms.%s", sarID)
	partKey := fmt.Sprintf("%s:part.%s", messageKey, partNumber)
	part := db.redis.Client.HMSet(partKey, []string{"part", partNumber, "message", message, "id", messageID})
	if part.Err() != nil {
		return part.Err()
	}
	mess := db.redis.Client.SAdd(messageKey, partKey)
	if mess.Err() != nil {
		return mess.Err()
	}

	return nil
}

type smsPart struct {
	ID      string
	Message string
}

func (db *db) GetAllSMSParts(sarID string) (map[string]smsPart, error) {
	parts := map[string]smsPart{}
	messageKey := fmt.Sprintf("sms.%s", sarID)

	partKeys := db.redis.Client.SMembers(messageKey)
	if partKeys.Err() != nil {
		return parts, partKeys.Err()
	}

	for _, partKey := range partKeys.Val() {
		part := db.redis.Client.HGetAll(partKey)
		if part.Err() != nil {
			return map[string]smsPart{}, part.Err()
		}
		parts[part.Val()["part"]] = smsPart{
			Message: part.Val()["message"],
			ID:      part.Val()["id"],
		}
	}

	return parts, nil
}
