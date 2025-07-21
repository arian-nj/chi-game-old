package gamesession

import "gopkg.in/telebot.v4"

type GameSession struct {
	Bot *telebot.Bot

	IsChatOn    bool
	IsGameEnded bool

	Gametype  gametype.GameType
	GameState Game

	CreatedAt       time.Time
	ExpireDuaration time.Duration
}

func NewGameSession(allSession *AllSession, bot *telebot.Bot, gameType gametype.GameType, gameState Game) *GameSession {
	gs := &GameSession{
		Bot:             bot,
		Gametype:        gameType,
		GameState:       gameState,
		IsChatOn:        true,
		CreatedAt:       time.Now(),
		ExpireDuaration: 2*time.Minute*2 + 30,
	}

	utils.RunBackgroundTask(func() {
		gs.MonitorGame(allSession)
	})

	return gs
}

func (g *GameSession) HandleChatMessage(bot telebot.API, senderID int, text string) error {
	players := g.GameState.Players()
	var senderPlayer *consoleplayer.ConsolePlayer
	var recieverPlayer *consoleplayer.ConsolePlayer

	for _, p := range players {
		if p.TgID != senderID {
			recieverPlayer = p
		} else {
			senderPlayer = p
		}
	}

	_, err := bot.Send(&telebot.User{ID: int64(recieverPlayer.TgID)},
		fmt.Sprintf("*_%s:_* %s", senderPlayer.Name, text), telebot.ModeMarkdownV2)
	return err
}

func (gs *GameSession) MonitorGame(allSession *AllSession) {
	<-gs.GameState.GetContext().Done()
	gs.IsGameEnded = true
	if gs.IsChatOn {
		expDur := 30 * time.Second
		gs.ExpireDuaration = time.Since(gs.CreatedAt) + expDur

		text := fmt.Sprintf("چت تا %d ثانیه دیگه بسته میشه", int(expDur.Seconds()))
		for _, player := range gs.GameState.Players() {
			_, err := gs.Bot.Send(player, text)
			if err != nil {
				slog.Error("can't send end game chat message", "err", err)
			}
		}
		time.Sleep(gs.ExpireDuaration)
	}
	gs.CleanAndDisconnect(allSession)
}

func (gs *GameSession) CleanAndDisconnect(allSession *AllSession) {
	if gs.IsChatOn {
		text := "چت قطع شد"
		for _, player := range gs.GameState.Players() {
			_, err := gs.Bot.Send(player, text, welcomeReplyKeyboard)
			if err != nil {
				slog.Error("can't send chat ended message", "err", err)
			}
		}
	}

	allSession.Mutex.Lock()
	defer allSession.Mutex.Unlock()

	for _, player := range gs.GameState.Players() {
		delete(allSession.Sessions, strconv.Itoa(player.TgID))
	}
}
