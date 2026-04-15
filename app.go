package main

import (
"context"
_ "embed"
"encoding/json"
"os"
"path/filepath"

"FasterRedux/pkg/launcher"
"github.com/getlantern/systray"
"github.com/wailsapp/wails/v2/pkg/runtime"
"golang.org/x/sys/windows/registry"
)

//go:embed build/windows/icon.ico
var trayIcon []byte

type Config struct {
	GtaPath        string   `json:"gta_path"`
	ReduxFolders   []string `json:"redux_folders"`
	ActiveRedux    string   `json:"active_redux"`
	AutoInject     bool     `json:"auto_inject"` // start monitoring ON launch   
	RunOnStartup   bool     `json:"run_on_startup"`
	StartWindowBox bool     `json:"start_window_box"` // false = start silently
}

type App struct {
ctx        context.Context
Config     Config
configPath string
launcher   *launcher.Launcher
}

func NewApp() *App {
cfgDir, _ := os.UserConfigDir()
appDir := filepath.Join(cfgDir, "FasterRedux")
os.MkdirAll(appDir, 0755)

return &App{
configPath: filepath.Join(appDir, "config.json"),
Config: Config{
ReduxFolders:   []string{},
AutoInject:     true, // defaults to true a.k.a active
RunOnStartup:   false,
StartWindowBox: true,
},
}
}

func (a *App) startup(ctx context.Context) {
a.ctx = ctx
a.loadConfig()

a.launcher = launcher.NewLauncher(func(msg string) {
runtime.EventsEmit(a.ctx, "status_update", msg)
})

if a.Config.AutoInject && a.Config.GtaPath != "" && a.Config.ActiveRedux != "" {
a.launcher.StartMonitor(a.Config.GtaPath, a.Config.ActiveRedux)
}

if !a.Config.StartWindowBox {
runtime.WindowHide(a.ctx)
}

go systray.Run(a.onTrayReady, func() {})
}

func (a *App) onTrayReady() {
systray.SetIcon(trayIcon)
systray.SetTitle("FasterRedux")
systray.SetTooltip("FasterRedux работает в фоне")

mShow := systray.AddMenuItem("Открыть интерфейс", "Показать главное окно")
systray.AddSeparator()
mQuit := systray.AddMenuItem("Полностью закрыть", "Остановить программу")

for {
select {
case <-mShow.ClickedCh:
runtime.WindowShow(a.ctx)
case <-mQuit.ClickedCh:
systray.Quit()
runtime.Quit(a.ctx)
return
}
}
}

func (a *App) shutdown(ctx context.Context) {
if a.launcher != nil {
a.launcher.StopMonitor()
}
systray.Quit()
}

func (a *App) loadConfig() {
b, err := os.ReadFile(a.configPath)
if err == nil {
json.Unmarshal(b, &a.Config)
}
}

func (a *App) saveConfig() {
b, _ := json.MarshalIndent(a.Config, "", "  ")
os.WriteFile(a.configPath, b, 0644)
}

func (a *App) GetConfig() Config {
return a.Config
}

func (a *App) SelectGtaFolder() string {
path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
Title: "Выберите папку с GTA 5",
})
if err == nil && path != "" {
a.SetGtaPath(path)
return path
}
return ""
}

func (a *App) SelectReduxFolder() string {
path, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
Title: "Выберите папку с Редуксом",
})
if err == nil && path != "" {
a.AddReduxFolder(path)
return path
}
return ""
}

func (a *App) SetGtaPath(p string) {
a.Config.GtaPath = p
a.saveConfig()
a.restartMonitorIfNeeded()
}

func (a *App) AddReduxFolder(p string) {
for _, f := range a.Config.ReduxFolders {
if f == p {
return
}
}
a.Config.ReduxFolders = append(a.Config.ReduxFolders, p)
a.saveConfig()
}

func (a *App) SetActiveRedux(p string) {
a.Config.ActiveRedux = p
a.saveConfig()
a.restartMonitorIfNeeded()
}

func (a *App) RemoveReduxFolder(p string) {
var newList []string
for _, f := range a.Config.ReduxFolders {
if f != p {
newList = append(newList, f)
}
}
a.Config.ReduxFolders = newList
if a.Config.ActiveRedux == p {
a.Config.ActiveRedux = ""
}
a.saveConfig()
a.restartMonitorIfNeeded()
}

func (a *App) ToggleAutoInject(enabled bool) {
a.Config.AutoInject = enabled
a.saveConfig()
if enabled {
a.restartMonitorIfNeeded()
} else {
a.launcher.StopMonitor()
}
}

func (a *App) SetRunOnStartup(enabled bool) {
a.Config.RunOnStartup = enabled
// Force auto-inject ON when enabling Startup so it has an effect
if enabled {
a.Config.AutoInject = true
}
a.saveConfig()

	key, err := registry.OpenKey(registry.CURRENT_USER, `Software\Microsoft\Windows\CurrentVersion\Run`, registry.SET_VALUE)
if err == nil {
defer key.Close()
exe, _ := os.Executable()
if enabled {
key.SetStringValue("FasterRedux", exe)
} else {
key.DeleteValue("FasterRedux")
}
}
a.restartMonitorIfNeeded()
}

func (a *App) SetStartWindowBox(enabled bool) {
a.Config.StartWindowBox = enabled
a.saveConfig()
}

func (a *App) restartMonitorIfNeeded() {
a.launcher.StopMonitor()
if a.Config.AutoInject && a.Config.GtaPath != "" && a.Config.ActiveRedux != "" {
a.launcher.StartMonitor(a.Config.GtaPath, a.Config.ActiveRedux)
}
}
