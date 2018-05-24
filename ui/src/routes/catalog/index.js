import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { SummaryList } from '../../components'

function CatalogPage({ catalog }) {
  const { summaryList } = catalog

  return (
    <SummaryList summaryList={summaryList} />
  )
}

CatalogPage.propTypes = {
  catalog: PropTypes.object,
}

export default connect(({ catalog }) => ({ catalog }))(CatalogPage)