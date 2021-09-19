package core

import "github.com/go-gl/gl/v2.1/gl"

type RenderBuffer struct {
	TextureId uint32
	bufferId  uint32
}

func NewRenderBuffer(w, h int32) *RenderBuffer {
	var frameBuffer uint32
	var textureId uint32

	gl.GenFramebuffers(1, &frameBuffer)
	gl.GenTextures(1, &textureId)

	gl.BindFramebuffer(gl.FRAMEBUFFER, frameBuffer)
	gl.BindTexture(gl.TEXTURE_2D, textureId)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, w, h, 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(0))
	return &RenderBuffer{
		bufferId:  frameBuffer,
		TextureId: textureId,
	}
}

func (rb *RenderBuffer) Bind() {}
func (rb *RenderBuffer) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (rb *RenderBuffer) Release() {}
