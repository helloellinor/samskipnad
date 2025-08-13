package main

import (
	"context"
	"samskipnad/pkg/sdk"
)

type CalculatorPlugin struct {
	*sdk.BasePlugin
}

func NewCalculatorPlugin() *CalculatorPlugin {
	return &CalculatorPlugin{
		BasePlugin: sdk.NewBasePlugin("calculator", "1.0.0"),
	}
}

func (p *CalculatorPlugin) Execute(ctx context.Context, params map[string]interface{}) (map[string]interface{}, error) {
	return map[string]interface{}{
		"status": "success",
		"plugin": p.Name(),
	}, nil
}

func main() {
	sdk.Serve(NewCalculatorPlugin())
}
