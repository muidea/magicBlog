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
	Start(bashURL, endpointID, authToken string) (string, bool)
	Stop()
	LoginAccount(account, password string) (model.AccountOnlineView, string, string, bool)
	LogoutAccount(authToken, sessionID string) bool
	StatusAccount(authToken, sessionID string) (model.AccountOnlineView, bool)
	CreateCatalog(name, description string, parent []model.Catalog, authToken, sessionID string) bool
	CreateArticle(title, content string, catalog []model.Catalog, authToken, sessionID string) bool
	FetchCatalog(name string) (model.CatalogDetailView, bool)
	QuerySummary(catalogID int) []model.SummaryView
	QueryCatalog(catalogID int) (model.CatalogDetailView, bool)
	QueryArticle(id int) (model.ArticleDetailView, bool)
	QueryLink(id int) (model.LinkDetailView, bool)
	QueryMedia(id int) (model.MediaDetailView, bool)
}

// NewCenterAgent 新建Agent
func NewCenterAgent() Agent {
	return &center{}
}

type center struct {
	httpClient *http.Client
	baseURL    string
	endpointID string
	authToken  string
	sessionID  string
}

func (s *center) Start(bashURL, endpointID, authToken string) (string, bool) {
	s.httpClient = &http.Client{}
	s.baseURL = bashURL
	s.endpointID = endpointID
	s.authToken = authToken

	sessionID, ok := s.verify()
	if !ok {
		return "", false
	}

	s.sessionID = sessionID
	log.Print("start centerAgent ok")
	return sessionID, true
}

func (s *center) Stop() {

}

func (s *center) verify() (string, bool) {
	type verifyResult struct {
		common_result.Result
		SessionID string `json:"sessionID"`
	}

	result := &verifyResult{}
	url := fmt.Sprintf("%s/%s/%s?authToken=%s", s.baseURL, "authority/endpoint/verify", s.endpointID, s.authToken)
	log.Print(url)
	response, err := s.httpClient.Get(url)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return "", false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("verify failed, statusCode:%d", response.StatusCode)
		return "", false
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed, err:%s", err.Error())
		return "", false
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return "", false
	}

	if result.ErrorCode == common_result.Success {
		return result.SessionID, true
	}

	log.Printf("verify failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return "", false
}

func (s *center) LoginAccount(account, password string) (model.AccountOnlineView, string, string, bool) {
	type loginParam struct {
		Account  string `json:"account"`
		Password string `json:"password"`
	}

	type loginResult struct {
		common_result.Result
		OnlineUser model.AccountOnlineView `json:"onlineUser"`
		SessionID  string                  `json:"sessionID"`
	}

	param := loginParam{Account: account, Password: password}
	result := &loginResult{}
	data, err := json.Marshal(param)
	if err != nil {
		log.Printf("marshal login param failed, err:%s", err.Error())
		return result.OnlineUser, "", "", false
	}

	bufferReader := bytes.NewBuffer(data)
	url := fmt.Sprintf("%s/%s", s.baseURL, "cas/user/")
	request, err := http.NewRequest("POST", url, bufferReader)
	if err != nil {
		log.Printf("construct request failed, url:%s, err:%s", url, err.Error())
		return result.OnlineUser, "", "", false
	}

	request.Header.Set("content-type", "application/json")
	response, err := s.httpClient.Do(request)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return result.OnlineUser, "", "", false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("login failed, statusCode:%d", response.StatusCode)
		return result.OnlineUser, "", "", false
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed, err:%s", err.Error())
		return result.OnlineUser, "", "", false
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return result.OnlineUser, "", "", false
	}

	if result.ErrorCode == common_result.Success {
		return result.OnlineUser, result.OnlineUser.AuthToken, result.SessionID, true
	}

	log.Printf("login failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return result.OnlineUser, "", "", false
}

func (s *center) LogoutAccount(authToken, sessionID string) bool {
	type logoutResult struct {
		common_result.Result
	}

	result := &logoutResult{}
	url := fmt.Sprintf("%s/%s/?authToken=%s&sessionID=%s", s.baseURL, "cas/user", authToken, sessionID)
	response, err := s.httpClient.Get(url)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("logout account failed, statusCode:%d", response.StatusCode)
		return false
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed, err:%s", err.Error())
		return false
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return false
	}

	if result.ErrorCode == common_result.Success {
		return true
	}

	log.Printf("logout account failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return false
}

func (s *center) StatusAccount(authToken, sessionID string) (model.AccountOnlineView, bool) {
	type statusResult struct {
		common_result.Result
		OnlineUser model.AccountOnlineView `json:"onlineUser"`
		SessionID  string                  `json:"sessionID"`
	}

	result := &statusResult{}
	url := fmt.Sprintf("%s/%s/?authToken=%s&sessionID=%s", s.baseURL, "cas/user", authToken, sessionID)
	response, err := s.httpClient.Get(url)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return result.OnlineUser, false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("status account failed, statusCode:%d", response.StatusCode)
		return result.OnlineUser, false
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed, err:%s", err.Error())
		return result.OnlineUser, false
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return result.OnlineUser, false
	}

	if result.ErrorCode == common_result.Success {
		return result.OnlineUser, true
	}

	log.Printf("status account failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return result.OnlineUser, false
}

func (s *center) CreateCatalog(name, description string, parent []model.Catalog, authToken, sessionID string) bool {
	type createParam struct {
		Name        string          `json:"name"`
		Description string          `json:"description"`
		Catalog     []model.Catalog `json:"catalog"`
	}

	type createResult struct {
		common_result.Result
		Catalog model.SummaryView `json:"catalog"`
	}

	param := createParam{Name: name, Description: description, Catalog: parent}
	data, err := json.Marshal(param)
	if err != nil {
		log.Printf("marshal create param failed, err:%s", err.Error())
		return false
	}

	bufferReader := bytes.NewBuffer(data)
	url := fmt.Sprintf("%s/%s?authToken=%s&sessionID=%s", s.baseURL, "content/catalog/", authToken, sessionID)
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
		log.Printf("read respose data failed, err:%s", err.Error())
		return false
	}

	result := &createResult{}
	err = json.Unmarshal(content, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return false
	}

	if result.ErrorCode == common_result.Success {
		return true
	}

	log.Printf("create catalog failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return false
}

func (s *center) CreateArticle(title, content string, catalog []model.Catalog, authToken, sessionID string) bool {
	type createParam struct {
		Name    string          `json:"name"`
		Content string          `json:"content"`
		Catalog []model.Catalog `json:"catalog"`
	}

	type createResult struct {
		common_result.Result
		Article model.SummaryView `json:"article"`
	}

	param := createParam{Name: title, Content: content, Catalog: catalog}
	data, err := json.Marshal(param)
	if err != nil {
		log.Printf("marshal create param failed, err:%s", err.Error())
		return false
	}

	bufferReader := bytes.NewBuffer(data)
	url := fmt.Sprintf("%s/%s?authToken=%s&sessionID=%s", s.baseURL, "content/article/", authToken, sessionID)
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
		log.Printf("create article failed, statusCode:%d", response.StatusCode)
		return false
	}

	contentData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed, err:%s", err.Error())
		return false
	}

	result := &createResult{}
	err = json.Unmarshal(contentData, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return false
	}

	if result.ErrorCode == common_result.Success {
		return true
	}

	log.Printf("create article failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return false
}

func (s *center) FetchCatalog(name string) (model.CatalogDetailView, bool) {
	type fetchResult struct {
		common_result.Result
		Catalog model.CatalogDetailView `json:"catalog"`
	}

	result := &fetchResult{}
	url := fmt.Sprintf("%s/%s?name=%s&authToken=%s&sessionID=%s", s.baseURL, "content/catalog/", name, s.authToken, s.sessionID)
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
		log.Printf("read respose data failed, err:%s", err.Error())
		return result.Catalog, false
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return result.Catalog, false
	}

	if result.ErrorCode == common_result.Success {
		return result.Catalog, true
	}

	log.Printf("fetch catalog failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return result.Catalog, false
}

func (s *center) QuerySummary(catalogID int) []model.SummaryView {
	type queryResult struct {
		common_result.Result
		Summary []model.SummaryView `json:"summary"`
	}

	result := &queryResult{Summary: []model.SummaryView{}}
	url := fmt.Sprintf("%s/%s/%d?authToken=%s&sessionID=%s", s.baseURL, "content/summary", catalogID, s.authToken, s.sessionID)
	response, err := s.httpClient.Get(url)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return result.Summary
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("query summary failed, statusCode:%d", response.StatusCode)
		return result.Summary
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed, err:%s", err.Error())
		return result.Summary
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return result.Summary
	}

	if result.ErrorCode == common_result.Success {
		return result.Summary
	}

	log.Printf("query summary failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return result.Summary
}

func (s *center) QueryCatalog(catalogID int) (model.CatalogDetailView, bool) {
	type queryResult struct {
		common_result.Result
		Catalog model.CatalogDetailView `json:"catalog"`
	}

	result := &queryResult{}
	url := fmt.Sprintf("%s/%s/%d?authToken=%s&sessionID=%s", s.baseURL, "content/catalog", catalogID, s.authToken, s.sessionID)
	response, err := s.httpClient.Get(url)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return result.Catalog, false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("query catalog failed, statusCode:%d", response.StatusCode)
		return result.Catalog, false
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed, err:%s", err.Error())
		return result.Catalog, false
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return result.Catalog, false
	}

	if result.ErrorCode == common_result.Success {
		return result.Catalog, true
	}

	log.Printf("query catalog failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return result.Catalog, false
}

func (s *center) QueryArticle(id int) (model.ArticleDetailView, bool) {
	type queryResult struct {
		common_result.Result
		Article model.ArticleDetailView `json:"article"`
	}

	result := &queryResult{}
	url := fmt.Sprintf("%s/%s/%d?authToken=%s&sessionID=%s", s.baseURL, "content/article", id, s.authToken, s.sessionID)
	response, err := s.httpClient.Get(url)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return result.Article, false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("query article failed, statusCode:%d", response.StatusCode)
		return result.Article, false
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed, err:%s", err.Error())
		return result.Article, false
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return result.Article, false
	}

	if result.ErrorCode == common_result.Success {
		return result.Article, true
	}

	log.Printf("query article failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return result.Article, false
}

func (s *center) QueryLink(id int) (model.LinkDetailView, bool) {
	type queryResult struct {
		common_result.Result
		Link model.LinkDetailView `json:"link"`
	}

	result := &queryResult{}
	url := fmt.Sprintf("%s/%s/%d?authToken=%s&sessionID=%s", s.baseURL, "content/link", id, s.authToken, s.sessionID)
	response, err := s.httpClient.Get(url)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return result.Link, false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("query link failed, statusCode:%d", response.StatusCode)
		return result.Link, false
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed, err:%s", err.Error())
		return result.Link, false
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return result.Link, false
	}

	if result.ErrorCode == common_result.Success {
		return result.Link, false
	}

	log.Printf("query link failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return result.Link, false
}

func (s *center) QueryMedia(id int) (model.MediaDetailView, bool) {
	type queryResult struct {
		common_result.Result
		Media model.MediaDetailView `json:"media"`
	}

	result := &queryResult{}
	url := fmt.Sprintf("%s/%s/%d?authToken=%s&sessionID=%s", s.baseURL, "content/media", id, s.authToken, s.sessionID)
	response, err := s.httpClient.Get(url)
	if err != nil {
		log.Printf("post request failed, err:%s", err.Error())
		return result.Media, false
	}

	if response.StatusCode != http.StatusOK {
		log.Printf("query media failed, statusCode:%d", response.StatusCode)
		return result.Media, false
	}

	content, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("read respose data failed, err:%s", err.Error())
		return result.Media, false
	}

	err = json.Unmarshal(content, result)
	if err != nil {
		log.Printf("unmarshal data failed, err:%s", err.Error())
		return result.Media, false
	}

	if result.ErrorCode == common_result.Success {
		return result.Media, true
	}

	log.Printf("query media failed, errorCode:%d, reason:%s", result.ErrorCode, result.Reason)
	return result.Media, false
}
