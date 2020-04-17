package route

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/muidea/magicBlog/config"
	"github.com/muidea/magicBlog/model"
	cmsClient "github.com/muidea/magicCMS/client"
	cmsModel "github.com/muidea/magicCMS/model"
	casClient "github.com/muidea/magicCas/client"
	casModel "github.com/muidea/magicCas/model"
	commonCommon "github.com/muidea/magicCommon/common"
	commonDef "github.com/muidea/magicCommon/def"
	"github.com/muidea/magicCommon/foundation/net"
	"github.com/muidea/magicCommon/foundation/util"
)

const currentCatalog = "current_catalog"
const archiveCatalog = "archive_catalog"
const systemCatalog = "system_catalog"

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

func (s *Registry) queryBlogCommon(clnt cmsClient.Client) (catalogs []*cmsModel.CatalogLite, archives []*cmsModel.CatalogLite, articleList []*cmsModel.ArticleView, err error) {
	blogCatalog, blogErr := clnt.QueryCatalogTree(s.cmsCatalog, 2)
	if blogErr != nil {
		err = blogErr
		return
	}

	var archiveTree *cmsModel.CatalogTree
	var systemTree *cmsModel.CatalogTree
	catalogs = []*cmsModel.CatalogLite{}
	archives = []*cmsModel.CatalogLite{}
	for _, cv := range blogCatalog.Subs {
		switch cv.Name {
		case currentCatalog:
			s.currentCatalog = cv.Lite()
		case archiveCatalog:
			archiveTree = cv
		case systemCatalog:
			systemTree = cv
		default:
			catalogs = append(catalogs, cv.Lite())
		}
	}

	if s.currentCatalog == nil {
		catalogView, catalogErr := clnt.CreateCatalog(currentCatalog, "auto create current catalog", blogCatalog.Lite())
		if catalogErr != nil {
			err = catalogErr
			return
		}

		s.currentCatalog = catalogView.Lite()
	}

	if archiveTree != nil {
		for _, cv := range archiveTree.Subs {
			archives = append(archives, cv.Lite())
		}
	} else {
		_, catalogErr := clnt.CreateCatalog(archiveCatalog, "auto create archive catalog", blogCatalog.Lite())
		if catalogErr != nil {
			err = catalogErr
			return
		}
	}

	if systemTree != nil {
		articleList, err = s.queryArticleList(clnt, systemTree.Lite(), nil)
	} else {
		_, catalogErr := clnt.CreateCatalog(systemCatalog, "auto create system catalog", blogCatalog.Lite())
		if catalogErr != nil {
			err = catalogErr
			return
		}
	}

	return
}

func (s *Registry) queryBlogPost(filter *filter, clnt cmsClient.Client) (fileName string, content interface{}, err error) {
	articleList, articleErr := s.queryArticleList(clnt, s.currentCatalog, filter.pageFilter)
	if articleErr != nil {
		err = articleErr
		return
	}

	var articlePtr *cmsModel.ArticleView
	for _, val := range articleList {
		fileName := fmt.Sprintf("%s.html", val.Title)
		if val.ID == filter.postID && fileName == filter.fileName {
			articlePtr = val
			break
		}
	}

	if articlePtr == nil {
		fileName = "404.html"
		return
	}

	info := map[string]interface{}{}
	commentList, commentErr := s.queryComments(clnt, articlePtr.ID, filter.pageFilter)
	if commentErr != nil {
		err = fmt.Errorf("queryComments failed,err:%s", commentErr.Error())
		return
	}
	info["Content"] = articlePtr
	info["Comments"] = commentList

	fileName = "post.html"
	content = info
	return
}

func (s *Registry) filterBlogArchive(filter *filter, archives []*cmsModel.CatalogLite, clnt cmsClient.Client) (fileName string, content interface{}, err error) {
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
		articlePtr, articleErr := s.queryArticle(clnt, filter.postID)
		if articleErr != nil {
			err = articleErr
			return
		}

		fileName = "post.html"
		content = articlePtr
		return
	}

	articleList, articleErr := s.queryArticleList(clnt, archivePtr, filter.pageFilter)
	if articleErr != nil {
		err = articleErr
		return
	}

	fileName = "index.html"
	content = articleList
	return
}

func (s *Registry) filterBlogCatalog(filter *filter, catalogs []*cmsModel.CatalogLite, clnt cmsClient.Client) (fileName string, content interface{}, err error) {
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
		articlePtr, articleErr := s.queryArticle(clnt, filter.postID)
		if articleErr != nil {
			err = articleErr
			return
		}

		fileName = "post.html"
		content = articlePtr
		return
	}

	articleList, articleErr := s.queryArticleList(clnt, catalogPtr, filter.pageFilter)
	if articleErr != nil {
		err = articleErr
		return
	}

	fileName = "index.html"
	content = articleList
	return
}

func (s *Registry) queryBlogLogin(filter *filter, articles []*cmsModel.ArticleView, clnt cmsClient.Client) (fileName string, content interface{}, err error) {
	var articlePtr *cmsModel.ArticleView
	for _, val := range articles {
		fileName := fmt.Sprintf("%s.html", val.Title)
		if fileName == filter.fileName {
			articlePtr = val
			break
		}
	}

	fileName = "login.html"
	content = articlePtr

	return
}

func (s *Registry) queryBlogAbout(filter *filter, articles []*cmsModel.ArticleView, clnt cmsClient.Client) (fileName string, content interface{}, err error) {
	var articlePtr *cmsModel.ArticleView
	for _, val := range articles {
		fileName := fmt.Sprintf("%s.html", val.Title)
		if fileName == filter.fileName {
			articlePtr = val
			break
		}
	}

	if articlePtr == nil {
		fileName = "404.html"
		return
	}

	fileName = "about.html"
	content = articlePtr

	return
}

func (s *Registry) queryBlogContact(filter *filter, articles []*cmsModel.ArticleView, clnt cmsClient.Client) (fileName string, content interface{}, err error) {
	var articlePtr *cmsModel.ArticleView
	for _, val := range articles {
		fileName := fmt.Sprintf("%s.html", val.Title)
		if fileName == filter.fileName {
			articlePtr = val
			break
		}
	}

	if articlePtr == nil {
		fileName = "404.html"
		return
	}

	info := map[string]interface{}{}
	commentList, commentErr := s.queryComments(clnt, articlePtr.ID, filter.pageFilter)
	if commentErr != nil {
		err = fmt.Errorf("queryComments failed,err:%s", commentErr.Error())
		return
	}
	info["Content"] = articlePtr
	info["Comments"] = commentList

	fileName = "contact.html"
	content = info
	return
}

func (s *Registry) filterBlogPostList(filter *filter, clnt cmsClient.Client) (ret []*cmsModel.ArticleView, err error) {
	articleList, articleErr := s.queryArticleList(clnt, s.currentCatalog, filter.pageFilter)
	if articleErr != nil {
		err = articleErr
		return
	}

	ret = articleList
	return
}

func (s *Registry) queryBlogPostEdit(filter *filter, clnt cmsClient.Client) (fileName string, content interface{}, err error) {
	if filter.action != "update_post" || filter.postID <= 0 {
		err = fmt.Errorf("illegal action, action:%s, postID:%d", filter.action, filter.postID)
		return
	}

	articleView, articleErr := s.queryArticle(clnt, filter.postID)
	if articleErr != nil {
		err = articleErr
		return
	}

	info := map[string]interface{}{}
	info["ID"] = articleView.ID
	info["Title"] = articleView.Title
	info["Content"] = articleView.Content

	catalogs := ""
	for _, val := range articleView.Catalog {
		if val.Name == s.currentCatalog.Name {
			continue
		}

		catalogs = fmt.Sprintf("%s,%s", catalogs, val.Name)
	}
	info["Catalog"] = strings.Trim(catalogs, ",")

	content = info
	fileName = "edit.html"
	return
}

func (s *Registry) deleteBlogPost(filter *filter, clnt cmsClient.Client) (err error) {
	if filter.action != "delete_post" || filter.postID <= 0 {
		err = fmt.Errorf("illegal action, action:%s, postID:%d", filter.action, filter.postID)
		return
	}

	_, err = s.deleteArticle(clnt, filter.postID)
	return
}

func (s *Registry) deleteBlogCatalog(filter *filter, catalogs []*cmsModel.CatalogLite, clnt cmsClient.Client) (err error) {
	if filter.action != "delete_catalog" || filter.catalogID <= 0 {
		err = fmt.Errorf("illegal action, action:%s, catalogId:%d", filter.action, filter.catalogID)
		return
	}

	var catalogPtr *cmsModel.CatalogLite
	for _, val := range catalogs {
		if val.ID == filter.catalogID {
			catalogPtr = val
			break
		}
	}
	if catalogPtr == nil {
		err = fmt.Errorf("illegal action, action:%s, catalogId:%d", filter.action, filter.catalogID)
		return
	}

	articleList, articleErr := s.queryArticleList(clnt, catalogPtr, nil)
	if articleErr != nil {
		err = articleErr
		return
	}
	if len(articleList) > 0 {
		err = fmt.Errorf("delete catalog failed, catalog is busy. action:%s, catalogId:%d", filter.action, filter.catalogID)
		return
	}

	_, err = s.deleteCatalog(clnt, filter.catalogID)
	return
}

func (s *Registry) queryArticle(clnt cmsClient.Client, id int) (ret *cmsModel.ArticleView, err error) {
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
	if pageFilter == nil {
		pageFilter = &util.PageFilter{PageSize: 10, PageNum: 1}
	}

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
	for _, val := range items {
		if val == currentCatalog || val == archiveCatalog {
			continue
		}

		cv, exist := catalogMapInfo[val]
		if !exist {
			catalogMapInfo[val] = nil
			newCatalogItems = append(newCatalogItems, val)
		} else if cv != nil {
			if val != archiveCatalog {
				ret = append(ret, cv)
			}
		}
	}

	for _, val := range newCatalogItems {
		newCatalog, newErr := clnt.CreateCatalog(val, "auto create blog catalog", blogCatalog.Lite())
		if newErr != nil {
			err = newErr
			return
		}
		if val != archiveCatalog {
			ret = append(ret, newCatalog.Lite())
		}
	}

	ret = append(ret, s.currentCatalog)

	return
}

// PostBlog post blog
func (s *Registry) PostBlog(res http.ResponseWriter, req *http.Request) {
	type postParam struct {
		ID      int    `json:"id"`
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

		memo := ""
		if param.ID > 0 {
			_, err = s.updateArticle(cmsClient, param.ID, param.Title, param.Content, catalogList)
			if err != nil {
				result.ErrorCode = commonDef.Failed
				result.Reason = "提交Blog失败, 更新出错"
				break
			}

			memo = fmt.Sprintf("更新Blog%s", param.Title)
		} else {
			_, err = s.createArticle(cmsClient, param.Title, param.Content, catalogList)
			if err != nil {
				result.ErrorCode = commonDef.Failed
				result.Reason = "提交Blog失败, 保存出错"
				break
			}

			memo = fmt.Sprintf("新建Blog%s", param.Title)
		}

		s.recordPostBlog(res, req, memo)

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

// PostComment post comment
func (s *Registry) PostComment(res http.ResponseWriter, req *http.Request) {
	type postParam struct {
		Name    string `json:"name"`
		EMail   string `json:"email"`
		Message string `json:"message"`
		Origin  string `json:"origin"`
		Host    int    `json:"host"`
	}

	type postResult struct {
		commonDef.Result
		Redirect string `json:"redirect"`
	}

	result := &postResult{}
	for {
		param := &postParam{}
		err := net.ParseJSONBody(req, param)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数"
			break
		}

		if param.Name == "" || param.EMail == "" || param.Origin == "" || param.Host == 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数,输入参数为空"
			break
		}

		cmsClient, cmsErr := s.getCMSClient()
		if cmsErr != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "留言失败, 系统出错"
			break
		}
		defer cmsClient.Release()

		_, err = cmsClient.CreateComment(param.Message, param.Name, &cmsModel.Unit{UID: param.Host, UType: cmsModel.ARTICLE}, 0)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "留言失败, 保存出错"
			break
		}

		result.ErrorCode = commonDef.Success
		result.Redirect = param.Origin
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

// ReplyComment reply comment
func (s *Registry) ReplyComment(res http.ResponseWriter, req *http.Request) {
	type postParam struct {
		Message string `json:"message"`
		Origin  string `json:"origin"`
		Host    int    `json:"host"`
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

		if param.Origin == "" || param.Host == 0 {
			result.ErrorCode = commonDef.Failed
			result.Reason = "非法参数,输入参数为空"
			break
		}

		cmsClient := cmsClient.NewClient(s.cmsService)
		defer cmsClient.Release()

		cmsClient.BindSession(sessionInfo)

		authPtr, _ := curSession.GetOption(commonCommon.AuthAccount)
		accountPtr := authPtr.(*casModel.AccountView)

		_, err = cmsClient.CreateComment(param.Message, accountPtr.Account, &cmsModel.Unit{UID: param.Host, UType: cmsModel.COMMENT}, 0)
		if err != nil {
			result.ErrorCode = commonDef.Failed
			result.Reason = "回复失败, 保存出错"
			break
		}

		result.ErrorCode = commonDef.Success
		result.Redirect = param.Origin
		break
	}

	block, err := json.Marshal(result)
	if err == nil {
		res.Write(block)
		return
	}

	res.WriteHeader(http.StatusExpectationFailed)
}

func (s *Registry) createArticle(clnt cmsClient.Client, title, content string, catalogs []*cmsModel.CatalogLite) (ret *cmsModel.ArticleView, err error) {
	blogArticle, blogErr := clnt.CreateArticle(title, content, catalogs)
	if blogErr != nil {
		err = blogErr
		log.Printf("CreateArticle failed, err:%s", err.Error())
		return
	}

	ret = blogArticle
	return
}

func (s *Registry) deleteArticle(clnt cmsClient.Client, id int) (ret *cmsModel.ArticleView, err error) {
	blogArticle, blogErr := clnt.DeleteArticle(id)
	if blogErr != nil {
		err = blogErr
		log.Printf("DeleteArticle failed, err:%s", err.Error())
		return
	}

	ret = blogArticle
	return
}

func (s *Registry) updateArticle(clnt cmsClient.Client, id int, title, content string, catalogs []*cmsModel.CatalogLite) (ret *cmsModel.ArticleView, err error) {
	blogArticle, blogErr := clnt.UpdateArticle(id, title, content, catalogs)
	if blogErr != nil {
		err = blogErr
		log.Printf("UpdateArticle failed, err:%s", err.Error())
		return
	}

	ret = blogArticle
	return
}

func (s *Registry) deleteCatalog(clnt cmsClient.Client, id int) (ret *cmsModel.CatalogView, err error) {
	blogCatalog, blogErr := clnt.DeleteCatalog(id)
	if blogErr != nil {
		err = blogErr
		log.Printf("DeleteCatalog failed, err:%s", err.Error())
		return
	}

	ret = blogCatalog
	return
}

func (s *Registry) queryComments(clnt cmsClient.Client, id int, pageFilter *util.PageFilter) (ret []*model.CommentView, err error) {
	blogComment, _, blogErr := clnt.FilterComment(&cmsModel.Unit{UID: id, UType: cmsModel.ARTICLE}, nil)
	if blogErr != nil {
		err = blogErr
		log.Printf("FilterComment failed, err:%s", err.Error())
		return
	}

	ret = []*model.CommentView{}
	for _, val := range blogComment {
		view := &model.CommentView{ID: val.ID, Content: val.Content, Creater: val.Creater, CreateDate: val.CreateDate}
		replyComment, _, replyErr := clnt.FilterComment(val.Unit(), nil)

		view.Reply = []interface{}{}
		if replyErr == nil {
			for _, sv := range replyComment {
				rv := &model.CommentView{ID: sv.ID, Content: sv.Content, Creater: sv.Creater, CreateDate: sv.CreateDate}
				view.Reply = append(view.Reply, rv)
			}
		}

		ret = append(ret, view)
	}

	return
}
