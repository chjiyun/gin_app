package asyncUtil

import (
	"context"
	"fmt"
	"sync"
)

type Group struct {
	wg      sync.WaitGroup
	ch      chan token
	cancel  func()
	errOnce sync.Once
	err     error
}
type token struct {
}

func (g *Group) done() {
	// 释放一个值，通道解除阻塞，新的任务进来
	if g.ch != nil {
		<-g.ch
	}
	g.wg.Done()
}

// Go 开启一个协程
func (g *Group) Go(fn func() error) {
	if g.ch != nil {
		g.ch <- token{}
	}
	g.wg.Add(1)
	go func() {
		defer g.done()

		if err := fn(); err != nil {
			// 只捕获第一个错误，无论执行多少次
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	}()
}
func (g *Group) Wait() error {
	g.wg.Wait()
	//cancel 一旦被调用，<-ctx.Done()就有信号
	if g.cancel != nil {
		g.cancel()
	}
	return g.err
}

// NewGroup 生成一个Group
//
// n 控制同时并发的数量，-1表示不限制并发数
func NewGroup(ctx context.Context, n int) (*Group, context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	g := Group{cancel: cancel}
	if n < 0 {
		g.ch = nil
		return &g, ctx
	}
	if len(g.ch) != 0 {
		panic(fmt.Sprintf("there are %v goroutines still alive", len(g.ch)))
	}
	//有缓冲通道
	g.ch = make(chan token, n)
	return &g, ctx
}
