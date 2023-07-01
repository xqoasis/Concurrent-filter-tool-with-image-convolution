// Package png allows for loading png images and applying
// image flitering effects on them.
package png

import (
	"image/color"
)

//color value struct
type Colors struct{
	Red float64
	Green float64
	Blue float64
	Alpha float64
}

func getOtherKernel(K string) [9]float64{
	var kernel [9]float64
	switch K{
	case "S":
		kernel = [9]float64{0.0, -1.0, 0.0, -1.0, 5.0, -1.0, 0.0, -1.0, 0.0}
	case "E":
		kernel = [9]float64{-1, -1, -1, -1, 8, -1, -1, -1, -1}
	case "B":
		kernel = [9]float64{1/9.0, 1 / 9, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0, 1 / 9.0}
	}
	return kernel
}

// Grayscale applies a grayscale filtering effect to the image
func (img *ImageTask) Grayscale(minY int, maxY int) {
	// Bounds returns defines the dimensions of the image. Always
	// use the bounds Min and Max fields to get out the width
	// and height for the image
	bounds := img.out.Bounds()
	for y := minY; y < maxY; y++ {
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

}

func (img *ImageTask) OtherEffect(K string, minY int, maxY int) {
	kernel := getOtherKernel(K)
	bounds := img.out.Bounds()
	for y := minY; y < maxY; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			img.Convolution(x, y, kernel)
		}
	}
}

func (img *ImageTask) Convolution(x int, y int, kernel [9]float64) {
	// not the mini worker's partition, but the whole img's max
	imgMaxX := img.out.Bounds().Max.X
	imgMaxY := img.out.Bounds().Max.Y

	newColor := Colors{Red: 0.0, Green: 0.0, Blue: 0.0, Alpha: 0.0}
	count := 0
	for yr := -1; yr <= 1; yr++ {
		for xr := -1; xr <= 1; xr++ {
			// default: zero-padding
			filterArgs := Colors{Red: 0.0, Green: 0.0, Blue: 0.0, Alpha: 0.0}
			cntx := x + xr
			cnty := y + yr
			if (cntx >= 0) && (cnty>=0) && (cntx <= imgMaxX) && (cnty <= imgMaxY) {
				r, g, b, a := img.in.At(cntx, cnty).RGBA()
				filterArgs = Colors{Red: float64(r), Green: float64(g), Blue: float64(b), Alpha: float64(a)}
				
			}
			newColor.Red += filterArgs.Red * kernel[count]
			newColor.Green += filterArgs.Green * kernel[count]
			newColor.Blue += filterArgs.Blue * kernel[count]

			if count == 4 {
				newColor.Alpha = filterArgs.Alpha
			}
			count ++
		}
	}
	img.out.Set(x, y, color.RGBA64{clamp(newColor.Red), clamp(newColor.Green), clamp(newColor.Blue), clamp(newColor.Alpha)})
}
func (img *ImageTask) SwapPointer() {
	temp := img.out
	img.out = img.in
	img.in = temp
}