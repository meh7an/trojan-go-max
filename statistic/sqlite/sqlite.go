package sqlite

import (
	"context"
	"database/sql"
	"time"

	// SQLite driver — pure Go, no CGO required
	_ "modernc.org/sqlite"

	"github.com/p4gefau1t/trojan-go/common"
	"github.com/p4gefau1t/trojan-go/config"
	"github.com/p4gefau1t/trojan-go/log"
	"github.com/p4gefau1t/trojan-go/statistic"
	"github.com/p4gefau1t/trojan-go/statistic/memory"
)

const Name = "SQLITE"

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
				"UPDATE users SET upload=upload+?, download=download+? WHERE password=?",
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
			log.Debug("SQLite daemon exiting...")
			return
		}
	}
}

func initDatabase(db *sql.DB) error {
	// WAL mode allows concurrent reads alongside writes without blocking.
	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		return err
	}
	// Busy timeout prevents immediate failure when another process holds a write lock.
	if _, err := db.Exec("PRAGMA busy_timeout=5000"); err != nil {
		return err
	}
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS users (
		id       INTEGER PRIMARY KEY AUTOINCREMENT,
		username TEXT    NOT NULL DEFAULT '',
		password TEXT    NOT NULL,
		quota    INTEGER NOT NULL DEFAULT 0,
		download INTEGER NOT NULL DEFAULT 0,
		upload   INTEGER NOT NULL DEFAULT 0
	)`)
	if err != nil {
		return err
	}
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_users_password ON users (password)")
	return err
}

func NewAuthenticator(ctx context.Context) (statistic.Authenticator, error) {
	cfg := config.FromContext(ctx, Name).(*Config)
	db, err := sql.Open("sqlite", cfg.SQLite.Path)
	if err != nil {
		return nil, common.NewError("failed to open SQLite database").Base(err)
	}
	if err := initDatabase(db); err != nil {
		return nil, common.NewError("failed to initialize SQLite database schema").Base(err)
	}
	memoryAuth, err := memory.NewAuthenticator(ctx)
	if err != nil {
		return nil, err
	}
	a := &Authenticator{
		db:             db,
		ctx:            ctx,
		updateDuration: time.Duration(cfg.SQLite.CheckRate) * time.Second,
		Authenticator:  memoryAuth.(*memory.Authenticator),
	}
	go a.updater()
	log.Debug("sqlite authenticator created")
	return a, nil
}

func init() {
	statistic.RegisterAuthenticatorCreator(Name, NewAuthenticator)
}
