package micro

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"golang.org/x/time/rate"
)

func TestRedisRateLimiter_Wait(t *testing.T) {
	Convey("Test", t, func() {
		r := rate.NewLimiter(rate.Every(time.Second), 1)
		_ = r
	})
}
