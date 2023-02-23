package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadENV(key, file string) string {
	err := godotenv.Load(file)
	HandleErr("EnvLoad Error", err)
	return os.Getenv(key)
}

func WriteENV(key, value, file string) {
	env, _ := godotenv.Unmarshal(fmt.Sprintf("%s=%s", key, value))
	err := godotenv.Write(env, file)
	if err != nil {
		log.Println("There was an error writing to the dotenv file")
	}
}
