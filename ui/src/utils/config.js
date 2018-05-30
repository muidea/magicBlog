const APIV1 = '/api/v1'

module.exports = {
  name: 'MagicBlog',
  prefix: 'magicBlog',
  footerText: 'MagicBlog © 2017 muidea.com',
  api: {
    indexPage: `${APIV1}/`,
    catalogPage: `${APIV1}/catalog/`,
    singleCatalogPage: `${APIV1}/catalog/:id`,
    contactPage: `${APIV1}/contact/`,
    aboutPage: `${APIV1}/about/`,
    articlePage: `${APIV1}/content/:id`,
    userStatus: `${APIV1}/maintain/status`,
    userLogin: `${APIV1}/maintain/login/`,
    userLogout: `${APIV1}/maintain/logout/`,
    noFoundPage: `${APIV1}/404.html`,
  },
}
