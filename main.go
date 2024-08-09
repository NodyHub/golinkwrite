package main

import (
	"archive/tar"
	"log/slog"
	"os"

	"github.com/alecthomas/kong"
)

type CLI struct {
	Input   string `arg:"" name:"input" help:"Input file."`
	Target  string `arg:"" name:"target" help:"Target destination in the filesystem."`
	Output  string `arg:"" name:"output" help:"Output file."`
	Verbose bool   `short:"v" name:"verbose" help:"Enable verbose output."`
	Version kong.VersionFlag
}

func main() {
	cli := CLI{}
	// ctx := kong.Parse(&cli,
	kong.Parse(&cli,
		kong.Name("golinkwrite"),
		kong.Description("Create a tar archive containing a symbolic link to a provided target and a provided file."),
		kong.UsageOnError(),
		kong.Vars{"version": "0.1.0"},
	)

	// Check for verbose output
	logLevel := slog.LevelError
	if cli.Verbose {
		logLevel = slog.LevelDebug
	}

	// setup logger
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))

	// log parameters
	logger.Debug("command line  parameters", "cli", cli)

	// get file mode from input file
	fileInfo, err := os.Stat(cli.Input)
	if err != nil {
		logger.Error("failed to get file info", "error", err)
		os.Exit(1)
	}

	// read the input file
	content, err := os.ReadFile(cli.Input)
	if err != nil {
		logger.Error("failed to read input file", "error", err)
		os.Exit(1)
	}

	// prepare the tar file
	out, err := os.Create(cli.Output)
	if err != nil {
		logger.Error("failed to create output file", "error", err)
		os.Exit(1)
	}

	// create writer add a file and a directory
	writer := tar.NewWriter(out)
	defer writer.Close()

	// add a symlink to the tar file
	header := &tar.Header{
		Name:     cli.Input,
		Linkname: cli.Output,
		Mode:     int64(fileInfo.Mode().Perm()),
		Typeflag: tar.TypeSymlink,
		Size:     0,
	}
	err = writer.WriteHeader(header)
	if err != nil {
		logger.Error("failed to write header to tar file", "error", err)
		os.Exit(1)
	}

	// add a file
	header = &tar.Header{
		Name: cli.Input,
		Mode: int64(fileInfo.Mode().Perm()),
		Size: int64(len(content)),
	}
	err = writer.WriteHeader(header)
	if err != nil {
		logger.Error("failed to write header to tar file", "error", err)
		os.Exit(1)
	}
	_, err = writer.Write([]byte(content))
	if err != nil {
		logger.Error("failed to write to tar file", "error", err)
		os.Exit(1)
	}
}
