package commonapp

import (
	"context"
	"sync"

	"github.com/arian-nj/chibazi/database"
	"github.com/arian-nj/chibazi/internals/config"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/telebot.v4"
)

type CommonApp struct {
	Config  *config.Config
	Queries *database.Queries
	Conn    *pgxpool.Pool
	Bot     *telebot.Bot
	Wg      *sync.WaitGroup
}

func NewCommon(conf *config.Config) *CommonApp {
	return &CommonApp{
		Config: conf,
		Wg:     &sync.WaitGroup{},
	}

}

func (c *CommonApp) ConfigureDatabase() error {

	conn, err := pgxpool.New(context.Background(), c.Config.DatabseUrl)
	if err != nil {
		return err
	}

	c.Queries = database.New(conn)
	c.Conn = conn
	return nil
}
