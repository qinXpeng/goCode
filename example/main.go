package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/qinXpeng/goCode/capt"
)

const (
	dx = 150
	dy = 50
)

func main() {

	err := capt.ReadFonts("fonts", ".ttf")
	if err != nil {
		fmt.Println(err)
		return
	}
	http.HandleFunc("/", Index)
	http.HandleFunc("/get/", Get)
	fmt.Println("服务已启动...")
	err = http.ListenAndServe(":8800", nil)
	if err != nil {
		log.Fatal(err)
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("tpl/index.html")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(w, nil)
}
func Get(w http.ResponseWriter, r *http.Request) {

	captchaImage, err := capt.NewCaptchaImage(dx, dy, capt.RandLightColor())

	captchaImage.DrawNoise(capt.CaptchaComplexLower)
	captchaImage.DrawTextNoise(capt.CaptchaComplexLower)
	captchaImage.DrawText(capt.RandText(6))
	captchaImage.DrawBorder(capt.ColorToRGB(0x17A7A7A))
	captchaImage.DrawHollowLine()
	if err != nil {
		fmt.Println(err)
	}

	captchaImage.SaveImage(w, capt.ImageFormatGif)
}
