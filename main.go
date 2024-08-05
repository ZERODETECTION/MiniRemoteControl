package main

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/vova616/screenshot"
	"image/png"
	"net/http"
	"os"
	"os/exec"
	"time"
)

const (
	defaultScreenshotDir  = "./screenshots/"
	defaultServerAddr     = ":8080"
	defaultInterval       = 2 * time.Second
	defaultMaxScreenshots = 10
	defaultCommandEnabled = false
)

var (
	screenshotDir      string
	serverAddr         string
	interval           time.Duration
	maxScreenshots     int
	commandEnabled     bool
)

func init() {
	// Define command-line arguments
	flag.StringVar(&screenshotDir, "dir", defaultScreenshotDir, "Directory to save screenshots")
	flag.StringVar(&serverAddr, "addr", defaultServerAddr, "Address and port of the HTTP server")
	flag.DurationVar(&interval, "interval", defaultInterval, "Interval between screenshots (e.g., 2s, 5m)")
	flag.IntVar(&maxScreenshots, "max", defaultMaxScreenshots, "Maximum number of screenshots to keep before deletion")
	flag.BoolVar(&commandEnabled, "enable-command", defaultCommandEnabled, "Enable command execution endpoint")

	// Parse command-line arguments
	flag.Parse()
}

func takeScreenshot() ([]byte, error) {
	// Capture the screenshot of the entire screen
	img, err := screenshot.CaptureScreen()
	if err != nil {
		return nil, err
	}

	// Encode the image to PNG format
	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func saveScreenshot(data []byte, filename string) error {
	// Create and write the screenshot data to a file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}

func cleanOldScreenshots(dir string, max int) error {
	// List all files in the directory
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	if len(files) <= max {
		return nil
	}

	// Delete old files if the number exceeds the maximum allowed
	for i := 0; i < len(files)-max; i++ {
		err := os.Remove(dir + "/" + files[i].Name())
		if err != nil {
			return err
		}
	}

	return nil
}

func handleCommand(w http.ResponseWriter, r *http.Request) {
	// Parse query parameters
	command := r.URL.Query().Get("command")
	param := r.URL.Query().Get("param")

	if command == "" {
		http.Error(w, "Missing command parameter", http.StatusBadRequest)
		return
	}

	// Prepare the command to be executed
	cmd := exec.Command("cmd", "/c", command)

	// If there are parameters, add them to the command
	if param != "" {
		cmd.Args = append(cmd.Args, param)
	}

	// Set up buffers to capture command output and error
	var out bytes.Buffer
	cmd.Stdout = &out
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	// Execute the command
	err := cmd.Run()
	if err != nil {
		http.Error(w, fmt.Sprintf("Command execution error: %v\n%s", err, stderr.String()), http.StatusInternalServerError)
		return
	}

	// Return command output
	fmt.Fprintf(w, "Command output:\n%s", out.String())
}

func startHTTPServer() {
	// Serve the screenshots directory over HTTP
	http.Handle("/", http.FileServer(http.Dir(screenshotDir)))

	// Register the command execution handler if enabled
	if commandEnabled {
		http.HandleFunc("/sendcommand", handleCommand)
	}

	fmt.Printf("Server running at http://localhost%s\n", serverAddr)
	err := http.ListenAndServe(serverAddr, nil)
	if err != nil {
		fmt.Printf("Error starting the server: %v\n", err)
	}
}

func main() {
	// Create the directory for screenshots if it doesn't exist
	err := os.MkdirAll(screenshotDir, os.ModePerm)
	if err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}

	// Start the HTTP server in a separate goroutine
	go startHTTPServer()

	count := 0
	for {
		// Capture and save the screenshot
		screenshotData, err := takeScreenshot()
		if err != nil {
			fmt.Printf("Error capturing screenshot: %v\n", err)
			continue
		}

		filename := fmt.Sprintf("%s/screenshot_%d.png", screenshotDir, count)
		err = saveScreenshot(screenshotData, filename)
		if err != nil {
			fmt.Printf("Error saving screenshot: %v\n", err)
			continue
		}

		fmt.Printf("Screenshot saved: %s\n", filename)
		count++

		// Clean up old screenshots if the number exceeds the maximum
		err = cleanOldScreenshots(screenshotDir, maxScreenshots)
		if err != nil {
			fmt.Printf("Error cleaning old screenshots: %v\n", err)
		}

		// Wait for the specified interval before taking the next screenshot
		time.Sleep(interval)
	}
}
