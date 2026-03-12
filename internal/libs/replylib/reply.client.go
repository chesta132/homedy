package replylib

import (
	"homedy/config"

	"github.com/chesta132/goreply/reply"
)

var Client = reply.NewClient(reply.Client{
	CodeAliases: CodeAliases,
	Transformer: transformer,
	DebugMode:   config.IsEnvDev(),
})
