import React from 'react'
import PropTypes from 'prop-types'
import { SummaryList } from '../SummaryPanel'
import { ArticleView } from '../ArticlePanel'
import EditBar from './EditBar'

function ContentPanel({ contentData }) {
  const { summary, content } = contentData
  const { type } = summary

  const getContent = (typeValue, value) => {
    if (typeValue === 'article') {
      return <ArticleView article={value} />
    } else if (typeValue === 'catalog') {
      return <SummaryList summaryList={value} />
    } else {
      return <div>aaa</div>
    }
  }

  const getBar = (typeValue) => {
    if (typeValue === 'catalog') {
      return <EditBar />
    }
  }

  return (
    <div>
      { getContent(type, content) }
      { getBar(type) }
    </div>
  )
}

ContentPanel.propTypes = {
  contentData: PropTypes.object,
}

export default ContentPanel
