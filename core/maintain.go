package core

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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
			log.Print("logout failed, illegal authToken or sessionID")
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
	Deep    int         `json:"deep"`
	SubItem interface{} `json:"subItem"`
}

func (s *Blog) fetchSubItem(id, curDeep int) []itemInfo {
	itemList := []itemInfo{}

	subItem := s.centerAgent.QuerySummary(id)
	for _, val := range subItem {
		info := itemInfo{}
		info.ID = val.ID
		info.Name = val.Name
		info.Type = val.Type
		info.Deep = curDeep + 1

		if val.Type == model.CATALOG {
			subList := s.fetchSubItem(val.ID, curDeep+1)
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
			log.Print("summaryAction, get summry failed, illegal authToken or sessionID")
			result.ErrorCode = common_result.Failed
			result.Reason = "无效Token或会话"
			break
		}

		curDeep := 0
		for _, val := range s.blogContent {
			info := itemInfo{}
			info.ID = val.ID
			info.Name = val.Name
			info.Type = val.Type
			info.Deep = curDeep

			if val.Type == model.CATALOG {
				subList := s.fetchSubItem(val.ID, curDeep)
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
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Catalog     []model.Catalog `json:"catalog"`
	}

	type catalogResult struct {
		common_result.Result
		Catalog model.SummaryView `json:"catalog"`
	}

	param := &catalogParam{}
	result := catalogResult{}
	for {
		authToken := req.URL.Query().Get(common.AuthTokenID)
		sessionID := req.URL.Query().Get(common.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("catalogCreateAction, create catalog failed, illegal authToken or sessionID")
			result.ErrorCode = common_result.Failed
			result.Reason = "无效Token或会话"
			break
		}

		err := net.ParsePostJSON(req, param)
		if err != nil {
			log.Printf("catalogCreateAction, ParsePostJSON failed, err:%s", err.Error())
			result.ErrorCode = common_result.Failed
			result.Reason = "非法请求"
			break
		}

		catalog, ok := s.centerAgent.CreateCatalog(param.Name, param.Description, param.Catalog, authToken, sessionID)
		if !ok {
			log.Print("catalogCreateAction, create catalog failed")
			result.ErrorCode = common_result.Failed
			result.Reason = "新建分类失败"
			break
		}

		result.ErrorCode = common_result.Success
		result.Catalog = catalog
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Blog) catalogUpdateAction(res http.ResponseWriter, req *http.Request) {
	log.Print("catalogUpdateAction")

	type catalogParam struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Catalog     []model.Catalog `json:"catalog"`
	}

	type catalogResult struct {
		common_result.Result
		Catalog model.SummaryView `json:"catalog"`
	}

	param := &catalogParam{}
	result := catalogResult{}
	for {
		authToken := req.URL.Query().Get(common.AuthTokenID)
		sessionID := req.URL.Query().Get(common.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("catalogUpdateAction, update catalog failed, illegal authToken or sessionID")
			result.ErrorCode = common_result.Failed
			result.Reason = "无效Token或会话"
			break
		}

		_, value := net.SplitRESTAPI(req.URL.Path)
		id, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("catalogUpdateAction, update catalog failed, illegal id, id:%s, err:%s", value, err.Error())
			result.ErrorCode = common_result.Failed
			result.Reason = "非法参数"
			break
		}

		err = net.ParsePostJSON(req, param)
		if err != nil {
			log.Printf("catalogUpdateAction, ParsePostJSON failed, err:%s", err.Error())
			result.ErrorCode = common_result.Failed
			result.Reason = "非法请求"
			break
		}

		catalog, ok := s.centerAgent.UpdateCatalog(id, param.Name, param.Description, param.Catalog, authToken, sessionID)
		if !ok {
			log.Print("catalogUpdateAction, update catalog failed")
			result.ErrorCode = common_result.Failed
			result.Reason = "更新分类失败"
			break
		}

		result.ErrorCode = common_result.Success
		result.Catalog = catalog
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Blog) catalogQueryAction(res http.ResponseWriter, req *http.Request) {
	log.Print("catalogQueryAction")

	type catalogResult struct {
		common_result.Result
		Content model.CatalogDetailView `json:"content"`
	}

	result := catalogResult{}
	for {
		authToken := req.URL.Query().Get(common.AuthTokenID)
		sessionID := req.URL.Query().Get(common.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("catalogQueryAction, query catalog failed, illegal authToken or sessionID")
			result.ErrorCode = common_result.Failed
			result.Reason = "无效Token或会话"
			break
		}
		_, value := net.SplitRESTAPI(req.URL.Path)
		id, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("catalogQueryAction, query catalog failed, illegal id, id:%s, err:%s", value, err.Error())
			result.ErrorCode = common_result.Failed
			result.Reason = "非法参数"
			break
		}

		catalog, ok := s.centerAgent.QueryCatalog(id)
		if !ok {
			log.Print("catalogQueryAction, query catalog failed, illegal id or no exist")
			result.ErrorCode = common_result.NoExist
			result.Reason = "对象不存在"
			break
		}

		result.Content = catalog
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
		Title   string          `json:"title"`
		Content string          `json:"content"`
		Catalog []model.Catalog `json:"catalog"`
	}

	type articleResult struct {
		common_result.Result
		Article model.SummaryView `json:"article"`
	}

	param := &articleParam{}
	result := articleResult{}
	for {
		authToken := req.URL.Query().Get(common.AuthTokenID)
		sessionID := req.URL.Query().Get(common.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("articleCreateAction, create article failed, illegal authToken or sessionID")
			result.ErrorCode = common_result.Failed
			result.Reason = "无效Token或会话"
			break
		}

		err := net.ParsePostJSON(req, param)
		if err != nil {
			log.Printf("articleCreateAction, ParsePostJSON failed, err:%s", err.Error())
			result.ErrorCode = common_result.Failed
			result.Reason = "非法请求"
			break
		}

		article, ok := s.centerAgent.CreateArticle(param.Title, param.Content, param.Catalog, authToken, sessionID)
		if !ok {
			log.Print("articleCreateAction, create article failed")
			result.ErrorCode = common_result.Failed
			result.Reason = "新建文章失败"
			break
		}

		result.ErrorCode = common_result.Success
		result.Article = article
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Blog) articleUpdateAction(res http.ResponseWriter, req *http.Request) {
	log.Print("articleUpdateAction")

	type articleParam struct {
		Title   string          `json:"title"`
		Content string          `json:"content"`
		Catalog []model.Catalog `json:"catalog"`
	}

	type articleResult struct {
		common_result.Result
		Article model.SummaryView `json:"article"`
	}

	param := &articleParam{}
	result := articleResult{}
	for {
		authToken := req.URL.Query().Get(common.AuthTokenID)
		sessionID := req.URL.Query().Get(common.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("articleUpdateAction, update article failed, illegal authToken or sessionID")
			result.ErrorCode = common_result.Failed
			result.Reason = "无效Token或会话"
			break
		}

		_, value := net.SplitRESTAPI(req.URL.Path)
		id, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("articleUpdateAction, update article failed, illegal id, id:%s, err:%s", value, err.Error())
			result.ErrorCode = common_result.Failed
			result.Reason = "非法参数"
			break
		}

		err = net.ParsePostJSON(req, param)
		if err != nil {
			log.Printf("articleUpdateAction, ParsePostJSON failed, err:%s", err.Error())
			result.ErrorCode = common_result.Failed
			result.Reason = "非法请求"
			break
		}

		article, ok := s.centerAgent.UpdateArticle(id, param.Title, param.Content, param.Catalog, authToken, sessionID)
		if !ok {
			log.Print("articleUpdateAction, update article failed")
			result.ErrorCode = common_result.Failed
			result.Reason = "更新文章失败"
			break
		}

		result.ErrorCode = common_result.Success
		result.Article = article
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Blog) catalogDeleteAction(res http.ResponseWriter, req *http.Request) {
	log.Print("catalogDeleteAction")

	type catalogResult struct {
		common_result.Result
	}

	result := catalogResult{}
	for {
		authToken := req.URL.Query().Get(common.AuthTokenID)
		sessionID := req.URL.Query().Get(common.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("catalogDeleteAction, delete catalog failed, illegal authToken or sessionID")
			result.ErrorCode = common_result.Failed
			result.Reason = "无效Token或会话"
			break
		}

		_, value := net.SplitRESTAPI(req.URL.Path)
		id, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("catalogDeleteAction, strconv.Atoi failed, err:%s", err.Error())
			result.ErrorCode = common_result.IllegalParam
			result.Reason = "非法参数"
			break
		}

		ok := s.centerAgent.DeleteCatalog(id, authToken, sessionID)
		if !ok {
			log.Printf("catalogDeleteAction, delete catalog failed, id=%d", id)
			result.ErrorCode = common_result.Failed
			result.Reason = "删除对象失败"
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

func (s *Blog) articleDeleteAction(res http.ResponseWriter, req *http.Request) {
	log.Print("articleDeleteAction")

	type articleResult struct {
		common_result.Result
	}

	result := articleResult{}
	for {
		authToken := req.URL.Query().Get(common.AuthTokenID)
		sessionID := req.URL.Query().Get(common.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("articleDeleteAction, delete article failed, illegal authToken or sessionID")
			result.ErrorCode = common_result.Failed
			result.Reason = "无效Token或会话"
			break
		}

		_, value := net.SplitRESTAPI(req.URL.Path)
		id, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("articleDeleteAction, strconv.Atoi failed, err:%s", err.Error())
			result.ErrorCode = common_result.IllegalParam
			result.Reason = "非法参数"
			break
		}

		ok := s.centerAgent.DeleteArticle(id, authToken, sessionID)
		if !ok {
			log.Printf("articleDeleteAction, delete article failed, illegal id, id=%d", id)
			result.ErrorCode = common_result.Failed
			result.Reason = "删除对象失败"
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
