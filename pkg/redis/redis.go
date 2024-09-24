package redis_db

import (
	"context"
	"encoding/json"

	"github.com/go-redis/redis/v8"
)

type RequestStatus struct {
	IsDepartment  bool `json:"is_department"`
	IsTittle      bool `json:"is_tittle"`
	IsDiscription bool `json:"is_discription"`
}

type Request struct {
	UserId      int64         `json:"user_id"`
	ChatId      int64         `json:"chat_id"`
	Department  string        `json:"department"`
	Tittle      string        `json:"tittle"`
	Discription string        `json:"discription"`
	Status      RequestStatus `json:"status"`
}

func SaveRequest(rdb *redis.Client, requestId string, request Request) error {
	data, err := json.Marshal(request)
	if err != nil {
		return err
	}

	return rdb.Set(context.Background(), requestId, data, 0).Err()
}

func GetRequest(rdb *redis.Client, requestId string) (Request, error) {
	data, err := rdb.Get(context.Background(), requestId).Bytes()
	if err != nil {
		return Request{}, err
	}

	var request Request
	err = json.Unmarshal(data, &request)
	if err != nil {
		return Request{}, err
	}

	return request, nil
}

func DeleteRequest(rdb *redis.Client, requestId string) error {
	exists, err := rdb.Exists(context.Background(), requestId).Result()
	if err != nil {
		return err
	}
	if exists > 0 {
		_, err = rdb.Del(context.Background(), requestId).Result()
		if err != nil {
			return err
		}
	}
	return nil
}
