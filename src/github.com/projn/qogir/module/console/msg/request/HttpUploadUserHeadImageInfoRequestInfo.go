package request

import "mime/multipart"

type HttpUploadUserHeadImageInfoRequestInfo struct {
	XAccessToken string
	UserId string
	File *multipart.FileHeader
}
