package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"strings"

	docopt "github.com/docopt/docopt-go"
	"github.com/gen2brain/go-fitz"
	"github.com/nfnt/resize"
	"github.com/schollz/progressbar"
	"github.com/ungerik/go-cairo"
)

// WIDTH is based on A4 ratio
var WIDTH = 210 * 4

// HEIGHT is based on A4 ratio
var HEIGHT = 297 * 4

func main() {
	usage := `Desc.
	Usage:
	  pngfy convert DIRECTORY [--target=TARGET][--width=WIDTH][--height=HEIGHT]
	  pngfy FILE [--target=TARGET][--width=WIDTH][--height=HEIGHT]
	  pngfy -h | --help
	Arguments:
		DIRECTORY         	Directory containing the pdf files to be converted
		FILE			  	Single pdf file to be converted
	Options:
	  -h --help                     	Show this screen.
	  --target=TARGET					Target directory for results.
	  --width=WIDTH						Width for rescaling the images.
	  --height=HEIGHT					Height for rescaling the images.`

	arguments, _ := docopt.ParseArgs(usage, nil, "1.0")

	// operators and parameters
	targetPath, _ := arguments["DIRECTORY"].(string)
	targetFile, _ := arguments["FILE"].(string)
	target, _ := arguments["--target"].(string)
	_width, _ := arguments["--width"].(string)
	_height, _ := arguments["--height"].(string)
	convert := arguments["convert"].(bool)
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
	targetDir := path.Join(dir, "output")
	if target != "" && target != "." {
		fmt.Println(target)
		targetDir = target
	}

	if !convert {
		convertPages(targetFile, targetDir, uint(width), uint(height))
		os.Exit(0)
	}

	files := getFiles(targetPath)
	fmt.Printf("Saving images to: %s\n", targetDir)
	fmt.Printf("Processing %d files\n", len(files))
	bar := progressbar.New(len(files))
	for _, file := range files {
		convertPages(file, targetDir, uint(width), uint(height))
		bar.Add(1)
	}
	fmt.Println()
}

func convertPages(file, targetDir string, width, height uint) {
	dir, fName := path.Split(file)
	parentDir := strings.Split(strings.Trim(dir, string(os.PathSeparator)), string(os.PathSeparator))
	parentDirName := parentDir[len(parentDir)-1]
	targetFileDir := path.Join(targetDir)
	os.MkdirAll(targetFileDir, 0770)
	pages := pdf2Surface(file, width, height)
	for n, page := range pages {
		// fmt.Println(fmt.Sprintf("%s/%s_%05d.png", targetFileDir, parentDirName, n))
		page.WriteToPNG(fmt.Sprintf("%s/%s_%s_%05d.png", targetFileDir, parentDirName, fName, n))
	}
}

func getFiles(filePath string) []string {

	var paths []string
	files, err := ioutil.ReadDir(filePath)
	if err != nil {
		panic(err)
	}

	for _, f := range files {
		if path.Ext(f.Name()) == ".pdf" {
			paths = append(paths, path.Join(filePath, f.Name()))
		}
	}
	return paths
}

func pdf2Surface(path string, width, height uint) []*cairo.Surface {
	doc, err := fitz.New(path)
	if err != nil {
		fmt.Println(path, err)
	}
	var pages = make([]*cairo.Surface, doc.NumPage())
	for n := 0; n < doc.NumPage(); n++ {
		img, err := doc.Image(n)
		if err != nil {
			fmt.Println(path, n+1, err)
		}

		resized := resize.Resize(width, height, img, resize.Lanczos2)
		surface := cairo.NewSurfaceFromImage(resized)
		pages[n] = surface
	}
	return pages
}
