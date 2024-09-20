package main

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"log/slog"
	"os"
	"time"

	"github.com/alecthomas/kong"
)

var CLI struct {
	Input   string `arg:"" name:"input" help:"Input file."`
	Target  string `arg:"" name:"target" help:"Target destination in the filesystem."`
	Output  string `arg:"" name:"output" help:"Output file."`
	Type    string `short:"t" name:"type" default:"tar" help:"Type of the archive. (tar, zip)"`
	Verbose bool   `short:"v" name:"verbose" help:"Enable verbose output."`
}

func main() {
	// ctx := kong.Parse(&cli,
	kong.Parse(&CLI,
		kong.Name("golinkwrite"),
		kong.Description("Create a tar archive containing a provided file and a symlink that points to the write destination."),
		kong.UsageOnError(),
		kong.Vars{"version": "0.1.0"},
	)

	// Check for verbose output
	logLevel := slog.LevelError
	if CLI.Verbose {
		logLevel = slog.LevelDebug
	}

	// setup logger
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))

	// log parameters
	logger.Debug("command line  parameters", "cli", CLI)

	// get file mode from input file
	fileInfo, err := os.Stat(CLI.Input)
	if err != nil {
		logger.Error("failed to get file info", "error", err)
		os.Exit(1)
	}
	logger.Debug("input permissions", "perm", fileInfo.Mode().Perm())

	// read the input file
	content, err := os.ReadFile(CLI.Input)
	if err != nil {
		logger.Error("failed to read input file", "error", err)
		os.Exit(1)
	}
	logger.Debug("input size", "size", len(content))

	// prepare the tar file
	out, err := os.Create(CLI.Output)
	if err != nil {
		logger.Error("failed to create output file", "error", err)
		os.Exit(1)
	}

	switch CLI.Type {

	case "tar":
		if err := createTar(out, fileInfo, content); err != nil {
			panic(fmt.Errorf("failed to create tar file: %w", err))
		}

	case "zip":
		if err := createZip(out, fileInfo, content); err != nil {
			panic(fmt.Errorf("failed to create zip file: %w", err))
		}

	default:
		panic(fmt.Errorf("unsupported archive type: %s", CLI.Type))

	}

	logger.Info("archive created", "output", CLI.Output)
}

func createTar(out io.Writer, fi fs.FileInfo, content []byte) error {

	// create writer add a file and a directory
	writer := tar.NewWriter(out)
	defer func() {
		if err := writer.Close(); err != nil {
			panic(fmt.Errorf("failed to close tar writer: %w", err))
		}
	}()

	// add a symlink to the tar file
	header := &tar.Header{
		Name:     CLI.Input,
		Linkname: CLI.Target,
		Mode:     int64(fi.Mode().Perm()),
		Typeflag: tar.TypeSymlink,
		Size:     0,
	}
	if err := writer.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write header to tar file: %w", err)
	}

	// add a file header
	header = &tar.Header{
		Name: CLI.Input,
		Mode: int64(fi.Mode().Perm()),
		Size: int64(len(content)),
	}
	if err := writer.WriteHeader(header); err != nil {
		return fmt.Errorf("failed to write header to tar file: %w", err)
	}

	// add the file content
	if _, err := writer.Write([]byte(content)); err != nil {
		return fmt.Errorf("failed to write file to tar file: %w", err)
	}

	return nil
}

func createZip(out io.Writer, fi fs.FileInfo, content []byte) error {

	// create zip writer
	writer := zip.NewWriter(out)
	defer func() {
		if err := writer.Close(); err != nil {
			panic(fmt.Errorf("failed to close zip writer: %w", err))
		}
	}()

	// create a new file header
	zipHeader := &zip.FileHeader{
		Name:     CLI.Input,
		Method:   zip.Store,
		Modified: time.Now(),
	}
	zipHeader.SetMode(os.ModeSymlink | 0755)

	// create a new file writer
	fw, err := writer.CreateHeader(zipHeader)
	if err != nil {
		return fmt.Errorf("failed to create zip header for symlink %s: %s", CLI.Input, err)
	}

	// write the symlink to the zip archive
	if _, err := fw.Write([]byte(CLI.Target)); err != nil {
		return fmt.Errorf("failed to write symlink target %s to zip archive: %s", CLI.Target, err)
	}

	// create a new file header
	zipHeader, err = zip.FileInfoHeader(fi)
	if err != nil {
		return fmt.Errorf("failed to create file header: %s", err)
	}

	// set the name of the file
	zipHeader.Name = CLI.Input

	// set the method of compression
	zipHeader.Method = zip.Deflate

	// create a new file writer
	zw, err := writer.CreateHeader(zipHeader)
	if err != nil {
		return fmt.Errorf("failed to create zip file header: %s", err)
	}

	// write the file to the zip archive

	// create reader for byte slice
	reader := bytes.NewReader(content)

	if _, err := io.Copy(zw, reader); err != nil {
		return fmt.Errorf("failed to write file to zip archive: %s", err)
	}

	return nil
}
