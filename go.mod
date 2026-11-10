module vc_spammer

go 1.25.0

require (
	github.com/NaeCqde/dca v0.0.0-20251007184203-c569c9368d38
	github.com/NaeCqde/disgolf v0.0.0-20251006093216-68255c69d3a6
	github.com/bwmarrin/discordgo v0.29.0
	github.com/goccy/go-yaml v1.18.0
)

require (
	github.com/gorilla/websocket v1.5.3 // indirect
	github.com/jonas747/ogg v0.0.0-20161220051205-b4f6f4cf3757 // indirect
	golang.org/x/crypto v0.52.0 // indirect
	golang.org/x/net v0.54.0 // indirect
	golang.org/x/sys v0.45.0 // indirect
)

replace github.com/bwmarrin/discordgo => ../vc_loopback/discordgo
