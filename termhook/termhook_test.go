package termhook

import (
	"bitbucket.org/futonredemption/goutil/atomics"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestAddCallback(t *testing.T) {
	assert := assert.New(t)

	manager := newSignalManagerForTest()

	signalCaught := atomics.NewAtomicBool()

	manager.addCallback(func(sig os.Signal) {
		signalCaught.Set(true)
	})
	manager.startListening()

	assert.False(signalCaught.Get(), "signalCaught should be false.")

	simulateSignalOnManager(manager)

	assert.True(signalCaught.Get(), "signalCaught should be true.")
}

func TestMultipleCallbacks(t *testing.T) {
	assert := assert.New(t)

	manager := newSignalManagerForTest()

	signalCaughtOne := atomics.NewAtomicBool()
	signalCaughtTwo := atomics.NewAtomicBool()

	manager.addCallback(func(sig os.Signal) {
		signalCaughtOne.Set(true)
	})
	manager.addCallback(func(sig os.Signal) {
		signalCaughtTwo.Set(true)
	})

	manager.startListening()

	assert.False(signalCaughtOne.Get(), "signalCaughtOne should be false.")
	assert.False(signalCaughtTwo.Get(), "signalCaughtTwo should be false.")

	simulateSignalOnManager(manager)

	assert.True(signalCaughtOne.Get(), "signalCaughtOne should be true.")
	assert.True(signalCaughtTwo.Get(), "signalCaughtTwo should be true.")
}

func TestStopListening(t *testing.T) {
	assert := assert.New(t)

	manager := newSignalManagerForTest()

	signalCaught := atomics.NewAtomicBool()

	manager.addCallback(func(sig os.Signal) {
		signalCaught.Set(true)
	})
	manager.startListening()
	manager.stopListening()

	assert.True(manager.isclosed.Get(), "manager.isclosed should be true")
	assert.False(signalCaught.Get(), "signal caught should not be set.")
}

func TestGlobalSignalCallbacks(t *testing.T) {
	assert := assert.New(t)

	signalCaught := atomics.NewAtomicBool()

	AddWithSignal(func(sig os.Signal) {
		signalCaught.Set(true)
	})

	assert.False(signalCaught.Get())

	simulateSignal()

	assert.True(signalCaught.Get())
}

func ExampleAdd() {
	signalCaught := atomics.NewAtomicBool()

	Add(func() {
		signalCaught.Set(true)
	})
	simulateSignal()

	fmt.Printf("Signal Caught: %t", signalCaught.Get())
	// Output: Signal Caught: true
}

func ExampleAddWithSignal() {
	signalCaught := atomics.NewAtomicBool()

	AddWithSignal(func(sig os.Signal) {
		signalCaught.Set(true)
	})
	simulateSignal()

	fmt.Printf("Signal Caught: %t", signalCaught.Get())
	// Output: Signal Caught: true
}

func newSignalManagerForTest() *signalManager {
	m := newSignalManager()
	m.intest.Set(true)
	return m
}

func simulateSignal() {
	simulateSignalOnManager(globalSignalManager)
}

func simulateSignalOnManager(manager *signalManager) {
	// Tell the signal manager to not kill the app.
	manager.intest.Set(true)
	log.Printf("Set Val (race condition?): %t", manager.intest.Get())
	manager.channel <- os.Interrupt
	<-manager.testchan
}
