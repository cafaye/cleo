package qa

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/cafaye/cleo/internal/taskstore"
)

type Adapter struct {
	store   *taskstore.Store
	repoKey string
	now     func() time.Time
}

func NewAdapter(store *taskstore.Store, repoKey string) *Adapter {
	return &Adapter{store: store, repoKey: strings.TrimSpace(repoKey), now: time.Now}
}

func (a *Adapter) Start(source string, ref string, goals string) (int64, error) {
	s, err := a.store.StartSession(context.Background(), source, ref, goals, a.now())
	if err != nil {
		return 0, err
	}
	return s.ID, nil
}

func (a *Adapter) LogIssue(sessionID int64, title string, details string, severity string) (int64, bool, error) {
	if severity == "" {
		severity = "medium"
	}
	if _, err := a.store.Session(context.Background(), sessionID); err != nil {
		return 0, false, err
	}
	key := dedupeKey(a.repoKey, title, details)
	task, created, err := a.store.UpsertOpenTask(context.Background(), taskstore.Task{
		SessionID: sessionID,
		RepoKey:   a.repoKey,
		Title:     title,
		Details:   details,
		Severity:  severity,
		DedupeKey: key,
	}, a.now())
	if err != nil {
		return 0, false, err
	}
	return task.ID, created, nil
}

func (a *Adapter) Finish(sessionID int64, verdict string) error {
	if err := validateVerdict(verdict); err != nil {
		return err
	}
	return a.store.FinishSession(context.Background(), sessionID, verdict, a.now())
}

func (a *Adapter) Report(sessionID int64) (string, error) {
	session, err := a.store.Session(context.Background(), sessionID)
	if err != nil {
		return "", err
	}
	tasks, err := a.store.TasksBySession(context.Background(), sessionID)
	if err != nil {
		return "", err
	}
	var lines []string
	lines = append(lines, fmt.Sprintf("QA session %d", session.ID))
	lines = append(lines, fmt.Sprintf("source=%s ref=%s verdict=%s", session.Source, session.Ref, emptyDefault(session.Verdict, "pending")))
	lines = append(lines, fmt.Sprintf("goals=%s", session.Goals))
	lines = append(lines, "")
	if len(tasks) == 0 {
		lines = append(lines, "No tasks logged.")
		return strings.Join(lines, "\n"), nil
	}
	lines = append(lines, "Tasks:")
	for _, task := range tasks {
		lines = append(lines, fmt.Sprintf("- #%d [%s] %s (status=%s occurrences=%d)", task.ID, task.Severity, task.Title, task.Status, task.Occurrences))
	}
	return strings.Join(lines, "\n"), nil
}

func dedupeKey(repoKey string, title string, details string) string {
	normalized := strings.ToLower(strings.TrimSpace(repoKey) + "|" + strings.TrimSpace(title) + "|" + strings.TrimSpace(details))
	sum := sha256.Sum256([]byte(normalized))
	return hex.EncodeToString(sum[:])
}

func validateVerdict(v string) error {
	switch strings.TrimSpace(v) {
	case "pass", "fail", "blocked":
		return nil
	default:
		return fmt.Errorf("--verdict must be pass|fail|blocked")
	}
}

func emptyDefault(v string, d string) string {
	if strings.TrimSpace(v) == "" {
		return d
	}
	return v
}
