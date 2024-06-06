package main

import (
	"fmt"
	"image"
	_ "image/color"
	_ "image/jpeg"
	"image/png"
	"os"

	"github.com/nfnt/resize"
)

// loadImage loads an image from a file and decodes it.
func loadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

// saveImage saves an image to a file in PNG format.
func saveImage(img image.Image, filename string) error {
	outFile, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer outFile.Close()

	err = png.Encode(outFile, img) // Use png.Encode to preserve transparency
	if err != nil {
		return err
	}
	return nil
}

// isTransparent checks if the pixel at (x, y) in the image is fully transparent.
func isTransparent(img image.Image, x, y int) bool {
	switch img := img.(type) {
	case *image.NRGBA:
		return img.NRGBAAt(x, y).A == 0
	case *image.RGBA:
		return img.RGBAAt(x, y).A == 0
	default:
		return false
	}
}

// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	} else {
// 		return b
// 	}
// }
// func min(a, b int) int {
// 	if a > b {
// 		return b
// 	} else {
// 		return a
// 	}
// }

func ImageResizing(img image.Image, xScale uint, yScale uint) image.Image {
	//0 in x or y scale causes fallback based on aspect ratio
	resizedImg := resize.Resize(uint(img.Bounds().Dx())*xScale, uint(img.Bounds().Dy())*yScale, img, resize.Lanczos2)

	out, err := os.Create("output.png")
	if err != nil {
		panic(err)
	}
	out.Close()
	return resizedImg
}

func main() {
	// Load the foreground image
	foregroundImage, err := loadImage("google.png")

	if err != nil {
		fmt.Printf("Failed to load image: %v\n", err)
		return
	}
	backgroundImage, err := loadImage("background.jpg")
	if err != nil {
		panic(err)
	}
	//resize to make foreground compatible with background
	foregroundImageHeight := foregroundImage.Bounds().Dy()
	// foregroundImageWidth := foregroundImage.Bounds().Dx()
	backgroundImageHieght := backgroundImage.Bounds().Dy()
	backgroundImageWidth := backgroundImage.Bounds().Dx()

	xScale := uint(backgroundImageHieght / foregroundImageHeight)

	resizedForeground := ImageResizing(foregroundImage, xScale, 0)
	// Create a new image with the same bounds as the foreground image
	// foregroundImageRect := resizedForeground.Bounds()
	// backgroundImageRect := backgroundImage.Bounds()
	// x0 := min(foregroundImageRect.Min.X, backgroundImageRect.Min.X)
	// x1 := max(foregroundImageRect.Max.X, backgroundImageRect.Max.X)
	// y0 := min(foregroundImage.Bounds().Min.Y, backgroundImage.Bounds().Min.Y)
	// y1 := max(foregroundImage.Bounds().Max.Y, backgroundImage.Bounds().Max.Y)
	newImage := image.NewNRGBA(image.Rect(0, 0, backgroundImageWidth, backgroundImageHieght))

	// Iterate over each pixel in the foreground image
	for x := 0; x < newImage.Rect.Dx(); x++ {
		for y := 0; y < newImage.Rect.Dx(); y++ {
			if isTransparent(resizedForeground, x, y) {
				//fallback on background pixel
				newImage.Set(x, y, backgroundImage.At(x, y))
			} else {
				// Copy the pixel from the foreground image
				newImage.Set(x, y, resizedForeground.At(x, y))
			}
		}
	}

	// Save the new image
	err = saveImage(newImage, "Merged.png")
	if err != nil {
		fmt.Printf("Failed to save image: %v\n", err)
		return
	}

	// defer out.Close()

	// jpeg.Encode(out, resizedImg,nil)

	fmt.Println("Image processing completed successfully")
}
