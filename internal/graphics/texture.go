package graphics

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"

	"os"

	"github.com/go-gl/gl/v4.6-core/gl"
)

const (
	SKYBOX_TEXTURE = 0
	BLOCKS_TEXTURE = 1
)

func createTexture(id uint32, xtype uint32) uint32 {
	var texture uint32
	gl.GenTextures(1, &texture)
	gl.ActiveTexture(gl.TEXTURE0 + id)
	gl.BindTexture(xtype, texture)
	return texture
}

func imageToRGBA(img image.Image) *image.RGBA {
	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, img.Bounds().Min, draw.Src)
	return rgba
}

/*
Takes multiple images of width and height "width" and stitches them together in a grid with "columns" columns
Returns the stitched image
*/
func stitchImages(rows, columns, width int, imgs ...image.Image) *image.RGBA {
	r := image.Rectangle{image.Point{0, 0}, image.Point{width * columns, width * rows}}

	rgba := image.NewRGBA(r)

	for i, img := range imgs {
		offset := image.Point{(i % columns) * width, (i / columns) * width}
		draw.Draw(rgba, image.Rect(offset.X, offset.Y, offset.X+width, offset.Y+width), img, image.Point{0, 0}, draw.Src)
	}
	return rgba
}

func setTextureInterpolation(target uint32, param int32) {
	gl.TexParameteri(target, gl.TEXTURE_MIN_FILTER, param)
	gl.TexParameteri(target, gl.TEXTURE_MAG_FILTER, param)
}

func setTextureWrap(target uint32, param int32) {
	gl.TexParameteri(target, gl.TEXTURE_WRAP_S, param)
	gl.TexParameteri(target, gl.TEXTURE_WRAP_T, param)
	gl.TexParameteri(target, gl.TEXTURE_WRAP_R, param)
}

/*
Load the cubemap texture in the "id" tray
Assumes a .png file
*/
func loadCubemap(path string, id uint32) error {
	createTexture(id, gl.TEXTURE_CUBE_MAP)

	faces := []string{"right", "left", "top", "bottom", "front", "back"}

	for i, face := range faces {
		imgFile, err := os.Open(fmt.Sprintf("%s/%s.png", path, face))
		if err != nil {
			return fmt.Errorf("os.Open(): %w", err)
		}
		img, err := png.Decode(imgFile)
		imgFile.Close()
		if err != nil {
			return fmt.Errorf("image.Decode(): %w", err)
		}

		rgba := imageToRGBA(img)

		gl.TexImage2D(
			gl.TEXTURE_CUBE_MAP_POSITIVE_X+uint32(i),
			0,
			gl.RGBA,
			int32(rgba.Rect.Size().X),
			int32(rgba.Rect.Size().Y),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			gl.Ptr(rgba.Pix),
		)
	}

	setTextureInterpolation(gl.TEXTURE_CUBE_MAP, gl.LINEAR)
	setTextureWrap(gl.TEXTURE_CUBE_MAP, gl.CLAMP_TO_EDGE)

	return nil
}

/*
Creates and loads the texture atlas in the "id" tray
"Path" is the path to a folder containing .png files
*/
func loadTextureAtlas(path string, id uint32, resolution int) (map[string]uint8, error) {
	createTexture(id, gl.TEXTURE_2D)

	textures, err := os.ReadDir(path)
	if err != nil {
		return nil, fmt.Errorf("os.ReadDir(): %w", err)
	}

	atlasMap := make(map[string]uint8)
	images := make([]image.Image, 0)
	for i, texture := range textures {
		imgFile, err := os.Open(fmt.Sprintf("%s/%s", path, texture.Name()))
		if err != nil {
			return nil, fmt.Errorf("os.Open(): %w", err)
		}
		img, err := png.Decode(imgFile)
		imgFile.Close()
		if err != nil {
			return nil, fmt.Errorf("image.Decode(): %w", err)
		}
		atlasMap[texture.Name()] = uint8(i)

		images = append(images, imageToRGBA(img))
	}

	atlasImage := stitchImages(16, 16, resolution, images...)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(atlasImage.Rect.Size().X),
		int32(atlasImage.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(atlasImage.Pix),
	)

	setTextureInterpolation(gl.TEXTURE_2D, gl.NEAREST)
	setTextureWrap(gl.TEXTURE_2D, gl.CLAMP_TO_EDGE)

	return atlasMap, nil
}
