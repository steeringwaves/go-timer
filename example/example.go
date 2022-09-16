package main

import (
	"context"
	"log"
	"time"

	"github.com/steeringwaves/go-timer"
)

type TimerTest struct {
	Messages chan int
}

func (client *TimerTest) modTimerWithContextLoop() {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		waitTimeout := 1 * time.Second
		timer := timer.NewTimerWithContext(ctx, waitTimeout)
		count := 0

		defer func() {
			timer.Stop()
		}()

		for {
			timeout := (123456789*count*count)%5000 + 10 // Semi-random yet predictable timeouts for both goroutines
			// log.Printf("modTimerWithContextLoop sleep for %d ms\n", timeout)
			time.Sleep(time.Duration(timeout * int(time.Millisecond)))

			timer.Reset(waitTimeout)

			select {
			case client.Messages <- count:
				// log.Printf("modTimerWithContextLoop message: %d\n", count)
				count++
				continue
			case <-timer.C:
				// log.Printf("modTimerWithContextLoop: Timed out\n")
				count++
				continue
			case <-ctx.Done():
				// log.Printf("modTimerWithContextLoop: context cancelled\n")
				return
			}
		}
	}()
}

func (client *TimerTest) modTimerLoop() {
	go func() {
		waitTimeout := 1 * time.Second
		timer := timer.NewTimer(waitTimeout)
		count := 0

		defer func() {
			timer.Stop()
		}()

		for {
			timeout := (123456789*count*count)%5000 + 10 // Semi-random yet predictable timeouts for both goroutines
			// log.Printf("modTimerLoop sleep for %d ms\n", timeout)
			time.Sleep(time.Duration(timeout * int(time.Millisecond)))

			timer.Reset(waitTimeout)

			select {
			case client.Messages <- count:
				// log.Printf("modTimerLoop message: %d\n", count)
				count++
				continue
			case <-timer.C:
				// log.Printf("modTimerLoop: Timed out\n")
				count++
				continue
			}
		}
	}()
}

func (client *TimerTest) goodTimerLoop() {
	go func() {
		waitTimeout := 1 * time.Second
		timer := time.NewTimer(waitTimeout)
		count := 0

		defer func() {
			if nil != timer {
				timer.Stop()
				timer = nil
			}
		}()

		for {
			timeout := (123456789*count*count)%5000 + 10 // Semi-random yet predictable timeouts for both goroutines
			// log.Printf("goodTimerLoop sleep for %d ms\n", timeout)
			time.Sleep(time.Duration(timeout * int(time.Millisecond)))

			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
					//nowait
				}
			}
			timer.Reset(waitTimeout)

			select {
			case client.Messages <- count:
				// log.Printf("goodTimerLoop message: %d\n", count)
				count++
				continue
			case <-timer.C:
				// log.Printf("goodTimerLoop: Timed out\n")
				count++
				continue
			}
		}
	}()
}

func (client *TimerTest) badTimerLoop() {
	go func() {
		waitTimeout := 1 * time.Second
		timer := time.NewTimer(waitTimeout)
		count := 0

		defer func() {
			if nil != timer {
				timer.Stop()
				timer = nil
			}
		}()

		for {
			timeout := (123456789*count*count)%5000 + 10 // Semi-random yet predictable timeouts for both goroutines
			// log.Printf("badTimerLoop sleep for %d ms\n", timeout)
			time.Sleep(time.Duration(timeout * int(time.Millisecond)))

			timer.Reset(waitTimeout)
			select {
			case client.Messages <- count:
				// log.Printf("badTimerLoop message: %d\n", count)
				count++
				continue
			case <-timer.C:
				// log.Printf("badTimerLoop: Timed out\n")
				count++
				continue
			}
		}
	}()
}

func New() *TimerTest {
	var client TimerTest
	client.Messages = make(chan int)

	return &client
}

func main() {
	cModWithContext := New()
	cMod := New()
	cGood := New()
	cBad := New()

	cModWithContext.modTimerWithContextLoop()
	cMod.modTimerLoop()
	cGood.goodTimerLoop()
	cBad.badTimerLoop()

	ctx := context.Background()

	go func() {
		waitTimeout := 5 * time.Second
		timer := timer.NewTimer(waitTimeout)

		for {
			timer.Reset(waitTimeout)
			select {
			case msg := <-cMod.Messages:
				log.Printf("modTimerWithContext response: %d\n", msg)
				continue
			case <-timer.C:
				// log.Printf("ModTimer: Timed out\n")
				continue
			}
		}
	}()

	go func() {
		waitTimeout := 5 * time.Second
		timer := timer.NewTimer(waitTimeout)

		for {
			timer.Reset(waitTimeout)
			select {
			case msg := <-cMod.Messages:
				log.Printf("ModTimer response: %d\n", msg)
				continue
			case <-timer.C:
				// log.Printf("ModTimer: Timed out\n")
				continue
			}
		}
	}()

	go func() {
		waitTimeout := 5 * time.Second
		timer := time.NewTimer(waitTimeout)

		for {
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
					//nowait
				}
			}
			timer.Reset(waitTimeout)
			select {
			case msg := <-cBad.Messages:
				log.Printf("BadTimer response: %d\n", msg)
				continue
			case <-timer.C:
				// log.Printf("BadTimer: Timed out\n")
				continue
			}
		}
	}()

	go func() {
		for {
			waitTimeout := 5 * time.Second
			timer := time.NewTimer(waitTimeout)

			for {
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
						//nowait
					}
				}
				timer.Reset(waitTimeout)

				select {
				case msg := <-cGood.Messages:
					log.Printf("GoodTimer response: %d\n", msg)
					continue
				case <-timer.C:
					// log.Printf("Goodimer: Timed out\n")
					continue
				}
			}
		}
	}()

	<-ctx.Done()
}
