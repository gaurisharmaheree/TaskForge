# TaskForge

A background job processing system built with **Go and Redis**.

TaskForge takes jobs from an API, puts them into a Redis queue, and processes them in the background using workers.

## How it works

```text
User
 ↓
Go API
 ↓
Redis Queue
 ↓
Worker
 ↓
Process Job

The API adds jobs to Redis, while the worker picks them up and processes them separately.

V1 Features
Go HTTP API for adding jobs
Redis-based job queue
Background worker
Concurrent job processing using goroutines
Basic retry handling
Basic metrics and logging
Example Job
{
  "type": "send_email",
  "retries": 3,
  "payload": {
    "to": "user@gmail.com",
    "subject": "Testing TaskForge"
  }
}

Job types are not limited to email. Other examples could be image processing, PDF generation, file processing, etc.

Tech Stack
Go
Redis
Goroutines
HTTP
What's Next?

TaskForge will be improved step by step with:

Exponential Backoff
Jitter
Dead Letter Queue (DLQ)
Stale Job Recovery
Graceful Shutdown
Multiple Workers

These features will be added gradually to handle real-world problems such as repeated failures, stuck jobs, and worker crashes.