package launcher

import (
"context"
"io/fs"
"log"
"os"
"path/filepath"
"strings"
"time"

"github.com/shirou/gopsutil/v3/process"
)

type Launcher struct {
GtaPath      string
ReduxPath    string
IsMonitoring bool
cancel       context.CancelFunc
onEvent      func(string)
}

func NewLauncher(onEvent func(string)) *Launcher {
return &Launcher{
onEvent: onEvent,
}
}

func (l *Launcher) StartMonitor(gtaPath, reduxPath string) {
if l.IsMonitoring {
return
}
l.GtaPath = gtaPath
l.ReduxPath = reduxPath
l.IsMonitoring = true

ctx, cancel := context.WithCancel(context.Background())
l.cancel = cancel

go func() {
l.onEvent("Ожидание запуска GTA 5...")
injected := false

for {
select {
case <-ctx.Done():
l.IsMonitoring = false
return
default:
time.Sleep(2 * time.Second)

if l.GtaPath == "" || l.ReduxPath == "" {
continue
}

isRunning := l.isGTA5Running()

if isRunning && !injected {
l.onEvent("GTA 5 обнаружена! Подмена файлов...")
err := l.SmartInject()
if err != nil {
l.onEvent("Ошибка подмены: " + err.Error())
} else {
l.onEvent("Редукс успешно применён!")
injected = true
}
} else if !isRunning && injected {
l.onEvent("GTA 5 закрыта. Восстановление файлов...")
l.Restore()
injected = false
l.onEvent("Ожидание запуска GTA 5...")
}
}
}
}()
}

func (l *Launcher) StopMonitor() {
if l.cancel != nil {
l.cancel()
}
if l.IsMonitoring {
l.Restore()
}
l.IsMonitoring = false
l.onEvent("Мониторинг остановлен")
}

func (l *Launcher) isGTA5Running() bool {
procs, err := process.Processes()
if err != nil {
return false
}
for _, p := range procs {
name, err := p.Name()
if err == nil && strings.ToLower(name) == "gta5.exe" {
return true
}
}
return false
}

func (l *Launcher) SmartInject() error {
log.Println("Injecting:", l.ReduxPath, "into", l.GtaPath)
return filepath.WalkDir(l.ReduxPath, func(path string, d fs.DirEntry, err error) error {
if err != nil || d.IsDir() {
return nil
}

ext := strings.ToLower(filepath.Ext(path))
if ext != ".rpf" && ext != ".asi" {
return nil
}

relPath, err := filepath.Rel(l.ReduxPath, path)
if err != nil {
return nil
}

targetPath := filepath.Join(l.GtaPath, relPath)
bakPath := targetPath + ".bak"

if _, err := os.Stat(targetPath); err == nil {
if _, err := os.Stat(bakPath); os.IsNotExist(err) {
os.Rename(targetPath, bakPath)
}
}

os.Remove(targetPath)
fileDir := filepath.Dir(targetPath)
os.MkdirAll(fileDir, 0755)

err = os.Link(path, targetPath)
if err != nil {
log.Println("Failed to link:", err)
}
return nil
})
}

func (l *Launcher) Restore() {
log.Println("Restoring original files...")
filepath.WalkDir(l.GtaPath, func(path string, d fs.DirEntry, err error) error {
if err != nil || d.IsDir() {
return nil
}

if strings.HasSuffix(path, ".bak") {
origPath := strings.TrimSuffix(path, ".bak")
os.Remove(origPath)
os.Rename(path, origPath)
}
return nil
})
}
