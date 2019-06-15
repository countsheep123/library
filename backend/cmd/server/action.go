package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/countsheep123/library/api"
	"github.com/countsheep123/library/db"
	"github.com/getsentry/raven-go"
	"go.uber.org/zap"
)

func action() error {
	readDB, err := sql.Open("postgres", pgReadDB)
	if err != nil {
		raven.CaptureError(err, nil)
		zap.S().Error(err)
		return err
	}
	defer readDB.Close()

	writeDB, err := sql.Open("postgres", pgWriteDB)
	if err != nil {
		raven.CaptureError(err, nil)
		zap.S().Error(err)
		return err
	}
	defer writeDB.Close()

	dbHandler, err := db.New(readDB, writeDB)
	if err != nil {
		raven.CaptureError(err, nil)
		zap.S().Error(err)
		return err
	}

	server, err := api.New(dbHandler, staticPath)
	if err != nil {
		raven.CaptureError(err, nil)
		zap.S().Error(err)
		return err
	}

	httpHandler, err := server.Handler()
	if err != nil {
		raven.CaptureError(err, nil)
		zap.S().Error(err)
		return err
	}

	srv := &http.Server{
		Addr:    listenAddr,
		Handler: httpHandler,
	}

	idleConnsClosed := make(chan struct{})

	go func() {
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		close(quit)

		zap.S().Info("shutting down in " + fmt.Sprint(gsSec))

		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(gsSec))
		defer cancel()

		srv.SetKeepAlivesEnabled(false)
		if err := srv.Shutdown(ctx); err != nil {
			zap.S().Fatal(err)
			return
		}
		close(idleConnsClosed)
		zap.S().Info("server exited")
	}()

	zap.S().Info("server started")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		zap.S().Fatal(err)
	}

	<-idleConnsClosed

	return nil
}
