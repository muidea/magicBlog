package core

import (
	"encoding/json"
	"log"
	"net/http"

	common_const "muidea.com/magicCommon/common"
	common_def "muidea.com/magicCommon/def"
	"muidea.com/magicCommon/foundation/net"
)

func (s *Blog) statusAction(res http.ResponseWriter, req *http.Request) {
	log.Print("statusAction")

	result := common_def.StatusAccountResult{}
	for {
		authToken := req.URL.Query().Get(common_const.AuthToken)
		sessionID := req.URL.Query().Get(common_const.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("statusAccount failed, illegal authToken or sessionID")
			result.ErrorCode = common_def.Failed
			result.Reason = "无效Token或会话"
			break
		}

		userView, sessionID, ok := s.centerAgent.StatusAccount(authToken, sessionID)
		if !ok {
			log.Print("statusAccount failed, illegal authToken or sessionID")
			result.ErrorCode = common_def.Failed
			result.Reason = "无效Token或会话"
			break
		}

		result.OnlineEntry = userView
		result.SessionID = sessionID
		result.ErrorCode = common_def.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Blog) loginAction(res http.ResponseWriter, req *http.Request) {
	log.Print("loginAction")

	param := &common_def.LoginAccountParam{}
	result := common_def.LoginAccountResult{}
	for {
		err := net.ParsePostJSON(req, param)
		if err != nil {
			log.Printf("ParsePostJSON failed, err:%s", err.Error())
			result.ErrorCode = common_def.Failed
			result.Reason = "非法请求"
			break
		}

		userView, authToken, sessionID, ok := s.centerAgent.LoginAccount(param.Account, param.Password)
		if !ok {
			log.Print("login failed, illegal account or password")
			result.ErrorCode = common_def.Failed
			result.Reason = "无效账号或密码"
			break
		}

		result.OnlineEntry = userView
		result.AuthToken = authToken
		result.SessionID = sessionID
		result.ErrorCode = common_def.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Blog) logoutAction(res http.ResponseWriter, req *http.Request) {
	log.Print("logoutAction")

	result := common_def.LogoutAccountResult{}
	for {
		authToken := req.URL.Query().Get(common_const.AuthToken)
		sessionID := req.URL.Query().Get(common_const.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("logout failed, illegal authToken or sessionID")
			result.ErrorCode = common_def.Failed
			result.Reason = "无效Token或会话"
			break
		}

		ok := s.centerAgent.LogoutAccount(authToken, sessionID)
		if !ok {
			log.Print("logout failed, illegal authToken or sessionID")
			result.ErrorCode = common_def.Failed
			result.Reason = "非法Token或会话"
			break
		}

		result.ErrorCode = common_def.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}
