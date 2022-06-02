// Package png allows for loading png images and applying
// image flitering effects on them
package png

import (
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
	"fmt"
)

// The Image represents a structure for working with PNG images.
// From Professor Samuels: You are allowed to update this and change it as you wish!
type ImageTask struct {
	in     *image.RGBA64   //The original pixels before applying the effect
	out    *image.RGBA64   //The updated pixels after applying teh effect
	Bounds image.Rectangle //The size of the image
	dataDir string
	inPath string
	outPath string
	effects []interface{}
}

func (i *ImageTask) GetIn() *image.RGBA64{
	return i.in
}

func (i *ImageTask) GetOut() *image.RGBA64{
	return i.out
}

func (i *ImageTask) GetDataDir() string{
	return i.dataDir
}

func (i *ImageTask) GetOutPath() string{
	return i.outPath
}

func (i *ImageTask) GetMaxY() int {
	return i.out.Bounds().Max.Y
}

func (i *ImageTask) GetMinY() int {
	return i.out.Bounds().Min.Y
}

func (i *ImageTask) GetEffects() []string{
	var effects []string
	for _, effect := range(i.effects) {
		effects = append(effects, string(effect.(string)))
	}
	return effects
}

func (i *ImageTask) SetInPointer(newPointer *image.RGBA64) {
	i.in = newPointer
}

func (i *ImageTask) SetDataDir(dataDir string) {
	i.dataDir = dataDir
}

func (i *ImageTask) SetInPath(inPath string) {
	i.inPath = inPath
}

func (i *ImageTask) SetOutPath(outPath string) {
	i.outPath = outPath
}

func (i *ImageTask) SetEffects(effects []interface{}) {
	i.effects = effects
}

//
// Public functions
//

// Load returns a Image that was loaded based on the filePath parameter
// From Professor Samuels:  You are allowed to modify and update this as you wish
func Load(filePath string) (*ImageTask, error) {

	inReader, err := os.Open(filePath)

	if err != nil {
		return nil, err
	}
	defer inReader.Close()

	inOrig, err := png.Decode(inReader)

	if err != nil {
		return nil, err
	}

	bounds := inOrig.Bounds()

	outImg := image.NewRGBA64(bounds)
	inImg := image.NewRGBA64(bounds)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := inOrig.At(x, y).RGBA()
			inImg.Set(x, y, color.RGBA64{uint16(r), uint16(g), uint16(b), uint16(a)})
		}
	}
	task := &ImageTask{}
	task.in = inImg
	task.out = outImg
	task.Bounds = bounds
	return task, nil
}

// Save saves the image to the given file
// From Professor Samuels:  You are allowed to modify and update this as you wish
func (img *ImageTask) Save(filePath string) error {

	outWriter, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer outWriter.Close()

	err = png.Encode(outWriter, img.out)
	if err != nil {
		return err
	}
	return nil
}

func (img *ImageTask) SaveToPNG(dataDir string, outPath string) {
	//Saves the image to a new file
	fullOutPath := fmt.Sprintf("../data/out/%v_%v", dataDir, outPath)
	err := img.Save(fullOutPath)	
	
	if err != nil {
		panic(err)
	}
}

//clamp will clamp the comp parameter to zero if it is less than zero or to 65535 if the comp parameter
// is greater than 65535.
func clamp(comp float64) uint16 {
	return uint16(math.Min(65535, math.Max(0, comp)))
}

func (img *ImageTask) Swap() {
	img.in = img.out
}