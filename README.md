[![go](https://github.com/danmurf/time-tracker/actions/workflows/go.yml/badge.svg?branch=master)](https://github.com/danmurf/time-tracker/actions/workflows/go.yml)
[![codeql-analysis](https://github.com/danmurf/time-tracker/actions/workflows/codeql-analysis.yml/badge.svg?branch=master)](https://github.com/danmurf/time-tracker/actions/workflows/codeql-analysis.yml)

# Time Tracker

Track task duration easily from the command line.

## To install
```shell
go install github.com/danmurf/time-tracker@latest
```

## To start a task
```shell
time-tracker start my-task
```

## To finish the task
```shell
time-tracker finish my-task
```

## To get the duration of the last completed task
```shell
time-tracker lastDuration my-task
```