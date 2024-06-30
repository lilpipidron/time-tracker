package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/lilpipidron/time-tracker/internal/config"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	fmt.Println(*config.MustLoad())
}
