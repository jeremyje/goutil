package termhook

import (
	"github.com/jeremyje/goutil/atomics"
	"os"
)

type signalManager struct {
	intest       *atomics.AtomicBool
	isclosed     *atomics.AtomicBool
	channel      chan os.Signal
	callbackList []SignalCallback
	testchan     chan bool
}

var globalSignalManager *signalManager

func newSignalManager() *signalManager {
	manager := &signalManager{
		intest:       atomics.NewAtomicBool(),
		channel:      make(chan os.Signal, 1),
		isclosed:     atomics.NewAtomicBool(),
		callbackList: []SignalCallback{},
		testchan:     make(chan bool),
	}

	addTerminatingSignals(manager.channel)
	return manager
}

func (sm *signalManager) startListening() {
	go func() {
		if !sm.isclosed.Get() {
			for sig := range sm.channel {
				for _, callback := range sm.callbackList {
					callback(sig)
				}
				if !sm.intest.Get() {
					os.Exit(0xf)
				} else {
					sm.testchan <- true
				}
			}
		}
	}()
}

func (sm *signalManager) stopListening() {
	sm.isclosed.Set(true)
	close(sm.channel)
}

func (sm *signalManager) addCallback(callback SignalCallback) {
	sm.callbackList = append(sm.callbackList, callback)
}

func init() {
	globalSignalManager = newSignalManager()
	globalSignalManager.startListening()
}
