package core

import (
	"encoding/json"
	"log"
	"net/http"

	"muidea.com/magicCenter/application/common"
	common_result "muidea.com/magicCommon/common"
	"muidea.com/magicCommon/foundation/net"
	"muidea.com/magicCommon/model"
)

func (s *Blog) statusAction(res http.ResponseWriter, req *http.Request) {
	log.Print("statusAction")

	type statusResult struct {
		common_result.Result
		OnlineUser model.AccountOnlineView `json:"onlineUser"`
	}

	result := statusResult{}
	for {
		authToken := req.URL.Query().Get(common.AuthTokenID)
		sessionID := req.URL.Query().Get(common.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("statusAccount failed, illegal authToken or sessionID")
			result.ErrorCode = common_result.Failed
			result.Reason = "无效Token或会话"
			break
		}

		userView, ok := s.centerAgent.StatusAccount(authToken, sessionID)
		if !ok {
			log.Print("statusAccount failed, illegal authToken or sessionID")
			result.ErrorCode = common_result.Failed
			result.Reason = "无效Token或会话"
			break
		}

		result.OnlineUser = userView
		result.ErrorCode = common_result.Success
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

	type loginParam struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}
	type loginResult struct {
		common_result.Result
		OnlineUser model.AccountOnlineView `json:"onlineUser"`
		AuthToken  string                  `json:"authToken"`
		SessionID  string                  `json:"sessionID"`
	}

	param := &loginParam{}
	result := loginResult{}
	for {
		err := net.ParsePostJSON(req, param)
		if err != nil {
			log.Printf("ParsePostJSON failed, err:%s", err.Error())
			result.ErrorCode = common_result.Failed
			result.Reason = "非法请求"
			break
		}

		userView, authToken, sessionID, ok := s.centerAgent.LoginAccount(param.Account, param.Password)
		if !ok {
			log.Print("login failed, illegal account or password")
			result.ErrorCode = common_result.Failed
			result.Reason = "无效账号或密码"
			break
		}

		result.OnlineUser = userView
		result.AuthToken = authToken
		result.SessionID = sessionID
		result.ErrorCode = common_result.Success
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

	type logoutResult struct {
		common_result.Result
	}

	result := logoutResult{}
	for {
		authToken := req.URL.Query().Get(common.AuthTokenID)
		sessionID := req.URL.Query().Get(common.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("statusAccount failed, illegal authToken or sessionID")
			result.ErrorCode = common_result.Failed
			result.Reason = "无效Token或会话"
			break
		}

		ok := s.centerAgent.LogoutAccount(authToken, sessionID)
		if !ok {
			log.Print("logout failed, illegal authToken or sessionID")
			result.ErrorCode = common_result.Failed
			result.Reason = "非法Token或会话"
			break
		}

		result.ErrorCode = common_result.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

type itemInfo struct {
	ID      int         `json:"id"`
	Name    string      `json:"name"`
	Type    string      `json:"type"`
	SubItem interface{} `json:"subItem"`
}

func (s *Blog) fetchSubItem(id int) []itemInfo {
	itemList := []itemInfo{}

	subItem := s.centerAgent.QuerySummary(id)
	for _, val := range subItem {
		info := itemInfo{}
		info.ID = val.ID
		info.Name = val.Name
		info.Type = val.Type

		if val.Type == model.CATALOG {
			subList := s.fetchSubItem(val.ID)
			info.SubItem = subList
		}

		itemList = append(itemList, info)
	}

	return itemList
}

func (s *Blog) summaryAction(res http.ResponseWriter, req *http.Request) {
	log.Print("summaryAction")

	type summaryResult struct {
		common_result.Result
		ItemList []itemInfo `json:"itemList"`
	}

	result := summaryResult{}
	for {
		authToken := req.URL.Query().Get(common.AuthTokenID)
		sessionID := req.URL.Query().Get(common.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("statusAccount failed, illegal authToken or sessionID")
			result.ErrorCode = common_result.Failed
			result.Reason = "无效Token或会话"
			break
		}

		for _, val := range s.blogContent {
			info := itemInfo{}
			info.ID = val.ID
			info.Name = val.Name
			info.Type = val.Type

			if val.Type == model.CATALOG {
				subList := s.fetchSubItem(val.ID)
				info.SubItem = subList
			}

			result.ItemList = append(result.ItemList, info)
		}

		result.ErrorCode = common_result.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	log.Print("summaryAction, json.Marshal, failed, err:" + err.Error())
	http.Redirect(res, req, "/default/index.html", http.StatusMovedPermanently)
}

func (s *Blog) catalogCreateAction(res http.ResponseWriter, req *http.Request) {
	log.Print("catalogCreateAction")

	type catalogParam struct {
		Name        string        `json:"name"`
		Description string        `json:"description"`
		Parent      model.Catalog `json:"parent"`
	}

	type catalogResult struct {
		common_result.Result
	}

	param := &catalogParam{}
	result := catalogResult{}
	for {
		authToken := req.URL.Query().Get(common.AuthTokenID)
		sessionID := req.URL.Query().Get(common.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("create catalog failed, illegal authToken or sessionID")
			result.ErrorCode = common_result.Failed
			result.Reason = "无效Token或会话"
			break
		}

		err := net.ParsePostJSON(req, param)
		if err != nil {
			log.Printf("ParsePostJSON failed, err:%s", err.Error())
			result.ErrorCode = common_result.Failed
			result.Reason = "非法请求"
			break
		}

		ok := s.centerAgent.CreateCatalog(param.Name, param.Description, []model.Catalog{param.Parent}, authToken, sessionID)
		if !ok {
			log.Print("login failed, illegal account or password")
			result.ErrorCode = common_result.Failed
			result.Reason = "新建分类失败"
			break
		}

		result.ErrorCode = common_result.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Blog) articleCreateAction(res http.ResponseWriter, req *http.Request) {
	log.Print("articleCreateAction")

	type articleParam struct {
		Title   string        `json:"title"`
		Content string        `json:"content"`
		Catalog model.Catalog `json:"catalog"`
	}

	type articleResult struct {
		common_result.Result
	}

	param := &articleParam{}
	result := articleResult{}
	for {
		authToken := req.URL.Query().Get(common.AuthTokenID)
		sessionID := req.URL.Query().Get(common.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("create article failed, illegal authToken or sessionID")
			result.ErrorCode = common_result.Failed
			result.Reason = "无效Token或会话"
			break
		}

		err := net.ParsePostJSON(req, param)
		if err != nil {
			log.Printf("ParsePostJSON failed, err:%s", err.Error())
			result.ErrorCode = common_result.Failed
			result.Reason = "非法请求"
			break
		}

		ok := s.centerAgent.CreateArticle(param.Title, param.Content, []model.Catalog{param.Catalog}, authToken, sessionID)
		if !ok {
			log.Print("login failed, illegal account or password")
			result.ErrorCode = common_result.Failed
			result.Reason = "无效账号或密码"
			break
		}

		result.ErrorCode = common_result.Success
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}
