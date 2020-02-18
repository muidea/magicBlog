package config

import "flag"

var endpointName = "magicBlog"
var cmsService = "http://127.0.0.1:8880"
var casService = "http://127.0.0.1:8081"
var userService = "http://127.0.0.1:8082"
var fileService = "http://127.0.0.1:8083"
var batisService = "http://127.0.0.1:8080"
var cmsCatalog = 192

func init() {
	flag.StringVar(&endpointName, "EndpointName", endpointName, "magicBlog endpoint name.")
	flag.StringVar(&cmsService, "CMSService", cmsService, "magicCMS service address")
	flag.IntVar(&cmsCatalog, "CMSCatalog", cmsCatalog, "magicBlog's cms catalog id")
	flag.StringVar(&casService, "CasService", casService, "magicCas service address")
	flag.StringVar(&userService, "UserService", userService, "magicUser service address")
	flag.StringVar(&fileService, "FileService", fileService, "magicFile service address")
	flag.StringVar(&batisService, "BatisService", batisService, "magicBatis service address")
}

//EndpointName endpointName
func EndpointName() string {
	return endpointName
}

//CasService cas service addr
func CasService() string {
	return casService
}

//UserService user service addr
func UserService() string {
	return userService
}

//FileService file service addr
func FileService() string {
	return fileService
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
