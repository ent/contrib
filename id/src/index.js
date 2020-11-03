import React from 'react';
import ReactDOM from 'react-dom';
import IDToolMain from './IDToolMain';
import {BrowserRouter} from 'react-router-dom';

ReactDOM.render(
  <React.StrictMode>
    <BrowserRouter>
      <IDToolMain />
    </BrowserRouter>
  </React.StrictMode>,
  document.getElementById('root')
);

