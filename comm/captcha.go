package comm

import (
	"image/color"

	"github.com/mojocn/base64Captcha"
	"github.com/satori/go.uuid"
)

var store = base64Captcha.DefaultMemStore

//configJsonBody json request body.
type configJsonBody struct {
	Id            string
	CaptchaType   string
	VerifyValue   string
	DriverAudio   *base64Captcha.DriverAudio
	DriverString  *base64Captcha.DriverString
	DriverChinese *base64Captcha.DriverChinese
	DriverMath    *base64Captcha.DriverMath
	DriverDigit   *base64Captcha.DriverDigit
}

func DriverStringFunc() (id, b64s string, err error) {
	e := configJsonBody{}
	uid := uuid.NewV4()
	e.Id = uid.String()
	e.DriverString = base64Captcha.NewDriverString(46, 140, 2, 2, 4, "234567890abcdefghjkmnpqrstuvwxyz", &color.RGBA{240, 240, 246, 246}, []string{"wqy-microhei.ttc"})
	driver := e.DriverString.ConvertFonts()
	cap := base64Captcha.NewCaptcha(driver, store)
	return cap.Generate()
}

func DriverDigitFunc() (id, b64s string, err error) {
	e := configJsonBody{}
	uid := uuid.NewV4()
	e.Id = uid.String()
	e.DriverDigit = base64Captcha.DefaultDriverDigit
	driver := e.DriverDigit
	cap := base64Captcha.NewCaptcha(driver, store)
	return cap.Generate()
}
