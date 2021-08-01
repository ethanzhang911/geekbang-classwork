package main

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"net/http"
	"os"
	"os/signal"
)

// startServer启动一个http server
// 因为要传给group.go函数内，所以要返回error
func startServer(s *http.Server) error {
	http.HandleFunc("/ping", handlerPing)
	return s.ListenAndServe()
}

func handlerPing(w http.ResponseWriter,req * http.Request){
	w.Write([]byte("pong"))
}


func main() {
	ctx, cancel := context.WithCancel(context.Background())
	group, groupErrCtx := errgroup.WithContext(ctx)
	// 这里要确保换缓冲通道数量至少为1
	c := make(chan os.Signal, 1)
	signal.Notify(c)

	httpServer := &http.Server{
		Addr:  "127.0.0.1:8080",
	}

	group.Go(func() error{
		fmt.Println("http服务已启动")
		return startServer(httpServer)
	})

	group.Go(func() error {
		<-groupErrCtx.Done()
		fmt.Println("http服务已关闭")
		return httpServer.Shutdown(groupErrCtx)
	})

	group.Go(func() error {
		for{
			select {
			case <-groupErrCtx.Done():
				return groupErrCtx.Err()
			case <-c:
				fmt.Println("signal信号监测到，服务停止")
				cancel()
				return nil
			}
		}
	})

	if err:= group.Wait();  err !=nil  {
       fmt.Println("errorgroup error",err)
	}
	fmt.Println("所有group均退出")

}

