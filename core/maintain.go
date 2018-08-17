package core

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	common_const "muidea.com/magicCommon/common"
	common_def "muidea.com/magicCommon/def"
	"muidea.com/magicCommon/foundation/net"
	"muidea.com/magicCommon/model"
)

type itemInfo struct {
	ID      int         `json:"id"`
	Name    string      `json:"name"`
	Type    string      `json:"type"`
	Deep    int         `json:"deep"`
	SubItem interface{} `json:"subItem"`
}

func (s *Blog) fetchSubItem(id, curDeep int, authToken, sessionID string) []itemInfo {
	itemList := []itemInfo{}

	summary := model.CatalogUnit{ID: id, Type: model.CATALOG}
	subItem := s.centerAgent.QuerySummaryContent(summary, authToken, sessionID)
	for _, val := range subItem {
		info := itemInfo{}
		info.ID = val.ID
		info.Name = val.Name
		info.Type = val.Type
		info.Deep = curDeep + 1

		if val.Type == model.CATALOG {
			subList := s.fetchSubItem(val.ID, curDeep+1, authToken, sessionID)
			info.SubItem = subList
		}

		itemList = append(itemList, info)
	}

	return itemList
}

func (s *Blog) summaryAction(res http.ResponseWriter, req *http.Request) {
	log.Print("summaryAction")

	type summaryResult struct {
		common_def.Result
		ItemList []itemInfo `json:"itemList"`
	}

	result := summaryResult{}
	for {
		authToken := req.URL.Query().Get(common_const.AuthToken)
		sessionID := req.URL.Query().Get(common_const.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("summaryAction, get summry failed, illegal authToken or sessionID")
			result.ErrorCode = common_def.Failed
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
				subList := s.fetchSubItem(val.ID, curDeep, authToken, sessionID)
				info.SubItem = subList
			}

			result.ItemList = append(result.ItemList, info)
		}

		result.ErrorCode = common_def.Success
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

	param := &common_def.CreateCatalogParam{}
	result := common_def.CreateCatalogResult{}
	for {
		authToken := req.URL.Query().Get(common_const.AuthToken)
		sessionID := req.URL.Query().Get(common_const.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("catalogCreateAction, create catalog failed, illegal authToken or sessionID")
			result.ErrorCode = common_def.Failed
			result.Reason = "无效Token或会话"
			break
		}

		err := net.ParsePostJSON(req, param)
		if err != nil {
			log.Printf("catalogCreateAction, ParsePostJSON failed, err:%s", err.Error())
			result.ErrorCode = common_def.Failed
			result.Reason = "非法请求"
			break
		}

		catalog, ok := s.centerAgent.CreateCatalog(param.Name, param.Description, param.Catalog, authToken, sessionID)
		if !ok {
			log.Print("catalogCreateAction, create catalog failed")
			result.ErrorCode = common_def.Failed
			result.Reason = "新建分类失败"
			break
		}

		result.ErrorCode = common_def.Success
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

	param := &common_def.UpdateCatalogParam{}
	result := common_def.UpdateCatalogResult{}
	for {
		authToken := req.URL.Query().Get(common_const.AuthToken)
		sessionID := req.URL.Query().Get(common_const.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("catalogUpdateAction, update catalog failed, illegal authToken or sessionID")
			result.ErrorCode = common_def.Failed
			result.Reason = "无效Token或会话"
			break
		}

		_, value := net.SplitRESTAPI(req.URL.Path)
		id, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("catalogUpdateAction, update catalog failed, illegal id, id:%s, err:%s", value, err.Error())
			result.ErrorCode = common_def.Failed
			result.Reason = "非法参数"
			break
		}

		err = net.ParsePostJSON(req, param)
		if err != nil {
			log.Printf("catalogUpdateAction, ParsePostJSON failed, err:%s", err.Error())
			result.ErrorCode = common_def.Failed
			result.Reason = "非法请求"
			break
		}

		catalog, ok := s.centerAgent.UpdateCatalog(id, param.Name, param.Description, param.Catalog, authToken, sessionID)
		if !ok {
			log.Print("catalogUpdateAction, update catalog failed")
			result.ErrorCode = common_def.Failed
			result.Reason = "更新分类失败"
			break
		}

		result.ErrorCode = common_def.Success
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

	result := common_def.QueryCatalogResult{}
	for {
		authToken := req.URL.Query().Get(common_const.AuthToken)
		sessionID := req.URL.Query().Get(common_const.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("catalogQueryAction, query catalog failed, illegal authToken or sessionID")
			result.ErrorCode = common_def.Failed
			result.Reason = "无效Token或会话"
			break
		}
		_, value := net.SplitRESTAPI(req.URL.Path)
		id, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("catalogQueryAction, query catalog failed, illegal id, id:%s, err:%s", value, err.Error())
			result.ErrorCode = common_def.Failed
			result.Reason = "非法参数"
			break
		}

		catalog, ok := s.centerAgent.QueryCatalog(id, authToken, sessionID)
		if !ok {
			log.Print("catalogQueryAction, query catalog failed, illegal id or no exist")
			result.ErrorCode = common_def.NoExist
			result.Reason = "对象不存在"
			break
		}

		result.Catalog = catalog
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

func (s *Blog) articleCreateAction(res http.ResponseWriter, req *http.Request) {
	log.Print("articleCreateAction")

	param := &common_def.CreateArticleParam{}
	result := common_def.CreateArticleResult{}
	for {
		authToken := req.URL.Query().Get(common_const.AuthToken)
		sessionID := req.URL.Query().Get(common_const.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("articleCreateAction, create article failed, illegal authToken or sessionID")
			result.ErrorCode = common_def.Failed
			result.Reason = "无效Token或会话"
			break
		}

		err := net.ParsePostJSON(req, param)
		if err != nil {
			log.Printf("articleCreateAction, ParsePostJSON failed, err:%s", err.Error())
			result.ErrorCode = common_def.Failed
			result.Reason = "非法请求"
			break
		}

		article, ok := s.centerAgent.CreateArticle(param.Title, param.Content, param.Catalog, authToken, sessionID)
		if !ok {
			log.Print("articleCreateAction, create article failed")
			result.ErrorCode = common_def.Failed
			result.Reason = "新建文章失败"
			break
		}

		result.ErrorCode = common_def.Success
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

	param := &common_def.UpdateArticleParam{}
	result := common_def.UpdateArticleResult{}
	for {
		authToken := req.URL.Query().Get(common_const.AuthToken)
		sessionID := req.URL.Query().Get(common_const.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("articleUpdateAction, update article failed, illegal authToken or sessionID")
			result.ErrorCode = common_def.Failed
			result.Reason = "无效Token或会话"
			break
		}

		_, value := net.SplitRESTAPI(req.URL.Path)
		id, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("articleUpdateAction, update article failed, illegal id, id:%s, err:%s", value, err.Error())
			result.ErrorCode = common_def.Failed
			result.Reason = "非法参数"
			break
		}

		err = net.ParsePostJSON(req, param)
		if err != nil {
			log.Printf("articleUpdateAction, ParsePostJSON failed, err:%s", err.Error())
			result.ErrorCode = common_def.Failed
			result.Reason = "非法请求"
			break
		}

		article, ok := s.centerAgent.UpdateArticle(id, param.Title, param.Content, param.Catalog, authToken, sessionID)
		if !ok {
			log.Print("articleUpdateAction, update article failed")
			result.ErrorCode = common_def.Failed
			result.Reason = "更新文章失败"
			break
		}

		result.ErrorCode = common_def.Success
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

	result := common_def.DestroyCatalogResult{}
	for {
		authToken := req.URL.Query().Get(common_const.AuthToken)
		sessionID := req.URL.Query().Get(common_const.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("catalogDeleteAction, delete catalog failed, illegal authToken or sessionID")
			result.ErrorCode = common_def.Failed
			result.Reason = "无效Token或会话"
			break
		}

		_, value := net.SplitRESTAPI(req.URL.Path)
		id, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("catalogDeleteAction, strconv.Atoi failed, err:%s", err.Error())
			result.ErrorCode = common_def.IllegalParam
			result.Reason = "非法参数"
			break
		}

		ok := s.centerAgent.DeleteCatalog(id, authToken, sessionID)
		if !ok {
			log.Printf("catalogDeleteAction, delete catalog failed, id=%d", id)
			result.ErrorCode = common_def.Failed
			result.Reason = "删除对象失败"
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

func (s *Blog) articleDeleteAction(res http.ResponseWriter, req *http.Request) {
	log.Print("articleDeleteAction")

	result := common_def.DestoryArticleResult{}
	for {
		authToken := req.URL.Query().Get(common_const.AuthToken)
		sessionID := req.URL.Query().Get(common_const.SessionID)
		if len(authToken) == 0 || len(sessionID) == 0 {
			log.Print("articleDeleteAction, delete article failed, illegal authToken or sessionID")
			result.ErrorCode = common_def.Failed
			result.Reason = "无效Token或会话"
			break
		}

		_, value := net.SplitRESTAPI(req.URL.Path)
		id, err := strconv.Atoi(value)
		if err != nil {
			log.Printf("articleDeleteAction, strconv.Atoi failed, err:%s", err.Error())
			result.ErrorCode = common_def.IllegalParam
			result.Reason = "非法参数"
			break
		}

		ok := s.centerAgent.DeleteArticle(id, authToken, sessionID)
		if !ok {
			log.Printf("articleDeleteAction, delete article failed, illegal id, id=%d", id)
			result.ErrorCode = common_def.Failed
			result.Reason = "删除对象失败"
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
