package ccw

import (
	"os"

	"github.com/zserge/lorca"
)

func BasicChromium(userdir string, private bool, args ...string) (lorca.UI, error) {
	if !private {
		os.MkdirAll(userdir, os.ModePerm)
	} else {
		add := true
		for _, arg := range args {
			if arg == "--incognito" {
				add = false
			}
		}
		if add {
			args = append(args, "--incognito")
		}
	}
	return lorca.New("", userdir, 0, 0, args...)
}

func ExtendedChromium(userdir string, private bool, extensiondirs []string, args ...string) (lorca.UI, error) {
	var extensionArgs []string
	for _, extension := range extensiondirs {
		extensionArgs = append(extensionArgs, "--load-extension")
		extensionArgs = append(extensionArgs, extension)
	}
	args = append(args, extensionArgs...)
	return BasicChromium(userdir, private, args...)
}

var CHROMIUM, ERROR = BasicChromium("basic", true, "--headless")

func Run() error {
	if ERROR != nil {
		return ERROR
	}
	defer CHROMIUM.Close()
	<-CHROMIUM.Done()
	return nil
}
