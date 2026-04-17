package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"

	appcore "opencode-cli-proxy/internal/app"
	"opencode-cli-proxy/internal/config"
)

func main() {
	a := app.NewWithID("com.opencode.proxy.desktop")
	w := a.NewWindow("OpenCode Go Gateway")
	w.Resize(fyne.NewSize(760, 560))

	configPath := widget.NewEntry()
	configPath.SetText("configs/config.example.yaml")

	status := widget.NewMultiLineEntry()
	status.SetMinRowsVisible(12)
	status.Disable()
	status.SetText("Ready.\n")

	serverHost := widget.NewEntry()
	serverHost.SetText("127.0.0.1")

	serverPort := widget.NewEntry()
	serverPort.SetText("18080")

	opencodeBinary := widget.NewEntry()
	opencodeBinary.SetText("opencode")
	opencodeBinary.SetPlaceHolder("opencode")

	opencodeAttach := widget.NewEntry()
	opencodeAttach.SetPlaceHolder("http://127.0.0.1:4096 (optional)")

	gatewayKey := widget.NewEntry()
	gatewayKey.SetText("sk-gw-demo")

	opencodeModel := widget.NewEntry()
	opencodeModel.SetText("opencode-go/glm-5.1")
	opencodeModel.SetPlaceHolder("opencode-go/glm-5.1")

	allowedModels := widget.NewEntry()
	allowedModels.SetText("opencode-go/glm-5.1")
	allowedModels.SetPlaceHolder("comma-separated models")

	openaiURL := widget.NewEntry()
	openaiURL.Disable()

	var mu sync.Mutex
	var gateway *appcore.Gateway
	var running bool

	appendStatus := func(msg string) {
		status.SetText(status.Text + msg + "\n")
	}

	parseCSV := func(input string) []string {
		parts := strings.Split(input, ",")
		result := make([]string, 0, len(parts))
		for _, part := range parts {
			item := strings.TrimSpace(part)
			if item != "" {
				result = append(result, item)
			}
		}
		return result
	}

	updateOpenAIURL := func() {
		openaiURL.SetText(fmt.Sprintf("http://%s:%s/v1", strings.TrimSpace(serverHost.Text), strings.TrimSpace(serverPort.Text)))
	}
	updateOpenAIURL()
	serverHost.OnChanged = func(string) { updateOpenAIURL() }
	serverPort.OnChanged = func(string) { updateOpenAIURL() }

	applyConfig := func() (*config.Config, error) {
		cfg, err := config.Load(configPath.Text)
		if err != nil {
			return nil, err
		}
		cfg.Server.Host = strings.TrimSpace(serverHost.Text)
		fmt.Sscanf(strings.TrimSpace(serverPort.Text), "%d", &cfg.Server.Port)
		if binary := strings.TrimSpace(opencodeBinary.Text); binary != "" {
			cfg.Upstream.Binary = binary
		}
		cfg.Upstream.Attach = strings.TrimSpace(opencodeAttach.Text)

		modelName := strings.TrimSpace(opencodeModel.Text)
		allowed := parseCSV(allowedModels.Text)
		if modelName != "" {
			cfg.Models = map[string]string{modelName: modelName}
			if len(allowed) == 0 {
				allowed = []string{modelName}
			}
		}
		if key := strings.TrimSpace(gatewayKey.Text); key != "" {
			cfg.Keys = map[string]config.KeyConfig{
				key: {
					Account:       "default",
					AllowedModels: allowed,
				},
			}
		}
		return cfg, cfg.Validate()
	}

	startButton := widget.NewButton("Start Gateway", func() {
		mu.Lock()
		defer mu.Unlock()
		if running {
			appendStatus("Gateway is already running.")
			return
		}
		cfg, err := applyConfig()
		if err != nil {
			dialog.ShowError(err, w)
			appendStatus("Failed to load config: " + err.Error())
			return
		}
		gateway = appcore.NewGateway(cfg)
		running = true
		appendStatus("Starting gateway at http://" + gateway.Address())
		go func() {
			if err := gateway.Start(); err != nil && err.Error() != "http: Server closed" {
				mu.Lock()
				running = false
				mu.Unlock()
				appendStatus("Gateway stopped with error: " + err.Error())
			}
		}()
	})

	stopButton := widget.NewButton("Stop Gateway", func() {
		mu.Lock()
		defer mu.Unlock()
		if !running || gateway == nil {
			appendStatus("Gateway is not running.")
			return
		}
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := gateway.Stop(ctx); err != nil {
			dialog.ShowError(err, w)
			appendStatus("Failed to stop gateway: " + err.Error())
			return
		}
		running = false
		appendStatus("Gateway stopped.")
	})

	form := widget.NewForm(
		widget.NewFormItem("Config file", configPath),
		widget.NewFormItem("Listen host", serverHost),
		widget.NewFormItem("Listen port", serverPort),
		widget.NewFormItem("Opencode binary", opencodeBinary),
		widget.NewFormItem("Attach server", opencodeAttach),
		widget.NewFormItem("Gateway API key", gatewayKey),
		widget.NewFormItem("Default model", opencodeModel),
		widget.NewFormItem("Allowed models", allowedModels),
		widget.NewFormItem("Client Base URL", openaiURL),
	)

	content := container.NewBorder(
		container.NewVBox(
			widget.NewLabel("Native desktop manager for a local opencode → OpenAI gateway."),
			widget.NewLabel("Use the Base URL, API key, and opencode model in Cherry Studio / Cursor / NextChat."),
		),
		container.NewHBox(startButton, stopButton),
		nil,
		nil,
		container.NewVBox(form, widget.NewLabel("Status"), status),
	)

	w.SetContent(content)
	w.SetCloseIntercept(func() {
		mu.Lock()
		defer mu.Unlock()
		if running && gateway != nil {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			_ = gateway.Stop(ctx)
		}
		w.Close()
	})
	w.ShowAndRun()
}
