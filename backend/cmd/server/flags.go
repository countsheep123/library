package main

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

var (
	listenAddr string
	sentryDSN  string
	gsSec      int64
	staticPath string
	pgReadDB   string
	pgWriteDB  string
)

var (
	flags = []cli.Flag{
		cli.StringFlag{
			Name:        "listen-addr",
			EnvVar:      "LISTEN_ADDR",
			Usage:       "Listen Addr",
			Value:       ":8081",
			Destination: &listenAddr,
		},
		cli.StringFlag{
			Name:        "sentry-dsn",
			EnvVar:      "SENTRY_DSN",
			Usage:       "Sentry DSN",
			Destination: &sentryDSN,
		},
		cli.Int64Flag{
			Name:        "graceful-shutdown-sec",
			EnvVar:      "GRACEFUL_SHUTDOWN_SEC",
			Usage:       "Graceful Shutdown Sec",
			Value:       5,
			Destination: &gsSec,
		},
		cli.StringFlag{
			Name:        "static-path",
			EnvVar:      "STATIC_PATH",
			Usage:       "Static Path",
			Value:       "../frontend/dist",
			Destination: &staticPath,
		},
		cli.StringFlag{
			Name:        "pg-read-db",
			EnvVar:      "PG_READ_DB",
			Usage:       "postgres read db",
			Value:       "postgres://user:pass@localhost/library_db?sslmode=disable",
			Destination: &pgReadDB,
		},
		cli.StringFlag{
			Name:        "pg-write-db",
			EnvVar:      "PG_WRITE_DB",
			Usage:       "postgres write db",
			Value:       "postgres://user:pass@localhost/library_db?sslmode=disable",
			Destination: &pgWriteDB,
		},
	}
)

func validate() error {
	if len(listenAddr) == 0 {
		return fmt.Errorf("invalid listen addr: %s", listenAddr)
	}
	if len(pgReadDB) == 0 {
		return fmt.Errorf("invalid pg read db: %s", pgReadDB)
	}
	if len(pgWriteDB) == 0 {
		return fmt.Errorf("invalid pg write db: %s", pgWriteDB)
	}
	return nil
}
