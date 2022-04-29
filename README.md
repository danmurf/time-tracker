# Time Tracker

Track task duration easily from the command line.

## To install
```shell
go install github.com/danmurf/time-tracker@latest
```

## To start a task
```shell
tt start my-task
```

## To finish the task
```shell
tt finish my-task
```

---
### Project To-Do List
- [ ] Record when a task has fished
- [ ] Don't allow a task to start if it is already in progress
- [ ] Store all observed task names for expose for autocompletion
- [ ] Tidy up bootstrap process
- [ ] Tests for task starter
- [ ] Configurable sqlite DB location