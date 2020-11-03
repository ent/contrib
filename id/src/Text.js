/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

 import {makeStyles} from '@material-ui/styles';
 import classNames from 'classnames';
 import React from 'react';

 const useStyles = makeStyles(() => ({
    text: {
        fontFamily: '"Roboto", "Helvetica", "Arial", sans-serif',
        fontWeight: 400,
        fontSize: '14px',
        lineHeight: 1.43,
        letterSpacing: '0.25px',
    },
 }));

 type Props = $ReadOnly<{|
    children: ?React.Node,
    className?: string,
 |}>;

 const Text = ({children, className, ...rest}: Props) => {
    const classes = useStyles();
    return (
        <span
        {...rest}
        className={classNames(classes.text, className)}>
        {children}
        </span>
    );
 }

 export default Text;