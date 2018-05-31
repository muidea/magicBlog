import React from 'react'
import PropTypes from 'prop-types'
import { Row, Button } from 'antd'
import styles from './EditBar.less'

function EditBar() {
  return (
    <Row type="flex" justify="center">
      <Button type="primary" className={styles.button}>新增分类</Button>
      <Button type="primary" className={styles.button}>新增文章</Button>
    </Row>
  )
}

EditBar.propTypes = {
  onAddCatalog: PropTypes.func,
  onAddArticle: PropTypes.func,
}

export default EditBar
