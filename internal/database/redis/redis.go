package redis

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"

	"github.com/redis/go-redis/v9"
	"github.com/teragrammer/payment-gateway-wrapper/internal/config"
)

var (
	redisClient redis.UniversalClient
)

func ConnectRedis() redis.UniversalClient {
	if redisClient == nil {
		cfg := config.Load()

		// Load the self-signed CA certificate
		caCert, err := os.ReadFile(cfg.RedisCA)
		if err != nil {
			return nil
		}

		// Load the client certificate
		clientCert, err := os.ReadFile(cfg.RedisClientCA)
		if err != nil {
			return nil
		}

		// Load the client private key
		clientKey, err := os.ReadFile(cfg.RedisClientKEY)
		if err != nil {
			return nil
		}

		// Create a certificate pool and add the CA cert
		certPool := x509.NewCertPool()
		certPool.AppendCertsFromPEM(caCert)

		// Load the client cert and key
		clientTLSCert, err := tls.X509KeyPair(clientCert, clientKey)
		if err != nil {
			log.Fatalf("failed to load client cert and key: %v", err)
		}

		// Configure TLS
		insecure := false
		if cfg.Env == "development" {
			insecure = true
		}
		tlsConfig := &tls.Config{
			RootCAs:            certPool,                         // Trust the self-signed certificate
			Certificates:       []tls.Certificate{clientTLSCert}, // Use the client cert and key
			InsecureSkipVerify: insecure,
		}

		// Configure the Redis client with TLS
		options := &redis.UniversalOptions{
			Addrs:     []string{cfg.RedisAddr}, // Replace with your Redis server address and port
			Password:  cfg.RedisPass,           // Set if your Redis server requires a password
			DB:        0,                       // Default database
			TLSConfig: tlsConfig,               // TLS config for secure connection
		}

		// Create a Redis client
		redisClient = redis.NewUniversalClient(options)

		// Ping Redis to check the connection
		ctx := context.Background()
		_, err = redisClient.Ping(ctx).Result()
		if err != nil {
			log.Fatalf("could not connect to Redis: %v", err)
		}
	}

	return redisClient
}

func CloseRedis() {
	if redisClient != nil {
		err := redisClient.Close()
		if err != nil {
			log.Fatalf("Error closing Redis connection: %v", err)
		}
		log.Println("Redis connection closed")
	}
}
