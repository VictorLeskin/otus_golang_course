package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewCache(t *testing.T) {
	t0 := NewCache(3)
	t1 := t0.(*lruCache)

	assert.Equal(t, 3, t1.capacity)
	assert.Equal(t, 0, t1.queue.Len())
	assert.Equal(t, 0, len(t1.items))
}

func Test_lruCase_Set(t *testing.T) {
	//  A
	t0 := NewCache(3)
	k := t0.Set("A", 111)

	assert.False(t, k)
	t1 := t0.(*lruCache)
	assert.Equal(t, 3, t1.capacity)
	assert.Equal(t, 1, t1.queue.Len())
	assert.Equal(t, 1, len(t1.items))

	ki := t1.queue.Front().Value.(cacheItem)
	assert.Equal(t, Key("A"), ki.key)
	assert.Equal(t, 111, ki.value)

	v, ok := t1.items["A"]
	assert.True(t, ok)
	k1 := v.Value.(cacheItem)
	assert.Equal(t, Key("A"), k1.key)
	assert.Equal(t, 111, k1.value)

	//  B
	t0.Set("B", 222)

	assert.False(t, k)
	assert.Equal(t, 3, t1.capacity)
	assert.Equal(t, 2, t1.queue.Len())
	assert.Equal(t, 2, len(t1.items))

	ki = t1.queue.Front().Value.(cacheItem)
	assert.Equal(t, Key("B"), ki.key)
	assert.Equal(t, 222, ki.value)

	_, ok = t1.items["B"]
	assert.True(t, ok)

	//  A
	k = t0.Set("A", 111111)
	assert.True(t, k)
	assert.Equal(t, 3, t1.capacity)
	assert.Equal(t, 2, t1.queue.Len())
	assert.Equal(t, 2, len(t1.items))

	ki = t1.queue.Front().Value.(cacheItem)
	assert.Equal(t, Key("A"), ki.key)
	assert.Equal(t, 111111, ki.value)

	_, ok = t1.items["A"]
	assert.True(t, ok)

	//  C
	k = t0.Set("C", 333)

	assert.False(t, k)
	assert.Equal(t, 3, t1.capacity)
	assert.Equal(t, 3, t1.queue.Len())
	assert.Equal(t, 3, len(t1.items))

	ki = t1.queue.Front().Value.(cacheItem)
	assert.Equal(t, Key("C"), ki.key)
	assert.Equal(t, 333, ki.value)

	_, ok = t1.items["C"]
	assert.True(t, ok)

	//  D
	k = t0.Set("D", 444)

	assert.False(t, k)
	assert.Equal(t, 3, t1.capacity)
	assert.Equal(t, 3, t1.queue.Len())
	assert.Equal(t, 3, len(t1.items))

	ki = t1.queue.Front().Value.(cacheItem)
	assert.Equal(t, Key("D"), ki.key)
	assert.Equal(t, 444, ki.value)

	_, ok = t1.items["D"]
	assert.True(t, ok)
	_, ok = t1.items["C"]
	assert.True(t, ok)
	_, ok = t1.items["A"]
	assert.True(t, ok)
}

func Test_lruCase_Get(t *testing.T) {
	//  A
	t0 := NewCache(3)
	t1 := t0.(*lruCache)

	_ = t0.Set("A", 111)
	_ = t0.Set("B", 222)

	ki := t1.queue.Front().Value.(cacheItem)
	assert.Equal(t, Key("B"), ki.key)

	res, ok := t0.Get("A")
	assert.True(t, ok)
	assert.Equal(t, 111, res)
	ki = t1.queue.Front().Value.(cacheItem)
	assert.Equal(t, Key("A"), ki.key)

	res, ok = t0.Get("B")
	assert.True(t, ok)
	assert.Equal(t, 222, res)
	ki = t1.queue.Front().Value.(cacheItem)
	assert.Equal(t, Key("B"), ki.key)

	_, ok = t0.Get("C")
	assert.False(t, ok)
}

func Test_lruCase_Clear(t *testing.T) {
	//  A
	t0 := NewCache(3)
	t1 := t0.(*lruCache)

	assert.Equal(t, 3, t1.capacity)
	assert.Equal(t, 0, t1.queue.Len())
	assert.Equal(t, 0, len(t1.items))

	t0.Clear()
	assert.Equal(t, 0, t1.queue.Len())
	assert.Equal(t, 0, len(t1.items))

	_ = t0.Set("A", 111)
	_ = t0.Set("B", 222)
	assert.Equal(t, 2, t1.queue.Len())
	assert.Equal(t, 2, len(t1.items))
	res, ok := t0.Get("A")
	assert.True(t, ok)
	assert.Equal(t, 111, res)

	t0.Clear()
	assert.Equal(t, 0, t1.queue.Len())
	assert.Equal(t, 0, len(t1.items))
	_, ok = t0.Get("A")
	assert.False(t, ok)
}

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic", func(_ *testing.T) {
		// Write me
	})
}

func TestCacheMultithreading(t *testing.T) {
	t.Skip() // Remove me if task with asterisk completed.

	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
