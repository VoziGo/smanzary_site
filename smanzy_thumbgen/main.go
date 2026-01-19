package main

import (
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/fsnotify/fsnotify"
)

// Config
var targetSizes = []struct {
	dirName string
	width   int
	height  int
}{
	{"320x200", 320, 200},
	{"800x600", 800, 600},
}

const uploadDir = "./uploads"

var forceRegen bool

func main() {
	flag.BoolVar(&forceRegen, "regenerate", false, "Scans folder and regenerates all thumbnails, then exits.")
	flag.Parse()

	setupDirectories()

	// 1. Run Garbage Collector (Clean up orphans on restart)
	fmt.Println(">>> Running Garbage Collector...")
	runGarbageCollector()

	if forceRegen {
		fmt.Println(">>> Mode: REGENERATE (Processing all existing files)")
		performRegeneration()
		fmt.Println(">>> Regeneration complete.")
		return
	}

	fmt.Println(">>> Mode: WATCHER (Waiting for changes in ./uploads)")
	startWatcher()
}

// -------------------------
// Core Logic
// -------------------------

func processFile(path string) {
	// Debounce: Wait for file write/lock to release
	time.Sleep(500 * time.Millisecond)

	fileName := filepath.Base(path)
	if strings.HasPrefix(fileName, ".") {
		return
	}

	if isImage(fileName) {
		fmt.Printf("[Create/Update] Image: %s\n", fileName)
		processImage(path, fileName)
	} else if isVideo(fileName) {
		fmt.Printf("[Create/Update] Video: %s\n", fileName)
		processVideo(path, fileName)
	}
}

func deleteThumbnailsFor(originalName string) {
	fmt.Printf("[Delete] Cleaning thumbs for: %s\n", originalName)

	baseName := strings.TrimSuffix(originalName, filepath.Ext(originalName))
	thumbName := baseName + ".jpg"

	for _, size := range targetSizes {
		thumbPath := filepath.Join(uploadDir, size.dirName, thumbName)

		// Check if it exists before trying to delete (avoid error logs)
		if _, err := os.Stat(thumbPath); err == nil {
			if err := os.Remove(thumbPath); err != nil {
				log.Printf("   Error deleting thumb %s: %v", thumbPath, err)
			} else {
				fmt.Printf("   -> Deleted %s\n", thumbPath)
			}
		}
	}
}

// -------------------------
// Processing Functions
// -------------------------

func processImage(path, originalName string) {
	img, err := imaging.Open(path)
	if err != nil {
		log.Printf("   Error opening image: %v", err)
		return
	}
	saveThumbnails(img, originalName)
}

func processVideo(path, originalName string) {
	cmd := exec.Command("ffmpeg", "-ss", "00:00:01", "-i", path, "-vframes", "1", "-f", "image2", "-pipe:1")
	var buffer bytes.Buffer
	cmd.Stdout = &buffer
	cmd.Stderr = nil // Suppress noisy ffmpeg logs

	if err := cmd.Run(); err != nil {
		// Fallback to frame 0
		cmd = exec.Command("ffmpeg", "-i", path, "-vframes", "1", "-f", "image2", "-pipe:1")
		cmd.Stdout = &buffer
		if err := cmd.Run(); err != nil {
			log.Printf("   Error extracting video frame: %v", err)
			return
		}
	}

	img, _, err := image.Decode(&buffer)
	if err != nil {
		log.Printf("   Error decoding video buffer: %v", err)
		return
	}
	saveThumbnails(img, originalName)
}

func saveThumbnails(img image.Image, originalName string) {
	baseName := strings.TrimSuffix(originalName, filepath.Ext(originalName))
	outputName := baseName + ".jpg"

	for _, size := range targetSizes {
		resizedImg := imaging.Fit(img, size.width, size.height, imaging.Lanczos)
		outputPath := filepath.Join(uploadDir, size.dirName, outputName)

		out, err := os.Create(outputPath)
		if err != nil {
			log.Printf("   Error creating output file: %v", err)
			continue
		}

		// Encode
		if err := jpeg.Encode(out, resizedImg, &jpeg.Options{Quality: 80}); err != nil {
			log.Printf("   Error encoding jpeg: %v", err)
		}
		out.Close()
	}
}

// -------------------------
// Watcher
// -------------------------

func startWatcher() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				// Ignore events inside the thumbnail folders to prevent loops
				if strings.Contains(event.Name, "320x200") || strings.Contains(event.Name, "800x600") {
					continue
				}

				// Handle CREATION / UPDATE
				if event.Op&fsnotify.Create == fsnotify.Create || event.Op&fsnotify.Write == fsnotify.Write {
					processFile(event.Name)
				}

				// Handle DELETION / RENAME (Moving a file out is treated as a rename)
				if event.Op&fsnotify.Remove == fsnotify.Remove || event.Op&fsnotify.Rename == fsnotify.Rename {
					deleteThumbnailsFor(filepath.Base(event.Name))
				}

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("Watcher error:", err)
			}
		}
	}()

	if err := watcher.Add(uploadDir); err != nil {
		log.Fatal(err)
	}
	<-done
}

// -------------------------
// Garbage Collector & Utils
// -------------------------

func runGarbageCollector() {
	// 1. Index all currently valid "base names" in the upload folder
	validBasenames := make(map[string]bool)

	files, err := os.ReadDir(uploadDir)
	if err != nil {
		log.Printf("GC Error reading uploads: %v", err)
		return
	}

	for _, f := range files {
		if f.IsDir() || strings.HasPrefix(f.Name(), ".") {
			continue
		}
		if isImage(f.Name()) || isVideo(f.Name()) {
			// Store "myvideo" from "myvideo.mp4"
			base := strings.TrimSuffix(f.Name(), filepath.Ext(f.Name()))
			validBasenames[base] = true
		}
	}

	// 2. Scan thumbnail folders and delete orphans
	for _, size := range targetSizes {
		thumbDir := filepath.Join(uploadDir, size.dirName)
		thumbs, err := os.ReadDir(thumbDir)
		if err != nil {
			continue
		}

		for _, t := range thumbs {
			if t.IsDir() {
				continue
			}

			// Thumbnails are always .jpg. Get the base name.
			// e.g. "myvideo.jpg" -> "myvideo"
			thumbBase := strings.TrimSuffix(t.Name(), filepath.Ext(t.Name()))

			// If "myvideo" is not in our valid list, delete the thumbnail
			if !validBasenames[thumbBase] {
				fullPath := filepath.Join(thumbDir, t.Name())
				fmt.Printf("[GC] Deleting orphan: %s\n", fullPath)
				os.Remove(fullPath)
			}
		}
	}
}

func performRegeneration() {
	entries, err := os.ReadDir(uploadDir)
	if err != nil {
		log.Fatalf("Failed to read dir: %v", err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			processFile(filepath.Join(uploadDir, entry.Name()))
		}
	}
}

func setupDirectories() {
	for _, size := range targetSizes {
		_ = os.MkdirAll(filepath.Join(uploadDir, size.dirName), 0755)
	}
}

func isImage(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".bmp"
}

func isVideo(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	return ext == ".mp4" || ext == ".mov" || ext == ".avi" || ext == ".mkv" || ext == ".webm"
}
