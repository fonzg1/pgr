package pgr

import (
	"context"
	"errors"
	"io"
	"sync"
	"time"

	"github.com/mattn/go-colorable"
)

type Poller struct {
	duration time.Duration
	bars     []*Bar

	mu      sync.RWMutex
	running bool
	out     io.Writer
}

func NewPoller(d time.Duration) *Poller {
	return &Poller{duration: d, out: colorable.NewColorableStdout()}
}

func (p *Poller) Add(bars ...*Bar) *Poller {
	if len(bars) == 0 {
		return p
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.bars = append(p.bars, bars...)
	return p
}

var errCantChangeOut = errors.New("cannot change out while running")

func (p *Poller) SetOut(out io.Writer) error {
	p.mu.RLock()
	if p.running {
		return errCantChangeOut
	}
	p.mu.Lock()
	defer p.mu.Unlock()
	p.out = out
	return nil
}

var ErrCanceled = errors.New("canceled by context")

func (p *Poller) Show(ctx context.Context) error {
	p.mu.Lock()
	p.running = true
	p.mu.Unlock()
	defer func() {
		p.mu.Lock()
		p.running = false
		p.mu.Unlock()
	}()

	termSave(p.out)
	for {
		select {
		case <-ctx.Done():
			return ErrCanceled
		case <-time.NewTimer(p.duration).C:
			termRestore(p.out)
			if err := p.poll(); err == nil {
				return nil
			} else if err != errUnfinished {
				return err
			}
		}
	}
}

var errUnfinished = errors.New("not finished yet")

// poll() returns nil error if all bars are finished.
func (p *Poller) poll() (err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	for i := range p.bars {
		termClearLine(p.out)
		if err := p.bars[i].tmpl.Execute(p.out, p.bars[i]); err != nil {
			return err
		}
		if _, err := p.out.Write([]byte{byte('\n')}); err != nil {
			return err
		}
		if p.bars[i].Current() < p.bars[i].Total() {
			err = errUnfinished
		}
	}
	return err
}

func (p *Poller) SetDuration(d time.Duration) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.duration = d
}
