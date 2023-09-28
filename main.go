package main

import (
	"ProcGuardProject/GuardTask"
	"ProcGuardProject/config"
	"ProcGuardProject/utils"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {

	config.InitConfig()
	if err := utils.InitLogger(config.AppConfigC.Setting.LogDir, config.AppConfigC.Setting.LogMaxSize, config.AppConfigC.Setting.LogInterval, "daemon"); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}

	utils.Info("Start Task：")
	var wg sync.WaitGroup
	for taskName, task := range config.AppConfigC.Tasks {
		wg.Add(1)
		utils.Info("INFO: Start %s", taskName)
		go func(taskName string, task config.TaskConfig) {
			defer wg.Done()
			GuardTask.GuardStartTask(taskName, task)
		}(taskName, task)
	}
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	<-signalChannel
	GuardTask.StopAllTasks()
	fmt.Println("The program has terminated and all tasks have stopped")
}
