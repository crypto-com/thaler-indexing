package syncservice

import "sync"

type RWSerialUint64 interface {
	Get() uint64
	// Set the underlying value if the new number is larger. Return true
	// if it is updated
	SetIfLarger(uint64) bool
}

type DefaultRWSerialUint64 struct {
	sync.RWMutex

	value uint64
}

func NewDefaultRWSerialUint64(initValue uint64) *DefaultRWSerialUint64 {
	return &DefaultRWSerialUint64{
		value: initValue,
	}
}

func (rwSerial *DefaultRWSerialUint64) Get() uint64 {
	rwSerial.RLock()
	defer rwSerial.RUnlock()

	return rwSerial.value
}

func (rwSerial *DefaultRWSerialUint64) SetIfLarger(value uint64) bool {
	current := rwSerial.Get()
	if current == value {
		return false
	}

	rwSerial.Lock()
	defer rwSerial.Unlock()

	if value > rwSerial.value {
		rwSerial.value = value
	}
	return true
}

// Send a true value to the provided channel when value is updated
// It does not block when the notify channel is blocked. It is advised to have
// a buffered channel such that downstream will get notified on next read anyway
type ChRWSerialUint64 struct {
	sync.RWMutex

	onUpdatedCh chan<- bool
	value       uint64
}

func NewChRWSerialUint64(onUpdatedCh chan<- bool, initValue uint64) *ChRWSerialUint64 {
	return &ChRWSerialUint64{
		onUpdatedCh: onUpdatedCh,
		value:       initValue,
	}
}

func (rwSerial *ChRWSerialUint64) Get() uint64 {
	rwSerial.RLock()
	defer rwSerial.RUnlock()

	return rwSerial.value
}

func (rwSerial *ChRWSerialUint64) SetIfLarger(value uint64) bool {
	current := rwSerial.Get()
	if current == value {
		return false
	}

	rwSerial.Lock()

	if value > rwSerial.value {
		rwSerial.value = value
		rwSerial.Unlock()

		select {
		case rwSerial.onUpdatedCh <- true:
		default:
		}
	} else {
		rwSerial.Unlock()
	}
	return true
}
