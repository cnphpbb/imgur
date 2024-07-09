package imgur

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/bmp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
	"golang.org/x/image/tiff"
	"golang.org/x/image/webp"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"io"
	"io/ioutil"
	"math"
)

// GetImageType returns the extension, mimetype and corresponding decoder function of the image.
// It judges the image by its first few bytes called magic number.
func GetImageType(bytes []byte) (ext string, mimetype string, decoder func(r io.Reader) (image.Image, error), err error) {
	if len(bytes) < 2 {
		err = ErrSourceImageNotSupport
		return
	}

	if bytes[0] == 0xFF && bytes[1] == 0xD8 {
		ext = "jpg"
		mimetype = "image/jpeg"
		decoder = jpeg.Decode
	}

	if len(bytes) >= 4 && bytes[0] == 0x89 && bytes[1] == 0x50 && bytes[2] == 0x4E && bytes[3] == 0x47 {
		ext = "png"
		mimetype = "image/png"
		decoder = png.Decode
	}

	if bytes[0] == 0x42 && bytes[1] == 0x4D {
		ext = "bmp"
		mimetype = "image/x-ms-bmp"
		decoder = bmp.Decode
	}

	if (bytes[0] == 0x49 && bytes[1] == 0x49) || (bytes[0] == 0x4D && bytes[1] == 0x4D) {
		ext = "tiff"
		mimetype = "image/tiff"
		decoder = tiff.Decode
	}

	if bytes[0] == 0x52 && bytes[1] == 0x49 {
		ext = "webp"
		mimetype = "image/webp"
		decoder = webp.Decode
	}

	if ext == "" {
		err = ErrSourceImageNotSupport
	}

	return
}

// Color2Hex converts a color.Color to its hex string representation.
func Color2Hex(c color.Color) string {
	r, g, b, _ := c.RGBA()
	return fmt.Sprintf("#%02X%02X%02X", uint8(r>>8), uint8(g>>8), uint8(b>>8))
}

// Image2RGBA converts an image to RGBA.
func Image2RGBA(img image.Image) *image.RGBA {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Over)
	return rgba
}

func fixp(x, y float64) fixed.Point26_6 {
	return fixed.Point26_6{fix(x), fix(y)}
}

func fix(x float64) fixed.Int26_6 {
	return fixed.Int26_6(math.Round(x * 64))
}

func unfix(x fixed.Int26_6) float64 {
	const shift, mask = 6, 1<<6 - 1
	if x >= 0 {
		return float64(x>>shift) + float64(x&mask)/64
	}
	x = -x
	if x >= 0 {
		return -(float64(x>>shift) + float64(x&mask)/64)
	}
	return 0
}

func Radians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

// LoadFontFace is a helper function to load the specified font file with
// the specified point size. Note that the returned `font.Face` objects
// are not thread safe and cannot be used in parallel across goroutines.
// You can usually just use the Context.LoadFontFace function instead of
// this package-level function.
func LoadFontFace(path string, points float64) (font.Face, error) {
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	f, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, err
	}
	face := truetype.NewFace(f, &truetype.Options{
		Size: points,
		// Hinting: font.HintingFull,
	})
	return face, nil
}

func getOpenTypeFontFace(fontFilePath string, fontSize, dpi float64) (*font.Face, error) {
	fontData, fontFileReadErr := ioutil.ReadFile(fontFilePath)
	if fontFileReadErr != nil {
		return nil, fontFileReadErr
	}
	otfFont, parseErr := opentype.Parse(fontData)
	if parseErr != nil {
		return nil, parseErr
	}
	otfFace, newFaceErr := opentype.NewFace(otfFont, &opentype.FaceOptions{
		Size: fontSize,
		DPI:  dpi,
	})
	if newFaceErr != nil {
		return nil, newFaceErr
	}
	return &otfFace, nil
}
