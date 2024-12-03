package configs

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Mode            string
	ApiAddress      string
	AccountAddress  string
	OrdersAddress   string
	ProductsAddress string
	ReviewsAddress  string
	Dsn             string
	JWTSecret       string
}

func LoadConfig() *Config {
	err := godotenv.Load("./configs/.env")
	if err != nil {
		panic(err)
	}

	return &Config{
		Mode:            os.Getenv("Mode"),
		Dsn:             os.Getenv("DSN"),
		ApiAddress:      os.Getenv("ApiAddress"),
		AccountAddress:  os.Getenv("AccountAddress"),
		OrdersAddress:   os.Getenv("ApiAddress"),
		ProductsAddress: os.Getenv("AccountAddress"),
		ReviewsAddress:  os.Getenv("ApiAddress"),
	}
}
