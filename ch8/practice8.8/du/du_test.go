package du

import (
	"testing"
	"time"
)

func TestListen(t *testing.T) {
	duListen([]string{"/Users/pipasese/Documents/code", "/Users/pipasese/Documents/miCode"}, 10*time.Second)
}
