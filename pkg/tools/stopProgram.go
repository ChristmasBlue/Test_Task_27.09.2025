package tools

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

// StopContext ожидает сигнала для вызова функции cancel(), которая отвечает за остановку сервиса
func StopContext(ctx *context.Context, cancel context.CancelFunc) {
	signCh := make(chan os.Signal, 1)
	signal.Notify(signCh, syscall.SIGINT)

	<-signCh
	cancel()
}
