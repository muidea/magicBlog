import React from 'react'
import PropTypes from 'prop-types'
import { Row } from 'antd'

function EditArticle() {
  return (
    <Row type="flex" justify="center">
      <div>EditArticle</div>
    </Row>
  )
}

EditArticle.propTypes = {
  onAddCatalog: PropTypes.func,
  onAddArticle: PropTypes.func,
}

export default EditArticle
