### Supported Payment Gateways
- [Xendit](https://docs.xendit.co/docs/overview)
- [PayMongo](https://developers.paymongo.com/docs/introduction)

### Database sample usage
```GO
// Example of interacting with MongoDB
collection := mongoClient.Database("mydb").Collection("mycollection")
_, err := collection.InsertOne(context.TODO(), map[string]string{"name": "test"})
if err != nil {
	log.Fatalf("Error inserting into MongoDB: %v", err)
}

// Example of interacting with Redis
err = redisClient.Set(context.Background(), "mykey", "myvalue", 10*time.Second).Err()
if err != nil {
	log.Fatalf("Error setting Redis key: %v", err)
}
```