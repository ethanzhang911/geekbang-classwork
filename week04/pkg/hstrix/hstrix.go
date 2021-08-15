package hstrix

import (
	"container/ring"
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Hstrix interface {
	Add() error
	Run()
}

// hstrixByEthan 结构体
// limitCount是最大限流访问请求数
// limitBucket是滑动窗口的个数
// curCount是当前总的请求数
type hstrixByEthan struct {
	limitCount  int
	limitBucket int
	interval    int
	head        *ring.Ring
	curCount    int32
	context     context.Context
}

// NewHstrixByEthan  构造器，返回hstrix接口类型
func NewHstrixByEthan(c context.Context, limitCount int, limitBucket int, interval int) Hstrix {
	return &hstrixByEthan{
		limitCount:  limitCount,
		limitBucket: limitBucket,
		interval:    interval,
		curCount:    0,
		context:     c,
	}
}

// Add 调用Add方法，将当前累计请求数增加1
func (h *hstrixByEthan) Add() error {
	// 我这里采取先加1，然后再判断，如果超过限流值，就返回-1
	n := atomic.AddInt32(&h.curCount, 1)
	if n > int32(h.limitCount) {
		atomic.AddInt32(&h.curCount, -1)
		return errors.New("the request has exceeded")
	} else {
		mut := sync.Mutex{}
		mut.Lock()
		pos := h.head.Prev()
		val := pos.Value.(int)
		val++
		pos.Value = val
		mut.Unlock()
		return nil
	}
}

func (h *hstrixByEthan) Run() {
	h.head = ring.New(h.limitBucket)
	// 初始化环形队列
	for i := 0; i < h.limitBucket; i++ {
		h.head.Value = 0
		h.head = h.head.Next()
	}
	go func() {
		timer := time.NewTicker(time.Second * 1)
		defer timer.Stop()
		// 按照interval属性进行轮询计算出当前滑动窗口中的总请求数
		//select {
		//case <-timer.C:
		//	{
		for range timer.C {
			// 这部是关键，讲当前指针内容从滑动窗口剔除，用于反应真实当前总的请求数
			// atomic.AddInt32要求参数是int32，所以这里直接先转换
			subCount := int32(0 - h.head.Value.(int))
			atomic.AddInt32(&h.curCount, subCount)
			h.head.Value = 0
			h.head = h.head.Next()
			fmt.Printf("当前请求数为%d，当前限流数为%d\n", atomic.LoadInt32(&h.curCount),h.limitCount)
		}
		//case <-h.context.Done():
		//	{
		//		fmt.Printf("%s\n", "滑动计时器结束")
		//		return
		//	}
		//}
	}()
}
