package GuardTask

import (
	"ProcGuardProject/config"
	"ProcGuardProject/utils"
	"fmt"
	"os"
	"os/exec"
	"sync"
	"time"
)

type TaskState struct {
	IsRunning bool
	Process   *os.Process
}

var taskStates = make(map[string]*TaskState)
var taskStatesLock sync.Mutex

func GuardStartTask(taskName string, taskConfig config.TaskConfig) {
	logDir := "./tasklogs/"
	if err := utils.InitLogger(logDir, taskConfig.LogMaxSize, taskConfig.LogInterval, taskName); err != nil {
		return
	}
	stdoutLogFile := taskName + ".stdout.log"
	stderrLogFile := taskName + ".stderr.log"
	if err := utils.InitLoggers(logDir, taskConfig.LogMaxSize, taskConfig.LogInterval, stdoutLogFile, stderrLogFile); err != nil {
		return
	}
	utils.Info("Starting task: %s\n", taskName)
	utils.Info("Start Process: %s\n", taskConfig.StartProcess)
	utils.Info("Process Args: %s\n", taskConfig.ProcessArgs)
	utils.Info("Working Directory: %s\n", taskConfig.WorkingDir)
	utils.Info("Environment Variables:")
	for key, value := range taskConfig.Environment {
		utils.Info("%s=%s\n", key, value)
	}
	utils.Info("Delay: %v\n", taskConfig.Delay)

	for {
		cmd := exec.Command(taskConfig.StartProcess, taskConfig.ProcessArgs)
		cmd.Dir = taskConfig.WorkingDir
		cmd.Env = os.Environ()
		cmd.Stdout = utils.StdoutLogger.Writer()
		cmd.Stderr = utils.StderrLogger.Writer()

		for key, value := range taskConfig.Environment {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
		}

		if err := cmd.Start(); err != nil {
			utils.Error("Failed to start task %s: %v\n", taskName, err)
			return
		}

		taskStatesLock.Lock()
		taskStates[taskName] = &TaskState{
			IsRunning: true,
			Process:   cmd.Process,
		}
		taskStatesLock.Unlock()

		utils.Info("Task %s started\n", taskName)

		if err := cmd.Wait(); err != nil {
			utils.Error("Task %s exited with error: %v\n", taskName, err)
		} else {
			utils.Info("Task %s completed successfully\n", taskName)
		}

		taskStatesLock.Lock()
		taskStates[taskName] = &TaskState{
			IsRunning: false,
		}
		taskStatesLock.Unlock()

		utils.Info("Restarting task %s in %d seconds...\n", taskName, taskConfig.Delay)
		time.Sleep(time.Duration(taskConfig.Delay) * time.Second)
	}
}

func StopAllTasks() {
	taskStatesLock.Lock()
	defer taskStatesLock.Unlock()

	for taskName, taskState := range taskStates {
		if taskState.IsRunning {
			if err := taskState.Process.Kill(); err != nil {
				utils.Info("Failed to stop task %s: %v\n", taskName, err)
			} else {
				utils.Info("Task %s stopped\n", taskName)
			}
		}
	}
}
