package imgur

import (
	"bytes"
	"github.com/golang/freetype/raster"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"io/ioutil"
	"math"
	"net/http"
	"os"
)

type ImageManager struct {
}

// Load an image from source.
// source can be a file path, a URL, a base64 encoded string, an *os.File, an image.Image or a byte slice.
func Load(source interface{}) *Image {
	switch source.(type) {
	case string:
		return loadFromString(source.(string))
	case *os.File:
		return LoadFromFile(source.(*os.File))
	case image.Image:
		return LoadFromImage(source.(image.Image))
	case []byte:
		return loadFromString(string(source.([]byte)))
	case *Image:
		return LoadFromImgo(source.(*Image))
	default:
		i := &Image{}
		i.addError(ErrSourceNotSupport)
		return i
	}
}

// loadFromString loads an image when the source is a string.
func loadFromString(source string) (i *Image) {
	i = &Image{}

	if len(source) == 0 {
		i.addError(ErrSourceStringIsEmpty)
		return
	}

	if len(source) > 4 && source[:4] == "http" {
		return LoadFromUrl(source)
	} else if len(source) > 10 && source[:10] == "data:image" {
		return LoadFromBase64(source)
	} else {
		return LoadFromPath(source)
	}
}

// LoadFromUrl loads an image when the source is an url.
func LoadFromUrl(url string) (i *Image) {
	i = &Image{}

	// Get the image response from the url.
	resp, err := http.Get(url)
	if err != nil {
		i.addError(err)
		return
	}
	defer resp.Body.Close()

	// Read the image data from the response.
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	// Get the extension, mimetype and corresponding decoder function of the image.
	ext, mime, decoder, err := GetImageType(bodyBytes[:8])
	if err != nil {
		i.addError(err)
		return
	}

	// Decode the image.
	file := bytes.NewReader(bodyBytes)
	img, err := decoder(file)
	if err != nil {
		i.addError(ErrSourceNotSupport)
		return
	}

	return &Image{
		image:     Image2RGBA(img),
		width:     img.Bounds().Dx(),
		height:    img.Bounds().Dy(),
		extension: ext,
		mimetype:  mime,
	}
}

// LoadFromPath loads an image from a path.
func LoadFromPath(path string) (i *Image) {
	i = &Image{}

	file, err := os.Open(path)
	if err != nil {
		i.addError(err)
		return
	}
	defer file.Close()

	return LoadFromFile(file)
}

// LoadFromFile loads an image from a file.
func LoadFromFile(file *os.File) (i *Image) {
	i = &Image{}

	// Read the first 8 bytes of the image.
	buf := make([]byte, 8)
	_, err := file.Read(buf)
	if err != nil {
		i.addError(err)
		return
	}

	// After reading the first 8 bytes, we seek back to the beginning of the file.
	_, err = file.Seek(0, 0)
	if err != nil {
		i.addError(err)
		return
	}

	// Get the extension, mimetype and corresponding decoder function of the image.
	ext, mime, decoder, err := GetImageType(buf)
	if err != nil {
		i.addError(err)
		return
	}

	// Decode the image.
	img, err := decoder(file)
	if err != nil {
		i.addError(err)
		return
	}

	// Set the image properties.
	stat, _ := file.Stat()

	return &Image{
		image:     Image2RGBA(img),
		width:     img.Bounds().Dx(),
		height:    img.Bounds().Dy(),
		extension: ext,
		mimetype:  mime,
		filesize:  stat.Size(),
	}
}

// LoadFromImage loads an image from an instance of image.Image.
func LoadFromImage(img image.Image) (i *Image) {
	i = &Image{}

	if img == nil {
		i.addError(ErrSourceImageIsNil)
		return
	}

	var formatName string
	switch img.(type) {
	case *image.NRGBA: // png
		formatName = "png"
	case *image.RGBA: // bmp, tiff
		formatName = "png"
	case *image.YCbCr: // jpeg, webp
		formatName = "jpg"
	default:
		i.addError(ErrSourceImageNotSupport)
		return
	}

	return &Image{
		image:     Image2RGBA(img),
		width:     img.Bounds().Dx(),
		height:    img.Bounds().Dy(),
		extension: formatName,
		mimetype:  "image/" + formatName,
	}
}

// LoadFromImgo loads an image from an instance of Image.
func LoadFromImgo(i *Image) *Image {
	return i
}

// Canvas create a new empty image.
func Canvas(width, height int, fillColor ...color.Color) *Image {
	var c color.Color
	if len(fillColor) == 0 {
		c = color.Transparent
	} else {
		c = fillColor[0]
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))
	for x := 0; x < img.Bounds().Dx(); x++ {
		for y := 0; y < img.Bounds().Dy(); y++ {
			img.Set(x, y, c)
		}
	}

	return &Image{
		image:     img,
		width:     img.Bounds().Dx(),
		height:    img.Bounds().Dy(),
		extension: "png",
		mimetype:  "image/png",
	}
}

func flattenPath(p raster.Path) [][]Point {
	var result [][]Point
	var path []Point
	var cx, cy float64
	for i := 0; i < len(p); {
		switch p[i] {
		case 0:
			if len(path) > 0 {
				result = append(result, path)
				path = nil
			}
			x := unfix(p[i+1])
			y := unfix(p[i+2])
			path = append(path, Point{x, y})
			cx, cy = x, y
			i += 4
		case 1:
			x := unfix(p[i+1])
			y := unfix(p[i+2])
			path = append(path, Point{x, y})
			cx, cy = x, y
			i += 4
		case 2:
			x1 := unfix(p[i+1])
			y1 := unfix(p[i+2])
			x2 := unfix(p[i+3])
			y2 := unfix(p[i+4])
			points := QuadraticBezier(cx, cy, x1, y1, x2, y2)
			path = append(path, points...)
			cx, cy = x2, y2
			i += 6
		case 3:
			x1 := unfix(p[i+1])
			y1 := unfix(p[i+2])
			x2 := unfix(p[i+3])
			y2 := unfix(p[i+4])
			x3 := unfix(p[i+5])
			y3 := unfix(p[i+6])
			points := CubicBezier(cx, cy, x1, y1, x2, y2, x3, y3)
			path = append(path, points...)
			cx, cy = x3, y3
			i += 8
		default:
			panic("bad path")
		}
	}
	if len(path) > 0 {
		result = append(result, path)
	}
	return result
}

func dashPath(paths [][]Point, dashes []float64, offset float64) [][]Point {
	var result [][]Point
	if len(dashes) == 0 {
		return paths
	}
	if len(dashes) == 1 {
		dashes = append(dashes, dashes[0])
	}
	for _, path := range paths {
		if len(path) < 2 {
			continue
		}
		previous := path[0]
		pathIndex := 1
		dashIndex := 0
		segmentLength := 0.0

		// offset
		if offset != 0 {
			var totalLength float64
			for _, dashLength := range dashes {
				totalLength += dashLength
			}
			offset = math.Mod(offset, totalLength)
			if offset < 0 {
				offset += totalLength
			}
			for i, dashLength := range dashes {
				offset -= dashLength
				if offset < 0 {
					dashIndex = i
					segmentLength = dashLength + offset
					break
				}
			}
		}

		var segment []Point
		segment = append(segment, previous)
		for pathIndex < len(path) {
			dashLength := dashes[dashIndex]
			point := path[pathIndex]
			d := previous.Distance(point)
			maxd := dashLength - segmentLength
			if d > maxd {
				t := maxd / d
				p := previous.Interpolate(point, t)
				segment = append(segment, p)
				if dashIndex%2 == 0 && len(segment) > 1 {
					result = append(result, segment)
				}
				segment = nil
				segment = append(segment, p)
				segmentLength = 0
				previous = p
				dashIndex = (dashIndex + 1) % len(dashes)
			} else {
				segment = append(segment, point)
				previous = point
				segmentLength += d
				pathIndex++
			}
		}
		if dashIndex%2 == 0 && len(segment) > 1 {
			result = append(result, segment)
		}
	}
	return result
}

func rasterPath(paths [][]Point) raster.Path {
	var result raster.Path
	for _, path := range paths {
		var previous fixed.Point26_6
		for i, point := range path {
			f := point.Fixed()
			if i == 0 {
				result.Start(f)
			} else {
				dx := f.X - previous.X
				dy := f.Y - previous.Y
				if dx < 0 {
					dx = -dx
				}
				if dy < 0 {
					dy = -dy
				}
				if dx+dy > 8 {
					// TODO: this is a hack for cases where two points are
					// too close - causes rendering issues with joins / caps
					result.Add1(f)
				}
			}
			previous = f
		}
	}
	return result
}

func dashed(path raster.Path, dashes []float64, offset float64) raster.Path {
	return rasterPath(dashPath(flattenPath(path), dashes, offset))
}
