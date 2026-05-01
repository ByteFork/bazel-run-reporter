package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	fileMode = 0o644
)

var (
	logger *log.Logger

	GitVersion = "0.0.1"
	GitCommit  = ""
	BuildDate  = ""
)

type Config struct {
	TestlogsDir    string
	OutputFile     string
	PostRunCmd     string
	PostRunTimeout time.Duration
	ShowVersion    bool
	Silent         bool
}

func main() {
	cfg := parseFlags()

	if cfg.ShowVersion {
		printVersion()
		return
	}

	if cfg.Silent {
		logger = log.New(io.Discard, "", 0)
	} else {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	var testXMLFiles []string

	walk, err := filepath.EvalSymlinks(cfg.TestlogsDir)
	if err != nil {
		logger.Fatalf("Error evaluating symlinks: %v", err)
	}

	err = filepath.Walk(walk, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		if filepath.Base(path) == "test.xml" {
			testXMLFiles = append(testXMLFiles, path)
		}

		return nil
	})
	if err != nil {
		logger.Fatalf("Error walking through bazel-testlogs: %v", err)
	}

	if len(testXMLFiles) == 0 {
		logger.Println("No test.xml files found.")

		return
	}

	logger.Printf("Found %d test.xml files.", len(testXMLFiles))

	suites := make([]TestSuites, 0, len(testXMLFiles))

	for _, file := range testXMLFiles {
		data, readErr := os.ReadFile(file)
		if readErr != nil {
			logger.Printf("Error reading test.xml file %s: %v", file, readErr)
			continue
		}

		parsed, parseErr := Parse(data)
		if parseErr != nil {
			logger.Printf("Error parsing test.xml file %s: %v", file, parseErr)
			continue
		}

		suites = append(suites, parsed)
	}

	m := MergeTestSuites(suites...)

	output, err := xml.MarshalIndent(m, "", "  ")
	if err != nil {
		logger.Fatalf("Error marshaling merged XML: %v", err)
	}

	if err := os.WriteFile(cfg.OutputFile, append([]byte(xml.Header), output...), fileMode); err != nil {
		logger.Fatalf("Error writing merged XML to file: %v", err)
	}

	logger.Printf("Tests written to %s", cfg.OutputFile)

	if cfg.PostRunCmd != "" {
		var c CommandHook

		if err := c.Set(cfg.PostRunCmd); err != nil {
			logger.Printf("Error setting post-run command: %v", err)
			return
		}

		logger.Printf("Running post-run command: %s", c.String())

		if err := c.Execute(cfg.PostRunTimeout); err != nil {
			logger.Printf("Error running post-run command: %v", err)
		}
	}
}

func parseFlags() *Config {
	var cfg Config

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "bazel-run-reporter version %s\n\n", GitVersion)
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: bazel-run-reporter [options]\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&cfg.TestlogsDir, "testlogs-dir", "bazel-testlogs", "Directory containing test.xml files")
	flag.StringVar(&cfg.OutputFile, "output-file", "results.xml", "Output file for merged test results")
	flag.StringVar(&cfg.PostRunCmd, "post-run", "", "Command to run after the tests results merged")
	flag.DurationVar(&cfg.PostRunTimeout, "post-run-timeout", time.Minute, "Timeout for the post-run command")
	flag.BoolVar(&cfg.ShowVersion, "version", false, "Show version information")
	flag.BoolVar(&cfg.Silent, "silent", false, "Silent mode (suppress output)")

	flag.Parse()

	return &cfg
}

func printVersion() {
	fmt.Printf("bazel-run-reporter version %s\n", GitVersion)

	if GitCommit != "" {
		fmt.Printf("GitCommit: %s\n", GitCommit)
	}

	if BuildDate != "" {
		fmt.Printf("BuildDate: %s\n", BuildDate)
	}
}
