package main

import (
	_ "embed"
	"os"
	"strings"

	"github.com/getlantern/systray"

	"os/exec"
	"regexp"
)

//go:embed icon.png
var icon []byte
var listExp = regexp.MustCompile(`Account:\s(.*)\n`)
var accOutput = regexp.MustCompile(`(\d+)`)

func main() {
	systray.Run(onReady, onExit)
}

func list_accounts() []string {
	cmd := exec.Command("cloak", "list")
	out, err := cmd.Output()
	if err != nil {
		return nil
	}

	matches := listExp.FindAllStringSubmatch(string(out), -1)
	accounts := make([]string, len(matches))
	for i, match := range matches {
		println("\"", match[1], "\"")
		accounts[i] = match[1]
	}
	return accounts
}

func onReady() {
	systray.SetIcon(icon)
	systray.SetTitle("cloak Tray")
	systray.SetTooltip("cloak Tray")

	accounts := list_accounts()
	for _, account := range accounts {
		item := systray.AddMenuItem(account, account)
		go func() {
			for {
				<-item.ClickedCh
				out, err := exec.Command("cloak", "view", account).Output()
				if err != nil {
					panic("cloak view failed")
				}

				out = accOutput.Find(out)

				sessionType, _ := os.LookupEnv("XDG_SESSION_TYPE")
				if sessionType == "wayland" {
					println("wl-copy")
					exec.Command("wl-copy", "-n", string(out)).Run()
				} else {
					command := exec.Command("xclip", "-selection", "clipboard")
					command.Stdin = strings.NewReader(string(out))
					command.Run()
				}
			}
		}()
	}

	mQuit := systray.AddMenuItem("Quit", "Quit the app")
	go func() {
		<-mQuit.ClickedCh
		systray.Quit()
	}()
}

func onExit() {

}
