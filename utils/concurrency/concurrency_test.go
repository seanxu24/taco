/**
 * @Author: Sean
 * @Date: 2022/5/19 14:55
 */

package concurrency

import (
	"context"
	"github.com/seanxu24/taco/utils/internal"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	assert := internal.NewAssert(t, "TestGenerate")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := NewChannel()
	intStream := c.Generate(ctx, 1, 2, 3)

	//for v := range intStream {
	//	t.Log(v) // 1,2,3
	//}

	assert.Equal(1, <-intStream)
	assert.Equal(2, <-intStream)
	assert.Equal(3, <-intStream)
}

func TestRepeat(t *testing.T) {
	assert := internal.NewAssert(t, "TestRepeat")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := NewChannel()
	intStream := c.Take(ctx, c.Repeat(ctx, 1, 2), 5)

	//for v := range intStream {
	//	t.Log(v) // 1,2,1,2,1
	//}

	assert.Equal(1, <-intStream)
	assert.Equal(2, <-intStream)
	assert.Equal(1, <-intStream)
	assert.Equal(2, <-intStream)
	assert.Equal(1, <-intStream)
}

func TestRepeatFn(t *testing.T) {
	assert := internal.NewAssert(t, "TestRepeatFn")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fn := func() any {
		s := 'a'
		return s
	}
	c := NewChannel()
	dataStream := c.Take(ctx, c.RepeatFn(ctx, fn), 3)

	//for v := range dataStream {
	//	t.Log(v) // a, a, a
	//}

	assert.Equal('a', <-dataStream)
	assert.Equal('a', <-dataStream)
	assert.Equal('a', <-dataStream)
}

func TestTake(t *testing.T) {
	assert := internal.NewAssert(t, "TestTake")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	numbers := make(chan any, 5)
	numbers <- 1
	numbers <- 2
	numbers <- 3
	numbers <- 4
	numbers <- 5
	defer close(numbers)

	c := NewChannel()
	intStream := c.Take(ctx, numbers, 3)

	assert.Equal(1, <-intStream)
	assert.Equal(2, <-intStream)
	assert.Equal(3, <-intStream)
}

func TestFanIn(t *testing.T) {
	assert := internal.NewAssert(t, "TestFanIn")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := NewChannel()
	channels := make([]<-chan any, 3)

	for i := 0; i < 3; i++ {
		channels[i] = c.Take(ctx, c.Repeat(ctx, i), 3)
	}

	mergedChannel := c.FanIn(ctx, channels...)

	for val := range mergedChannel {
		t.Logf("\t%d\n", val)
	}

	assert.Equal(1, 1)
}

func TestOr(t *testing.T) {
	assert := internal.NewAssert(t, "TestOr")

	sig := func(after time.Duration) <-chan any {
		c := make(chan interface{})
		go func() {
			defer close(c)
			time.Sleep(after)
		}()
		return c
	}

	start := time.Now()

	c := NewChannel()
	<-c.Or(
		sig(1*time.Second),
		sig(2*time.Second),
		sig(3*time.Second),
		sig(4*time.Second),
		sig(5*time.Second),
	)

	t.Logf("done after %v", time.Since(start))

	assert.Equal(1, 1)
}

func TestOrDone(t *testing.T) {
	assert := internal.NewAssert(t, "TestOrDone")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := NewChannel()
	intStream := c.Take(ctx, c.Repeat(ctx, 1), 3)

	var res any
	for val := range c.OrDone(ctx, intStream) {
		t.Logf("%v", val)
		res = val
	}

	assert.Equal(1, res)
}

func TestTee(t *testing.T) {
	assert := internal.NewAssert(t, "TestTee")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := NewChannel()
	intStream := c.Take(ctx, c.Repeat(ctx, 1), 4)

	out1, out2 := c.Tee(ctx, intStream)
	for val := range out1 {
		val1 := val
		val2 := <-out2
		//t.Log("val1 is",val1)
		//t.Log("val2 is",val2)
		assert.Equal(1, val1)
		assert.Equal(1, val2)
	}
}

func TestBridge(t *testing.T) {
	assert := internal.NewAssert(t, "TestBridge")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := NewChannel()
	genVals := func() <-chan <-chan any {
		chanStream := make(chan (<-chan any))
		go func() {
			defer close(chanStream)

			for i := 0; i < 10; i++ {
				stream := make(chan any, 1)
				stream <- i
				close(stream)
				chanStream <- stream
			}
		}()

		return chanStream
	}

	index := 0
	for val := range c.Bridge(ctx, genVals()) {
		t.Logf("%v ", val) //0 1 2 3 4 5 6 7 8 9
		assert.Equal(index, val)
		index++
	}
}
