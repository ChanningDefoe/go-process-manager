# Go Process Manager Example

A quick example of a go process that can handle multiple interfaces and shut them down if it encounters a signal interrupt, signal terminate, or an error.

Test file [main.go file here](/example/gin/main.go).

Run the test file with:

```bash
go run example/gin/main.go
```

It does the following:
- Spins up 2 servers on localhost:8080 and localhost:8081.
- The servers has the endpoints `/ping` and `/error` for throwing a signal terminate.
- There is an ErrorClass struct that throws an error after 10 seconds.

The terminal logs will exhibit the behavior. Hit `Cmd⌘+C` on macos to throw a signal interrupt and observe the behavior. 
