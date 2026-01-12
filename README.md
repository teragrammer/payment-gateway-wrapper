### Supported Payment Gateways

- [Xendit](https://docs.xendit.co/docs/overview): Comprehensive suite of APIs and no code tools, you can seamlessly accept payments, send payouts, and manage your business finances with ease.
- [PayMongo](https://developers.paymongo.com/docs/introduction): PayMongo is a payment gateway that allows businesses to easily accept payments from their customers online via credit/debit card (straight), BillEase (buy now, pay later), e-wallets, and online banking.

### Getting Started

- Clone the repository

```
$ git clone https://github.com/teragrammer/payment-gateway-wrapper
$ cd payment-gateway-wrapper
```

- Configure your .env (.env.example)

- Initialize Docker

```
$ docker compose up --build
$ sh scripts/run-docker.sh
```

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

### Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository.
2. Create a new branch (git checkout -b feature/your-feature).
3. Commit your changes (git commit -m 'Add your feature').
4. Push to the branch (git push origin feature/your-feature).
5. Open a Pull Request.
   Please ensure your code follows the project's coding standards and includes relevant tests.

### Hire Me

```
If you like this project and need help with development, customization, or integration, feel free to reach out!

I’m available for freelance work, consulting, and collaboration.

Thank you for checking out Payment Gateway Wrapper API!
Feel free to contribute or open issues.
```

### License

This project is licensed under the MIT License. See the LICENSE file for details.