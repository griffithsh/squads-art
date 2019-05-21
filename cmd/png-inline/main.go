package main

import (
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func fmtBytes(b []byte) string {
	var result string
	for _, b := range b {
		result += fmt.Sprintf("0x%x,", b)
	}
	return result
}

func tplFile(filename string) string {
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	if len(b) == 0 {
		panic("no bytes in file " + filename)
	}

	return fmt.Sprintf(`{name: "%s", b: []byte{ %s }},`, filepath.Base(filename), fmtBytes(b))

}

func tplGo(packageName string, filenames []string) string {
	tpl := `// This file was autogenerated; DO NOT EDIT
package %s

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
)

// Images stores decoded Images keyed by their filename as inline compiled resources.
var Images map[string]image.Image

func init() {
	Images = map[string]image.Image{}

	for _, file := range %s {
		decoded, err := png.Decode(bytes.NewReader(file.b))
		if err != nil {
			panic(fmt.Sprintf("png.Decode %s: %s", file.name, err))
		}
		Images[file.name] = decoded
	}
}
`

	return fmt.Sprintf(tpl, packageName, tplFiles(filenames), "%s", "%v")
}

func tplFiles(filenames []string) string {
	tpl := `[]struct {
		name string
		b    []byte
	}{%s
	}`
	var data string
	for _, f := range filenames {
		data += "\n" + tplFile(f)
	}
	return fmt.Sprintf(tpl, data)
}

func filesInDir(dir string) []string {
	result := []string{}

	infos, err := ioutil.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	for _, info := range infos {
		if strings.ToLower(filepath.Ext(info.Name())) != ".png" {
			continue
		}
		absDir, err := filepath.Abs(dir)
		if err != nil {
			panic(err)
		}
		result = append(result, path.Join(absDir, info.Name()))
	}
	return result
}

func main() {
	inFlag := flag.String("i", "./", "i, in configures a directory of pngs to inline")
	outDirFlag := flag.String("o", "./static/images.go", "o, output configures a file to write the inlined png data to")
	packageFlag := flag.String("p", "static", "p, package customises the package name that the inlined image file uses")
	flag.Parse()

	inDir := *inFlag
	outDir := *outDirFlag
	outPackage := *packageFlag

	filenames := filesInDir(inDir)

	data := tplGo(outPackage, filenames)

	b, err := format.Source([]byte(data))
	if err != nil {
		panic(fmt.Errorf("format.Source: %s", err))
	}
	if err := os.MkdirAll(filepath.Dir(outDir), 0755); err != nil {
		panic(err)
	}
	if err := ioutil.WriteFile(path.Clean(outDir), b, 0644); err != nil {
		panic(err)
	}
}
