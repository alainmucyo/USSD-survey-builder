package service

import (
	"survey-ussd/store/kafka/producer"
	"survey-ussd/store/redis"
	"time"
)

type Service struct {
	cache    *redis.Cache
	producer *producer.Producer
}

func New(cache *redis.Cache, producer *producer.Producer) *Service {
	return &Service{cache: cache, producer: producer}
}

// CacheSomething Cache something
func (s *Service) CacheSomething(sessionId string, something string) error {
	err := s.cache.SetValue(sessionId+"-something", something, 10*time.Minute)
	return err
}

// GetSomething Get something from Redis
func (s *Service) GetSomething(sessionId string) (string, error) {
	something, err := s.cache.GetValue(sessionId + "-receiver-phone-number")
	return something, err
}
