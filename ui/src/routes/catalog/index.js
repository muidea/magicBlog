import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { SummaryView } from '../common'

function CatalogPage({ catalog, dispatch }) {
  const { summaryList } = catalog

  const onSelect = (value) => {
    dispatch({ type: 'app/redirectContent', payload: { ...value } })
  }

  return (
    <SummaryView summaryList={summaryList} readOnly onSelect={onSelect} />
  )
}

CatalogPage.propTypes = {
  catalog: PropTypes.object,
  dispatch: PropTypes.func,
}

export default connect(({ catalog, dispatch }) => ({ catalog, dispatch }))(CatalogPage)
