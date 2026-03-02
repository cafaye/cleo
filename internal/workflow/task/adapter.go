package task

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cafaye/cleo/internal/taskstore"
)

type Adapter struct {
	store *taskstore.Store
	now   func() time.Time
}

func NewAdapter(store *taskstore.Store) *Adapter {
	return &Adapter{store: store, now: time.Now}
}

func (a *Adapter) List(status string) (string, error) {
	tasks, err := a.store.ListTasks(context.Background(), status)
	if err != nil {
		return "", err
	}
	if len(tasks) == 0 {
		return "No tasks found.", nil
	}
	var b strings.Builder
	for _, task := range tasks {
		fmt.Fprintf(&b, "#%d [%s] %s status=%s occurrences=%d\n", task.ID, task.Severity, task.Title, task.Status, task.Occurrences)
	}
	return strings.TrimRight(b.String(), "\n"), nil
}

func (a *Adapter) Show(id int64) (string, error) {
	task, err := a.store.Task(context.Background(), id)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Task #%d\nrepo=%s\nstatus=%s severity=%s\noccurrences=%d\n\ntitle: %s\n\ndetails:\n%s", task.ID, task.RepoKey, task.Status, task.Severity, task.Occurrences, task.Title, task.Details), nil
}

func (a *Adapter) Claim(id int64) error {
	return a.store.UpdateTaskStatus(context.Background(), id, "in_progress", a.now())
}

func (a *Adapter) Close(id int64) error {
	return a.store.UpdateTaskStatus(context.Background(), id, "closed", a.now())
}
