import React from 'react'
import PropTypes from 'prop-types'
import { Row, Button } from 'antd'
import styles from './ContentBar.less'

function ContentBar({ currentItem, onAdd }) {
  const { id, name, type } = currentItem

  return (
    <Row type="flex" justify="center">
      <Button type="primary" className={styles.button} onClick={() => onAdd({ id, name, type, data: { type: 'catalog' } })}>新增分类</Button>
      <Button type="primary" className={styles.button} onClick={() => onAdd({ id, name, type, data: { type: 'article' } })}>新增文章</Button>
    </Row>
  )
}

ContentBar.propTypes = {
  currentItem: PropTypes.object,
  onAdd: PropTypes.func,
}

export default ContentBar
