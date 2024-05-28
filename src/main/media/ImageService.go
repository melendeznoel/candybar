package media

import (
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"log"
	UtilityService "main/helper"
	"math"
	"net/http"
	"strconv"
	"sync"
	"time"
)

// Compare will
func Compare(images []Figure) ImageCompareResult {
	figures := getImages(images)

	if len(figures) == 0 {
		return ImageCompareResult{}
	}

	readImages(figures)

	buildRelatives(figures)

	defineRanks(figures)

	return ImageCompareResult{ID: UtilityService.FetchNewID(), Images: figures}
}

func readImages(figures []Figure) {
	for idx, figure := range figures {
		if figure.image == nil {
			if img, err := UtilityService.ToImage(figure.Data); err == nil {
				figures[idx].image = img

			} else {
				log.Fatal(err)
			}
		}
	}
}

func colorCompare(a, b color.Color) (diff int64) {
	r1, g1, b1, a1 := a.RGBA()
	r2, g2, b2, a2 := b.RGBA()

	diff += int64(math.Abs(float64(r1 - r2)))
	diff += int64(math.Abs(float64(g1 - g2)))
	diff += int64(math.Abs(float64(b1 - b2)))
	diff += int64(math.Abs(float64(a1 - a2)))

	return diff
}

func boundsCompare(a, b image.Rectangle) bool {
	return a.Min.X == b.Min.X && a.Min.Y == b.Min.Y && a.Max.X == b.Max.X && a.Max.Y == b.Max.Y
}

func resizeImages(imgs []image.Image) {
	x := 0
	y := 0

	for idx := range imgs {
		img := imgs[idx]

		imgSize := img.Bounds().Size()

		if imgSize.X < x && imgSize.Y < y {

		}
	}

	/*	m := resize.Resize(100, 100, img, resize.NearestNeighbor)

		out, err := os.Create("test_qsized.jpg")

		if err != nil {
			log.Fatal(err)
		}

		defer out.Close()

		// write new image to file
		errors := jpeg.Encode(out, m, nil)

		if errors != nil {
			fmt.Println("somet thing goin wrong ")
			os.Exit(0)
		}*/
}

// func resizeImage(img image.Image, size, image.Point) image.Image  {

// 	return
// }

// getImages will download the image data
func getImages(figures []Figure) []Figure {
	var wg sync.WaitGroup

	wg.Add(len(figures))

	for idx, figure := range figures {
		go func(i int, figure Figure) {
			defer wg.Done()

			res, err := http.Get(figures[i].URL)

			if err != nil {
				log.Fatal(err)

			} else {
				defer res.Body.Close()

				body, err := ioutil.ReadAll(res.Body)

				if err != nil {
					log.Fatal(err)

				} else {
					figures[i].ID = UtilityService.FetchNewID()
					figures[i].Data = body
					figures[i].Type = res.Header.Get("Content-Type")
					figures[i].image = nil

					lastModified := res.Header.Get("last-modified")

					if lmt, e := time.Parse(time.RFC1123, lastModified); e == nil {
						figures[i].LastModified = lmt

					} else {
						log.Fatal(e)
					}

					dataLength := res.Header.Get("Content-Length")

					if dl, err := strconv.ParseInt(dataLength, 10, 64); err == nil {
						figures[i].DataLength = dl

					} else {
						log.Fatal(err)
					}
				}
			}
		}(idx, figure)
	}

	wg.Wait()

	return figures
}

func buildRelatives(figures []Figure) {
	var wg sync.WaitGroup

	wg.Add(len(figures))

	for i := range figures {
		go func(entityIdx int) {
			defer wg.Done()

			for idx := range figures {
				if figures[idx].ID != figures[entityIdx].ID && !isRelated(figures[idx].ID, figures[entityIdx]) {
					if figures[idx].image == nil {
						if img, err := UtilityService.ToImage(figures[idx].Data); err == nil {
							figures[idx].image = img

						} else {
							log.Fatal(err)
						}
					}

					// var akinResult int64
					// var resizedImage image.Image

					if boundsCompare(figures[entityIdx].image.Bounds(), figures[idx].image.Bounds()) {
						akinResult := akin(figures[entityIdx], figures[idx])

						figures[entityIdx].Relatives = append(figures[entityIdx].Relatives, Relative{ID: figures[idx].ID, Percentage: akinResult, Rank: 2})
					} else {

						if figures[idx].image.Bounds().In(figures[entityIdx].image.Bounds()) {
							// resizedImage = resizeImage(figures[entityIdx].image.Bounds(), figures[entityIdx].image.Bounds().Size())

						} else {

							// resizedImage = resizeImage(figures[entityIdx].image.Bounds(), figures[entityIdx].image.Bounds().Size())

						}
					}
				}
			}
		}(i)
	}

	wg.Wait()
}

func defineRanks(figures []Figure) {
	// todo:  determine rank by sorting the percentage on the relatives
}

func isRelated(id string, figure Figure) bool {
	for _, r := range figure.Relatives {
		if r.ID == id {
			return true
		}
	}

	return false
}

func akin(targetFigure, sourceFigure Figure) (diffs int64) {
	for y := targetFigure.image.Bounds().Min.Y; y < targetFigure.image.Bounds().Max.Y; y++ {
		for x := targetFigure.image.Bounds().Min.X; x < targetFigure.image.Bounds().Max.X; x++ {
			diffs += colorCompare(targetFigure.image.Bounds().At(x, y), sourceFigure.image.Bounds().At(x, y))
		}
	}

	return diffs
}

func FastCompare(img1, img2 *image.RGBA) (int64, error) {
	if img1.Bounds() != img2.Bounds() {
		return 0, fmt.Errorf("image bounds not equal: %+v, %+v", img1.Bounds(), img2.Bounds())
	}

	accumError := int64(0)

	for i := 0; i < len(img1.Pix); i++ {
		accumError += int64(sqDiffUInt8(img1.Pix[i], img2.Pix[i]))
	}

	return int64(math.Sqrt(float64(accumError))), nil
}

func sqDiffUInt8(x, y uint8) uint64 {
	d := uint64(x) - uint64(y)
	return d * d
}
