import React from 'react';
import { connect } from 'dva';
import styles from './Catalogs.css';

function Catalogs() {
  return (
    <div className={styles.normal}>
      Route Component: Catalogs
    </div>
  );
}

function mapStateToProps() {
  return {};
}

export default connect(mapStateToProps)(Catalogs);
