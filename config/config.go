package config

import "flag"

var endpointName = "magicBlog"
var cmsIdentityID = "magicBlog"
var cmsAuthToken = "magicBlog"
var cmsService = "http://127.0.0.1:8880"
var casService = "http://127.0.0.1:8081"
var batisService = "http://127.0.0.1:8080"
var cmsCatalog = 192

func init() {
	flag.StringVar(&endpointName, "EndpointName", endpointName, "magicBlog endpoint name.")
	flag.StringVar(&cmsIdentityID, "IdentityID", cmsIdentityID, "magicCMS identity id.")
	flag.StringVar(&cmsAuthToken, "AuthToken", cmsAuthToken, "magicCMS authority token.")
	flag.StringVar(&cmsService, "CMSService", cmsService, "magicCMS service address")
	flag.IntVar(&cmsCatalog, "CMSCatalog", cmsCatalog, "magicBlog's cms catalog id")
	flag.StringVar(&casService, "CasService", casService, "magicCas service address")
	flag.StringVar(&batisService, "BatisService", batisService, "magicBatis service address")
}

//EndpointName endpointName
func EndpointName() string {
	return endpointName
}

//IdentityID endpointName
func IdentityID() string {
	return cmsIdentityID
}

//AuthToken endpointName
func AuthToken() string {
	return cmsAuthToken
}

//CasService cas service addr
func CasService() string {
	return casService
}

//BatisService batis service addr
func BatisService() string {
	return batisService
}

// CMSService cms service addr
func CMSService() string {
	return cmsService
}

// CMSCatalog cms catalog id
func CMSCatalog() int {
	return cmsCatalog
}
