# smanzy_thumbgen

`smanzy_thumbgen` is a high-performance, automated thumbnail generation service written in Go. It monitors a directory for new media files and instantly generates optimized thumbnails.

## Overview

The service acts as a background worker that watches the `./uploads` directory. When an image or video is added, it automatically generates thumbnails in predefined sizes. It also maintains directory hygiene by removing thumbnails for deleted files (Garbage Collection).

## Features

- **Real-time File Watching**: Utilizes `fsnotify` to detect file system events immediately.
- **Filesystem-based Triggers**: Control the service by "touching" special files in the upload directory.
- **Garbage Collection**: Automatically purges orphaned thumbnails on startup and via triggers.
- **Docker Ready**: Includes a lightweight Alpine-based Dockerfile with `ffmpeg`.

### Supported Formats

- **Images**: JPG, JPEG, PNG, GIF, BMP
- **Videos**: MP4, MOV, AVI, MKV, WEBM (Extracts a frame at 00:00:01)
- **HEIC**: Full support for HEIC/HEIF files (via ffmpeg)

### Output Sizes

Thumbnails are generated in the following subdirectories within `uploads`:
- `320x200/`: Small thumbnails
- `800x600/`: Medium thumbnails

## Usage

### Running with Docker (Recommended)

The service is pre-configured for Docker. It requires `ffmpeg` which is included in the provided Dockerfile.

```bash
# Build the image
docker build -t thumbgen .

# Run the container
# Mount your local uploads directory to /app/uploads
docker run -d \
  -v $(pwd)/uploads:/app/uploads \
  --name thumbgen \
  thumbgen
```

### Running Locally

**Prerequisites:**
- Go 1.23 or higher
- `ffmpeg` installed and available in your system PATH.

```bash
# Install dependencies
go mod download

# Run the service
go run main.go
```

## Configuration & Control

### Command Line Flags

| Flag | Description |
|------|-------------|
| `--regenerate` | Scans the `uploads` folder, regenerates thumbnails for ALL files, and then exits. |

Example:
```bash
go run main.go --regenerate
```

### Triggers (Runtime Control)

You can control the running service by creating empty files in the `./uploads` directory. The service detects these files, performs the action, and then deletes the trigger file.

- **Force Regeneration**:
  ```bash
  touch uploads/.trigger_regenerate
  ```
  *Effect: Re-processes every file in the uploads directory.*

- **Force Garbage Collection**:
  ```bash
  touch uploads/.trigger_gc
  ```
  *Effect: Scans for and deletes thumbnails that no longer have a corresponding source file.*

## Development

### Project Structure

- `main.go`: Application entry point and core logic.
- `Dockerfile`: Multi-stage build (Go builder -> Alpine runner with ffmpeg).
- `go.mod / go.sum`: Go module definitions.

### Dependencies

- `github.com/disintegration/imaging`: High-quality image processing.
- `github.com/fsnotify/fsnotify`: Cross-platform file system notifications.
