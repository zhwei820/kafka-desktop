package main

import (
	"context"
	"fmt"
	"kafka_desktop/backend/conf"
	"os"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/zhwei820/gconv"
)

// App struct
type App struct {
	ctx  context.Context
	conf *conf.Config
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in startup", r)
			wd, _ := os.Getwd()
			runtime.MessageDialog(ctx, runtime.MessageDialogOptions{
				Type:          "error",
				Title:         wd,
				Message:       gconv.String(r),
				Buttons:       nil,
				DefaultButton: "Confirm",
				CancelButton:  "Cancel",
			})
		}
	}()

	a.ctx = ctx
	a.conf = conf.New()
	a.conf.LoadAll()
	if len(a.conf.Instance) == 0 {
		a.conf.Instance["local"] = &conf.LocalKafkaConfig
		a.conf.Configs().ConfigNames = map[string]bool{"local": true}
		a.conf.SaveAll()
	}

	runtime.WindowMaximise(ctx)
}

func (a *App) shutdown(ctx context.Context) {
	a.conf.SaveAll()
	time.Sleep(500 * time.Millisecond)
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
