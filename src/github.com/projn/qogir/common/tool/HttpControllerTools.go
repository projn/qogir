package tool

import (
	"../../bean"
	"../define"
	"../exception"
	"../initialize"
	"../msg/response"
	"../util"
	"github.com/kataras/iris"
	"reflect"
)

func DealHttpRequest(url string, paramObj interface{}, bodyObj interface{}, ctx iris.Context) {

	ctx.Header(define.HEADER_CONTENT_TYPE, define.CONTENT_TYPE_APPLICATION_JSON_UTF_8);
	if url == "" || paramObj == nil ||ctx == nil {
		bean.Logger.Error("Invaild param.")

		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(MakeHttpErrorResponseMsgInfo(exception.RESULT_INVAILD_PARAM_ERROR,
			exception.GetCommonErrorDecription(exception.RESULT_INVAILD_PARAM_ERROR)))
		return
	}

	bean.Logger.Infof("Request url(%s).", ctx.Request().RequestURI)

	method := ctx.Method()
	requestServiceInfo, e :=initialize.GetRequestServiceInfo(url, method)
	if e != nil {
		bean.Logger.Errorf("Invaild request service info, uri(%s), method(%s).",
			ctx.Request().RequestURI, method)

		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(MakeHttpErrorResponseMsgInfo(exception.RESULT_INVAILD_PARAM_ERROR,
			exception.GetCommonErrorDecription(exception.RESULT_INVAILD_PARAM_ERROR)))
		return
	}

	if requestServiceInfo.Type == "" || requestServiceInfo.Type != initialize.SERVICE_TYPE_HTTP {
		bean.Logger.Error("Invaild request service type info, type(" + requestServiceInfo.Type + ").")

		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(MakeHttpErrorResponseMsgInfo(exception.RESULT_INVAILD_PARAM_ERROR,
			exception.GetCommonErrorDecription(exception.RESULT_INVAILD_PARAM_ERROR)))
		return
	}

	if requestServiceInfo.AuthorizationFilter != nil {
		e := requestServiceInfo.AuthorizationFilter(ctx, requestServiceInfo.UserRoleNameList)
		if e!= nil {
			bean.Logger.Error("Check authorization info error,error info(" + e.Error() + ").")

			ctx.StatusCode(iris.StatusUnauthorized)
			ctx.JSON(MakeHttpErrorResponseMsgInfo(exception.RESULT_INVAILD_USER_TOKEN_ERROR,
				exception.GetCommonErrorDecription(exception.RESULT_INVAILD_USER_TOKEN_ERROR)))
			return
		}
	}

	httpRequestInfo, e := util.ConvertHttpRequestInfo(ctx, paramObj, bodyObj)
	if e != nil {
		bean.Logger.Error("Convert param error,error info(" + e.Error() + ").")

		ctx.StatusCode(iris.StatusBadRequest)
		ctx.JSON(MakeHttpErrorResponseMsgInfo(exception.RESULT_ANALYSE_REQUEST_ERROR,
			exception.GetCommonErrorDecription(exception.RESULT_ANALYSE_REQUEST_ERROR)))
		return
	}

	if httpRequestInfo.ParamObj!= nil {
		e = util.CheckParam(httpRequestInfo.ParamObj)
		if e != nil {
			bean.Logger.Error("Check param error,error info(" + e.Error() + ").")

			ctx.StatusCode(iris.StatusBadRequest)
			ctx.JSON(MakeHttpErrorResponseMsgInfo(exception.RESULT_ANALYSE_REQUEST_ERROR,
				exception.GetCommonErrorDecription(exception.RESULT_ANALYSE_REQUEST_ERROR)))
			return
		}
	}

	httpResponseInfo := requestServiceInfo.HttpServiceFunction(httpRequestInfo)

	ctx.Header(define.HEADER_CONTENT_TYPE, define.CONTENT_TYPE_APPLICATION_JSON_UTF_8);
	if httpResponseInfo.HeaderInfoMap != nil {
		for k, v := range httpResponseInfo.HeaderInfoMap {
			ctx.Header(k, v)
		}
	}

	if httpResponseInfo.Msg != nil {
		objType := reflect.TypeOf(httpResponseInfo.Msg)
		objValue := reflect.ValueOf(httpResponseInfo.Msg)
		if objType.Kind() == reflect.Slice {
			ctx.Write(objValue.Bytes())
		} else {
			ctx.JSON(httpResponseInfo.Msg)
		}
	} else {
		ctx.JSON(MakeHttpErrorResponseMsgInfo(exception.RESULT_OK,
			exception.GetCommonErrorDecription(exception.RESULT_OK)))
	}

 }

func MakeHttpErrorResponseMsgInfo(errorCode string, errorDescription string) response.HttpErrorResponseMsgInfo {
	httpErrorResponseMsgInfo := response.HttpErrorResponseMsgInfo{errorCode, errorDescription}
	return httpErrorResponseMsgInfo
}