package main

import (
	"github.com/oneworkspacestudio/my-news/internal/news"
	"github.com/oneworkspacestudio/my-news/internal/webui"
	"github.com/oneworkspacestudio/one-workspace/pkg/apphost"
	"os"
	"path/filepath"
	"strings"
)

var version = "0.1.0-preview.1"

func main() {
	hasData := false
	for _, v := range os.Args[1:] {
		if v == "--data" || strings.HasPrefix(v, "--data=") {
			hasData = true
		}
	}
	if !hasData {
		root, e := os.UserConfigDir()
		if e != nil {
			apphost.Report(e)
			os.Exit(1)
		}
		os.Args = append(os.Args, "--data", filepath.Join(root, "OneWorkspace", "MyNews"))
	}
	if e := apphost.Main(apphost.Config{ID: "news", Name: "My News", Version: version, Listen: "127.0.0.1:9030"}, webui.Files(), func(h *apphost.Host) (apphost.Application, error) { return news.New(h) }); e != nil {
		apphost.Report(e)
		os.Exit(1)
	}
}
