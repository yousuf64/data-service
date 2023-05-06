package workrunner

import (
	"github.com/hashicorp/logutils"
	"log"
	"os"
	"sync"
	"testing"
	"time"
)

func init() {
	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: "WARN",
		Writer:   os.Stdout,
	}
	log.SetOutput(filter)
}

func TestWorkRunner_Run(t *testing.T) {
	wr := New[string, string]()

	var wg sync.WaitGroup
	iterCount := 10_000
	wg.Add(iterCount)

	for i := 0; i < iterCount; i++ {
		go func() {
			_ = <-wr.Observe("911_365", func() string {
				time.Sleep(time.Nanosecond * 1)
				return "111"
			})

			//log.Println(response)
			wg.Done()
		}()
	}

	wg.Wait()
}
