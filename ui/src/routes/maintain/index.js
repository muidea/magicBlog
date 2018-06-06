import React from 'react'
import PropTypes from 'prop-types'
import { connect } from 'dva'
import { Row, Col } from 'antd'
import { ContentNav, ContentView } from '../common'
import styles from './index.less'

function MaintainPage({ maintain, dispatch }) {
  const { itemList, action } = maintain

  const onSelect = (value) => {
    const { id, type, name } = value

    dispatch({ type: 'maintain/refreshContent', payload: { command: 'view', id, type, name } })
  }

  const onAdd = (value) => {
    const { id, type, name } = value
    dispatch({ type: 'maintain/refreshContent', payload: { command: 'add', id, type, name } })
  }

  const onModify = (value) => {
    const { id, type, name } = value
    dispatch({ type: 'maintain/refreshContent', payload: { command: 'modify', id, type, name } })
  }

  const onDelete = (value) => {
    const { id, type, name } = value
    dispatch({ type: 'maintain/refreshContent', payload: { command: 'delete', id, type, name } })
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
        <ContentView contentData={action} onSelect={onSelect} onAdd={onAdd} onModify={onModify} onDelete={onDelete} onSubmit={onSubmit} />
      </Col>
    </Row>
  )
}

MaintainPage.propTypes = {
  dispatch: PropTypes.func,
}

export default connect(({ maintain, dispatch }) => ({ maintain, dispatch }))(MaintainPage)
