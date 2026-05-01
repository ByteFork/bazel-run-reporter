# Contributing

Thanks for helping improve Bazel Run Reporter.

## Before you start

Open an issue describing the change if it is non-trivial. Small fixes (typos, doc tweaks, obvious bug fixes) do not need one.

## Development

Requires Go 1.26+.

```bash
git clone https://github.com/ByteFork/bazel-run-reporter.git
cd bazel-run-reporter

go build -o bazel-run-reporter .
go test -race ./...
golangci-lint run ./...
```

## Pull request workflow

1. Fork and create a branch: `git checkout -b type/short-description` (e.g. `fix/post-run-timeout`, `feat/junit-merge-option`).
2. Keep commits focused. [Conventional Commits](https://www.conventionalcommits.org/) are preferred but not required.
3. Run `go test -race ./...` and `golangci-lint run ./...` locally. Both must pass.
4. Open the PR against `main`.
5. CI must be green and the PR needs one approval before merge.
6. PRs are squash-merged into `main`.

## Coding standards

- `gofmt` and `goimports` clean.
- `golangci-lint` must pass with 0 issues (see `.golangci.yml`).
- Tests accompany new functionality and bug fixes when behavior changes.
- No `TODO`/`FIXME` without a linked issue number.

## AI usage

AI assistance is welcome. Disclose it by including a `Co-Authored-By:` trailer in the commit message, or by mentioning the tool in the PR description. You are responsible for reviewing and understanding everything you submit.

## Reporting bugs and feature requests

Use GitHub Issues for bugs and feature requests. For security concerns, please email the maintainers directly instead of opening a public issue.

## License

Contributions are licensed under the [MIT License](LICENSE). By opening a PR you agree your contributions will be distributed under these terms.
