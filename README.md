SIS - URL Analyzer (Checkpoint 1)

A small synchronous in-memory Go application. A task is a unit of work: analyze one URL and produce a result.

How to run:
go run .
go run . https://aruskin.com/
go test ./...

Input: a URL string.
Validation: the URL must not be empty, must be parsable, must use http or https, must have a host, and the port (if given) must be between 1 and 65535. Invalid input is rejected with ErrInvalidInput.
Processing: the URL is parsed and the result contains scheme, host, port, secure flag, path, path depth, number of query parameters and fragment flag.
Result: a Result struct stored in the task.

Task lifecycle:
pending -> processing -> completed
processing -> failed

Allowed transitions:
pending to processing
processing to completed
processing to failed

completed and failed are final states. Any other transition, for example processing a completed task again, returns ErrInvalidTransition and leaves the task unchanged.

Errors: ErrInvalidInput, ErrNotFound, ErrInvalidTransition. They are returned as errors and checked with errors.Is.

Files:
task.go - Task, Status, Result, errors, transitions
analyzer.go - URL validation and analysis
manager.go - Analyzer interface, in-memory storage, Create, Get, Process
main.go - demo
*_test.go - unit tests

Notes: Manager is not safe for concurrent use yet. Concurrency is planned for Checkpoint 2.
