package agent

import (
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
	Start(bashURL, endpointID, authToken string) bool
	Stop()

	LoginAccount(account, password string) (model.AccountOnlineView, string, string, bool)
	LogoutAccount(authToken, sessionID string) bool
	StatusAccount(authToken, sessionID string) (model.AccountOnlineView, bool)

	QuerySummary(catalogID int) []model.SummaryView

	FetchCatalog(name string) (model.CatalogDetailView, bool)
	QueryCatalog(catalogID int) (model.CatalogDetailView, bool)
	CreateCatalog(name, description string, parent []model.Catalog, authToken, sessionID string) bool

	QueryArticle(id int) (model.ArticleDetailView, bool)
	CreateArticle(title, content string, catalog []model.Catalog, authToken, sessionID string) bool

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

func (s *center) Start(bashURL, endpointID, authToken string) bool {
	s.httpClient = &http.Client{}
	s.baseURL = bashURL
	s.endpointID = endpointID
	s.authToken = authToken

	sessionID, ok := s.verify()
	if !ok {
		return false
	}

	s.sessionID = sessionID
	log.Print("start centerAgent ok")
	return true
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