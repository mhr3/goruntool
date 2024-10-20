package main

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: goruntool module@version [flags]")
		os.Exit(1)
	}

	arg := os.Args[1]
	parts := strings.Split(arg, "@")
	if len(parts) != 2 {
		fmt.Println("Argument must be in the form module@version")
		os.Exit(1)
	}

	module := parts[0]
	version := parts[1]

	fnv32a := fnv.New32a()
	fnv32a.Write([]byte(module))
	fnv32a.Write([]byte("@"))
	fnv32a.Write([]byte(version))

	// Create a temporary directory
	tempDir, err := os.MkdirTemp("", fmt.Sprintf("goruntool-%x", fnv32a.Sum(nil)))
	if err != nil {
		fmt.Printf("Failed to create temporary directory: %v\n", err)
		os.Exit(1)
	}

	repo := strings.TrimSuffix(module, "/")
	subDir := ""
	cmd := "binary"

	if idx := strings.LastIndex(module, "/"); idx != -1 {
		cmd = module[idx+1:]
	}

	if strings.HasPrefix(module, "github.com/") {
		// extract the repository name
		repoParts := strings.Split(module, "/")
		if len(repoParts) >= 3 {
			repo = "https://" + strings.Join(repoParts[0:3], "/") + ".git"
			subDir = strings.Join(repoParts[3:], "/")
		}
	} else {
		repo = "https://" + repo + ".git"
	}

	var stdoutBuf, stderrBuf bytes.Buffer

	// Clone the repository
	cloneCmd := exec.Command("git", "clone", "--branch", version, repo, tempDir)
	cloneCmd.Stdout = &stdoutBuf
	cloneCmd.Stderr = &stderrBuf
	if err := cloneCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to clone repository: %v\n", err)
		stdOutStr := strings.TrimSpace(stdoutBuf.String())
		if len(stdOutStr) > 0 {
			fmt.Fprintf(os.Stderr, "%s\n", stdoutBuf.String())
		}
		stdErrStr := strings.TrimSpace(stderrBuf.String())
		if len(stdErrStr) > 0 {
			fmt.Fprintf(os.Stderr, "%s\n", stderrBuf.String())
		}
		os.Exit(1)
	}

	stdoutBuf.Reset()
	stderrBuf.Reset()

	// Build the module
	buildCmd := exec.Command("go", "build", "-o", cmd, "./"+subDir)
	buildCmd.Dir = tempDir
	buildCmd.Stdout = &stdoutBuf
	buildCmd.Stderr = &stderrBuf
	if err := buildCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to build module: %v\n", err)
		stdOutStr := strings.TrimSpace(stdoutBuf.String())
		if len(stdOutStr) > 0 {
			fmt.Fprintf(os.Stderr, "%s\n", stdoutBuf.String())
		}
		stdErrStr := strings.TrimSpace(stderrBuf.String())
		if len(stdErrStr) > 0 {
			fmt.Fprintf(os.Stderr, "%s\n", stderrBuf.String())
		}
		os.Exit(1)
	}

	// Run the built binary
	binaryPath := filepath.Join(tempDir, cmd)

	runCmd := exec.Command(binaryPath, os.Args[2:]...)
	runCmd.Stdout = os.Stdout
	runCmd.Stderr = os.Stderr
	if err := runCmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to run module: %v\n", err)
		os.Exit(1)
	}
}
