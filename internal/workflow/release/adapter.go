package release

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/cafaye/cleo/internal/ghcli"
	releaseruntime "github.com/cafaye/cleo/internal/workflow/release/runtime"
)

type Adapter struct {
	gh   *ghcli.Client
	repo string
	root string
}

func NewAdapter(owner, repo string) *Adapter {
	root, _ := os.Getwd()
	return &Adapter{gh: ghcli.New(), repo: owner + "/" + repo, root: root}
}

func (a *Adapter) CheckGitClean() error {
	out, err := runLocal("git", "status", "--porcelain")
	if err != nil {
		return err
	}
	if strings.TrimSpace(out) != "" {
		return fmt.Errorf("working tree is not clean")
	}
	return nil
}

func (a *Adapter) EnsureReleaseMissing(version string) error {
	_, err := a.gh.Run("release", "view", version, "--repo", a.repo)
	if err == nil {
		return fmt.Errorf("release %s already exists", version)
	}
	if strings.Contains(strings.ToLower(err.Error()), "release not found") {
		return nil
	}
	if strings.Contains(strings.ToLower(err.Error()), "http 404") {
		return nil
	}
	return err
}

func (a *Adapter) Cut(version string) error {
	if _, err := runLocal("git", "tag", version); err != nil {
		return err
	}
	_, err := runLocal("git", "push", "origin", version)
	return err
}

func (a *Adapter) Publish(version string, draft bool, generateNotes bool) error {
	args := []string{"release", "create", version, "--repo", a.repo, "--verify-tag"}
	if draft {
		args = append(args, "--draft")
	}
	if generateNotes {
		args = append(args, "--generate-notes")
	}
	if releaseruntime.DetectGo(a.root) {
		assets, err := releaseruntime.BuildGoReleaseArtifacts(version)
		if err != nil {
			return err
		}
		args = append(args, assets...)
	}
	_, err := a.gh.Run(args...)
	return err
}

func (a *Adapter) Verify(version string) error {
	out, err := a.gh.Run("release", "view", version, "--repo", a.repo, "--json", "tagName,url,isDraft,isPrerelease,assets")
	if err != nil {
		return err
	}
	if !releaseruntime.DetectGo(a.root) {
		return nil
	}
	var payload struct {
		Assets []struct {
			Name string `json:"name"`
		} `json:"assets"`
	}
	if err := ghcli.DecodeJSON(out, &payload); err != nil {
		return err
	}
	have := map[string]bool{}
	for _, a := range payload.Assets {
		have[strings.TrimSpace(a.Name)] = true
	}
	for _, name := range releaseruntime.ExpectedGoAssetNames(version) {
		if !have[name] {
			return fmt.Errorf("missing release asset: %s", name)
		}
	}
	return nil
}

func runLocal(name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("%s %s: %s", name, strings.Join(args, " "), strings.TrimSpace(string(out)))
	}
	return string(out), nil
}
