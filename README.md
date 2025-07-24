# Immich Uploader (IU)

A simple CLI tool written in Go for uploading media files to an Immich server.

## Overview

`iu` (immich-uploader) is a command-line utility that scans a directory for media files (images and videos) and uploads them to an [Immich][immich-project] server.
It preserves file creation and modification timestamps during the upload process.

## Features

- Recursively scans directories for media files
- Supports common image formats (jpg, jpeg, png, gif, bmp, tiff, webp)
- Supports common video formats (mp4, mov, avi, mkv, wmv, flv, webm)
- Preserves file creation and modification dates
- Simple command-line interface

## Prerequisites

- Go 1.x or higher
- An Immich server
- API key for your Immich server

## Installation

### Building from source

```bash
# For Linux
GOOS=linux GOARCH=amd64 go build -o iu-linux iu.go

# For macOS
GOOS=darwin GOARCH=amd64 go build -o iu-mac iu.go

# For Windows
GOOS=windows GOARCH=amd64 go build -o iu.exe iu.go
```

## Configuration

Before using the tool, you need to set the API key as an environment variable:

```bash
export IMAGEUP_API_KEY=your_api_key_here
```

## Usage

```bash
./iu /path/to/your/media/directory
```

By default the program uploads files to an Immich server running locally using the Immich project defaults (http://127.0.0.1:2283/api).
If the Immich server is located elsewhere, you can specify the server URL using the `--url` flag:

```bash
iu-linux --url=http://your-immich-server-url/api /path/to/your/media/directory
```

The tool will scan the specified directory recursively for supported media files and upload them to your Immich server.

## Example

```bash
export IMAGEUP_API_KEY=your_api_key_here
./iu ~/Pictures/Vacation2023
```


[immich-project]: https://immich.app/
