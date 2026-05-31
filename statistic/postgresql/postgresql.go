package postgresql

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// PostgreSQL driver
	_ "github.com/lib/pq"

	"github.com/p4gefau1t/trojan-go/common"
	"github.com/p4gefau1t/trojan-go/config"
	"github.com/p4gefau1t/trojan-go/log"
	"github.com/p4gefau1t/trojan-go/statistic"
	"github.com/p4gefau1t/trojan-go/statistic/memory"
)

const Name = "POSTGRESQL"

type Authenticator struct {
	*memory.Authenticator
	db             *sql.DB
	updateDuration time.Duration
	ctx            context.Context
}

func (a *Authenticator) updater() {
	for {
		for _, user := range a.ListUsers() {
			hash := user.Hash()
			// Traffic directions are swapped intentionally to match the server-side perspective.
			sent, recv := user.ResetTraffic()

			s, err := a.db.Exec(
				"UPDATE users SET upload=upload+$1, download=download+$2 WHERE password=$3",
				recv, sent, hash,
			)
			if err != nil {
				log.Error(common.NewError("failed to update traffic data in user table").Base(err))
				continue
			}
			if r, err := s.RowsAffected(); err != nil {
				if r == 0 {
					a.DelUser(hash)
				}
			}
		}
		log.Info("buffered data has been written into the database")

		rows, err := a.db.Query("SELECT password, quota, download, upload FROM users")
		if err != nil || rows.Err() != nil {
			log.Error(common.NewError("failed to pull data from the database").Base(err))
			time.Sleep(a.updateDuration)
			continue
		}
		for rows.Next() {
			var hash string
			var quota, download, upload int64
			if err := rows.Scan(&hash, &quota, &download, &upload); err != nil {
				log.Error(common.NewError("failed to scan row from query result").Base(err))
				break
			}
			if download+upload < quota || quota < 0 {
				a.AddUser(hash)
			} else {
				a.DelUser(hash)
			}
		}

		select {
		case <-time.After(a.updateDuration):
		case <-a.ctx.Done():
			log.Debug("PostgreSQL daemon exiting...")
			return
		}
	}
}

func connectDatabase(cfg *PostgreSQLConfig) (*sql.DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.Username,
		cfg.Password,
		cfg.ServerHost,
		cfg.ServerPort,
		cfg.Database,
		cfg.SSLMode,
	)
	return sql.Open("postgres", dsn)
}

func NewAuthenticator(ctx context.Context) (statistic.Authenticator, error) {
	cfg := config.FromContext(ctx, Name).(*Config)
	db, err := connectDatabase(&cfg.PostgreSQL)
	if err != nil {
		return nil, common.NewError("failed to connect to PostgreSQL server").Base(err)
	}
	memoryAuth, err := memory.NewAuthenticator(ctx)
	if err != nil {
		return nil, err
	}
	a := &Authenticator{
		db:             db,
		ctx:            ctx,
		updateDuration: time.Duration(cfg.PostgreSQL.CheckRate) * time.Second,
		Authenticator:  memoryAuth.(*memory.Authenticator),
	}
	go a.updater()
	log.Debug("postgresql authenticator created")
	return a, nil
}

func init() {
	statistic.RegisterAuthenticatorCreator(Name, NewAuthenticator)
}