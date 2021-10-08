package resource

import (
	"errors"
	"image"
	"image/draw"
	_ "image/jpeg"
	_ "image/png"
	"os"

	"github.com/ddomurad/goCraft/core"
	"github.com/go-gl/gl/v3.3-core/gl"
)

const (
	RT_TEXTURE core.ResourceType = "texture"
)

type TextureParams struct {
	FilePath         string
	NearestFiltering bool
}

type TextureData struct {
	Id uint32
}

func GetEmptyTexture(uri string) core.Resource {
	return core.Resource{
		Type:  RT_TEXTURE,
		Uri:   uri,
		Empty: true,
		Data: TextureData{
			Id: 0,
		},
	}
}

type FileTextureLoader struct{}

func (l FileTextureLoader) CanLoad(resourceType core.ResourceType, uri string, param core.LoaderParam) bool {
	if resourceType != RT_TEXTURE {
		return false
	}

	switch param.(type) {
	case TextureParams:
		return true
	default:
		return false
	}
}

func (l FileTextureLoader) Load(uri string, param core.LoaderParam) (core.Resource, error) {
	textureParams := param.(TextureParams)
	textureFile, err := os.Open(textureParams.FilePath)

	if err != nil {
		return GetEmptyTexture(uri), err
	}
	defer textureFile.Close()

	img, _, err := image.Decode(textureFile)
	if err != nil {
		return GetEmptyTexture(uri), err
	}

	rgba := image.NewRGBA(img.Bounds())
	draw.Draw(rgba, rgba.Bounds(), img, image.Pt(0, 0), draw.Src)

	if rgba.Stride != rgba.Rect.Size().X*4 { // TODO-cs: why?
		return GetEmptyTexture(uri), errors.New("unsported stride")
	}

	var textureId uint32
	gl.GenTextures(1, &textureId)
	gl.BindTexture(gl.TEXTURE_2D, textureId)
	defer gl.BindTexture(gl.TEXTURE_2D, 0)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_BASE_LEVEL, 0)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAX_LEVEL, 0)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA,
		int32(rgba.Rect.Size().X), int32(rgba.Rect.Size().Y),
		0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(rgba.Pix))

	if textureParams.NearestFiltering {
		gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	} else {
		gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameterf(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	}

	return core.Resource{
		Type:  RT_TEXTURE,
		Uri:   uri,
		Empty: false,
		Data: TextureData{
			Id: textureId,
		},
		Unload: func() {
			gl.DeleteTextures(1, &textureId)
		},
	}, nil
}

func NewFileTextureLoader() FileTextureLoader {
	return FileTextureLoader{}
}
