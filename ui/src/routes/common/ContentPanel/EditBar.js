import React from 'react'
import PropTypes from 'prop-types'
import { Row, Button } from 'antd'
import styles from './EditBar.less'

function EditBar({ onAddCatalog, onAddArticle, currentItem }) {
  const onAddNewCatalog = () => {
    onAddCatalog(currentItem)
  }

  const onAddNewArticle = () => {
    onAddArticle(currentItem)
  }

  return (
    <Row type="flex" justify="center">
      <Button type="primary" className={styles.button} onClick={onAddNewCatalog}>新增分类</Button>
      <Button type="primary" className={styles.button} onClick={onAddNewArticle}>新增文章</Button>
    </Row>
  )
}

EditBar.propTypes = {
  onAddCatalog: PropTypes.func,
  onAddArticle: PropTypes.func,
  currentItem: PropTypes.object,
}

export default EditBar
