package di

import (
	"context"
	"time"

	"github.com/RakhimovAns/CodeforcesStats/internal/config"
	diut "github.com/RakhimovAns/CodeforcesStats/pkg/di"
	"github.com/RakhimovAns/txmananger/postgres"
	txman "github.com/RakhimovAns/txmananger/tx_manager"
	"github.com/RakhimovAns/wrapper/pkg/closer"
	"github.com/exaring/otelpgx"
)

func (d *DI) Postgres(ctx context.Context) postgres.Postgres {
	return diut.Once(ctx, func(ctx context.Context) postgres.Postgres {
		pcfg := postgres.NewConfig(
			config.PostgresUsername(),
			config.PostgresPassword(),
			config.PostgresHost(),
			config.PostgresPort(),
			config.PostgresDatabase(),
		)

		pcfg.WithConnAmount(90)
		pcfg.WithMinConnAmount(20)
		pcfg.WithMaxConnLifetime(time.Hour)
		pcfg.WithMaxConnIdleTime(10 * time.Minute)
		pcfg.WithHealthCheckPeriod(2 * time.Minute)

		pcfg.WithTracer(otelpgx.NewTracer(
			otelpgx.WithIncludeQueryParameters(),
			otelpgx.WithTrimSQLInSpanName(),
		))

		pg, err := postgres.NewPostgres(ctx, d.Log(ctx), pcfg)
		if err != nil {
			d.mustExit(err)
		}

		if err := otelpgx.RecordStats(pg.Pool()); err != nil {
			d.mustExit(err)
		}

		closer.Add(func() error {
			d.Log(ctx).Info("shutting down postgres")

			pg.Close()

			return nil
		})

		return pg
	})
}

func (d *DI) TxManager(ctx context.Context) txman.TxManager {
	return diut.Once(ctx, func(ctx context.Context) txman.TxManager {
		return txman.New(d.Postgres(ctx))
	})
}
