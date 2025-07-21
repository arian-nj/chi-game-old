package sharedapp

import (
	"context"
	"sync"

	"github.com/arian-nj/chibazi/database"
	"github.com/arian-nj/chibazi/internals/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/telebot.v4"
)

type SharedApp struct {
	Config  *config.Config
	Queries *database.Queries
	Conn    *pgxpool.Pool
	Bot     *telebot.Bot
	Wg      *sync.WaitGroup
}

func NewSharedApp(conf *config.Config) *SharedApp {
	return &SharedApp{
		Config: conf,
		Wg:     &sync.WaitGroup{},
	}

}

func (c *SharedApp) ConfigureDatabase() error {
	conn, err := pgxpool.New(context.Background(), c.Config.DatabseUrl)
	if err != nil {
		return err
	}

	c.Queries = database.New(conn)
	c.Conn = conn
	return nil
}
