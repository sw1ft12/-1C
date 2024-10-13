package main

import (
    "1C/internal/config"
    "1C/internal/server"
    "context"
    "github.com/jackc/pgx/v5/pgxpool"
    "log"
)

func main() {
    cfg := config.GetConfig()
    pgConn, err := pgxpool.New(context.TODO(), cfg.PostgresConn)
    err = pgConn.Ping(context.TODO())
    if err != nil {
        log.Fatal(err)
    }
    s := server.NewServer(pgConn)
    err = s.Run(cfg.Address)
    if err != nil {
        log.Fatal(err)
    }
}
