package commonapp

import (
	"context"

	"github.com/arian-nj/ultrun/database"
	"github.com/arian-nj/ultrun/internals/config"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CommonApp struct {
	Config  *config.Config
	Queries *database.Queries
	Conn    *pgxpool.Pool
}

func NewCommon(conf *config.Config) *CommonApp {
	return &CommonApp{
		Config: conf,
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
