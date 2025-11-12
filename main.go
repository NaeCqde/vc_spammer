package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/NaeCqde/disgolf"
	"github.com/bwmarrin/discordgo"
)

func main() {
	slog.SetDefault(
		slog.New(
			slog.NewTextHandler(
				os.Stdout,
				&slog.HandlerOptions{Level: slog.LevelInfo},
			),
		),
	)

	slog.Info("Start")

	err := LoadConfig()
	if err != nil {
		panic(err)
	}

	slog.Info("DATA", "Prefix", config.Prefix, "OwnerIDs", config.OwnerIDs)

	var bots []*disgolf.Bot

	for _, token := range config.Tokens {
		bot, err := disgolf.New(token)
		if err != nil {
			panic(err)
		}

		bot.Identify.Intents = discordgo.IntentsAll

		bot.Router.Register(&disgolf.Command{
			Name:               "play",
			Description:        "音楽を再生します",
			MessageHandler:     disgolf.MessageHandlerFunc(HandlePlay),
			MessageMiddlewares: []disgolf.MessageHandler{disgolf.MessageHandlerFunc(HasOwnerMiddleware)},
		})

		bot.Router.Register(&disgolf.Command{
			Name:               "leave",
			Description:        "ボイスチャンネルから切断します",
			MessageHandler:     disgolf.MessageHandlerFunc(HandleVCLeave),
			MessageMiddlewares: []disgolf.MessageHandler{disgolf.MessageHandlerFunc(HasOwnerMiddleware)},
		})

		bot.AddHandler(bot.Router.MakeMessageHandler(&disgolf.MessageHandlerConfig{
			Prefixes:      []string{config.Prefix},
			MentionPrefix: true,
		}))

		bot.AddHandler(func(s *discordgo.Session, m *discordgo.Ready) {
			slog.Info("Bot is ready", "UserName", fmt.Sprintf("%v#%v", m.User.Username, m.User.Discriminator), "UserID", m.User.ID)
		})

		if err = bot.Open(); err != nil {
			panic(err)
		}

		bots = append(bots, bot)

		time.Sleep(time.Second * 5)
	}

	stopBot := make(chan os.Signal, 1)
	signal.Notify(stopBot, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)

	<-stopBot

	for _, bot := range bots {
		err = bot.Close()
		if err != nil {
			panic(err)
		}
	}

	slog.Info("End")
}
