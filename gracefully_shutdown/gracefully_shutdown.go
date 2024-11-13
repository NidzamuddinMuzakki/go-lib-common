package gracefully_shutdown

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

// Operation is a cleanup function on shutting down
type Operation func(ctx context.Context) error

// GracefullyShutdown waits for termination syscalls and doing clean up operations after received it
func GracefullyShutdown(ctx context.Context, timeout time.Duration, ops map[string]Operation) <-chan struct{} {
	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		// add any other syscalls that you want to be notified with
		signal.Notify(s, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
		<-s

		log.Println("shutting down")

		// set timeout for the ops to be done to prevent system hang
		timeoutFunc := time.AfterFunc(timeout, func() {
			log.Printf("timeout %d ms has been elapsed, force exit", timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		// Do the operations asynchronously to save time
		for key, op := range ops {
			wg.Add(1)
			inOp := op
			inKey := key
			go func() {
				defer wg.Done()

				log.Printf("cleaning up: %s", inKey)
				if err := inOp(ctx); err != nil {
					log.Printf("%s: clean up failed: %s", inKey, err.Error())
					return
				}

				log.Printf("%s was shutdown gracefully", inKey)
			}()
		}

		wg.Wait()

		close(wait)
	}()

	return wait
}
