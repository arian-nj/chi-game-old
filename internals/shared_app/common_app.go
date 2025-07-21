package sharedapp

import (
	"context"
	"sync"

	"github.com/arian-nj/chibazi/database"
	"github.com/arian-nj/chibazi/internals/config"
	gametype "github.com/arian-nj/chibazi/internals/game_type"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/telebot.v4"
)

type SharedApp struct {
	Config  *config.Config
	Queries *database.Queries
	Conn    *pgxpool.Pool
	Bot     *telebot.Bot
	Wg      *sync.WaitGroup

	GameSessions *AllSession
	MatchMaking  *MatchMaking
}

func NewSharedApp(conf *config.Config) *SharedApp {
	return &SharedApp{
		Config: conf,
		Wg:     &sync.WaitGroup{},

		GameSessions: &AllSession{
			Sessions: map[string]*GameSession{},
			Mutex:    sync.Mutex{},
		},
		MatchMaking: &MatchMaking{
			WaitingPlayers: map[gametype.GameType][]*Ticket{},
			Mutex:          sync.Mutex{},
		},
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
