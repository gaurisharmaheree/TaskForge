package worker

import (
	"TaskForge/internal/task"
	"fmt"
	"time"
)

func ProcessTask(taskToExecute task.Task) error {
	if taskToExecute.Payload == nil {
		return fmt.Errorf("payload is empty")
	}

	switch taskToExecute.Type {

	case "send_email":
		time.Sleep(2 * time.Second)
		fmt.Println("Sending email to", taskToExecute.Payload["to"], "with subject", taskToExecute.Payload["subject"])
		return nil

	case "resize_image":
		fmt.Println("Resizing image to x:", taskToExecute.Payload["new_x"], "y:", taskToExecute.Payload["new_y"])
		return nil

	case "generate_pdf":
		fmt.Println("Generating pdf...")
		return nil

	case "":
		return fmt.Errorf("task type is empty")

	default:
		return fmt.Errorf("unsupported task type: %s", taskToExecute.Type)
	}
}