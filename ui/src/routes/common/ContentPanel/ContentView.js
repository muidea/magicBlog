import React from 'react'
import PropTypes from 'prop-types'
import { SummaryView } from '../SummaryPanel'
import { ArticleView, ArticleEditor } from '../ArticlePanel'
import { CatalogEditor } from '../CatalogPanel'
import ContentBar from './ContentBar'

function ContentView({ contentData, onSelect, onAdd, onModify, onDelete, onSubmit }) {
  const { command, id, type, name, data } = contentData

  const onSubmitData = (value) => {
    console.log(value)
    onSubmit({ ...value, command })
  }

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
      return <SummaryView summaryList={summaryList} onSelect={onSelect} onModify={onModify} onDelete={onDelete} />
    }
  }

  const getAddContent = () => {
    const content = { ...data, catalog: [{ id, name }] }
    if (data.type === 'article') {
      return <ArticleEditor content={content} onSubmit={onSubmitData} />
    } else if (data.type === 'catalog') {
      return <CatalogEditor content={content} onSubmit={onSubmitData} />
    }
  }

  const getModifyContent = () => {
    const content = { ...data }
    if (type === 'article') {
      return <ArticleEditor content={content} onSubmit={onSubmitData} />
    } else {
      return <CatalogEditor content={content} onSubmit={onSubmitData} />
    }
  }

  const getBar = () => {
    if (type === 'catalog') {
      const item = { id, type, name }
      return <ContentBar currentItem={item} onAdd={onAdd} />
    }
  }

  return (
    <div>
      { getContent() }
      { (command === 'view') && getBar() }
    </div>
  )
}

ContentView.propTypes = {
  contentData: PropTypes.object,
  onSelect: PropTypes.func,
  onAdd: PropTypes.func,
  onModify: PropTypes.func,
  onDelete: PropTypes.func,
}

export default ContentView
