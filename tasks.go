package main

import (
	"sync"
)

const (
	// StateWaiting 等待状态（添加后默认）
	StateWaiting = "waiting"
	// StateCompleted 完成状态
	StateCompleted = "completed"
	// StateError 失败状态
	StateError = "failed"
	// StateNone 空状态（没有构建任务）
	StateNone = "none"
)

// 保护任务状态的互斥量
var taskStateMutex sync.Mutex

// 执行任务的缓冲通道
var taskChan = make(chan Task, *cl)

// 储存任务状态的 map
var taskState = make(map[string]string)

// Task 添加到队列的任务结构体
type Task struct {
	TplKey       string
	Subs         Subs
	RunnableList []makeFunc
}

//
type makeFunc func(string, Subs) (string, error)

// addMakeTask 添加一个生成任务
func addMakeTask(task Task) string {
	go func() {
		taskChan <- task
	}()
	hash := task.Subs.Hash(task.TplKey)
	updateTaskState(hash, StateWaiting)
	return hash
}

// updateTaskState 更新任务状态
// 状态更新操作加锁
func updateTaskState(hash, state string) {
	taskStateMutex.Lock()
	defer taskStateMutex.Unlock()

	taskState[hash] = state
}

// loadTaskState 读取任务状态
// 状态更新操作加锁
func loadTaskState(hash string) (state string) {
	taskStateMutex.Lock()
	defer taskStateMutex.Unlock()

	resultState, exists := taskState[hash]
	if !exists {
		state = StateNone
	} else {
		state = resultState
	}
	return
}

// asyncMakeAction 异步生成任务启动
// goroutine 函数
func asyncMakeAction() {
	for {
		task := <-taskChan
		var curTaskHash string
		var err error
		for _, f := range task.RunnableList {
			curTaskHash, err = f(task.TplKey, task.Subs)
		}
		if err != nil {
			updateTaskState(curTaskHash, StateError)
		} else {
			updateTaskState(curTaskHash, StateCompleted)
		}
	}
}
