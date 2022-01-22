package sqlite

import goutils "github.com/simonski/goutils"

func NewApp() *App {
	return &App{}
}

type App struct {
}

func (a *App) HandleInput(command string, cli *goutils.CLI) {

}
