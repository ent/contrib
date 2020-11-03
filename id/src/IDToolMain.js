/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import EntDetails from './EntDetails';
import EntSearchBar from './EntSearchBar';
import React from 'react';
import Text from './Text';
import {Route} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(theme => ({
  '@global': {
    body: {
      margin: 0,
    },
  },
  root: {
    display: 'flex',
  },
  content: {
    height: '100vh',
    width: '100vw',
    display: 'flex',
    flexDirection: 'column',
  },
  placeholderContainer: {
    display: 'flex',
    flexGrow: 1,
    alignItems: 'center',
    justifyContent: 'center',
    color: '#303846',
  },
}));

const IDToolMain = () => {
  const classes = useStyles();

  return (
      <div className={classes.root}>
          <div className={classes.content}>
            <EntSearchBar />
            <Route path="/id/:id" component={EntDetails} />
            <Route
              path="/id/"
              exact
              render={() => (
                <div className={classes.placeholderContainer}>
                  <Text>Enter an ID above to see its fields and edges</Text>
                </div>
              )}
            />
          </div>
      </div>
  );
}

export default IDToolMain;
