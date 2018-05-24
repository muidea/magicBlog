package core

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"muidea.com/magicCommon/foundation/net"
	"muidea.com/magicCommon/model"
	engine "muidea.com/magicEngine"
)

type route struct {
	pattern string
	method  string
	handler interface{}
}

func (s *route) Pattern() string {
	return s.pattern
}

func (s *route) Method() string {
	return s.method
}

func (s *route) Handler() interface{} {
	return s.handler
}

func newRoute(pattern, method string, handler interface{}) engine.Route {
	return &route{pattern: pattern, method: method, handler: handler}
}

// NewBlog 新建Blog
func NewBlog(centerServer, name, account, password string) (Blog, bool) {
	blog := Blog{centerAgent: NewCenterAgent()}

	agent := NewCenterAgent()
	if !agent.Start(centerServer, name, account, password) {
		return blog, false
	}
	blogCatalog, ok := agent.FetchCatalog(name)
	if !ok {
		ok = agent.CreateCatalog(name, "MagicBlog auto create catalog.")
		if !ok {
			log.Print("create blog root catalog failed.")
			return blog, false
		}
	}
	blogCatalog, ok = agent.FetchCatalog(name)
	if !ok {
		log.Print("fetch blog root ctalog failed.")
		return blog, false
	}

	blogContent := agent.QuerySummary(blogCatalog.ID)

	blog.centerAgent = agent
	blog.blogInfo = blogCatalog
	blog.blogContent = blogContent

	return blog, true
}

// Blog Blog对象
type Blog struct {
	centerAgent Agent
	blogInfo    model.CatalogDetailView
	blogContent []model.SummaryView
}

// Startup 启动
func (s *Blog) Startup(router engine.Router) {
	mainRoute := newRoute("/", "GET", s.mainPage)
	router.AddRoute(mainRoute)

	catalogListRoute := newRoute("/catalog/", "GET", s.catalogListPage)
	router.AddRoute(catalogListRoute)

	catalogRoute := newRoute("/catalog/:id", "GET", s.catalogPage)
	router.AddRoute(catalogRoute)

	contentRoute := newRoute("/content/:id", "GET", s.contentPage)
	router.AddRoute(contentRoute)

	aboutRoute := newRoute("/about/", "GET", s.aboutPage)
	router.AddRoute(aboutRoute)

	contactRoute := newRoute("/contact/", "GET", s.contactPage)
	router.AddRoute(contactRoute)

	noFoundRoute := newRoute("/404.html", "GET", s.noFoundPage)
	router.AddRoute(noFoundRoute)
}

// Teardown 销毁
func (s *Blog) Teardown() {
	if s.centerAgent != nil {
		s.centerAgent.Stop()
	}
}

func (s *Blog) getIndexView() (model.SummaryView, bool) {
	for _, v := range s.blogContent {
		if v.Name == "Index" && v.Type == model.CATALOG {
			return v, true
		}
	}

	return model.SummaryView{}, false
}

func (s *Blog) getCatalogView() (model.SummaryView, bool) {
	for _, v := range s.blogContent {
		if v.Name == "Catalog" && v.Type == model.CATALOG {
			return v, true
		}
	}

	return model.SummaryView{}, false
}

func (s *Blog) getAboutView() (model.SummaryView, bool) {
	for _, v := range s.blogContent {
		if v.Name == "About" && v.Type == model.ARTICLE {
			return v, true
		}
	}

	return model.SummaryView{}, false
}

func (s *Blog) getContactView() (model.SummaryView, bool) {
	for _, v := range s.blogContent {
		if v.Name == "Contact" && v.Type == model.ARTICLE {
			return v, true
		}
	}

	return model.SummaryView{}, false
}

func (s *Blog) get404View() (model.SummaryView, bool) {
	for _, v := range s.blogContent {
		if v.Name == "404" && v.Type == model.ARTICLE {
			return v, true
		}
	}

	return model.SummaryView{}, false
}

func (s *Blog) mainPage(res http.ResponseWriter, req *http.Request) {
	log.Print("mainPage")

	summary := []model.SummaryView{}
	indexView, ok := s.getIndexView()
	if ok {
		summary = s.centerAgent.QuerySummary(indexView.ID)
		block, err := json.Marshal(summary)
		if err == nil {
			res.Write(block)
			return
		}

		log.Print("mainPage, json.Marshal, failed, err:" + err.Error())
	}

	http.Redirect(res, req, "/default/index.html", http.StatusMovedPermanently)
}

func (s *Blog) catalogListPage(res http.ResponseWriter, req *http.Request) {
	log.Print("catalogListPage")

	summary := []model.SummaryView{}
	catalogView, ok := s.getCatalogView()
	if ok {
		summary = s.centerAgent.QuerySummary(catalogView.ID)
		block, err := json.Marshal(summary)
		if err == nil {
			res.Write(block)
			return
		}

		log.Print("catalogListPage, json.Marshal, failed, err:" + err.Error())
	}

	http.Redirect(res, req, "/default/catalog.html", http.StatusMovedPermanently)
}

func (s *Blog) catalogPage(res http.ResponseWriter, req *http.Request) {
	log.Print("catalogPage")

	_, value := net.SplitRESTAPI(req.URL.Path)
	id, err := strconv.Atoi(value)
	if err == nil {
		summary := s.centerAgent.QuerySummary(id)
		block, err := json.Marshal(summary)
		if err == nil {
			res.Write(block)
			return
		}

		log.Print("catalogPage, json.Marshal, failed, err:" + err.Error())
	} else {
		log.Print("catalogPage, strconv.atoi failed, err:" + err.Error())
	}

	http.Redirect(res, req, "/404.html", http.StatusMovedPermanently)
}

func (s *Blog) contentPage(res http.ResponseWriter, req *http.Request) {
	log.Print("contentPage")

	_, value := net.SplitRESTAPI(req.URL.Path)
	id, err := strconv.Atoi(value)
	if err == nil {
		article, ok := s.centerAgent.QueryArticle(id)
		if ok {
			block, err := json.Marshal(article)
			if err == nil {
				res.Write(block)
				return
			}

			log.Print("contentPage, json.Marshal, failed, err:" + err.Error())
		} else {
			log.Printf("contentPage, nofound article content, id:%d", id)
		}
	}

	http.Redirect(res, req, "/404.html", http.StatusMovedPermanently)
}

func (s *Blog) aboutPage(res http.ResponseWriter, req *http.Request) {
	log.Print("aboutPage")

	aboutView, ok := s.getAboutView()
	if ok {
		article, ok := s.centerAgent.QueryArticle(aboutView.ID)
		if ok {
			block, err := json.Marshal(article)
			if err == nil {
				res.Write(block)
				return
			}

			log.Print("aboutPage, json.Marshal, failed, err:" + err.Error())
		} else {
			log.Printf("aboutPage, nofound about content, id:%d", aboutView.ID)
		}
	}

	http.Redirect(res, req, "/default/about.html", http.StatusMovedPermanently)
}

func (s *Blog) contactPage(res http.ResponseWriter, req *http.Request) {
	log.Print("contactPage")

	contactView, ok := s.getContactView()
	if ok {
		article, ok := s.centerAgent.QueryArticle(contactView.ID)

		if ok {
			block, err := json.Marshal(article)
			if err == nil {
				res.Write(block)
				return
			}

			log.Print("contactPage, json.Marshal, failed, err:" + err.Error())
		} else {
			log.Printf("contactPage, nofound contact content, id:%d", contactView.ID)
		}
	}

	http.Redirect(res, req, "/default/contact.html", http.StatusMovedPermanently)
}

func (s *Blog) noFoundPage(res http.ResponseWriter, req *http.Request) {
	log.Print("noFoundPage")

	noFoundView, ok := s.get404View()
	if ok {
		article, ok := s.centerAgent.QueryArticle(noFoundView.ID)

		if ok {
			block, err := json.Marshal(article)
			if err == nil {
				res.Write(block)
				return
			}

			log.Print("noFoundPage, json.Marshal, failed, err:" + err.Error())
		} else {
			log.Printf("noFoundPage, nofound 404 content, id:%d", noFoundView.ID)
		}
	}

	http.Redirect(res, req, "/default/404.html", http.StatusMovedPermanently)
}
