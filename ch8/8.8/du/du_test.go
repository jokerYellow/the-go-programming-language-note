package du

import (
	"testing"
)

var dir = "/Users/pipasese/Documents/code"

func TestDu(t *testing.T) {
	du([]string{dir})
}

func TestDuProgress(t *testing.T) {
	duShowProgress([]string{dir})
}

func TestDuFaster(t *testing.T) {
	duShowProgressFaster([]string{dir})
}
