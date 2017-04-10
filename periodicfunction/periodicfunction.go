package periodicfunction

import (
	"github.com/jeremyje/goutil/atomics"
	"time"
)

type PeriodicFunction struct {
	Interval time.Duration
	Function func()
	ticker   *time.Ticker
	enabled  *atomics.AtomicBool
}

func (pf *PeriodicFunction) Start() {
	if pf.ticker == nil {
		pf.enabled = atomics.NewAtomicBool()
		pf.enabled.Set(true)
		pf.ticker = time.NewTicker(pf.Interval)
		go func() {
			for range pf.ticker.C {
				if pf.enabled.Get() {
					pf.Function()
				}
			}
		}()
	}
}

func (pf *PeriodicFunction) Stop() {
	pf.ticker.Stop()
	pf.enabled.Set(false)
	//close(pf.ticker.C)
	pf.ticker = nil
}
