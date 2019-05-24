package main

import (
	"encoding/json"
	"flag"
	"github.com/fogleman/gg"
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"image/jpeg"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"runtime/pprof"
	"strconv"
	"time"
)

type InputJsonData struct {
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Format   string `json:"format"`
	FrameNum []int  `json:"timestamps"`
}

type ImageWithPosition struct {
	img    image.Image
	row    int
	column int
}

// error handler
func check(err error) {
	if err != nil {
		panic(err)
	}
}

func generateFinImage(thumbnailWidth, thumbnailHeight, columnLength, rowWidth int) *image.RGBA {
	finRect := image.Rectangle{
		Min: image.ZP,
		Max: image.Point{X: thumbnailWidth * rowWidth, Y: thumbnailHeight * columnLength},
	}
	return image.NewRGBA(finRect)
}

func resizeImg(img image.Image, thumbnailSize uint) image.Image {
	rezImg := resize.Thumbnail(thumbnailSize, thumbnailSize, img, resize.NearestNeighbor)
	return rezImg
}

func getDataFromJsonFile(path string) InputJsonData {
	jsonData, err := ioutil.ReadFile(path)
	check(err)
	var data InputJsonData
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		panic(err)
	}
	return data
}

func divmod(numerator, denominator int) (quotient, remainder int) {
	quotient, remainder = numerator/denominator, numerator%denominator
	return
}

func getThumbnailWidthHeight(imgPath string, thumbnailSize uint) (thumbnailWidth, thumbnailHeight int) {
	img, err := gg.LoadJPG(imgPath)
	check(err)
	img = resizeImg(img, thumbnailSize)
	thumbnailWidth = img.Bounds().Dx()
	thumbnailHeight = img.Bounds().Dy()
	return
}

func handleSingleImage(imgPath string, thumbnailSize uint, row, column int) (imgp ImageWithPosition) {
	img, err := gg.LoadJPG(imgPath)
	check(err)
	img = resizeImg(img, thumbnailSize)
	thumbnailWidth := img.Bounds().Dx()
	thumbnailHeight := img.Bounds().Dy()
	imgp = ImageWithPosition{img, row * thumbnailWidth, column * thumbnailHeight}
	return
}

func generateJPG(imageFolderPath, jsonPath string, sampleInterval, columnLength int, thumbnailSize uint) *image.RGBA {
	jsonData := getDataFromJsonFile(jsonPath)
	frames := jsonData.FrameNum
	imgNum := len(frames)
	thumbnailWidth, thumbnailHeight := getThumbnailWidthHeight(path.Join(imageFolderPath, "0.jpeg"), thumbnailSize)
	row, column := 0, 0

	c := make(chan ImageWithPosition, imgNum/sampleInterval)

	for i := 0; i <= imgNum-sampleInterval; i += sampleInterval {
		row, column = divmod(i/sampleInterval, columnLength)
		go func(row, column, i int) {
			c <- handleSingleImage(path.Join(imageFolderPath, strconv.Itoa(i)+".jpeg"), thumbnailSize,
				row, column)
		}(row, column, i)
	}

	finImage := generateFinImage(thumbnailWidth, thumbnailHeight, columnLength, row+1)

	for i := 0; i < imgNum/sampleInterval; i++ {
		p := <-c
		rect := image.Rectangle{Min: image.Point{X: p.row, Y: p.column}, Max: image.Point{X: p.row, Y: p.column}.Add(p.img.Bounds().Size())}
		draw.Draw(finImage, rect, p.img, image.ZP, draw.Src)
	}

	return finImage
}

var cpuprofile = flag.String("cpuprofile", "cpu.prof", "write cpu profile to file.")
var memprofile = flag.String("memprofile", "mem.prof", "write memory profile to file")

func main() {

	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal("could not create CPU profile: ", err)
		}
		defer f.Close()
		if err := pprof.StartCPUProfile(f); err != nil {
			log.Fatal("could not start CPU profile: ", err)
		}
		defer pprof.StopCPUProfile()
	}

	start := time.Now()
	finImage := generateJPG("test_files",
		"test.json",
		20,
		25,
		100)
	out, err := os.Create("./JPG.jpg")
	check(err)
	err = jpeg.Encode(out, finImage, &jpeg.Options{Quality: 75})
	check(err)
	print("time consumption: ", time.Now().Sub(start).Seconds(), "\n")

	if *memprofile != "" {
		f, err := os.Create(*memprofile)
		if err != nil {
			log.Fatal("could not create memory profile: ", err)
		}
		defer f.Close()
		runtime.GC() // get up-to-date statistics
		if err := pprof.WriteHeapProfile(f); err != nil {
			log.Fatal("could not write memory profile: ", err)
		}
	}
}
