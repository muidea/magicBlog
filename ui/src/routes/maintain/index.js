import React from 'react'
import PropTypes from 'prop-types'
import queryString from 'query-string'
import { connect } from 'dva'
import { Row, Col } from 'antd'
import { ContentNav, ContentView } from '../common'
import styles from './index.less'

function MaintainPage({ maintain, dispatch }) {
  const { itemList, action } = maintain

  const onSelect = (value) => {
    const { id, type, name } = value
    const url = '/maintain?'.concat(queryString.stringify({ command: 'view', id, type, name }))

    dispatch({ type: 'maintain/redirectContent', payload: { url } })
  }

  const onSubmit = (value) => {
    dispatch({ type: 'maintain/submitContent', payload: { ...value } })
  }

  return (
    <Row type="flex" align="top">
      <Col md={6} lg={6} xl={6} className={styles.nav}>
        <ContentNav itemList={itemList} onSelect={onSelect} />
      </Col>
      <Col md={18} lg={18} xl={18}>
        <ContentView contentData={action} onSubmit={onSubmit} />
      </Col>
    </Row>
  )
}

MaintainPage.propTypes = {
  dispatch: PropTypes.func,
}

export default connect(({ maintain, dispatch }) => ({ maintain, dispatch }))(MaintainPage)
