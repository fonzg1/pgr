// +build ignore

package main

import (
	"context"
	"math"
	"strings"
	"time"

	"github.com/tyru/pgr"
)

func main() {
	b1 := pgr.NewBarFunc(math.MaxInt64, dash())
	b2 := pgr.NewBarFunc(math.MaxInt64, dash())
	b3 := pgr.NewBarFunc(math.MaxInt64, dash())
	go incBy(b1, 30*time.Millisecond)
	go incBy(b2, 20*time.Millisecond)
	go incBy(b3, 40*time.Millisecond)

	poller := pgr.NewPoller(100 * time.Millisecond).Add(b1)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		// Add new progress bar (b2)
		time.Sleep(2 * time.Second)
		poller.Add(b2)

		// Add new progress bar (b3)
		time.Sleep(2 * time.Second)
		poller.Add(b3)

		// Speed up 10x faster
		time.Sleep(2 * time.Second)
		poller.SetDuration(10 * time.Millisecond)

		// end
		time.Sleep(2 * time.Second)
		cancel()
	}()

	poller.Show(ctx)
}

func dash() pgr.FormatFunc {
	forward := true
	dash := 0
	const maxDash = 20
	return func(*pgr.Bar) string {
		var aa string
		if forward {
			if dash == 0 {
				aa = "┏( ^o^)┛"
			} else {
				aa = strings.Repeat("　", dash-1) + "三┏( ^o^)┛"
			}
		} else {
			if dash == 0 {
				aa = strings.Repeat("　", maxDash) + "┗(^o^ )┓"
			} else {
				aa = strings.Repeat("　", maxDash-dash) + "┗(^o^ )┓三"
			}
		}
		if dash >= maxDash {
			forward = !forward
			dash = 0
		} else {
			dash++
		}
		return aa
	}
}

func incBy(p *pgr.Bar, d time.Duration) {
	for {
		p.Inc()
		time.Sleep(d)
	}
}
