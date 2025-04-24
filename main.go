package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

const (
	fileMode = 0o644
)

var (
	testlogsDir string
	outputFile  string
	postRunCmd  string
	silent      bool
	showVersion bool
	logger      *log.Logger

	GitVersion = "0.0.1"
	GitCommit  = ""
	BuildDate  = ""
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "bazel-run-reporter version %s\n\n", GitVersion)
		fmt.Fprintf(flag.CommandLine.Output(), "Usage: bazel-run-reporter [options]\n\n")
		fmt.Fprintf(flag.CommandLine.Output(), "Options:\n")
		flag.PrintDefaults()
	}

	flag.StringVar(&testlogsDir, "testlogs-dir", "bazel-testlogs", "Directory containing test.xml files")
	flag.StringVar(&outputFile, "output-file", "results.xml", "Output file for merged test results")
	flag.StringVar(&postRunCmd, "post-run", "", "Command to run after the tests results merged")
	flag.BoolVar(&showVersion, "version", false, "Show version information")
	flag.BoolVar(&silent, "silent", false, "Silent mode (suppress output)")

	flag.Parse()

	if showVersion {
		version()
		return
	}

	if silent {
		logger = log.New(io.Discard, "", 0)
	} else {
		logger = log.New(os.Stdout, "", log.LstdFlags)
	}

	var testXMLFiles []string

	walk, err := filepath.EvalSymlinks(testlogsDir)
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

	if err := os.WriteFile(outputFile, append([]byte(xml.Header), output...), fileMode); err != nil {
		logger.Fatalf("Error writing merged XML to file: %v", err)
	}

	logger.Printf("Tests written to %s", outputFile)

	if postRunCmd != "" {
		var c command

		c.Set(postRunCmd)

		logger.Printf("Running post-run command: %s", c.String())

		if err := c.Execute(); err != nil {
			logger.Printf("Error running post-run command: %v", err)
		}
	}
}

func version() {
	fmt.Printf("bazel-run-reporter version %s\n", GitVersion)

	if GitCommit != "" {
		fmt.Printf("GitCommit: %s\n", GitCommit)
	}

	if BuildDate != "" {
		fmt.Printf("BuildDate: %s\n", BuildDate)
	}
}
