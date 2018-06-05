import React from 'react'
import PropTypes from 'prop-types'
import { SummaryView } from '../SummaryPanel'
import { ArticleView, ArticleEditor } from '../ArticlePanel'
import { CatalogEditor } from '../CatalogPanel'
import ContentBar from './ContentBar'

function ContentView({ contentData }) {
  const { command, id, type, name, data } = contentData

  const getContent = () => {
    if (command === 'add') {
      return getAddContent()
    } else if (command === 'modify') {
      return getModifyContent()
    } else {
      return getViewContent()
    }
  }

  const getViewContent = () => {
    if (type === 'article') {
      return <ArticleView article={data} />
    } else {
      let summaryList = []
      if (data !== null) {
        summaryList = data
      }
      return <SummaryView summaryList={summaryList} />
    }
  }

  const getAddContent = () => {
    if (type === 'article') {
      return <ArticleEditor content={data} />
    } else {
      return <CatalogEditor content={data} />
    }
  }

  const getModifyContent = () => {
    if (type === 'article') {
      return <ArticleEditor content={data} />
    } else {
      return <CatalogEditor content={data} />
    }
  }

  const getBar = (typeValue, item) => {
    if (typeValue === 'catalog') {
      return <ContentBar currentItem={item} />
    }
  }

  return (
    <div>
      { getContent() }
    </div>
  )
}

ContentView.propTypes = {
  contentData: PropTypes.object,
}

export default ContentView
