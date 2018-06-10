const APIV1 = '/api/v1'

module.exports = {
  name: 'MagicBlog',
  prefix: 'magicBlog',
  footerText: 'magicBlog © 2017 muidea.com',
  api: {
    indexQuery: `${APIV1}/`,
    catalogQuerySummary: `${APIV1}/catalog/`,
    catalogQuerySummaryByID: `${APIV1}/catalog/:id`,
    contactQuery: `${APIV1}/contact`,
    aboutQuery: `${APIV1}/about`,
    articleQuery: `${APIV1}/content/:id`,
    userStatus: `${APIV1}/maintain/status`,
    userLogin: `${APIV1}/maintain/login`,
    userLogout: `${APIV1}/maintain/logout`,
    summaryQuery: `${APIV1}/maintain/summary`,
    catalogCreate: `${APIV1}/maintain/catalog`,
    catalogQuery: `${APIV1}/maintain/catalog/:id`,
    articleCreate: `${APIV1}/maintain/article`,
    catalogDelete: `${APIV1}/maintain/catalog/:id`,
    articleDelete: `${APIV1}/maintain/article/:id`,
    catalogUpdate: `${APIV1}/maintain/catalog/:id`,
    articleUpdate: `${APIV1}/maintain/article/:id`,
    noFoundPage: `${APIV1}/404.html`,
  },
}
