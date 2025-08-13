package main

import (
	"context"
	"fmt"
	"log"
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
	operation, ok := params["operation"].(string)
	if !ok {
		return nil, fmt.Errorf("operation parameter is required")
	}

	a, aOk := params["a"].(float64)
	b, bOk := params["b"].(float64)
	if !aOk || !bOk {
		return nil, fmt.Errorf("parameters 'a' and 'b' must be numbers")
	}

	var result float64
	var err error

	switch operation {
	case "add":
		result = a + b
	case "subtract":
		result = a - b
	case "multiply":
		result = a * b
	case "divide":
		if b == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		result = a / b
	default:
		return nil, fmt.Errorf("unsupported operation: %s", operation)
	}

	return map[string]interface{}{
		"status":    "success",
		"plugin":    p.Name(),
		"version":   p.Version(),
		"operation": operation,
		"operands":  []float64{a, b},
		"result":    result,
		"message":   fmt.Sprintf("%.2f %s %.2f = %.2f", a, operation, b, result),
	}, err
}

func main() {
	plugin := NewCalculatorPlugin()
	// Avoid writing to stdout before handshake; use stderr if needed
	log.Printf("Starting Calculator Plugin %s v%s", plugin.Name(), plugin.Version())
	sdk.Serve(plugin)
}
