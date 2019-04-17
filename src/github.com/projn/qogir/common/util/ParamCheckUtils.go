package util

import (
	"encoding/json"
	"errors"
	"github.com/magiconair/properties"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	PARAM_TYPE_NULL = iota
	PARAM_TYPE_JSON
	PARAM_TYPE_IP_ADDRESS
	PARAM_TYPE_PATH
	PARAM_TYPE_REGEX
)

const PARAM_LIMIT = "ParamLimit"

type ParamLimit struct {
	Type      int
	Regex     string
	Value     []string
	MaxValue  string
	MinValue  string
	MaxLength int
	MinLength int
	Precision int
}

func CheckParam(obj interface{}) error {
	var e error = nil

	objType := reflect.TypeOf(obj)
	objValue := reflect.ValueOf(obj)
	for i := 0 ;i<objType.NumField(); i++ {

		field := objType.Field(i)
		paramLimitStr := field.Tag.Get(PARAM_LIMIT)

		if paramLimitStr=="" {
			continue
		}

		p := properties.MustLoadString(paramLimitStr)
		var paramLimit ParamLimit
		if err := p.Decode(&paramLimit); err != nil {
			return err
		}

		switch field.Type.Kind() {
		case reflect.Int:
			value := int(objValue.Field(i).Int())
			e = checkInt(value, paramLimit)
		case reflect.Float32:
			value := float32(objValue.Field(i).Float())
			e = checkFloat32(value, paramLimit);
		case reflect.Float64:
			value := float64(objValue.Field(i).Float())
			e = checkFloat64(value, paramLimit);
		case reflect.String:
			value := objValue.Field(i).String()
			e = checkString(value, paramLimit);
		case reflect.Ptr:
			for j:=0;j<objValue.Field(i).Len();j++ {
				subValue := objValue.Field(i).Slice(i, i+1)
				e = CheckParam(subValue)

				if e != nil {
					continue
				}
			}
		default:
			continue
		}

	}

	return e
}

func checkInt(value int, paramLimit ParamLimit) error {
	strVal := strconv.Itoa(value)
	if "" != paramLimit.MaxValue {
		maxVal, e := strconv.Atoi(paramLimit.MaxValue);
		if e != nil {
			return e
		}
		if value > maxVal {
			return errors.New("Field is larger than " + paramLimit.MaxValue + ", value(" + strVal + ").");
		}
	}
	if "" != paramLimit.MinValue {
		minVal, e := strconv.Atoi(paramLimit.MinValue);
		if e != nil {
			return e
		}
		if value < minVal {
			return errors.New("Field is smaller than " + paramLimit.MinValue + ", value(" + strVal + ").");
		}
	}

	if paramLimit.Value != nil && len(paramLimit.Value) != 0 {
		var exist bool = false
		for _, val := range paramLimit.Value {
			if val == strVal {
				exist = true
			}
		}

		if !exist {
			return errors.New("Field is invaild, value(" + strVal + ").");
		}
	}
	return nil
}

func checkFloat32(value float32, paramLimit ParamLimit) error {
	strVal := strconv.FormatFloat(float64(value), 'f', 6, 32)
	if "" != paramLimit.MaxValue {
		maxVal, e := strconv.ParseFloat(paramLimit.MaxValue, 32);
		if e != nil {
			return e
		}
		if value > float32(maxVal) {
			return errors.New("Field is larger than " + paramLimit.MaxValue + ", value(" + strVal + ").");
		}
	}
	if "" != paramLimit.MinValue {
		minVal, e := strconv.Atoi(paramLimit.MinValue);
		if e != nil {
			return e
		}
		if value < float32(minVal) {
			return errors.New("Field is smaller than " + paramLimit.MinValue + ", value(" + strVal + ").");
		}
	}

	if paramLimit.Precision != 0 {
		valSlice := strings.Split(strVal, ".");
		if len(valSlice) == 2 && len(valSlice[1]) != paramLimit.Precision {
			return errors.New("Field precision do not equal with " + strconv.Itoa(paramLimit.Precision) + ", value(" + strVal + ").");
		}
	}

	if paramLimit.Value != nil && len(paramLimit.Value) != 0 {
		var exist bool = false
		for _, val := range paramLimit.Value {
			floatVal, e := strconv.ParseFloat(val, 32)
			if e != nil {
				return e
			}
			if float32(floatVal) == value {
				exist = true
			}
		}

		if !exist {
			return errors.New("Field is invaild, value(" + strVal + ").");
		}
	}
	return nil
}

func checkFloat64(value float64, paramLimit ParamLimit) error {
	strVal := strconv.FormatFloat(float64(value), 'f', 6, 64)
	if "" != paramLimit.MaxValue {
		maxVal, e := strconv.ParseFloat(paramLimit.MaxValue, 64);
		if e != nil {
			return e
		}
		if value > float64(maxVal) {
			return errors.New("Field is larger than " + paramLimit.MaxValue + ", value(" + strVal + ").");
		}
	}
	if "" != paramLimit.MinValue {
		minVal, e := strconv.Atoi(paramLimit.MinValue);
		if e != nil {
			return e
		}
		if value < float64(minVal) {
			return errors.New("Field is smaller than " + paramLimit.MinValue + ", value(" + strVal + ").");
		}
	}

	if paramLimit.Precision != 0 {
		valSlice := strings.Split(strVal, ".");
		if len(valSlice) == 2 && len(valSlice[1]) != paramLimit.Precision {
			return errors.New("Field precision do not equal with " + strconv.Itoa(paramLimit.Precision) + ", value(" + strVal + ").");
		}
	}

	if paramLimit.Value != nil && len(paramLimit.Value) != 0 {
		var exist bool = false
		for _, val := range paramLimit.Value {
			floatVal, e := strconv.ParseFloat(val, 64)
			if e != nil {
				return e
			}
			if float64(floatVal) == value {
				exist = true
			}
		}

		if !exist {
			return errors.New("Field is invaild, value(" + strVal + ").");
		}
	}
	return nil
}

func checkString(value string, paramLimit ParamLimit) error {
	if paramLimit.MaxLength != -1 {
		if len(value) > paramLimit.MaxLength {
			return errors.New("Field length is larger than " + strconv.Itoa(paramLimit.MaxLength) + ", value(" + value + ").");
		}
	}
	if paramLimit.MinLength != -1 {
		if len(value) < paramLimit.MinLength {
			return errors.New("Field length is smaller than " + strconv.Itoa(paramLimit.MinLength) + ", value(" + value + ").");
		}
	}

	if paramLimit.Value != nil && len(paramLimit.Value) != 0 {
		var exist bool = false
		for _, val := range paramLimit.Value {
			if val == value {
				exist = true
			}
		}

		if !exist {
			return errors.New("Field is invaild, value(" + value + ").");
		}
	}

	if paramLimit.Type != PARAM_TYPE_NULL {
		switch paramLimit.Type {
		case PARAM_TYPE_JSON:
			if !isValidJson(value) {
				return errors.New("Field is not json, value(" + value + ").")
			}
		case PARAM_TYPE_IP_ADDRESS:
			if !isValidIpAddress(value) {
				return errors.New("Field is not ip, value(" + value + ").")
			}
		case PARAM_TYPE_PATH:
			if !isValidPath(value) {
				return errors.New("Field is not path, value(" + value + ").")
			}
		case PARAM_TYPE_REGEX:
			if paramLimit.Regex != "" {
				match, _ := regexp.MatchString(paramLimit.Regex, value)
				if !match {
					return errors.New("Field is not match regex, value(" + value + ").");
				}
			}
		default:
		}
	}
	return nil
}

func isValidIpAddress(ip string) bool {
	const regex string = "(2[5][0-5]|2[0-4]\\d|1\\d{2}|\\d{1,2})\\.(25[0-5]|2[0-4]\\d|1\\d{2}|\\d{1,2})\\.(25[0-5]|2[0-4]\\d|1\\d{2}|\\d{1,2})\\.(25[0-5]|2[0-4]\\d|1\\d{2}|\\d{1,2})";
	match, _ := regexp.MatchString(regex, ip)
	return match
}

func isValidPath(path string) bool {
	const regex string = "[a-zA-Z]:(\\\\([0-9a-zA-Z]+))+|(\\/([0-9a-zA-Z]+))+";
	match, _ := regexp.MatchString(regex, path)
	return match
}

func isValidJson(jsonStr string) bool {
	var data interface{}
	e := json.Unmarshal([]byte(jsonStr), &data)
	if e != nil {
		return true
	} else {
		return false
	}
}
