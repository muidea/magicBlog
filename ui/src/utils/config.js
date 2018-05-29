const APIV1 = '/api/v1'

module.exports = {
  name: 'MagicCenter',
  prefix: 'magicCenter',
  footerText: 'MagicBlog © 2017 muidea.com',
  api: {
    indexPage: `${APIV1}/`,
    catalogPage: `${APIV1}/catalog/`,
    singleCatalogPage: `${APIV1}/catalog/:id`,
    contactPage: `${APIV1}/contact/`,
    aboutPage: `${APIV1}/about/`,
    articlePage: `${APIV1}/content/:id`,
    maintainStatus: `${APIV1}/maintain/status`,
    maintainLogin: `${APIV1}/maintain/login/`,
    maintainLogout: `${APIV1}/maintain/logout/`,
    noFoundPage: `${APIV1}/404.html`,
  },
}
