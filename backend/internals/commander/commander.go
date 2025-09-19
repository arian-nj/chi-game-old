package commander

import (
	"sync"

	"github.com/arian-nj/chibazi/backend/internals/utils"
)

type Subscriber interface {
	Update(command Command)
}

type Command interface {
	Execute()
}

type Commander struct {
	Commands        []Command
	DoneCommands    []Command
	CommandNotifire chan any

	Subscribers []Subscriber

	mu sync.Mutex
}

func NewCommander() *Commander {
	return &Commander{
		Commands:        []Command{},
		DoneCommands:    []Command{},
		CommandNotifire: make(chan any, 6),

		Subscribers: []Subscriber{},
	}
}
func (commander *Commander) notif() {
	utils.RunBackgroundTask(func() {
		commander.CommandNotifire <- nil
	})
}

// Commands
func (commander *Commander) PushCommand(newCommand Command) {
	commander.Commands = append(commander.Commands, newCommand)
	commander.notif()
}

func (commander *Commander) InjectCommand(newAction Command) {
	commander.Commands = append([]Command{newAction}, commander.Commands...)
	commander.notif()
}
func (commander *Commander) PopCommand() Command {
	firstAction := commander.Commands[0]
	commander.Commands = commander.Commands[1:]
	return firstAction
}

func (commander *Commander) ApplyCommand(newCommand Command) {
	commander.mu.Lock()
	defer commander.mu.Unlock()

	newCommand.Execute()
	commander.Notify(newCommand)
	commander.DoneCommands = append(commander.DoneCommands, newCommand)
}

// Subscription
func (commander *Commander) Notify(command Command) {
	var wg sync.WaitGroup

	for _, sub := range commander.Subscribers {
		wg.Add(1)
		go func(s Subscriber) {
			defer wg.Done()
			sub.Update(command) // pass both state + action
		}(sub)
	}
	wg.Wait()
}

func (commander *Commander) Subscribe(subscriber Subscriber) {
	commander.Subscribers = append(commander.Subscribers, subscriber)
}

func (commander *Commander) Unsubscribe(targetSub Subscriber) {
	for i, sub := range commander.Subscribers {
		if sub == targetSub {
			commander.Subscribers = append(
				commander.Subscribers[:i],
				commander.Subscribers[i+1:]...,
			)
			return
		}
	}
}
