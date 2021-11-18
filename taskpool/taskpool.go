package taskpool

import (
	"github.com/panjf2000/ants/v2"
	"os/exec"
	"sync"
)

var pool *ants.Pool
var taskMap sync.Map

type Task struct {
	Command *exec.Cmd
	Do      func()
	Stop    func()
}

func Init() {
	pool, _ = ants.NewPool(10000)
}

func Run(key string, task Task) {
	taskMap.Store(key, task)
	pool.Submit(task.Do)
}

func Cancel(key string) {
	task, _ := taskMap.Load(key)
	task.(Task).Stop()
}
