# Gracefully Shutdown

## Introduction
A graceful shutdown in a process is when a process is turned off by the operating system (OS) is allowed to perform its tasks of safely shutting down processes and closing connections.

So it means that any cleanup task should be done before the application exits, whether it’s a server completing the ongoing requests, removing temporary files, etc.

Graceful shutdown is also one of [The Twelve-Factor App](https://12factor.net/disposability), which is Disposability.

What's got in this package.
1. GracefullyShutdown - used to waits for termination syscalls and doing clean up operations after received it.

## Using Package
```go 
// main.go

func main() {
    // Initialize master DB
    master := newDB()
    // Initialize slave DB
    slave := newDB()
	
    ...
    
    httpServer := http.NewServer(
        sentry,
        ...
    )
    httpServer.Serve(ctx)
    
    wait := gracefully_shutdown.GracefullyShutdown(ctx, time.Duration(5)*time.Second,
        map[string]gs.Operation{
            "masterdb": func(ctx context.Context) error {
                return master.Close()
            },
            "slavedb": func(ctx context.Context) error {
                return slave.Close()
            },
            "server": func(ctx context.Context) error {
                return httpServer.Shutdown(ctx)
            },
        }
    )
    <-wait	
}
```