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
    maintainPage: `${APIV1}/maintain/`,
    loginPage: `${APIV1}/login/`,
    noFoundPage: `${APIV1}/404.html`,
  },
}
