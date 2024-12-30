package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gobold"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

const (
	imageWidth  = 1920
	imageHeight = 1080
)

var (
	startDate = time.Date(2025, 1, 6, 0, 0, 0, 0, time.Local)
	endDate   = time.Date(2025, 2, 28, 0, 0, 0, 0, time.Local)
	skipCount = 17 + 1
	skipDate  = []time.Time{
		time.Date(2025, 1, 13, 0, 0, 0, 0, time.Local),
		time.Date(2025, 2, 11, 0, 0, 0, 0, time.Local),
		time.Date(2025, 2, 24, 0, 0, 0, 0, time.Local),
	}
)

func main() {
	os.Exit(run())
}

func run() int {
	for startDate.Before(endDate) {
		if !isSkipDate(&startDate) && startDate.Weekday() != time.Sunday && startDate.Weekday() != time.Saturday {
			virtualCount, actualCount := calcCount(startDate, endDate)

			createImage(startDate.Format("2006-01-02"), virtualCount, actualCount)

			if actualCount <= 0 {
				break
			}
		}

		startDate = startDate.AddDate(0, 0, 1)
	}

	return 0
}

func calcCount(startDate time.Time, endDate time.Time) (int, int) {
	var (
		virtualCount = 0
		actualCount  = -skipCount
	)
	for startDate.Before(endDate) {
		if !isSkipDate(&startDate) && startDate.Weekday() != time.Sunday && startDate.Weekday() != time.Saturday {
			virtualCount++
			actualCount++
		}

		startDate = startDate.AddDate(0, 0, 1)
	}

	return virtualCount, actualCount
}

func isSkipDate(date *time.Time) bool {
	for _, d := range skipDate {
		if d.Equal(*date) {
			return true
		}
	}
	return false
}

func createImage(date string, virtualCount int, actualCount int) {
	img := image.NewRGBA(image.Rect(0, 0, imageWidth, imageHeight))
	draw.Draw(img, img.Bounds(), &image.Uniform{C: color.White}, image.Point{}, draw.Src)

	f, err := opentype.Parse(gobold.TTF)
	if err != nil {
		log.Fatal(err)
	}

	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    120.0,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot:  fixed.Point26_6{X: fixed.I(imageWidth - 480), Y: fixed.I(240)},
	}
	d.DrawString(fmt.Sprintf("%d (%d)", virtualCount, actualCount))

	out, err := os.Create(fmt.Sprintf("out/%s.jpg", date))
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	if err := jpeg.Encode(out, img, nil); err != nil {
		log.Fatal(err)
	}
}
