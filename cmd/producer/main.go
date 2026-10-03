package main

import (
	"TaskForge/internal/task"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx = context.Background()

func connectRedis() *redis.Client {
	redisURL := os.Getenv("REDIS_URL")
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		log.Fatal("Could not parse Redis URL:", err)
	}
	return redis.NewClient(opt)
}

func main() {
	godotenv.Load()

	port := ":" + os.Getenv("PORT_PRODUCER")

	rdb = connectRedis()

	http.HandleFunc("/enqueue", enqueueHandler)

	log.Println("Producer starting on port", port)

	err := http.ListenAndServe(port, nil)
	if err != nil {
		log.Fatal("Error starting server:", err)
	}
}

func enqueueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.Error(w, "Only POST requests allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain")

	var newTask task.Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		http.Error(w, "Bad request: invalid JSON", http.StatusBadRequest)
		return
	}

	if newTask.Type == "" {
		http.Error(w, "Bad request: 'type' is required", http.StatusBadRequest)
		return
	}

	if newTask.Type == "send_email" {
		if newTask.Payload["to"] == nil || newTask.Payload["subject"] == nil {
			http.Error(w, "Bad request: send_email requires 'to' and 'subject' in payload", http.StatusBadRequest)
			return
		}
	}

	taskBytes, err := json.Marshal(newTask)
	if err != nil {
		http.Error(w, "Bad request: could not process task", http.StatusBadRequest)
		return
	}

	queueLength, err := rdb.RPush(ctx, "task_queue", taskBytes).Result()
	if err != nil {
		http.Error(w, "Internal server error: could not queue task", http.StatusInternalServerError)
		return
	}

	log.Println("Queue length is now:", queueLength)

	fmt.Fprintf(w, "Task of type '%s' added to the queue successfully", newTask.Type)
}