import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { SummaryList } from '../common'

function IndexPage({ index }) {
  const { summaryList } = index
  return (
    <SummaryList summaryList={summaryList} />
  )
}

IndexPage.propTypes = {
  index: PropTypes.object,
}

export default connect(({ index }) => ({ index }))(IndexPage)
