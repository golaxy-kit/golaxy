package launcher

import (
	"github.com/alecthomas/kingpin/v2"
	"kit.golaxy.org/golaxy/util/option"
)

// NewApp 创建应用
func NewApp(settings ...option.Setting[AppOptions]) App {
	app := &_App{
		options: option.Make(Option{}.Default(), settings...),
	}
	return app
}

// App 应用
type App interface {
	// Run 运行
	Run()
}

type _App struct {
	options AppOptions
}

// Run 运行
func (app *_App) Run() {
	var ptPath = kingpin.Flag("pt", "服务原型配置文件（*.json|*.xml）。").Default("./pt.json").ExistingFile()

	var runApp = kingpin.Command("run", "开始运行。").Alias("r").Default()
	var services = runApp.Flag("services", "需要启动的服务列表。").Strings()

	var printInfo = kingpin.Command("print", "打印信息。").Alias("p")
	var printPt = printInfo.Command("pt", "打印所有服务原型。")
	var printComp = printInfo.Command("comp", "打印所有组件。")

	cmd := kingpin.Parse()

	switch cmd {
	case runApp.FullCommand():
		app.runApp(*services, *ptPath)
		return
	case printInfo.FullCommand():
		return
	case printComp.FullCommand():
		app.printComp()
		return
	case printPt.FullCommand():
		app.printPt(*ptPath)
		return
	}

	for _, customCmd := range app.options.Commands {
		if cmd == customCmd.Clause.FullCommand() {
			customCmd.Run()
			return
		}
	}

	kingpin.Usage()
}
