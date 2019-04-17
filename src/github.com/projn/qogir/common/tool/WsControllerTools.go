package tool
import (
	"encoding/json"
	"github.com/kataras/iris"
	"../initialize"
	"../struct"
	"../msg/request"
	"../util"
	"../define"
	"../../bean"
	"github.com/kataras/iris/websocket"
	"reflect"
)


func DealWsRequest(textMsg string, paramObj interface{}, bodyObj interface{}, connection websocket.Connection) bool {
	if textMsg == "" || connection == nil {
		bean.Logger.Error("Invaild param.")
		return false
	}

	wsRequestMsgInfo := request.WsRequestMsgInfo{}
	err := json.Unmarshal([]byte(textMsg), wsRequestMsgInfo)
	if err != nil {
		bean.Logger.Error("Invaild web socket request msg,msg info(" + err.Error() + ").")
		return false
	}

	method := initialize.SERVICE_METHOD_WS_POST
	requestServiceInfo, e :=initialize.GetRequestServiceInfo(wsRequestMsgInfo.MsgId, method)
	if e != nil {
		bean.Logger.Errorf("Invaild request service info, msg id(%s), method(%s).",
			wsRequestMsgInfo.MsgId, method)
		return false
	}

	if requestServiceInfo.Type == "" || requestServiceInfo.Type != initialize.SERVICE_TYPE_WS {
		bean.Logger.Error("Invaild request service type info, type(" + requestServiceInfo.Type + ").")
		return false
	}

	if requestServiceInfo.AuthorizationFilter != nil {
		e := requestServiceInfo.AuthorizationFilter(connection, requestServiceInfo.UserRoleNameList)
		if e!= nil {
			bean.Logger.Error("Check authorization info error,error info(" + e.Error() + ").")
			return false
		}
	}


	wsRequestInfo, e := util.ConvertWsRequestInfo(connection, textMsg, paramObj, bodyObj)
	if e != nil {
		bean.Logger.Error("Convert param error,error info(" + e.Error() + ").")
		return false
	}

	if wsRequestInfo.ParamObj!= nil {
		e = util.CheckParam(wsRequestInfo.ParamObj)
		if e != nil {
			bean.Logger.Error("Check param error,error info(" + e.Error() + ").")
			return false
		}
	}

	wsResponseInfo := requestServiceInfo.WsServiceFunction(wsRequestInfo)
	if wsResponseInfo.Msg != nil {
		data, e :=json.Marshal(wsResponseInfo.Msg)
		if e != nil {
			bean.Logger.Error("Convert json info error,error info(" + e.Error() + ").")
			return false
		}
		connection.To(websocket.Broadcast).EmitMessage(data)
	}

	if wsResponseInfo.ExtendInfoMap != nil {
		for k, v := range wsResponseInfo.ExtendInfoMap {
			connection.SetValue(k, v)
		}

		agentId := wsResponseInfo.ExtendInfoMap[define.AGENT_ID_KEY];
		if agentId == "" {
			if (!WsSessionInfoMap.getInstance().addWebSocketSessionInfo(agentId, session)) {
			LOGGER.error("Add web socket session to pool error, pool size({}).",
			WsSessionInfoMap.getInstance().getPoolSize());
			return;
		}

		IAgentMasterInfoDao agentMasterInfoDao = InitializeBean.getBean(AgentMasterInfoDaoImpl.class);
		if(agentMasterInfoDao!=null) {
		String url = ServiceData.getMasterInfo().isServerSsl()? HTTP_URL_HEADER: HTTPS_URL_HEADER
		+ ServiceData.getMasterInfo().getServerIp() + ":"
		+ ServiceData.getMasterInfo().getServerPort() + "/"
		+ API_URL_HEADER + HTTP_API_SERVICE_SEND_WS_MSG;
		agentMasterInfoDao.setAgentMasterInfo(new AgentMasterInfo(agentId,
		ServiceData.getMasterInfo().getServerIp(),
		ServiceData.getMasterInfo().getServerPort(),url));
	}
	}
	}

}