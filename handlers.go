package main

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"slices"
	"strings"

	"github.com/NaeCqde/disgolf"
)

var (
	ChannelMentionPattern = regexp.MustCompile(`(?m)<#(\d+)>|https://(?:[A-Za-z]+\.)?discord.com/channels/\d+/(\d+)`)
)

func HasOwnerMiddleware(ctx *disgolf.MessageCtx) {
	if slices.Contains(config.OwnerIDs, ctx.Message.Author.ID) {
		ctx.Next()
	}
}

/*
func IgnoreSelfMiddleware(ctx *disgolf.MessageCtx) {
	if ctx.Message.Author.ID != ctx.State.User.ID && !ctx.Message.Author.Bot {
		ctx.Next()
	}
}
*/

func JoinVC(ctx *disgolf.MessageCtx) (*VoiceManager, error) {
	// 接続済みのボイスチャンネルがあるか確認
	voiceManager, err := GetManager(ctx.Message.GuildID, ctx.State.User.ID)
	var voiceChannelID string

	var mentionChannels []string
	matches := ChannelMentionPattern.FindAllStringSubmatch(ctx.Message.Content, -1)
	for _, v := range matches {
		if strings.HasPrefix(v[0], "<#") {
			mentionChannels = append(mentionChannels, v[1])
		} else {
			mentionChannels = append(mentionChannels, v[2])
		}
	}

	if len(mentionChannels) != 0 {
		voiceChannelID = mentionChannels[0]
	} else {
		for _, g := range ctx.State.Guilds {
			for _, v := range g.VoiceStates {
				if v.UserID == ctx.Message.Author.ID {
					voiceChannelID = v.ChannelID
				}
			}
		}
	}

	if voiceChannelID == "" {
		return nil, fmt.Errorf("voice channel is not selected")
	}

	// 既にボイスチャンネルに入っている
	if err == nil {
		if voiceManager.Channel.ID != voiceChannelID {
			//対象のボイスチャンネルが違う場合は移動する
			voiceChannel, err := ctx.Channel(voiceChannelID)
			if err != nil {
				return nil, err
			}

			err = ctx.VoiceStateUpdate(voiceManager.GuildID, voiceChannelID, true, true)
			if err != nil {
				return nil, err
			}

			voiceManager.Channel = voiceChannel
		}
		//同じボイスチャンネルなら何もしない

		return voiceManager, nil
	}

	voiceChannel, err := ctx.Channel(voiceChannelID)
	if err != nil {
		return nil, err
	}

	conn, err := ctx.ChannelVoiceJoin(context.TODO(), ctx.Message.GuildID, voiceChannelID, true, true)
	if err != nil {
		return nil, err
	}

	slog.Info("参加", "ChannelID", voiceChannelID)
	return NewManager(conn, ctx.Message.GuildID, ctx.State.User.ID, voiceChannel)
}

func HandleVCLeave(ctx *disgolf.MessageCtx) {
	voiceManager, err := GetManager(ctx.Message.GuildID, ctx.State.User.ID)

	if err == nil {
		ctx.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, "🆗")

		slog.Info("退出", "ChannelID", voiceManager.Channel.ID)
		voiceManager.Leave()
	}
}

func HandlePlay(ctx *disgolf.MessageCtx) {
	if len(ctx.Arguments) >= 1 {
		if audio, ok := config.Audio[ctx.Arguments[0]]; ok {
			// エラーが起きても無視
			ctx.MessageReactionAdd(ctx.Message.ChannelID, ctx.Message.ID, "🆗")

			mgr, err := JoinVC(ctx)
			if err != nil {
				slog.Error(err.Error())
				return
			}

			go mgr.Play(audio)
		}
	} else {
		text := "<name>    ：説明\n"
		for name, audio := range config.Audio {
			text += fmt.Sprintf("%-15s：%s\n", name, audio.Description)
		}

		text += fmt.Sprintf("\n構文：%splay <name>\n例：　%splay op", config.Prefix, config.Prefix)

		ctx.Reply(text, false)
	}
}
