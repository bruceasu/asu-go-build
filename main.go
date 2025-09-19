/*
  go run main.go -project <your_project_name> -output <output_directory>
*/
package main

import (
	"bytes"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	// Parse command line arguments
	// Example: go run main.go -project myapp -output bin
	projectName := flag.String("project", "", "project name")
	outputDir := flag.String("output", "bin", "output directory")
	flag.Parse()
	
	// If the project name argument is empty, try to get the module name from go.mod
	if *projectName == "" {
		// If still not available, use the current directory name
		wd, err := os.Getwd()
		if err != nil {
			fmt.Printf("Failed to get current directory: %v\n", err)
			os.Exit(1)
		}
		*projectName = filepath.Base(wd)
	}

	// Create the output directory
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		fmt.Printf("Failed to create output directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Building local executable...")

	// Get local GOOS and GOARCH from go env
	localGOOS, err := runCmdAndGetStdout("go", "env", "GOOS")
	if err != nil {
		fmt.Printf("Failed to get local GOOS: %v\n", err)
		os.Exit(1)
	}
	localGOARCH, err := runCmdAndGetStdout("go", "env", "GOARCH")
	if err != nil {
		fmt.Printf("Failed to get local GOARCH: %v\n", err)
		os.Exit(1)
	}

	localOutputName := fmt.Sprintf("%s-%s-%s", *projectName, localGOOS, localGOARCH)
	if localGOOS == "windows" {
		localOutputName += ".exe"
	}
	if err := buildWithEnv(*outputDir, localOutputName, localGOOS, localGOARCH); err != nil {
		fmt.Printf("Local build failed: %v\n", err)
		os.Exit(1)
	}

	// Specify target platforms for cross-compilation
	platforms := []string{
		"windows/amd64",
		"windows/386",
		"linux/amd64",
		// "darwin/amd64", // Uncomment if needed
	}

	fmt.Println("Starting cross-compilation...")
	for _, p := range platforms {
		parts := strings.Split(p, "/")
		if len(parts) != 2 {
			fmt.Printf("Invalid platform configuration: %s\n", p)
			os.Exit(1)
		}
		goos := parts[0]
		goarch := parts[1]

		outputName := fmt.Sprintf("%s-%s-%s", *projectName, goos, goarch)
		if goos == "windows" {
			outputName += ".exe"
		}

		if err := buildWithEnv(*outputDir, outputName, goos, goarch); err != nil {
			fmt.Printf("Cross-compilation failed for %s, error: %v\n", p, err)
			os.Exit(1)
		}
	}

	fmt.Println("Build complete.")
}

// buildWithEnv calls go build, sets GOOS/GOARCH, and outputs to the specified filename
func buildWithEnv(outputDir, outputName, goos, goarch string) error {
	cmd := exec.Command("go", "build", "-o", filepath.Join(outputDir, outputName))
	cmd.Env = append(os.Environ(),
		"GOOS="+goos,
		"GOARCH="+goarch,
	)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// runCmdAndGetStdout executes a command and returns its standard output (trimmed of trailing newlines)
func runCmdAndGetStdout(name string, args ...string) (string, error) {
	var out bytes.Buffer
	cmd := exec.Command(name, args...)
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", err
	}
	return strings.TrimSpace(out.String()), nil
}
