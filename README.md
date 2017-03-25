# go_ms

go_ms is a CLI-based Minesweeper and solver written in Go.

## Prerequisites

> Suggest using built-in means like `brew` or `apt-get` to avoid human error.

* Go - https://golang.org/doc/install#install
* (optional) Docker - https://docs.docker.com/engine/getstarted/

## Running go_ms

Run `go build` in the repo directory to get an executable to run.

or

Simply run `go run *.go` in the repo directory to avoid building an executable.

## Solver

This will house the strategy/logic needed to solve Minesweeper automatically. The current thought process is the following:

1. Identify all safe moves and mines based on immediate neighbors
2. Identify all safe moves based on known mines
3. Identify all safe moves by using neighbors' info
4. Identify least risky move
5. Blind click if no other options are found

The pass rate is a little above 90% (100,000 iterations) currently.
