package utg

import (
	"context"
	"log/slog"
	"time"

	slerr "github.com/RakhimovAns/logger/pkg/err"
	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/td/examples"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"golang.org/x/time/rate"
)

type BotConfig struct {
	Log            *slog.Logger
	PhoneNumber    string
	AppID          int
	AppHash        string
	SessionStorage telegram.SessionStorage
}

type Bot struct {
	log             *slog.Logger
	client          *telegram.Client
	waiter          *floodwait.Waiter
	updatesRecovery *updates.Manager
	flow            auth.Flow
}

func New(config BotConfig) *Bot {
	d := tg.NewUpdateDispatcher()

	log := config.Log.With(slog.String("instance", "user_bot"))

	waiter := floodwait.NewWaiter().WithCallback(func(ctx context.Context, wait floodwait.FloodWait) {
		log.Warn("flood wait", slog.Duration("wait", wait.Duration))
	})

	flow := auth.NewFlow(examples.Terminal{
		PhoneNumber: config.PhoneNumber,
	}, auth.SendCodeOptions{})

	updatesRecovery := updates.New(updates.Config{
		Handler: d,
	})

	client := telegram.NewClient(
		config.AppID,
		config.AppHash,
		telegram.Options{
			UpdateHandler: updatesRecovery,
			Middlewares: []telegram.Middleware{
				waiter,
				ratelimit.New(rate.Every(time.Millisecond*100), 5),
			},
			SessionStorage: config.SessionStorage,
		},
	)

	return &Bot{
		log:             log,
		client:          client,
		waiter:          waiter,
		flow:            flow,
		updatesRecovery: updatesRecovery,
	}
}

func (b *Bot) Client() *telegram.Client {
	return b.client
}

func (b *Bot) Run(ctx context.Context, cb func(ctx context.Context, client *telegram.Client) error) error {
	err := b.waiter.Run(ctx, func(ctx context.Context) error {
		if err := b.client.Run(ctx, func(ctx context.Context) error {
			if err := b.client.Auth().IfNecessary(ctx, b.flow); err != nil {
				return slerr.WithSource(err, "auth")
			}

			user, err := b.client.Self(ctx)
			if err != nil {
				return slerr.WithSource(err, "self")
			}

			b.log.Info("go user bot!", slog.Int64("user_id", user.ID), slog.String("name", user.FirstName))

			if err := cb(ctx, b.client); err != nil {
				return slerr.WithSource(err, "callback")
			}

			return b.updatesRecovery.Run(ctx, b.client.API(), user.ID, updates.AuthOptions{
				IsBot: user.Bot,
				OnStart: func(ctx context.Context) {
					b.log.Info("update recovery initialized and started, listening for events")
				},
			})
		}); err != nil {
			return slerr.WithSource(err, "waiter")
		}

		return nil
	})
	if err != nil {
		return slerr.WithSource(err)
	}

	return nil
}

type Storage struct {
	redis *redis.Client
}

func NewStorage(client *redis.Client) *Storage {
	return &Storage{
		redis: client,
	}
}

func (s *Storage) LoadSession(ctx context.Context) ([]byte, error) {
	session, err := s.redis.Get(ctx, "telegram:bot:user:session").Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		return nil, errors.Wrap(err, "get session")
	}

	return []byte(session), nil
}

func (s *Storage) StoreSession(ctx context.Context, data []byte) error {
	err := s.redis.Set(ctx, "telegram:bot:user:session", string(data), 0).Err()

	return err
}
