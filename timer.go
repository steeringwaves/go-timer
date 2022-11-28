package timer

import (
	"context"
	"time"
)

type Timer struct {
	C chan time.Time

	ctx    context.Context
	cancel context.CancelFunc
	timer  *time.Timer
}

func NewTimerWithContext(ctx context.Context, d time.Duration) *Timer {
	t := Timer{
		C:     make(chan time.Time, 1),
		timer: time.NewTimer(d),
	}

	t.ctx, t.cancel = context.WithCancel(ctx)

	go func() {
		defer func() {
			t.Stop()
		}()

		for {
			select {
			case when := <-t.timer.C:
				select {
				case t.C <- when:
				default:
					//dont wait
				}
			case <-t.ctx.Done():
				return
			}
		}
	}()

	return &t
}

func NewTimer(d time.Duration) *Timer {
	return NewTimerWithContext(context.Background(), d)
}

func (t *Timer) Stop() bool {
	t.cancel()
	stopped := true
	if !t.timer.Stop() {
		stopped = false
		select {
		case <-t.timer.C:
		default:
			//dont wait
		}
	}

	if nil != t.C {
		select {
		case <-t.C:
		default:
			//dont wait
		}

		// close(t.C)
		// t.C = nil
	}
	return stopped
}

func (t *Timer) Reset(d time.Duration) bool {
	if !t.timer.Stop() {
		select {
		case <-t.timer.C:
		default:
			//dont wait
		}
	}

	if nil != t.C {
		select {
		case <-t.C:
		default:
			//dont wait
		}

		// close(t.C)
		// t.C = nil
	}

	// t.C = make(chan time.Time, 1)

	return t.timer.Reset(d)
}
