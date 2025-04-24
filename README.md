# Bazel Run Reporter | [![CI](https://github.com/ByteFork/bazel-run-reporter/actions/workflows/ci.yml/badge.svg)](https://github.com/ByteFork/bazel-run-reporter/actions/workflows/ci.yml)

A command-line tool that collects and merges test results from Bazel test runs. This tool scans the bazel-testlogs directory (or any specified directory) for JUnit/XML test reports and combines them into a single consistent report file. Later combined test results can be reported with the `-post-run` flag

## Features

- Finds and merges multiple test.xml files from Bazel test outputs
- Preserves test suite structure and details while eliminating duplicates
- Properly handles test failures, errors, and skipped tests
- Generates a single XML file compatible with CI systems and test visualization tools
- Silent mode for CI/CD pipeline integration
- Post-run command execution for seamless integration with reporting services

## Usage

```bash
$ bazel-run-reporter -h
bazel-run-reporter version 0.0.1

Usage: bazel-run-reporter [options]

Options:
  -output-file string
    	Output file for merged test results (default "results.xml")
  -post-run string
    	Command to run after the tests results merged
  -silent
    	Silent mode (suppress output)
  -testlogs-dir string
    	Directory containing test.xml files (default "bazel-testlogs")
  -version
    	Show version information
```

### Container Image

```bash
$ podman run -v $(pwd)/testdata:/testdata ghcr.io/bytefork/bazel-run-reporter -testlogs-dir /testdata -output-file /testdata/merged.x
ml
2025/04/24 22:56:12 Found 2 test.xml files.
2025/04/24 22:56:12 Tests written to /testdata/merged.xml
```

_Example_

```bash
# Run tests with Bazel
$ bazel test //...

# Use Testmo CLI as post-run command
$ export TESTMO_TOKEN=********
$ export POST_RUN="testmo automation:run:submit \
  --instance https://<your-name>.testmo.net \
  --project-id 1 \
  --name \"New Test Run\" \
  --source \"service-a\" \
  --results results.xml"

# Merge results and upload to a reporting service
$ bazel-run-reporter -post-run "$POST_RUN"
```

## Installation

Download a binary from [Releases](https://github.com/ByteFork/bazel-run-reporter/releases) or install from sources with `go install`:

```bash
$ go install github.com/ByteFork/bazel-run-reporter@latest
```

## License

This repository is [MIT](LICENSE) licensed.
