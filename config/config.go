package config

import (
	"fmt"
	"log"
	"os"

	"github.com/Netflix/go-env"
	"github.com/joho/godotenv"
)

type Config struct {
	Env            string `env:"APP_ENV"`
	ServerPort     string `env:"APP_SERVER_PORT"`
	FirestoreCert  string `env:"APP_FIRESTORE_CERT"`
	PrivKey        string `env:"APP_HANGO_CERT_PRIVATE"`
	PublicKey      string `env:"APP_HANGO_CERT_PUBLIC"`
	MessageBridKey string `env:"APP_API_KEY_MESSAGE_BIRD"`
	StorageName    string `env:"APP_STORAGE_NAME"`
}

func ParseConfig() (cfg *Config) {
	godotenv.Load(".env")
	switch os.Getenv("APP_ENV") {
	case "local":
		godotenv.Overload(".env")
		log.Println("Running on local environments")
		break
	default:
		log.Println("Running on production environments")
		break
	}

	cfg = new(Config)
	_, err := env.UnmarshalFromEnviron(cfg)
	if err != nil {
		fmt.Printf("Cannot read config %s\n", err.Error())
		os.Exit(1)
	}
	return
}
