package ccw

import (
	"fmt"
	"log"
	"os"

	hashdir "github.com/sger/go-hashdir"
	"github.com/eyedeekay/lorca"
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
			extensionArgs = append(extensionArgs, "--load-extension="+extension)
		} else {
			log.Println("extension load warning,", err)
		}
	}
	args = append(args, extensionArgs...)
	return BasicChromium(userdir, private, args...)
}

func SecureExtendedChromium(userdir string, private bool, extensiondirs, extensionhashes []string, args ...string) (lorca.UI, error) {
	var extensionArgs []string
	for index, extension := range extensiondirs {
		if _, err := os.Stat(extension); err != nil {
			return nil, err
		}
		if hash, err := hashdir.Create(extension, "md5"); err != nil {
			if extensionhashes[index] == hash {
				extensionArgs = append(extensionArgs, "--load-extension="+extension)
			} else {
				return nil, fmt.Errorf("hash mismatch error on extension", extension, hash)
			}
		} else {
			return nil, fmt.Errorf("hash calculation error on extension", extension, err)
		}
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
