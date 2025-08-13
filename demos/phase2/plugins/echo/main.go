package main

import (
	"context"
	"fmt"
	"log"
	"samskipnad/pkg/sdk"
)

type EchoPlugin struct {
	*sdk.BasePlugin
}

func NewEchoPlugin() *EchoPlugin {
	return &EchoPlugin{
		BasePlugin: sdk.NewBasePlugin("echo", "1.0.0"),
	}
}

func (p *EchoPlugin) Execute(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	message, ok := params["message"].(string)
	if !ok {
		message = "Hello from Echo Plugin!"
	}

	return map[string]interface{}{
		"status":       "success",
		"plugin":       p.Name(),
		"version":      p.Version(),
		"echo":         message,
		"params_count": len(params),
		"params":       params,
		"timestamp":    fmt.Sprintf("%d", ctx.Value("timestamp")),
	}, nil
}

func main() {
	plugin := NewEchoPlugin()
	// Avoid writing to stdout before handshake; use stderr if needed
	log.Printf("Starting Echo Plugin %s v%s", plugin.Name(), plugin.Version())
	sdk.Serve(plugin)
}
