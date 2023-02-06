package environment

type Environment struct {
	KafkaBroker      string
	Port             string
	KafkaGroupId     string
	RedisURL         string
	RedisPassword    string
	IncomingRequest  string
	OutgoingResponse string
}

func New(
	kafkaBroker string,
	port string,
	kafkaGroupId string,
	redisURL string,
	redisPassword string,
	IncomingRequest string,
	OutgoingResponse string,
) *Environment {
	return &Environment{
		KafkaGroupId:     kafkaGroupId,
		KafkaBroker:      kafkaBroker,
		Port:             port,
		RedisURL:         redisURL,
		RedisPassword:    redisPassword,
		IncomingRequest:  IncomingRequest,
		OutgoingResponse: OutgoingResponse,
	}
}
