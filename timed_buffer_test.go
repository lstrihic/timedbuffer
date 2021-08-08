package timedbuffer

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestTimedBuffer_Add(t *testing.T) {
	var mu sync.Mutex
	var count = 0
	tb := NewTimedBuffer(15, 1, func(data []interface{}) {
		mu.Lock()
		defer mu.Unlock()
		count += len(data)
	})

	for i := 0; i < 100; i++ {
		err := tb.Add(fmt.Sprintf("data %d", i))
		assert.NoError(t, err)
	}

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, 90, count)
}

func TestTimedBuffer_AutoFlush(t *testing.T) {
	var mu sync.Mutex
	var count = 0
	tb := NewTimedBuffer(15, 0, func(data []interface{}) {
		mu.Lock()
		defer mu.Unlock()
		count += len(data)
	})

	for i := 0; i < 100; i++ {
		err := tb.Add(fmt.Sprintf("data %d", i))
		assert.NoError(t, err)
	}

	time.Sleep(200 * time.Millisecond)
	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, 100, count)
}

func TestTimedBuffer_Flush(t *testing.T) {
	var mu sync.Mutex
	var count = 0
	tb := NewTimedBuffer(500, 20, func(data []interface{}) {
		mu.Lock()
		defer mu.Unlock()
		count += len(data)
	})

	for i := 0; i < 20; i++ {
		err := tb.Add(fmt.Sprintf("data %d", i))
		assert.NoError(t, err)
	}

	err := tb.Flush()
	assert.NoError(t, err)

	mu.Lock()
	defer mu.Unlock()
	assert.Equal(t, 20, count)
}

func TestTimedBuffer_Close(t *testing.T) {
	tb := NewTimedBuffer(15, 0, func(data []interface{}) {

	})
	err := tb.Close()
	assert.NoError(t, err)

	err1 := tb.Add("msg")
	assert.Error(t, err1, "the buffer is already closed")

	err2 := tb.Close()
	assert.Error(t, err2, "the buffer is already closed")

	err3 := tb.Flush()
	assert.Error(t, err3, "the buffer is already closed")
}
