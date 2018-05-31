import React from 'react'
import PropTypes from 'prop-types'
import SummaryList from '../SummaryList'
import ArticleView from '../ArticleView'

function SummaryContent({ contentData }) {
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

  return (
      getContent(type, content)
  )
}

SummaryContent.propTypes = {
  contentData: PropTypes.object,
}

export default SummaryContent
