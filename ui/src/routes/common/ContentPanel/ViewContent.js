import React from 'react'
import PropTypes from 'prop-types'
import { SummaryList } from '../SummaryPanel'
import { ArticleView } from '../ArticlePanel'
import EditBar from './EditBar'

function ViewContent({ contentData, onAddCatalog, onAddArticle }) {
  const { currentItem, data } = contentData
  const { type } = currentItem

  const getContent = (typeValue, value) => {
    if (typeValue === 'article') {
      return <ArticleView article={value} />
    } else if (typeValue === 'catalog') {
      return <SummaryList summaryList={value} />
    } else {
      return <div>invalid typeValue</div>
    }
  }

  const getBar = (typeValue, item) => {
    if (typeValue === 'catalog') {
      return <EditBar onAddCatalog={onAddCatalog} onAddArticle={onAddArticle} currentItem={item} />
    }
  }

  return (
    <div>
      { getContent(type, data) }
      { getBar(type, currentItem) }
    </div>
  )
}

ViewContent.propTypes = {
  contentData: PropTypes.object,
  onAddCatalog: PropTypes.func,
  onAddArticle: PropTypes.func,
}

export default ViewContent
