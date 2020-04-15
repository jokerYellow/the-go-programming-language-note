package du

import (
	duOld "joker/goprogramlearn/ch8/8.8/du"
	"log"
	"time"
)

func duListen(roots []string, duration time.Duration) {
	for {
		select {
		case <-time.Tick(duration):
			for _, root := range roots {
				log.Print(root + ": ")
				duOld.DuFaster([]string{root})
			}
		}
	}
}
