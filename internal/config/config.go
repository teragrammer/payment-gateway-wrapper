package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/teragrammer/payment-gateway-wrapper/internal/utils"
)

type Config struct {
	Env                string
	Port               string
	MongoURI           string
	MongoDBName        string
	MongoTLSCAFILE     string
	MongoTLSCERTFILE   string
	MongoTLSKEYFILE    string
	RedisAddr          string
	RedisPass          string
	RedisCA            string
	RedisClientCA      string
	RedisClientKEY     string
	RateLimitAPIKeyReq int
	RateLimitAPIKeyInt int64
	RateLimitIPReq     int
	RateLimitIPInt     int64
	JWTSecret          string
	JWTExpiredDays     int
}

var (
	cfg *Config
)

func Load(options ...string) *Config {
	if cfg != nil {
		return cfg
	}

	env := "./.env"
	if len(options) > 0 {
		env = options[0]
	}

	_ = godotenv.Load(env)

	cfg = &Config{
		Env:                os.Getenv("ENV"),
		Port:               os.Getenv("PORT"),
		MongoURI:           os.Getenv("MONGO_URI"),
		MongoDBName:        os.Getenv("MONGO_DB_NAME"),
		MongoTLSCAFILE:     os.Getenv("MONGO_TLS_CA_FILE"),
		MongoTLSCERTFILE:   os.Getenv("MONGO_TLS_CERT_FILE"),
		MongoTLSKEYFILE:    os.Getenv("MONGO_TLS_KEY_FILE"),
		RedisAddr:          os.Getenv("REDIS_ADDR"),
		RedisPass:          os.Getenv("REDIS_PASS"),
		RedisCA:            os.Getenv("REDIS_CA"),
		RedisClientCA:      os.Getenv("REDIS_CLIENT_CA"),
		RedisClientKEY:     os.Getenv("REDIS_CLIENT_KEY"),
		RateLimitAPIKeyReq: utils.StringToInt(os.Getenv("RATE_LIMIT_API_KEY_REQ"), 100),
		RateLimitAPIKeyInt: utils.StringToInt64(os.Getenv("RATE_LIMIT_API_KEY_INT"), 1),
		RateLimitIPReq:     utils.StringToInt(os.Getenv("RATE_LIMIT_IP_REQ"), 100),
		RateLimitIPInt:     utils.StringToInt64(os.Getenv("RATE_LIMIT_IP_INT"), 1),
		JWTSecret:          os.Getenv("JWT_SECRET"),
		JWTExpiredDays:     utils.StringToInt(os.Getenv("JWT_EXPIRED_DAYS"), 1),
	}

	if cfg.Port == "" || cfg.MongoURI == "" || cfg.RedisAddr == "" {
		log.Fatal("missing required environment variables")
	}

	return cfg
}
