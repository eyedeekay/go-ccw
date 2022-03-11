package ccw

import (
	"fmt"
	"log"
	"os"
	"strings"

	hashdir "github.com/sger/go-hashdir"
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
	log.Println(args)
	return lorca.NewChromium("", userdir, 800, 600, args...)
}

func ExtendedChromium(userdir string, private bool, extensiondirs []string, args ...string) (lorca.UI, error) {
	var extensionArgs []string
	for _, extension := range extensiondirs {
		if _, err := os.Stat(extension); err == nil {
			extensionArgs = append(extensionArgs, extension)
		} else {
			log.Println("extension load warning,", err)
		}
	}
	finalExtensionArgs := "--load-extension=" + strings.Join(extensionArgs, ",")
	args = append(args, finalExtensionArgs)
	return BasicChromium(userdir, private, args...)
}

func SecureExtendedChromium(userdir string, private bool, extensiondirs, extensionhashes []string, args ...string) (lorca.UI, error) {
	var extensionArgs []string
	for index, extension := range extensiondirs {
		if _, err := os.Stat(extension); err != nil {
			return nil, err
		}
		if hash, err := hashdir.Create(extension, "sha256"); err == nil {
			if extensionhashes[index] == hash {
				extensionArgs = append(extensionArgs, extension)
			} else {
				return nil, fmt.Errorf("hash mismatch error on extension", extension, hash, extensionhashes[index])
			}
		} else {
			return nil, fmt.Errorf("hash calculation error on extension", extension, err)
		}
	}
	finalExtensionArgs := "--load-extension=" + strings.Join(extensionArgs, ",")
	args = append(args, finalExtensionArgs)
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
