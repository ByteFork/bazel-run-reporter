<div align="center" markdown="1">

# Bazel Run Reporter

[![Tests](https://github.com/ByteFork/bazel-run-reporter/actions/workflows/ci.yml/badge.svg)](https://github.com/ByteFork/bazel-run-reporter/actions/workflows/ci.yml)
[![Release](https://img.shields.io/github/v/release/ByteFork/bazel-run-reporter?sort=semver)](https://github.com/ByteFork/bazel-run-reporter/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/ByteFork/bazel-run-reporter)](https://goreportcard.com/report/github.com/ByteFork/bazel-run-reporter)


A command-line tool that collects and merges test results from Bazel test runs. This tool scans the bazel-testlogs directory (or any specified directory) for JUnit/XML test reports and combines them into a single consistent report file. Later combined test results can be reported with the `-post-run` flag
<br />

</div>

## Features

- Finds and merges multiple test.xml files from Bazel test outputs
- Preserves test suite structure and details while eliminating duplicates
- Properly handles test failures, errors, and skipped tests
- Generates a single XML file compatible with CI systems and test visualization tools
- Silent mode for CI/CD pipeline integration
- Post-run command execution for seamless integration with reporting services
- Configurable post-run timeout to prevent hanging reporting commands

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
  -post-run-timeout duration
      Timeout for the post-run command (default 1m0s)
  -silent
    	Silent mode (suppress output)
  -testlogs-dir string
    	Directory containing test.xml files (default "bazel-testlogs")
  -version
    	Show version information
```

### Container Image

```bash
$ podman run -v $(pwd)/testdata:/testdata ghcr.io/bytefork/bazel-run-reporter -testlogs-dir /testdata -output-file /testdata/merged.xml
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
$ bazel-run-reporter -post-run "$POST_RUN" -post-run-timeout 2m
```

## Installation

Install the latest release with the installer:

```bash
$ curl -fsSL https://install.bytefork.io/bazel-run-reporter | sh
```

Install a specific version or directory:

```bash
$ curl -fsSL https://install.bytefork.io/bazel-run-reporter | sh -s -- --version v0.0.1 --bindir ~/.local/bin
```

You can also run the installer directly from GitHub:

```bash
$ curl -fsSL https://raw.githubusercontent.com/ByteFork/bazel-run-reporter/main/install.sh | sh
```

Or download a binary from [Releases](https://github.com/ByteFork/bazel-run-reporter/releases), or install from source with `go install`:

```bash
$ go install github.com/ByteFork/bazel-run-reporter@latest
```

## License

This repository is [MIT](LICENSE) licensed.
