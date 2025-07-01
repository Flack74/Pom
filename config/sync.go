package config

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

type SyncProvider interface {
	Upload(localPath, remotePath string) error
	Download(remotePath, localPath string) error
	IsAvailable() bool
}

type GitHubSync struct {
	RepoURL string
	Token   string
}

type DropboxSync struct {
	AccessToken string
}

func (g *GitHubSync) IsAvailable() bool {
	_, err := exec.LookPath("git")
	return err == nil && g.RepoURL != ""
}

func (g *GitHubSync) Upload(localPath, remotePath string) error {
	if !g.IsAvailable() {
		return fmt.Errorf("git not available or repo URL not set")
	}

	// Simple git-based sync
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	syncDir := filepath.Join(configDir, "sync")
	
	// Clone or pull repo
	if _, err := os.Stat(syncDir); os.IsNotExist(err) {
		cmd := exec.Command("git", "clone", g.RepoURL, syncDir)
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to clone repo: %v", err)
		}
	} else {
		cmd := exec.Command("git", "pull")
		cmd.Dir = syncDir
		cmd.Run() // Ignore errors for pull
	}

	// Copy local data to sync dir
	cmd := exec.Command("cp", "-r", localPath, filepath.Join(syncDir, remotePath))
	if err := cmd.Run(); err != nil {
		return err
	}

	// Commit and push
	cmd = exec.Command("git", "add", ".")
	cmd.Dir = syncDir
	if err := cmd.Run(); err != nil {
		return err
	}

	cmd = exec.Command("git", "commit", "-m", fmt.Sprintf("Sync data %s", time.Now().Format("2006-01-02 15:04:05")))
	cmd.Dir = syncDir
	cmd.Run() // Ignore errors if no changes

	cmd = exec.Command("git", "push")
	cmd.Dir = syncDir
	return cmd.Run()
}

func (g *GitHubSync) Download(remotePath, localPath string) error {
	if !g.IsAvailable() {
		return fmt.Errorf("git not available or repo URL not set")
	}

	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	syncDir := filepath.Join(configDir, "sync")
	
	// Pull latest changes
	cmd := exec.Command("git", "pull")
	cmd.Dir = syncDir
	if err := cmd.Run(); err != nil {
		return err
	}

	// Copy from sync dir to local
	cmd = exec.Command("cp", "-r", filepath.Join(syncDir, remotePath), localPath)
	return cmd.Run()
}

func (d *DropboxSync) IsAvailable() bool {
	_, err := exec.LookPath("rclone")
	return err == nil && d.AccessToken != ""
}

func (d *DropboxSync) Upload(localPath, remotePath string) error {
	if !d.IsAvailable() {
		return fmt.Errorf("rclone not available or access token not set")
	}

	cmd := exec.Command("rclone", "copy", localPath, fmt.Sprintf("dropbox:%s", remotePath))
	return cmd.Run()
}

func (d *DropboxSync) Download(remotePath, localPath string) error {
	if !d.IsAvailable() {
		return fmt.Errorf("rclone not available or access token not set")
	}

	cmd := exec.Command("rclone", "copy", fmt.Sprintf("dropbox:%s", remotePath), localPath)
	return cmd.Run()
}

func GetSyncProvider(config Config) SyncProvider {
	switch config.CloudProvider {
	case "github":
		return &GitHubSync{
			RepoURL: os.Getenv("POM_GITHUB_REPO"),
			Token:   os.Getenv("POM_GITHUB_TOKEN"),
		}
	case "dropbox":
		return &DropboxSync{
			AccessToken: os.Getenv("POM_DROPBOX_TOKEN"),
		}
	default:
		return nil
	}
}

func SyncData(upload bool) error {
	config, err := LoadConfig()
	if err != nil {
		return err
	}

	if !config.CloudSync {
		return fmt.Errorf("cloud sync is disabled")
	}

	provider := GetSyncProvider(config)
	if provider == nil {
		return fmt.Errorf("no sync provider configured")
	}

	if !provider.IsAvailable() {
		return fmt.Errorf("sync provider not available")
	}

	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	if upload {
		return provider.Upload(configDir, "pom-data")
	} else {
		return provider.Download("pom-data", configDir)
	}
}