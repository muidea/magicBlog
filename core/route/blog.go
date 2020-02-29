package route

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	cmsClient "github.com/muidea/magicCMS/client"
	cmsModel "github.com/muidea/magicCMS/model"
	commonDef "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/net"
	"github.com/muidea/magicCommon/foundation/util"
)

const currentCatalog = "current_catalog"
const archiveCatalog = "archive_catalog"

func (s *Registry) filterPostList(res http.ResponseWriter, req *http.Request) interface{} {
	filter := commonDef.NewFilter([]string{"catalog"})
	filter.Decode(req)

	var catalogPtr *cmsModel.CatalogLite
	catalogStr, catalogOK := filter.ContentFilter.Items["catalog"]
	if catalogOK {
		val, err := strconv.Atoi(catalogStr)
		if err == nil {
			catalogPtr = &cmsModel.CatalogLite{ID: val}
		}
	}

	type filterResult struct {
		commonDef.Result
		Catalogs []*cmsModel.CatalogLite `json:"catalogs"`
		Archives []*cmsModel.CatalogLite `json:"archives"`
		Articles []*cmsModel.ArticleView `json:"articles"`
	}

	curSession := s.sessionRegistry.GetSession(res, req)

	sessionInfo := curSession.GetSessionInfo()
	result := &filterResult{
		Catalogs: []*cmsModel.CatalogLite{},
		Archives: []*cmsModel.CatalogLite{},
		Articles: []*cmsModel.ArticleView{},
	}

	cmsClient := cmsClient.NewClient(s.cmsService)
	defer cmsClient.Release()

	cmsClient.BindSession(sessionInfo)
	catalogList, catalogErr := s.queryCatalog(cmsClient)
	if catalogErr == nil {
		result.Catalogs = catalogList
	}
	archiveList, archiveErr := s.queryArchive(cmsClient)
	if archiveErr == nil {
		result.Archives = archiveList
	}
	articleList, articleErr := s.queryArticle(cmsClient, catalogPtr, filter.PageFilter)
	if articleErr == nil {
		result.Articles = articleList
	}

	return result
}

func (s *Registry) queryArticle(clnt cmsClient.Client, catalog *cmsModel.CatalogLite, pageFilter *util.PageFilter) (ret []*cmsModel.ArticleView, err error) {
	blogArticle, _, blogErr := clnt.FilterArticle(catalog, pageFilter)
	if blogErr != nil {
		err = blogErr
		log.Printf("FilterArticle failed, err:%s", err.Error())
		return
	}

	ret = blogArticle
	return
}

func (s *Registry) queryArchive(clnt cmsClient.Client) (ret []*cmsModel.CatalogLite, err error) {
	if s.archiveCatalog == nil {
		err = fmt.Errorf("empty archive blogs")
		return
	}

	archiveList, blogErr := clnt.QueryCatalogTree(s.archiveCatalog.ID, 1)
	if blogErr != nil {
		err = blogErr
		log.Printf("QueryCatalogTree failed, err:%s", err.Error())
		return
	}

	for _, cv := range archiveList.Subs {
		switch cv.Name {
		case currentCatalog:
			s.currentCatalog = cv.Lite()
		case archiveCatalog:
			s.archiveCatalog = cv.Lite()
		default:
			ret = append(ret, cv.Lite())
		}
	}

	return
}

func (s *Registry) queryCatalog(clnt cmsClient.Client) (ret []*cmsModel.CatalogLite, err error) {
	blogCatalog, blogErr := clnt.QueryCatalogTree(s.cmsCatalog, 1)
	if blogErr != nil {
		err = blogErr
		return
	}

	for _, cv := range blogCatalog.Subs {
		switch cv.Name {
		case currentCatalog:
			s.currentCatalog = cv.Lite()
		case archiveCatalog:
			s.archiveCatalog = cv.Lite()
		default:
			ret = append(ret, cv.Lite())
		}
	}

	return
}

func (s *Registry) getCatalogs(catalog string, clnt cmsClient.Client) (ret []*cmsModel.CatalogLite, err error) {
	blogCatalog, blogErr := clnt.QueryCatalogTree(s.cmsCatalog, 1)
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

	items := strings.Split(catalog, ",")
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
