package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	docopt "github.com/docopt/docopt-go"
	fitz "github.com/gen2brain/go-fitz"
	"github.com/nfnt/resize"
	"github.com/schollz/progressbar"
	cairo "github.com/ungerik/go-cairo"
)

var WIDTH = 300
var HEIGHT = 450

func main() {
	usage := `Desc.
	Usage:
	  pngfy DIRECTORY [--target=TARGET][--width=WIDTH][--height=HEIGHT]
	  pngfy -h | --help
	Arguments:
		DIRECTORY         Directory containing the pdf files to be converted
	Options:
	  -h --help                     	Show this screen.
	  --target=TARGET					Target directory for results.
	  --width=WIDTH						Width for rescaling the images.
	  --height=HEIGHT					Height for rescaling the images.`

	arguments, _ := docopt.ParseArgs(usage, nil, "1.0")

	// operators and parameters
	targetPath, _ := arguments["DIRECTORY"].(string)
	target, _ := arguments["--target"].(string)
	_width, _ := arguments["--width"].(string)
	_height, _ := arguments["--height"].(string)

	width, err := strconv.Atoi(_width)
	if err != nil {
		fmt.Printf("Invalid parameter %v for width. Expected int", width)
	}
	height, err := strconv.Atoi(_height)
	if err != nil {
		fmt.Printf("Invalid parameter %v for height. Expected int", height)
	}

	// default values
	if width == 0 {
		width = WIDTH
	}
	if height == 0 {
		height = HEIGHT
	}

	targetPath = strings.TrimRight(targetPath, string(os.PathSeparator))

	dir, _ := path.Split(targetPath)
	targetDir := path.Join(dir, "target")

	if target != "" {
		targetDir = target
	}

	files := getFiles(targetPath)
	fmt.Printf("Saving images to: %s\n", targetDir)
	fmt.Printf("Processing %d files\n", len(files))
	bar := progressbar.New(len(files))
	for _, file := range files {
		_, fName := path.Split(file)
		targetFileDir := path.Join(targetDir, strings.Split(fName, ".")[0])
		os.MkdirAll(targetFileDir, 0770)
		pages := GetPdfBytes(file, true, uint(width), uint(height))
		for n, page := range pages {
			page.WriteToPNG(fmt.Sprintf("%s/%d.png", targetFileDir, n))
		}

		bar.Add(1)
	}
	fmt.Println()
}

func getFiles(filePath string) []string {

	var paths []string
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		log.Fatal(err)
	}

	for _, f := range files {
		if path.Ext(f.Name()) == ".pdf" {
			paths = append(paths, path.Join(filePath, f.Name()))
		}
	}
	return paths
}

func GetPdfBytes(path string, width, height uint) []*cairo.Surface {
	doc, err := fitz.New(path)
	if err != nil {
		log.Fatal(err)
	}
	var pages = make([]*cairo.Surface, doc.NumPage())
	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.Image(n)
		if err != nil {
			log.Fatal(err)
		}

		resized := resize.Resize(width, height, img, resize.Lanczos2)
		surface := cairo.NewSurfaceFromImage(resized)
		pages[n] = surface
	}
	return pages
}
