package util

import (
	"../struct"
	"../define"
	"encoding/json"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
	"github.com/magiconair/properties"
	"reflect"
	"strconv"
)

const (
	PARAM_LOCATION_TYPE_PATH = iota
	PARAM_LOCATION_TYPE_QUERY
	PARAM_LOCATION_TYPE_HEADER
	PARAM_LOCATION_TYPE_BODY
)

const PARAM_LOCATION = "ParamLocation"

type ParamLocation struct {
	Type int
	File bool
}

const PARAM_NAME = "ParamName"

func ConvertHttpRequestInfo(ctx iris.Context, paramObj interface{}, bodyObj interface{}) (_struct.HttpRequestInfo, error) {

	objType := reflect.TypeOf(paramObj)
	objValue := reflect.ValueOf(paramObj)

	var err error = nil
	for i := 0; i < objType.NumField(); i++ {

		field := objType.Field(i)
		paramLocationStr := field.Tag.Get(PARAM_LOCATION)

		if paramLocationStr == "" {
			continue
		}

		p := properties.MustLoadString(paramLocationStr)
		var paramLocation ParamLocation
		if err := p.Decode(&paramLocation); err != nil {
			return _struct.HttpRequestInfo{},err
		}

		fieldName := field.Name;
		paramName := field.Tag.Get(PARAM_NAME)
		if paramName != "" {
			fieldName = paramName
		}

		if (paramLocation.Type == PARAM_LOCATION_TYPE_HEADER) {
			fieldValue := ctx.GetHeader(fieldName);
			err = setFieldValue(objValue.Field(i).Elem(), objValue.Field(i).Type(), fieldValue);
		} else if (paramLocation.Type == PARAM_LOCATION_TYPE_QUERY) {
			fieldValue := ctx.URLParam(fieldName);
			err = setFieldValue(objValue.Field(i).Elem(), objValue.Field(i).Type(), fieldValue);
		} else if (paramLocation.Type == PARAM_LOCATION_TYPE_PATH) {
			fieldValue := ctx.Params().Get(fieldName)
			err = setFieldValue(objValue.Field(i).Elem(), objValue.Field(i).Type(), fieldValue);
		} else if (paramLocation.Type == PARAM_LOCATION_TYPE_BODY) {
			if paramLocation.File == false {
				err := ctx.ReadJSON(bodyObj)
				if err == nil {
					objValue.Field(i).Elem().SetPointer(&bodyObj)
				}
			} else {
				file, info, err := ctx.FormFile("file")
				if err == nil {
					objValue.Field(i).Elem().SetPointer(info)
				}
				defer file.Close()
			}
		}
		if err !=nil {
			break
		}
 	}

	language := ctx.GetHeader(define.HEADER_LANGUAGE);
	return _struct.HttpRequestInfo{language,paramObj}, err

}

func ConvertWsRequestInfo(connection websocket.Connection, textMsg string, paramObj interface{}, bodyObj interface{}) (_struct.WsRequestInfo, error) {

	objType := reflect.TypeOf(paramObj)
	objValue := reflect.ValueOf(paramObj)

	var err error = nil
	for i := 0; i < objType.NumField(); i++ {

		field := objType.Field(i)
		paramLocationStr := field.Tag.Get(PARAM_LOCATION)

		if paramLocationStr == "" {
			continue
		}

		p := properties.MustLoadString(paramLocationStr)
		var paramLocation ParamLocation
		if err := p.Decode(&paramLocation); err != nil {
			return _struct.WsRequestInfo{},err
		}

		fieldName := field.Name;
		paramName := field.Tag.Get(PARAM_NAME)
		if paramName != "" {
			fieldName = paramName
		}

		if paramLocation.Type == PARAM_LOCATION_TYPE_HEADER || paramLocation.Type == PARAM_LOCATION_TYPE_QUERY {
			fieldValue := connection.GetValueString(fieldName);
			err = setFieldValue(objValue.Field(i).Elem(), objValue.Field(i).Type(), fieldValue);
		} else if paramLocation.Type == PARAM_LOCATION_TYPE_BODY {
			err := json.Unmarshal([]byte(textMsg), bodyObj)
			if err == nil {
				objValue.Field(i).Elem().SetPointer(&bodyObj)
			}
		}
		if err !=nil {
			break
		}
	}

	return _struct.WsRequestInfo{paramObj}, err

}

func ConvertMsgRequestInfo(msgText string, paramObj interface{}, bodyObj interface{}) (_struct.MsgRequestInfo, error) {

	objType := reflect.TypeOf(paramObj)
	objValue := reflect.ValueOf(paramObj)

	var err error = nil
	for i := 0; i < objType.NumField(); i++ {

		field := objType.Field(i)
		paramLocationStr := field.Tag.Get(PARAM_LOCATION)

		if paramLocationStr == "" {
			continue
		}

		p := properties.MustLoadString(paramLocationStr)
		var paramLocation ParamLocation
		if err := p.Decode(&paramLocation); err != nil {
			return _struct.MsgRequestInfo{},err
		}

		if paramLocation.Type == PARAM_LOCATION_TYPE_BODY {
			err := json.Unmarshal([]byte(msgText), bodyObj)
			if err == nil {
				objValue.Field(i).Elem().SetPointer(&bodyObj)
			}
		}
		if err !=nil {
			break
		}
	}

	return _struct.MsgRequestInfo{Msg: paramObj}, err

}

func ConvertRpcRequestInfo(msgText string, paramObj interface{}, bodyObj interface{}) (_struct.RpcRequestInfo, error) {

	objType := reflect.TypeOf(paramObj)
	objValue := reflect.ValueOf(paramObj)

	var err error = nil
	for i := 0; i < objType.NumField(); i++ {

		field := objType.Field(i)
		paramLocationStr := field.Tag.Get(PARAM_LOCATION)

		if paramLocationStr == "" {
			continue
		}

		p := properties.MustLoadString(paramLocationStr)
		var paramLocation ParamLocation
		if err := p.Decode(&paramLocation); err != nil {
			return _struct.RpcRequestInfo{},err
		}

		if paramLocation.Type == PARAM_LOCATION_TYPE_BODY {
			err := json.Unmarshal([]byte(msgText), bodyObj)
			if err == nil {
				objValue.Field(i).Elem().SetPointer(&bodyObj)
			}
		}
		if err !=nil {
			break
		}
	}

	return _struct.RpcRequestInfo{ paramObj}, err

}

func setFieldValue(value reflect.Value, _type reflect.Type, fieldValue string) error {
	var err error = nil

	switch _type.Kind() {
	case reflect.Int:
		fallthrough
	case reflect.Int32:
		fallthrough
	case reflect.Int64:
		if value.CanSet() {
			intValue, err := strconv.ParseInt(fieldValue, 10, 64)
			if err == nil {
				value.SetInt(intValue)
			}
		}
	case reflect.Float32:
		fallthrough
	case reflect.Float64:
		if value.CanSet() {
			floatValue, err := strconv.ParseFloat(fieldValue, 64)
			if err == nil {
				value.SetFloat(floatValue)
			}
		}
	case reflect.Bool:
		if value.CanSet() {
			boolValue, err := strconv.ParseBool(fieldValue)
			if err == nil {
				value.SetBool(boolValue)
			}
		}
	case reflect.String:
		if value.CanSet() {
			value.SetString(fieldValue)
		}
	}
	return err
}
