const APIV1 = '/api/v1'

module.exports = {
  name: 'MagicBlog',
  prefix: 'magicBlog',
  footerText: 'MagicBlog © 2017 muidea.com',
  api: {
    indexQuery: `${APIV1}/`,
    catalogQueryAll: `${APIV1}/catalog/`,
    catalogQuery: `${APIV1}/catalog/:id`,
    contactQuery: `${APIV1}/contact/`,
    aboutQuery: `${APIV1}/about/`,
    articleQuery: `${APIV1}/content/:id`,
    userStatus: `${APIV1}/maintain/status`,
    userLogin: `${APIV1}/maintain/login`,
    userLogout: `${APIV1}/maintain/logout`,
    summaryQuery: `${APIV1}/maintain/summary`,
    noFoundPage: `${APIV1}/404.html`,
  },
}
