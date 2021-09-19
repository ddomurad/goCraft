package render_buffer

import "github.com/go-gl/gl/v3.3-core/gl"

type RenderBuffer struct {
	TextureId uint32
	bufferId  uint32
}

func NewRenderBuffer(w, h int32, nearesFilter bool) *RenderBuffer {
	var frameBuffer uint32
	var textureId uint32

	gl.GenFramebuffers(1, &frameBuffer)
	gl.GenTextures(1, &textureId)

	gl.BindFramebuffer(gl.FRAMEBUFFER, frameBuffer)
	gl.BindTexture(gl.TEXTURE_2D, textureId)

	defer gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	defer gl.BindTexture(gl.TEXTURE_2D, 0)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, w, h, 0, gl.RGBA, gl.UNSIGNED_BYTE, nil)
	if nearesFilter {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	} else {
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	}

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, textureId, 0)

	return &RenderBuffer{
		bufferId:  frameBuffer,
		TextureId: textureId,
	}
}

func (rb *RenderBuffer) Bind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, rb.bufferId)
}

func (rb *RenderBuffer) Unbind() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
}

func (rb *RenderBuffer) Release() {
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.BindTexture(gl.TEXTURE_2D, 0)
	gl.DeleteFramebuffers(1, &rb.bufferId)
	gl.DeleteTextures(1, &rb.TextureId)
}
