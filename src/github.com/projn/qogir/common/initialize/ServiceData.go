package initialize

import (
	"github.com/kataras/iris/core/errors"
	"../struct"
)

const SERVICE_METHOD_HTTP_POST = "post"
const SERVICE_METHOD_HTTP_GET = "get"
const SERVICE_METHOD_HTTP_PUT = "put"
const SERVICE_METHOD_HTTP_DELETE = "delete"
const SERVICE_METHOD_WS_POST = "post"
const SERVICE_METHOD_MSG_NORMAL = "normal"
const SERVICE_METHOD_MSG_ORDER = "order"
const SERVICE_METHOD_MSG_BROADCAST = "broadcast"
const SERVICE_TYPE_HTTP = "http"
const SERVICE_TYPE_WS = "ws"
const SERVICE_TYPE_MSG = "msg"

type RequestServiceInfo struct {
	ServiceName string 
	Method string
	Type string
	UserRoleNameList []string
	AuthorizationFilter func(interface{}, []string) error
	HttpServiceFunction func( httpRequestInfo _struct.HttpRequestInfo ) _struct.HttpResponseInfo
	MsgServiceFunction func(bornTimestamp int64 , msgRequestInfo _struct.MsgRequestInfo) bool
	WsServiceFunction func(wsRequestInfo _struct.WsRequestInfo) _struct.WsResponseInfo
	RpcServiceFunction func(rpcRequestInfo _struct.RpcRequestInfo) _struct.RpcResponseInfo
}

var RequestServiceInfoMap map[string]map[string]RequestServiceInfo = nil

func AddRequestServiceInfo(url string, method string, requestServiceInfo RequestServiceInfo) {
	if RequestServiceInfoMap == nil {
		RequestServiceInfoMap = make(map[string]map[string]RequestServiceInfo)
	}

	dataMap := map[string]RequestServiceInfo{
		method : requestServiceInfo,
	}
	RequestServiceInfoMap[url]=dataMap
}

func GetRequestServiceInfo(url string, method string) (RequestServiceInfo, error) {
	ret := RequestServiceInfo{}

	if url == "" || method == "" {
		return ret, errors.New("Invaild prarm.")
	}

	return RequestServiceInfoMap[url][method], nil
}
