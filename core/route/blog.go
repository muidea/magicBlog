package route

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/muidea/magicBlog/config"
	cmsClient "github.com/muidea/magicCMS/client"
	cmsModel "github.com/muidea/magicCMS/model"
	casClient "github.com/muidea/magicCas/client"
	commonCommon "github.com/muidea/magicCommon/common"
	commonDef "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/net"
	"github.com/muidea/magicCommon/foundation/util"
)

const currentCatalog = "current_catalog"
const archiveCatalog = "archive_catalog"

func (s *Registry) verifyEndpoint() (ret *commonCommon.SessionInfo, err error) {
	casClient := casClient.NewClient(s.casService)
	defer casClient.Release()

	identityID := config.IdentityID()
	authToken := config.AuthToken()
	_, ret, err = casClient.VerifyEndpoint(identityID, authToken)

	return
}

func (s *Registry) getCMSClient() (ret cmsClient.Client, err error) {
	sessionInfo, sessionErr := s.verifyEndpoint()
	if sessionErr != nil {
		log.Printf("verifyEndpoint failed, err:%s", sessionErr.Error())
		err = sessionErr
		return
	}

	sessionInfo.Scope = commonCommon.ShareSession
	cmsClient := cmsClient.NewClient(s.cmsService)
	cmsClient.BindSession(sessionInfo)

	ret = cmsClient

	return
}

func (s *Registry) getCommonInfo(clnt cmsClient.Client) (catalogs []*cmsModel.CatalogLite, archives []*cmsModel.CatalogLite, err error) {
	blogCatalog, blogErr := clnt.QueryCatalogTree(s.cmsCatalog, 2)
	if blogErr != nil {
		err = blogErr
		return
	}

	var archiveTree *cmsModel.CatalogTree
	catalogs = []*cmsModel.CatalogLite{}
	archives = []*cmsModel.CatalogLite{}
	for _, cv := range blogCatalog.Subs {
		switch cv.Name {
		case currentCatalog:
			s.currentCatalog = cv
		case archiveCatalog:
			archiveTree = cv
		default:
			catalogs = append(catalogs, cv.Lite())
		}
	}
	if archiveTree != nil {
		for _, cv := range archiveTree.Subs {
			archives = append(archives, cv.Lite())
		}
	}

	return
}

func (s *Registry) filterArchive(filter *filter, archives []*cmsModel.CatalogLite, clnt cmsClient.Client) (fileName string, content interface{}, err error) {
	var archivePtr *cmsModel.CatalogLite
	for _, val := range archives {
		if filter.archiveName == val.Name {
			archivePtr = val
			break
		}
	}
	if archivePtr == nil {
		err = fmt.Errorf("illegal archive, name:%s", filter.archiveName)
		return
	}
	if filter.fileName != "" {
		articlePtr, articleErr := s.queryArticle(clnt, archivePtr, filter.pageID)
		if articleErr != nil {
			err = articleErr
			return
		}

		fileName = "post.html"
		content = articlePtr
		return
	}

	articleList, articleErr := s.queryArticleList(clnt, archivePtr, nil)
	if articleErr != nil {
		err = articleErr
		return
	}

	fileName = "index.html"
	content = articleList
	return
}

func (s *Registry) filterCatalog(filter *filter, catalogs []*cmsModel.CatalogLite, clnt cmsClient.Client) (fileName string, content interface{}, err error) {
	var catalogPtr *cmsModel.CatalogLite
	for _, val := range catalogs {
		if filter.catalogName == val.Name {
			catalogPtr = val
			break
		}
	}
	if catalogPtr == nil {
		err = fmt.Errorf("illegal catalog, name:%s", filter.catalogName)
		return
	}

	if filter.fileName != "" {
		articlePtr, articleErr := s.queryArticle(clnt, catalogPtr, filter.pageID)
		if articleErr != nil {
			err = articleErr
			return
		}

		fileName = "post.html"
		content = articlePtr
		return
	}

	articleList, articleErr := s.queryArticleList(clnt, catalogPtr, nil)
	if articleErr != nil {
		err = articleErr
		return
	}

	fileName = "index.html"
	content = articleList
	return
}

func (s *Registry) filterAbout(filter *filter, clnt cmsClient.Client) (ret *cmsModel.ArticleView, err error) {
	return
}

func (s *Registry) filterContact(filter *filter, clnt cmsClient.Client) (ret *cmsModel.ArticleView, err error) {
	return
}

func (s *Registry) filterPostList(filter *filter, clnt cmsClient.Client) (ret []*cmsModel.ArticleView, err error) {
	articleList, articleErr := s.queryArticleList(clnt, nil, nil)
	if articleErr != nil {
		err = articleErr
		return
	}

	ret = articleList
	return
}

func (s *Registry) queryArticle(clnt cmsClient.Client, catalog *cmsModel.CatalogLite, id int) (ret *cmsModel.ArticleView, err error) {
	blogArticle, blogErr := clnt.QueryArticle(id)
	if blogErr != nil {
		err = blogErr
		log.Printf("QueryArticle failed, err:%s", err.Error())
		return
	}

	ret = blogArticle
	return
}

func (s *Registry) queryArticleList(clnt cmsClient.Client, catalog *cmsModel.CatalogLite, pageFilter *util.PageFilter) (ret []*cmsModel.ArticleView, err error) {
	blogArticle, _, blogErr := clnt.FilterArticle(catalog, pageFilter)
	if blogErr != nil {
		err = blogErr
		log.Printf("FilterArticle failed, err:%s", err.Error())
		return
	}

	ret = blogArticle
	return
}

func (s *Registry) getCatalogs(catalog string, clnt cmsClient.Client) (ret []*cmsModel.CatalogLite, err error) {
	blogCatalog, blogErr := clnt.QueryCatalogTree(s.cmsCatalog, 2)
	if blogErr != nil {
		err = blogErr
		return
	}

	newCatalogItems := []string{}
	ret = []*cmsModel.CatalogLite{}
	catalogMapInfo := map[string]*cmsModel.CatalogLite{}
	for _, cv := range blogCatalog.Subs {
		catalogMapInfo[cv.Name] = cv.Lite()
	}

	items := strings.Split(strings.Trim(catalog, " "), ",")
	items = append(items, currentCatalog, archiveCatalog)
	for _, val := range items {
		cv, exist := catalogMapInfo[val]
		if !exist {
			catalogMapInfo[val] = nil
			newCatalogItems = append(newCatalogItems, val)
		} else if cv != nil {
			ret = append(ret, cv)
		}
	}

	for _, val := range newCatalogItems {
		newCatalog, newErr := clnt.CreateCatalog(val, "auto create catalog", blogCatalog.Lite())
		if newErr != nil {
			err = newErr
			return
		}

		ret = append(ret, newCatalog.Lite())
	}

	return
}

// PostBlog post blog
func (s *Registry) PostBlog(res http.ResponseWriter, req *http.Request) {
	type postParam struct {
		Title   string `json:"title"`
		Content string `json:"content"`
		Catalog string `json:"catalog"`
	}

	type postResult struct {
		commonDef.Result
		Redirect string `json:"redirect"`
	}

	curSession := s.sessionRegistry.GetSession(res, req)

	sessionInfo := curSession.GetSessionInfo()
	result := &postResult{}
	for {
		param := &postParam{}
		err := net.ParseJSONBody(req, param)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数"
			break
		}

		if param.Title == "" || param.Catalog == "" {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数,输入参数为空"
			break
		}

		cmsClient := cmsClient.NewClient(s.cmsService)
		defer cmsClient.Release()

		cmsClient.BindSession(sessionInfo)

		catalogList, catalogErr := s.getCatalogs(param.Catalog, cmsClient)
		if catalogErr != nil {
			log.Printf("getCatalogs failed, err:%s", catalogErr.Error())
			result.ErrorCode = commonDef.Failed
			result.Reason = "提交Blog失败, 查询分类出错"
			break
		}

		_, err = cmsClient.CreateArticle(param.Title, param.Content, catalogList)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "提交Blog失败, 保存出错"
			break
		}

		s.recordPostBlog(res, req, param.Title)

		result.ErrorCode = commonDef.Success
		result.Redirect = "/"
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}
