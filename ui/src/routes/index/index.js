import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { SummaryView } from '../common'

function IndexPage({ index, dispatch }) {
  const { summaryList } = index

  const onSelect = (value) => {
    dispatch({ type: 'app/redirectContent', payload: { ...value } })
  }

  return (
    <SummaryView summaryList={summaryList} readOnly onSelect={onSelect} />
  )
}

IndexPage.propTypes = {
  index: PropTypes.object,
  dispatch: PropTypes.func,
}

export default connect(({ index, dispatch }) => ({ index, dispatch }))(IndexPage)
