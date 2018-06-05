import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { SummaryView } from '../common'

function IndexPage({ index }) {
  const { summaryList } = index
  return (
    <SummaryView summaryList={summaryList} readOnly />
  )
}

IndexPage.propTypes = {
  index: PropTypes.object,
}

export default connect(({ index }) => ({ index }))(IndexPage)
