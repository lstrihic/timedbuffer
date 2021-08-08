## Timed buffer
Simple implementation of timed buffer in go..

See it in action:
```go
tb := NewTimedBuffer(15, 1, func(data []interface{}) {
    // TODO do something with data
})
defer tb.Close()

for i := 0; i < 100; i++ {
    err := tb.Add(fmt.Sprintf("data %d", i))
    if err != nil {
    	return
    }   
}
```
