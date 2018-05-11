import React from 'react';
import { Router, Route } from 'dva/router';
import IndexPage from './routes/IndexPage';

import Catalogs from './routes/Catalogs.js';

function RouterConfig({ history }) {
  return (
    <Router history={history}>
      <Route path="/" component={IndexPage} />
      <Route path="/catalogs" component={Catalogs} />
    </Router>
  );
}

export default RouterConfig;
