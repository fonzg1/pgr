// +build ignore

package main

import (
	"bytes"
	"context"
	"math"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/tyru/pgr"
)

func main() {
	b1 := pgr.NewBarFunc(math.MaxInt64, dash())
	b2 := pgr.NewBarFunc(math.MaxInt64, uwaaaa())
	b3 := pgr.NewBarFunc(math.MaxInt64, pugya())
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

func uwaaaa() pgr.FormatFunc {
	start := 0

	// constants
	cyan := color.New(color.FgCyan).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	b1 := cyan("▂")
	b2 := magenta("▅")
	b3 := blue("▇")
	b4 := yellow("▇")
	b5 := green("▓")
	b6 := red("▒")
	b7 := "░"
	wave := []string{b1, b2, b3, b4, b5, b6, b7, b6, b5, b4, b3, b2}
	const face = " ('ω')"

	return func(*pgr.Bar) string {
		var buf bytes.Buffer

		// wave 1
		j := start
		for i := 0; i < len(wave); i++ {
			buf.WriteString(wave[j])
			j = (j + 1) % len(wave)
		}

		buf.WriteString(face)

		// wave 2: reverse of wave 1
		j = start
		for i := 0; i < len(wave); i++ {
			buf.WriteString(wave[len(wave)-1-j])
			j = (j + 1) % len(wave)
		}

		start = (start + 1) % len(wave)
		return buf.String()
	}
}

func pugya() pgr.FormatFunc {
	i := 0

	// constants
	cyan := color.New(color.FgCyan).SprintFunc()
	magenta := color.New(color.FgMagenta).SprintFunc()
	blue := color.New(color.FgBlue).SprintFunc()
	yellow := color.New(color.FgYellow).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	red := color.New(color.FgRed).SprintFunc()
	rainbow := []func(...interface{}) string{cyan, magenta, blue, yellow, green, red}
	convert := func(s string, i int) string {
		var buf bytes.Buffer
		for j, c := range s {
			buf.WriteString(rainbow[(i+j)%len(rainbow)](string(c)))
		}
		return buf.String()
	}

	return func(*pgr.Bar) string {
		i++
		switch i % 4 {
		case 0:
			return convert("ﾌﾟｷﾞｬｰｰｰｰｰｰｰm9(^Д^)9mｰｰｰｰｰｰｰｰ", i)
		case 1:
			return convert("ﾌﾟｷﾞｬｰｰｰｰｰｰm9(^Д^)9mｰｰｰｰｰｰｰｰｰ", i)
		case 2:
			return convert("ﾌﾟｷﾞｬｰｰｰｰｰｰｰm9(^Д^)9mｰｰｰｰｰｰｰｰ", i)
		default:
			return convert("ﾌﾟｷﾞｬｰｰｰｰｰｰｰｰm9(^Д^)9mｰｰｰｰｰｰｰ", i)
		}
	}
}

func incBy(p *pgr.Bar, d time.Duration) {
	for {
		p.Inc()
		time.Sleep(d)
	}
}
