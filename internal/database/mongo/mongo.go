package mongo

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"log"
	"os"
	"time"

	"github.com/teragrammer/payment-gateway-wrapper/internal/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	mongoClient *mongo.Client
	mongoDB     *mongo.Database
	ctx         context.Context
	cancel      context.CancelFunc
)

func ConnectMongo(uri string) (*mongo.Client, context.Context, context.CancelFunc, error) {
	if mongoClient == nil {
		tlsConfig, err := loadTLSConfig()
		if err != nil {
			log.Fatalf("Error TSL MongoDB: %v", err)
			return nil, nil, nil, err
		}

		clientOpts := options.Client().
			ApplyURI(uri).
			SetTLSConfig(tlsConfig)

		ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		mongoClient, err = mongo.Connect(ctx, clientOpts)
		if err != nil {
			log.Fatalf("Error connecting to MongoDB: %v", err)
			return nil, nil, nil, err
		}

		if err := mongoClient.Ping(ctx, nil); err != nil {
			log.Fatalf("Error pinging MongoDB: %v", err)
			return nil, nil, nil, err
		}
	}

	return mongoClient, ctx, cancel, nil
}

func DefaultMongo() (*mongo.Database, error) {
	if mongoDB == nil {
		cfg := config.Load()
		client, _, _, err := ConnectMongo(cfg.MongoURI)

		if err != nil {
			return nil, err
		}

		mongoDB = client.Database(cfg.MongoDBName)
	}

	return mongoDB, nil
}

func CloseMongo() {
	if mongoClient != nil {
		if err := mongoClient.Disconnect(context.TODO()); err != nil {
			log.Fatalf("Error disconnecting MongoDB: %v", err)
		}
		log.Println("MongoDB connection closed")
	}
}

func loadTLSConfig() (*tls.Config, error) {
	cfg := config.Load()

	caFile := cfg.MongoTLSCAFILE
	certFile := cfg.MongoTLSCERTFILE
	keyFile := cfg.MongoTLSKEYFILE

	insecure := false
	if cfg.Env == "development" {
		insecure = true
	}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: insecure,
	}

	// Load CA cert
	if caFile != "" {
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return nil, err
		}

		caPool := x509.NewCertPool()
		caPool.AppendCertsFromPEM(caCert)
		tlsConfig.RootCAs = caPool
	}

	// Load client cert (mTLS)
	if certFile != "" && keyFile != "" {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return nil, err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	tlsConfig.MinVersion = tls.VersionTLS12
	return tlsConfig, nil
}
