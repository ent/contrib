/**
 * @generated
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 **/

 /**
 * @flow
 */

/* eslint-disable */

'use strict';

/*::
import type { ConcreteRequest } from 'relay-runtime';
export type EntDetailsQueryVariables = {|
  id: string
|};
export type EntDetailsQueryResponse = {|
  +vertex: ?{|
    +id: string,
    +type: string,
    +fields: $ReadOnlyArray<{|
      +name: string,
      +value: string,
      +type: string,
    |}>,
    +edges: $ReadOnlyArray<{|
      +name: string,
      +type: string,
      +ids: $ReadOnlyArray<string>,
    |}>,
  |}
|};
export type EntDetailsQuery = {|
  variables: EntDetailsQueryVariables,
  response: EntDetailsQueryResponse,
|};
*/


/*
query EntDetailsQuery(
  $id: ID!
) {
  vertex(id: $id) {
    id
    type
    fields {
      name
      value
      type
    }
    edges {
      name
      type
      ids
    }
  }
}
*/

const node/*: ConcreteRequest*/ = (function(){
var v0 = [
  {
    "defaultValue": null,
    "kind": "LocalArgument",
    "name": "id"
  }
],
v1 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "type",
  "storageKey": null
},
v2 = {
  "alias": null,
  "args": null,
  "kind": "ScalarField",
  "name": "name",
  "storageKey": null
},
v3 = [
  {
    "alias": null,
    "args": [
      {
        "kind": "Variable",
        "name": "id",
        "variableName": "id"
      }
    ],
    "concreteType": "Vertex",
    "kind": "LinkedField",
    "name": "vertex",
    "plural": false,
    "selections": [
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "id",
        "storageKey": null
      },
      (v1/*: any*/),
      {
        "alias": null,
        "args": null,
        "concreteType": "Field",
        "kind": "LinkedField",
        "name": "fields",
        "plural": true,
        "selections": [
          (v2/*: any*/),
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "value",
            "storageKey": null
          },
          (v1/*: any*/)
        ],
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "concreteType": "Edge",
        "kind": "LinkedField",
        "name": "edges",
        "plural": true,
        "selections": [
          (v2/*: any*/),
          (v1/*: any*/),
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "ids",
            "storageKey": null
          }
        ],
        "storageKey": null
      }
    ],
    "storageKey": null
  }
];
return {
  "fragment": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Fragment",
    "metadata": null,
    "name": "EntDetailsQuery",
    "selections": (v3/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "EntDetailsQuery",
    "selections": (v3/*: any*/)
  },
  "params": {
    "cacheID": "1c539c26a6c41aaacc0c449a2f231b02",
    "id": null,
    "metadata": {},
    "name": "EntDetailsQuery",
    "operationKind": "query",
    "text": "query EntDetailsQuery(\n  $id: ID!\n) {\n  vertex(id: $id) {\n    id\n    type\n    fields {\n      name\n      value\n      type\n    }\n    edges {\n      name\n      type\n      ids\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = '4be33521bcf990ae9df86c758a47fcde';

module.exports = node;
