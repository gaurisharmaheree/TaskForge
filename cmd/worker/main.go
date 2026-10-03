package main

import (
	"TaskForge/internal/logger"
	"TaskForge/internal/task"
	"TaskForge/internal/worker"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
)

var totalJobsInQueue int64
var jobsDone int = 0
var jobsFailed int = 0

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

	port := ":" + os.Getenv("PORT_WORKER")

	rdb := connectRedis()

	var wg sync.WaitGroup
	numWorkers := 3 // how many goroutines pull jobs concurrently

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go runWorker(rdb, &wg, i)
	}

	http.HandleFunc("/metrics", metricsHandler)

	log.Println("Worker starting on port", port)

	go func() {
		err := http.ListenAndServe(port, nil)
		if err != nil {
			log.Fatal("Error starting server:", err)
		}
	}()

	wg.Wait()
	fmt.Println("All workers finished executing")
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Only GET requests allowed", http.StatusMethodNotAllowed)
		return
	}

	metrics := task.Metrics{
		TotalJobsInQueue: totalJobsInQueue,
		JobsDone:         jobsDone,
		JobsFailed:       jobsFailed,
	}

	res, err := json.Marshal(metrics)
	if err != nil {
		http.Error(w, "Could not encode metrics", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, string(res))
}

func runWorker(rdb *redis.Client, wg *sync.WaitGroup, workerID int) {
	defer wg.Done()

	for {
		res, err := rdb.BLPop(ctx, 0, "task_queue").Result()
		if err != nil {
			log.Println("Error reading from Redis:", err)
			break
		}

		totalJobsInQueue, _ = rdb.LLen(ctx, "task_queue").Result()

		var taskToExecute task.Task
		err = json.Unmarshal([]byte(res[1]), &taskToExecute)
		if err != nil {
			log.Println("Could not parse task from Redis:", err)
			continue
		}

		err = worker.ProcessTask(taskToExecute)

		if err != nil {
			jobsFailed++
			logger.LogFailure(taskToExecute, err)
			log.Printf("[worker-%d] Task failed: %v", workerID, err)

			taskToExecute.Retries--
			if taskToExecute.Retries > 0 {
				retryBytes, _ := json.Marshal(taskToExecute)
				rdb.RPush(ctx, "task_queue", retryBytes)
				log.Printf("[worker-%d] Retrying task, retries left: %d", workerID, taskToExecute.Retries)
			} else {
				log.Printf("[worker-%d] Task permanently failed after all retries", workerID)
			}
		} else {
			jobsDone++
			logger.LogSuccess(taskToExecute)
			log.Printf("[worker-%d] Task completed successfully", workerID)
		}
	}
}