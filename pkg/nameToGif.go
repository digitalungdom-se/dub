package pkg

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"math"
	"net/http"

	"github.com/tfriedel6/canvas"
	"github.com/tfriedel6/canvas/backend/softwarebackend"
)

func NameToGif(name string, avatarURL string) (*bytes.Buffer, error) {
	gifAnimation := &gif.GIF{
		Image:     make([]*image.Paletted, 100),
		Delay:     make([]int, 100),
		LoopCount: -1,
	}

	gifAnimation.Image = gifAnimation.Image[:0]
	gifAnimation.Delay = gifAnimation.Delay[:0]

	backend := softwarebackend.New(800, 250)
	cv := canvas.New(backend)

	wint, hint := cv.Size()
	w, h := float64(wint), float64(hint)

	response, e := http.Get(avatarURL)
	if e != nil {
		return nil, e
	}

	avatarImage, _, _ := image.Decode(response.Body)
	response.Body.Close()

	backendTMP := softwarebackend.New(128, 128)
	cvTMP := canvas.New(backendTMP)

	wintTMP, hintTMP := cvTMP.Size()
	wTMP, hTMP := float64(wintTMP), float64(hintTMP)

	avatarCVTMP, _ := cvTMP.LoadImage(avatarImage)

	cvTMP.SetGlobalAlpha(0)
	cvTMP.BeginPath()
	cvTMP.Arc(wTMP/2, hTMP/2, hTMP/2, 0, math.Pi*2, true)
	cvTMP.Clip()
	cvTMP.Stroke()
	cvTMP.ClosePath()
	cvTMP.Fill()

	cvTMP.DrawImage(avatarCVTMP, 0, 0)

	avatarRGBA := cvTMP.GetImageData(0, 0, wintTMP, hintTMP)
	avatarIMG := avatarRGBA.SubImage(avatarRGBA.Rect)

	avatarCV, _ := cv.LoadImage(avatarIMG)

	cv.SetFillStyle("#3E5ABE")
	cv.FillRect(0, 0, w, h)

	cv.SetStrokeStyle("#3E5ABE")
	cv.StrokeRect(0, 0, w, h)

	var nameFont float64 = 70
	cv.SetFont("/usr/share/fonts/Unifont/Unifont.ttf", nameFont)

	for cv.MeasureText(name).Width > w-325 {
		nameFont -= 1
		cv.SetFont("/usr/share/fonts/Unifont/Unifont.ttf", nameFont)
	}

	for i := 0; i < len(name)+1; i++ {
		cv.SetFillStyle("#ffffff")

		cv.FillText(fmt.Sprintf(">%v_", name[:i]), 250, h/2+nameFont/2)

		rgbaImage := cv.GetImageData(0, 0, wint, hint)
		img := rgbaImage.SubImage(rgbaImage.Rect)

		cv.DrawImage(avatarCV, 25, 25, 200, 200)

		AddImageToGif(gifAnimation, &img, 15)

		cv.SetFillStyle("#3E5ABE")
		cv.FillRect(0, 0, w, h)
	}

	for i := 0; i < 20; i++ {
		cv.SetFillStyle("#ffffff")

		if len(gifAnimation.Image)%2 == 0 {
			cv.FillText(fmt.Sprintf(">%v", name), 250, h/2+nameFont/2)
		} else {
			cv.FillText(fmt.Sprintf(">%v_", name), 250, h/2+nameFont/2)
		}

		rgbaImage := cv.GetImageData(0, 0, wint, hint)
		img := rgbaImage.SubImage(rgbaImage.Rect)

		cv.DrawImage(avatarCV, 25, 25, 200, 200)

		AddImageToGif(gifAnimation, &img, 45)

		cv.SetFillStyle("#3E5ABE")
		cv.FillRect(0, 0, w, h)
	}

	cv.SetFillStyle("#ffffff")
	cv.FillText(fmt.Sprintf(">%v_", name), 250, h/2+nameFont/2)

	rgbaImage := cv.GetImageData(0, 0, wint, hint)
	img := rgbaImage.SubImage(rgbaImage.Rect)

	cv.DrawImage(avatarCV, 25, 25, 200, 200)

	AddImageToGif(gifAnimation, &img, 15)

	gifbuf := new(bytes.Buffer)

	gif.EncodeAll(gifbuf, gifAnimation)

	return gifbuf, nil
}

func AddImageToGif(gifImages *gif.GIF, img *image.Image, delay int) {
	opts := gif.Options{
		NumColors: 216,
		Drawer:    draw.FloydSteinberg,
	}
	b := (*img).Bounds()

	pimg := image.NewPaletted(b, palette.WebSafe[:opts.NumColors])
	if opts.Quantizer != nil {
		pimg.Palette = opts.Quantizer.Quantize(make(color.Palette, 0, opts.NumColors), *img)
	}

	opts.Drawer.Draw(pimg, b, *img, image.ZP)

	gifImages.Image = append(gifImages.Image, pimg)
	gifImages.Delay = append(gifImages.Delay, delay)
}
