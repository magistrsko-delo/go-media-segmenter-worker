package main

import (
	"github.com/joho/godotenv"
	"log"
	"main/Models"
	"main/Worker"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	Models.InitEnv()
}


func main() {

	worker := Worker.InitWorker()
	defer worker.RabbitMQ.Conn.Close()
	defer worker.RabbitMQ.Ch.Close()
	worker.Work()

}