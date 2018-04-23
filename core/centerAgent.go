package core

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	common_result "muidea.com/magicCommon/common"
	"muidea.com/magicCommon/model"
)

// Agent Center访问代理
type Agent interface {
	Start(bashURL, name, account, password string) bool
	Stop()
	CreateCatalog(name, description string) bool
	FetchCatalog(name string) (model.CatalogDetailView, bool)
	QuerySummary() []model.SummaryView
	QueryCatalog() []model.CatalogDetailView
	QuerySubCatalog(catalogID int) []model.CatalogDetailView
	QueryContent(id int) model.ArticleDetailView
}

// NewCenterAgent 新建Agent
func NewCenterAgent() Agent {
	return &center{}
}

type center struct {
	httpClient  *http.Client
	baseURL     string
	catalogName string
	account     string
	password    string

	onlineView model.AccountOnlineView
	sessionID  string
}

func (s *center) Start(bashURL, name, account, password string) bool {
	s.httpClient = &http.Client{}
	s.baseURL = bashURL
	s.catalogName = name
	s.account = account
	s.password = password

	if !s.login() {
		return false
	}

	log.Print("start centerAgent ok")
	return true
}

func (s *center) Stop() {

}

func (s *center) login() bool {
	type loginParam struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}

	type loginResult struct {
		common_result.Result
		OnlineUser model.AccountOnlineView `json:"onlineUser"`
		SessionID  string                  `json:"sessionID"`
	}

	param := loginParam{Account: s.account, Password: s.password}
	data, err := json.Marshal(param)
	if err != nil {
		log.Printf("marshal login param failed, err:%s", err.Error())
		return false
	}

	bufferReader := bytes.NewBuffer(data)
	url := fmt.Sprintf("%s/%s", s.baseURL, "cas/user/")
	request, err := http.NewRequest("POST", url, bufferReader)
	if err != nil {
		log.Printf("construct request failed, url:%s, err:%s", url, err.Error())
		return false
	}

	request.Header.Set("content-type", "application/json")
	response, err := s.httpClient.Do(request)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("login failed, statusCode:%d", response.StatusCode)
		return false
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed")
	}

	result := &loginResult{}
	err = json.Unmarshal(content, result)
	if result.ErrorCode == common_result.Success {
		s.onlineView = result.OnlineUser
		s.sessionID = result.SessionID
		return true
	}

	return false
}

func (s *center) CreateCatalog(name, description string) bool {
	type createParam struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Catalog     []model.Catalog `json:"catalog"`
	}

	type createResult struct {
		common_result.Result
		Catalog model.SummaryView `json:"catalog"`
	}

	param := createParam{Name: name, Description: description, Catalog: []model.Catalog{}}
	data, err := json.Marshal(param)
	if err != nil {
		log.Printf("marshal create param failed, err:%s", err.Error())
		return false
	}

	bufferReader := bytes.NewBuffer(data)
	url := fmt.Sprintf("%s/%s?authToken=%s&sessionID=%s", s.baseURL, "content/catalog/", s.onlineView.AuthToken, s.sessionID)
	request, err := http.NewRequest("POST", url, bufferReader)
	if err != nil {
		log.Printf("construct request failed, url:%s, err:%s", url, err.Error())
		return false
	}

	request.Header.Set("content-type", "application/json")
	response, err := s.httpClient.Do(request)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("create catalog failed, statusCode:%d", response.StatusCode)
		return false
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed")
	}

	result := &createResult{}
	err = json.Unmarshal(content, result)
	if result.ErrorCode == common_result.Success {
		return true
	}

	log.Printf("create catalog failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return false
}

func (s *center) FetchCatalog(name string) (model.CatalogDetailView, bool) {
	type fetchResult struct {
		common_result.Result
		Catalog model.CatalogDetailView `json:"catalog"`
	}

	result := &fetchResult{}
	url := fmt.Sprintf("%s/%s?name=%s&authToken=%s&sessionID=%s", s.baseURL, "content/catalog/", name, s.onlineView.AuthToken, s.sessionID)
	response, err := s.httpClient.Get(url)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return result.Catalog, false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("fetch catalog failed, statusCode:%d", response.StatusCode)
		return result.Catalog, false
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed")
	}

	err = json.Unmarshal(content, result)
	if result.ErrorCode == common_result.Success {
		return result.Catalog, true
	}

	return result.Catalog, false
}

func (s *center) QuerySummary() []model.SummaryView {
	viewList := []model.SummaryView{}
	return viewList
}

func (s *center) QueryCatalog() []model.CatalogDetailView {
	viewList := []model.CatalogDetailView{}

	return viewList
}

func (s *center) QuerySubCatalog(catalogID int) []model.CatalogDetailView {
	viewList := []model.CatalogDetailView{}

	return viewList
}

func (s *center) QueryContent(id int) model.ArticleDetailView {
	view := model.ArticleDetailView{}

	return view
}
