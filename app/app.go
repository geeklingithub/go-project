package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

// App 应用
type App struct {
	*Option                    //应用配置项
	ctx     context.Context    //应用上下文
	cancel  context.CancelFunc //上下文取消信号
}

// Server 服务接口
type Server interface {
	Start(context.Context)
	Stop(context.Context)
}

// Init 应用初始化
func Init(opts ...OptFunc) *App {
	//初始化配置项

	//默认配置
	o := &Option{
		closeSignals: []os.Signal{syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT},
	}

	//自定义配置
	for _, opt := range opts {
		opt(o)
	}

	//返回应用实例对象
	ctx, cancel := context.WithCancel(context.Background())
	return &App{
		Option: o,
		cancel: cancel,
		ctx:    ctx,
	}
}

// Start 应用启动
func (node *App) Start() {

	wg := &sync.WaitGroup{}
	for _, server := range node.servers {
		wg.Add(1)

		//服务关闭
		go func() {
			//应用关闭时,关闭服务
			<-node.ctx.Done()
			server.Stop(node.ctx)
		}()

		//服务启动
		server := server
		go func() {
			wg.Done()
			server.Start(node.ctx)
		}()
	}

	wg.Wait()

	//信号通知关服
	c := make(chan os.Signal, 1)
	signal.Notify(c, node.closeSignals...)
	go func() {
		<-c
		node.Stop()
	}()
}

// Stop 应用关闭
func (node *App) Stop() {

	node.cancel()
}
