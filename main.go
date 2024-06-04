package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"

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


func main() {
	// Load the foreground image
	foregroundImage, err := loadImage("google.png")
	if err != nil {
		fmt.Printf("Failed to load image: %v\n", err)
		return
	}

	// Create a new image with the same bounds as the foreground image
	newImage := image.NewNRGBA(foregroundImage.Bounds())

	// Iterate over each pixel in the foreground image
	for x := 0; x < foregroundImage.Bounds().Dx(); x++ {
		for y := 0; y < foregroundImage.Bounds().Dy(); y++ {
			if isTransparent(foregroundImage, x, y) {
				//fallback on background pixel
				newImage.Set(x, y, color.RGBA{R: 255, G: 255, B: 255, A: 255})
			} else {
				// Copy the pixel from the foreground image
				newImage.Set(x, y, foregroundImage.At(x, y))
			} 
		}
	}

	// Save the new image
	err = saveImage(newImage, "Merged.png")
	if err != nil {
		fmt.Printf("Failed to save image: %v\n", err)
		return
	}

	fmt.Println("Image processing completed successfully")
}
