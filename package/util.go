package util

import (
	"fmt"
	"os"
	"time"
)

func GenerateCaptureFileName() string {
	ts := time.Now().Format("2006-01-02_15-04-05")
	return "capture_" + ts + ".jsonl"
}

// Generates replay session folder and returns full path to it
func GenerateReplaySessionFolder(baseDir string) (string, error) {
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	folderName := fmt.Sprintf("replay_%s", timestamp)
	fullPath := fmt.Sprintf("%s/%s", baseDir, folderName)
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return "", err
	}
	return fullPath, nil
}

func GenerateReplayFileName(replayedFromFile string) string {
	//haven't put .jsonl at the end here because replayedFromFile has that in the end
	return "replay_" + replayedFromFile
}
