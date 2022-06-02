// Package png allows for loading png images and applying
// image flitering effects on them.
package png

import (
	"image"
	"image/color"
	"math"
)

// Grayscale applies a grayscale filtering effect to the image
func (img *ImageTask) Grayscale(minY, maxY int) *image.RGBA64 {

	// Bounds returns defines the dimensions of the image. Always
	// use the bounds Min and Max fields to get out the width
	// and height for the image
	bounds := img.out.Bounds()
	for y := minY; y < maxY; y++ {
	//for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			//Returns the pixel (i.e., RGBA) value at a (x,y) position
			// Note: These get returned as int32 so based on the math you'll
			// be performing you'll need to do a conversion to float64(..)
			r, g, b, a := img.in.At(x, y).RGBA()

			//Note: The values for r,g,b,a for this assignment will range between [0, 65535].
			//For certain computations (i.e., convolution) the values might fall outside this
			// range so you need to clamp them between those values.
			greyC := clamp(float64(r+g+b) / 3)

			//Note: The values need to be stored back as uint16 (I know weird..but there's valid reasons
			// for this that I won't get into right now).
			img.out.Set(x, y, color.RGBA64{greyC, greyC, greyC, uint16(a)})
		}
	}
	return img.out
}


func Convolve2d(kernel []float64) [][]float64 {
	// convert input 1D kernel to 2D and flip it
	var flippedKernel2D [][]float64
	h := int(math.Sqrt(float64(len(kernel))))

	for col := len(kernel); col > 0; col = col - h {
		flippedKernel2D = append(flippedKernel2D, kernel[col-h:col])
	}
	return flippedKernel2D
}

// ref: https://stackoverflow.com/questions/42677196/convolution-in-golang
func (img *ImageTask) ApplyEffect(kernel []float64, minY int, maxY int) {
	bounds := img.out.Bounds()
	var sumR float64
    var sumB float64
    var sumG float64
    var r uint32
    var g uint32
    var b uint32
	var a uint32

	matrice := Convolve2d(kernel)
	for y := minY; y < maxY; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {

            for i := -1; i <= 1; i++ {
                for j := -1; j <= 1; j++ {

                    var imageX int
                    var imageY int

                    imageX = x + i
                    imageY = y + j

                    r, g, b, a = img.in.At(imageX, imageY).RGBA()
					sumR = (sumR + (float64(r) * matrice[i+1][j+1]))
                    sumG = (sumG + (float64(g) * matrice[i+1][j+1]))
                    sumB = (sumB + (float64(b) * matrice[i+1][j+1]))
                }
            }

            img.out.Set(x, y, color.RGBA64{clamp(sumR), clamp(sumG), clamp(sumB), uint16(a)})
			sumR = 0
			sumB = 0
			sumG = 0
        }
    }
}

func (img *ImageTask) ExecuteEffects(effectType string, minY int, maxY int, waitgroup chan bool) {
	var kernel []float64
	if effectType == "G" {
		img.Grayscale(minY, maxY)
	} else {
		if effectType == "S" {
			kernel = []float64{0,-1,0,-1,5,-1,0,-1,0}
		} else if effectType == "E" {
			kernel = []float64{-1,-1,-1,-1,8,-1,-1,-1,-1}
		} else if effectType == "B" {
			kernel = []float64{1 / 9.0, 1 / 9, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0}
		}
		img.ApplyEffect(kernel, minY, maxY)
	}
	
	waitgroup <- true
}

func (img *ImageTask) ExecuteEffectsBSP(effectType string, minY int, maxY int) {
	var kernel []float64
	if effectType == "G" {
		img.Grayscale(minY, maxY)
	} else {
		if effectType == "S" {
			kernel = []float64{0,-1,0,-1,5,-1,0,-1,0}
		} else if effectType == "E" {
			kernel = []float64{-1,-1,-1,-1,8,-1,-1,-1,-1}
		} else if effectType == "B" {
			kernel = []float64{1 / 9.0, 1 / 9, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0}
		}
		img.ApplyEffect(kernel, minY, maxY)
	}
}
