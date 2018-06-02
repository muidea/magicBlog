import React from 'react'
import PropTypes from 'prop-types'
import { Button, Form, Input } from 'antd'
import styles from './index.less'

const FormItem = Form.Item
const { TextArea } = Input

const EditCatalog = ({
  onSubmit,
  contentItem,
  form: {
    getFieldDecorator,
    validateFieldsAndScroll,
  },
}) => {
  const handleOk = () => {
    validateFieldsAndScroll((errors, values) => {
      if (errors) {
        return
      }

      const { parent } = contentItem
      const { id, type } = parent

      values = {
        ...values,
        parent: { id, type },
      }

      onSubmit(values)
    })
  }

  return (
    <div className={styles.form}>
      <div>
        <form>
          <FormItem label="分类名" hasFeedback>
            {getFieldDecorator('name', {
              initialValue: contentItem.name,
              rules: [{ required: true }],
            })(<Input size="large" onPressEnter={handleOk} placeholder="输入分类名" />)}
          </FormItem>
          <FormItem label="描述">
            {getFieldDecorator('description', {
              initialValue: contentItem.description,
            })(<TextArea rows={3} cols={3} />)}
          </FormItem>
        </form>
      </div>
      <div>
        <Button type="primary" size="large" onClick={handleOk}>
            添加
        </Button>
      </div>
    </div>
  )
}

EditCatalog.propTypes = {
  form: PropTypes.object,
  onSubmit: PropTypes.func,
  contentItem: PropTypes.object,
}

export default Form.create()(EditCatalog)
