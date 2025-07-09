package commonapp

import (
	"context"
	"sync"

	"github.com/arian-nj/chibazi/database"
	"github.com/arian-nj/chibazi/games/dotbox_console"
	xoconsole "github.com/arian-nj/chibazi/games/xo_console"
	"github.com/arian-nj/chibazi/internals/config"
	gametype "github.com/arian-nj/chibazi/internals/game_type"
	"github.com/arian-nj/chibazi/internals/lobby"
	matchmaking "github.com/arian-nj/chibazi/internals/match_making"
	"github.com/jackc/pgx/v5/pgxpool"
	"gopkg.in/telebot.v4"
)

type CommonApp struct {
	Config      *config.Config
	Queries     *database.Queries
	Conn        *pgxpool.Pool
	Bot         *telebot.Bot
	Wg          *sync.WaitGroup
	MatchMaking *matchmaking.MatchMaking
	Lobby       *lobby.Lobby
}

func NewCommon(conf *config.Config) *CommonApp {
	return &CommonApp{
		Config: conf,
		Wg:     &sync.WaitGroup{},
		MatchMaking: &matchmaking.MatchMaking{
			WaitingPlayers: make(map[gametype.GameType][]*matchmaking.Ticket),
			Mutex:          sync.Mutex{},
		},
		Lobby: &lobby.Lobby{
			XOGames: make(map[string]*xoconsole.XOGame),
			DotBox:  make(map[string]*dotbox_console.DotBoxGame),
		},
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
