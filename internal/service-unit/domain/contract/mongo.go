package contract

// basic recored schema for mongo stateful unit
type MongoRecord struct {
	ID      int    `bson:"id"`
	Payload string `bson:"payload"`
}
