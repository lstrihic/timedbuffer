package timedbuffer

import (
	"errors"
	"time"
)

//TimedBuffer represents timed buffer.
type TimedBuffer interface {
	//Add add item to buffer
	//receives:
	//- data: the item that will be inserted into buffer.
	//returns error, if any.
	Add(data interface{}) error
	//Flush flush data from buffer, if any.
	//returns error, if any.
	Flush() error
	//Close close the buffer.
	//returns error, if any.
	Close() error
}

type timedBuffer struct {
	add           chan interface{}
	flush         chan struct{}
	data          []interface{}
	cb            func([]interface{})
	done          chan struct{}
	index         uint
	size          uint
	flushInterval uint
}

//NewTimedBuffer crete new instance of timed buffer.
//receives:
//- size: the size of the buffer
//- flushInterval: the time in seconds after which will buffer be flushed
//- cb: callback function
//returns new instance of TimedBuffer.
func NewTimedBuffer(size uint, flushInterval uint, cb func([]interface{})) TimedBuffer {
	tb := &timedBuffer{
		add:           make(chan interface{}),
		flush:         make(chan struct{}),
		done:          make(chan struct{}),
		data:          make([]interface{}, size),
		index:         0,
		size:          size,
		flushInterval: flushInterval,
		cb:            cb,
	}
	go tb.startBuffer()
	return tb
}

//Add add data to buffer.
func (tb *timedBuffer) Add(data interface{}) error {
	select {
	case <-tb.done:
		return errors.New("the buffer is already closed")
	default:
		tb.add <- data
	}
	return nil
}

//Flush flush data from buffer if any.
func (tb *timedBuffer) Flush() error {
	select {
	case <-tb.done:
		return errors.New("the buffer is already closed")
	default:
		tb.flush <- struct{}{}
	}
	return nil
}

//Close close the buffer.
func (tb *timedBuffer) Close() error {
	select {
	case <-tb.done:
		return errors.New("the buffer is already closed")
	default:
		close(tb.done)
	}
	return nil
}

func (tb *timedBuffer) flushData() {
	if tb.index != 0 {
		tb.cb(tb.data[:tb.index])
		tb.index = 0
	}
}

func (tb *timedBuffer) startBuffer() {
	interval := time.Duration(tb.flushInterval) * time.Second
	timer := time.NewTimer(interval)
	defer timer.Stop()
	defer close(tb.flush)
	defer close(tb.add)
	for {
		select {
		case data := <-tb.add:
			if tb.index != 0 && tb.index%tb.size == 0 {
				tb.flushData()
				timer.Reset(interval)
			}
			tb.data[tb.index] = data
			tb.index++
		case <-tb.flush:
			tb.flushData()
			timer.Reset(interval)
		case <-timer.C:
			tb.flushData()
			timer.Reset(interval)
		case <-tb.done:
			return
		}
	}
}
