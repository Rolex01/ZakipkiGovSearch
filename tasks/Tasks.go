package tasks

import (
	"os/exec"
	"path"
	"log"
)

var tasksRoot string = "./"

func SetTasksRoot(newTasksRoot string) {
	tasksRoot = newTasksRoot
}

func GetTask(taskName string) *exec.Cmd {
	taskPath := path.Join(tasksRoot, taskName + ".sh")
	log.Println("Tasks run:", taskPath)
	taskRunner := exec.Command("sh", taskPath)
	return taskRunner
}
