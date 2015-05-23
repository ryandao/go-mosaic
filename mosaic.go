package mosaic

import (
	"github.com/nfnt/resize"
	"image"
	"image/draw"
	"math"
)

// Wrapper for a tile image with average color attached
type TileImage struct {
	img      image.Image
	avgColor Color
}

// Struct for holding information of an image square.
type ImageSq struct {
	bounds   image.Rectangle
	avgColor Color
}

// A combination of RGB values in floating point numbers.
type Color struct {
	r float64
	g float64
	b float64
}

// Generate a mosaic image based on a target image and a list of tile images.
// `sqSize` specifies the tile size of the output image.
func Mosaic(target image.Image, tiles []image.Image, sqSize int) (image.Image, error) {
	resizedtiles := resizeimages(tiles, sqSize)
	var tileImgs []TileImage
	for _, tile := range resizedtiles {
		tileImgs = append(tileImgs, TileImage{tile, avgColor(tile)})
	}

	squares := colorProfile(target, sqSize)
	x0 := 0
	y0 := 0
	x1 := sqSize * len(squares)
	y1 := sqSize * len(squares)
	dest := image.NewRGBA(image.Rect(x0, y0, x1, y1))

	// Draw the tile image on top of the newly created one
	for y := 0; y < len(squares)-1; y++ {
		for x := 0; x < len(squares)-1; x++ {
			closest := closestTileByColor(squares[y][x].avgColor, tileImgs)
			bounds := image.Rect(x*sqSize, y*sqSize, (x+1)*sqSize, (y+1)*sqSize)
			draw.Draw(dest, bounds, closest.img, image.ZP, draw.Src)
		}
	}

	return dest, nil
}

// Resize a list of images to the given size.
func resizeimages(imgs []image.Image, size int) []image.Image {
	var result []image.Image

	for _, img := range imgs {
		m := resize.Resize(uint(size), uint(size), img, resize.NearestNeighbor)
		result = append(result, m)
	}

	return result
}

// Generate a color profile for an image.
// The profile is a 2-D array of image squares with their average color values.
func colorProfile(img image.Image, sqSize int) [][]ImageSq {
	bounds := img.Bounds()
	numSq := (bounds.Max.X-bounds.Min.X)/sqSize + 1

	result := make([][]ImageSq, numSq)
	for i := range result {
		result[i] = make([]ImageSq, numSq)
	}

	for j := 0; j < numSq; j++ {
		for i := 0; i < numSq; i++ {
			y := j * sqSize
			x := i * sqSize
			var sumR, sumG, sumB, count float64

			for y1 := y; y1 < y+sqSize && y1 < bounds.Max.Y; y1++ {
				for x1 := x; x1 < x+sqSize && x1 < bounds.Max.X; x1++ {
					r, g, b, _ := img.At(x1, y1).RGBA()
					sumR += float64(r)
					sumG += float64(g)
					sumB += float64(b)
					count += 1
				}
			}

			avgColor := Color{sumR / count, sumG / count, sumB / count}
			result[j][i] = ImageSq{image.Rect(x, y, x+sqSize, y+sqSize), avgColor}
		}
	}

	return result
}

// Calculate the distance between 2 colors.
// Similar to the distance between 2 points in 3D space.
func colorDistance(c1 Color, c2 Color) float64 {
	return math.Sqrt((c1.r-c2.r)*(c1.r-c2.r) + (c1.g-c2.g)*(c1.g-c2.g) + (c1.b-c2.b)*(c1.b-c2.b))
}

// Find the tile that has average color closest to a given color.
func closestTileByColor(c Color, tileImgs []TileImage) *TileImage {
	var closest *TileImage
	min := math.MaxFloat64

	for i, _ := range tileImgs {
		dist := colorDistance(c, tileImgs[i].avgColor)
		if dist < min {
			closest = &tileImgs[i]
			min = dist
		}
	}

	return closest
}

// Calculate average RGB color for an image.
func avgColor(m image.Image) Color {
	var sumR, sumG, sumB, count float64
	bounds := m.Bounds()

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := m.At(x, y).RGBA()
			sumR += float64(r)
			sumG += float64(g)
			sumB += float64(b)
			count += 1
		}
	}

	return Color{sumR / count, sumG / count, sumB / count}
}
