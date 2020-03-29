package route

import (
	"fmt"
	"net/http"
	"path"
	"strconv"
	"strings"
)

type filter struct {
	pageID      int
	fileName    string
	catalogName string
	archiveName string
}

func (s *filter) isArchive() bool {
	return s.archiveName != "" && s.catalogName == ""
}

func (s *filter) isCatalog() bool {
	return s.catalogName != "" && s.archiveName == ""
}

func (s *filter) decode(req *http.Request) error {
	filePath, fileName := path.Split(req.URL.EscapedPath())

	str := req.URL.Query().Get("pageid")
	if str != "" {
		val, err := strconv.Atoi(str)
		if err != nil {
			return err
		}

		s.pageID = val
	}

	s.fileName = fileName

	items := strings.Split(strings.Trim(filePath, "/"), "/")
	itemSize := len(items)
	// /view/xxx.html
	if itemSize == 1 {
		return nil
	}

	// /view/post/xxx.html
	if itemSize == 2 {
		val := items[1]
		switch val {
		case "post":
		default:
			return fmt.Errorf("illegal path, url:%s", filePath)
		}

		return nil
	}

	// /view/catalog/xxx/
	// /view/catalog/xxx/xxx.html
	// /view/archive/xxx/
	// /view/archive/xxx/xxx.html
	if itemSize == 3 {
		val := items[1]
		name := items[2]
		switch val {
		case "catalog":
			s.catalogName = name
		case "archive":
			s.archiveName = name
		default:
			return fmt.Errorf("illegal path, url:%s", filePath)
		}

		return nil
	}

	return fmt.Errorf("illegal path, url:%s", filePath)
}
