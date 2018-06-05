import React from 'react'
import PropTypes from 'prop-types'
import queryString from 'query-string'
import { Row, Button } from 'antd'
import styles from './ContentBar.less'

function ContentBar({ currentItem }) {
  const AddUrl = (item) => {
    const { id, name, type } = item
    item = { action: 'add', id, name, type }

    return '/maintain?'.concat(queryString.stringify(item))
  }

  return (
    <Row type="flex" justify="center">
      <Button type="primary" className={styles.button} href={AddUrl(currentItem)}>新增分类</Button>
      <Button type="primary" className={styles.button} href={AddUrl(currentItem)}>新增文章</Button>
    </Row>
  )
}

ContentBar.propTypes = {
  currentItem: PropTypes.object,
}

export default ContentBar
