package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	jobinterface "williamfeng323/mooncake-duty/src/interfaces/job"

	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"

	middleware "williamfeng323/mooncake-duty/src/infrastructure/middlewares"
	webInterface "williamfeng323/mooncake-duty/src/interfaces/web"
)

func main() {
	router := gin.Default()
	router.Use(middleware.Logger())

	webInterface.RegisterAccountRoute(router)
	webInterface.RegisterProjectRoute(router)
	srv := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}

	go func() {
		// service connections
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	cronJobs := createCronJobs()
	go cronJobs.Start()
	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	jobCtx := cronJobs.Stop()
	if jobCtx != nil {
		<-jobCtx.Done()
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")
}

func createCronJobs() *cron.Cron {
	logger := cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))
	cronJobs := cron.New(
		cron.WithChain(cron.SkipIfStillRunning(logger), cron.Recover(logger)),
		cron.WithLogger(logger))

	jobinterface.RegisterIssueJobs(cronJobs)
	return cronJobs
}
