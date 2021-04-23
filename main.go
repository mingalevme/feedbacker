package main

import (
	"context"
	"fmt"
	_ "github.com/lib/pq"
	"github.com/mingalevme/feedbacker/internal/app"
	"github.com/mingalevme/feedbacker/internal/http"
	"github.com/mingalevme/feedbacker/pkg/envvarbag"
	"github.com/pkg/errors"
	netHTTP "net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	envVarBag := envvarbag.New()
	var env app.Env = app.NewEnv(envVarBag)
	env.Build()

	if err := env.Dispatcher().Run(); err != nil {
		panic(errors.Wrap(err, "dispatcher: starting"))
	}

	address := envVarBag.Get("HTTP_LISTEN_ADDRESS", "0.0.0.0:8080")
	var s http.Server = http.NewEchoServer(address, env)
	go func() {
		fmt.Printf("Listening on http://%s\n", s.GetAddr())
		err := s.ListenAndServe()
		if !errors.Is(err, netHTTP.ErrServerClosed) {
			panic(errors.Wrap(err, "http server: listening and serving"))
		}
	}()

	WaitForShutdown(func() {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		wg := sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("Shutting down the HTTP-server ...")
			if err := s.Shutdown(ctx); err != nil {
				fmt.Println("Error while shutting down the HTTP-server", err)
			} else {
				fmt.Println("The HTTP-server has been shut down")
			}
		}()
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("Shutting down the dispatcher ...")
			if err := env.Dispatcher().Stop(); err != nil {
				fmt.Println("Error while shutting down the dispatcher", err)
			} else {
				fmt.Println("The dispatcher has been shut down")
			}
		}()
		wg.Wait()
		// Closing container after the server and dispatcher to gracefully finish the queue
		wg.Add(1)
		go func() {
			defer wg.Done()
			fmt.Println("Closing the container ...")
			env.Close()
			fmt.Println("Container has been closed")
		}()
		wg.Wait()
	})
}

func HandleSignal(signals ...os.Signal) chan os.Signal {
	c := make(chan os.Signal, 1)
	signal.Notify(c, signals...)
	return c
}

func WaitForShutdown(shutdowner func()) {
	count := 0
	for {
		sig := <-HandleSignal(os.Interrupt)
		fmt.Println("Signal received", sig)
		count++
		if count > 1 {
			fmt.Println("Exit")
			os.Exit(1)
		}
		fmt.Println("Shutting down the app ...")
		go func() {
			shutdowner()
			fmt.Println("The app has been successfully shut down")
			os.Exit(1)
		}()
	}
}
