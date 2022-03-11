package ccw

import (
	"archive/zip"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strings"

	//hashdir "github.com/sger/go-hashdir"

	"github.com/zserge/lorca"
)

func addFiles(w *zip.Writer, basePath, baseInZip string) {
	// Open the Directory
	files, err := ioutil.ReadDir(basePath)
	if err != nil {
		fmt.Println(err)
	}

	for _, file := range files {
		//fmt.Println(basePath + file.Name())
		if !file.IsDir() {
			dat, err := ioutil.ReadFile(filepath.Join(basePath, file.Name()))
			if err != nil {
				fmt.Println(err)
			}

			// Add some files to the archive.
			f, err := w.Create(filepath.Join(baseInZip, file.Name()))
			if err != nil {
				fmt.Println(err)
			}
			_, err = f.Write(dat)
			if err != nil {
				fmt.Println(err)
			}
		} else if file.IsDir() {

			// Recurse
			newBase := filepath.Join(basePath, file.Name(), "/")
			//fmt.Println("Recursing and Adding SubDir: " + file.Name())
			//fmt.Println("Recursing and Adding SubDir: " + newBase)

			addFiles(w, newBase, baseInZip+file.Name()+"/")
		}
	}
}

func ZipAndHashDir(dir string) (string, error) {
	// first, zip the directory at dir
	zipfile := dir + ".zip"
	os.Remove(zipfile)
	// Get a Buffer to Write To
	outFile, err := os.Create(zipfile)
	if err != nil {
		fmt.Println(err)
	}
	defer outFile.Close()

	// Create a new zip archive.
	w := zip.NewWriter(outFile)

	// Add some files to the archive.
	addFiles(w, dir, "")

	if err != nil {
		fmt.Println(err)
	}

	// Make sure to check the error on Close.
	err = w.Close()
	if err != nil {
		fmt.Println(err)
	}
	// hash the resulting .zip file
	hasher := sha256.New()
	s, err := ioutil.ReadFile(zipfile)
	hasher.Write(s)
	if err != nil {
		log.Fatal(err)
	}
	//finally, remove the zip file
	os.Remove(zipfile)

	return hex.EncodeToString(hasher.Sum(nil)), nil

}

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
	var err error
	for index, extension := range extensiondirs {
		if _, err = os.Stat(extension); err != nil {
			return nil, err
		}
		//if hash, ierr := hashdir.Create(extension, "sha256"); ierr == nil {
		if hash, ierr := ZipAndHashDir(extension); ierr == nil {
			log.Println("extension hash check,", index, extension, extensionhashes[index], hash)
			if extensionhashes[index] == hash {
				extensionArgs = append(extensionArgs, extension)
			} else {
				err = fmt.Errorf("extension hash mismatch on %s, expected %s, got %s, %s", extension, extensionhashes[index], hash, err)
				//return nil, fmt.Errorf("hash mismatch error on extension", extension, hash, extensionhashes[index])
			}
		} else {
			err = fmt.Errorf("hash calculation error on extension %s %s %s", extension, ierr, err)
			//return nil,
		}
	}
	if err != nil {
		return nil, err
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
