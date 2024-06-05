package main

import (
	"context"
	"database/sql"
	"github.com/heroiclabs/nakama-common/runtime"
)

func InitModule(ctx context.Context, logger runtime.Logger, db *sql.DB, nk runtime.NakamaModule, initializer runtime.Initializer) error {
	logger.Info("Go module loaded")

	if err := initializer.RegisterRpc("GetContent", GetContent); err != nil {
		return err
	}
	return nil
}
