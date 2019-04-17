package controller

import (
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"projn.com/qogir/common/struct"
	"projn.com/qogir/common/define"
	"projn.com/qogir/module/console/msg/request"
	"projn.com/qogir/module/console/msg/request/type"

	"github.com/kataras/iris"
)

func RegisterHandler(app *iris.Application) error {
	if app == nil {
		return errors.New("Application is invaild.");
	}

	app.Get("/user/verification-code", GetVerificationCodeInfo)
	app.Post("/user/login", Login)
	app.Get("/user/logout/{user_id}", Logout)
	app.Get("/user/headimage/{user_id}", UploadUserHeadImageInfo)

	return nil
}

func GetVerificationCodeInfo(ctx iris.Context) {

}

func Login(ctx iris.Context)  {
	locale := ctx.GetHeader(define.HEADER_LANGUAGE)

	var loginRequestInfo _type.LoginRequestInfo
	error := ctx.ReadJSON(loginRequestInfo)

	httpLoginRequestInfo := request.HttpLoginRequestInfo{loginRequestInfo}

	httpRequestInfo:=_struct.HttpRequestInfo{locale, httpLoginRequestInfo}
}

func Logout(ctx iris.Context) {
	locale := ctx.GetHeader(define.HEADER_LANGUAGE)

	xAccessToken := ctx.GetHeader("x-access-token")
	userId := ctx.Params().Get("user_id")

	httpLogoutRequestInfo := request.HttpLogoutRequestInfo{xAccessToken, userId}

	httpRequestInfo:=_struct.HttpRequestInfo{locale, httpLogoutRequestInfo}
}

func UploadUserHeadImageInfo(ctx iris.Context) {
	locale := ctx.GetHeader(define.HEADER_LANGUAGE)

	xAccessToken := ctx.GetHeader("x-access-token")
	userId := ctx.Params().Get("user_id")
	file, info, err := ctx.FormFile("file")
	if err != nil {
		ctx.StatusCode(iris.StatusInternalServerError)
		ctx.HTML("Error while uploading: <b>" + err.Error() + "</b>")
		return
	}
	defer file.Close()

	httpUploadUserHeadImageInfoRequestInfo := request.HttpUploadUserHeadImageInfoRequestInfo{xAccessToken, userId, info}

	httpRequestInfo:=_struct.HttpRequestInfo{locale, httpUploadUserHeadImageInfoRequestInfo}

	GetVerificationCodeInfoService

}

func saveUploadedFile(fh *multipart.FileHeader, destDirectory string) (int64, error) {
	src, err := fh.Open()
	if err != nil {
		return 0, err
	}
	defer src.Close()
	out, err := os.OpenFile(filepath.Join(destDirectory, fh.Filename),
		os.O_WRONLY|os.O_CREATE, os.FileMode(0666))
	if err != nil {
		return 0, err
	}
	defer out.Close()
	return io.Copy(out, src)
}