package managers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestShutdownManagerWillEndWaitOnError(t *testing.T) {
	sm := NewShutdownManager()
	sm.Start(func() error {
		return assert.AnError
	}, func() error {
		return nil
	})

	err := sm.Wait()
	if assert.Error(t, err) {
		assert.EqualError(t, err, assert.AnError.Error())
	}
}

func TestShutdownManagerWillEndWaitOnCancel(t *testing.T) {
	c := make(chan bool, 1)
	s := make(chan bool, 1)

	sm := NewShutdownManager()
	sm.Start(func() error {
		s <- true
		<-c
		return nil
	}, func() error {
		c <- true
		return nil
	})

	<-s
	sm.Cancel()

	err := sm.Wait()
	assert.Nil(t, err)
}

func TestShutdownManagerWillEndWaitOnSigTerm(t *testing.T) {
	c := make(chan bool, 1)
	s := make(chan bool, 1)

	sm := NewShutdownManager()
	sm.Start(func() error {
		s <- true
		<-c
		return nil
	}, func() error {
		c <- true
		return nil
	})

	<-s
	stopChan <- os.Interrupt

	err := sm.Wait()
	assert.Nil(t, err)
}
