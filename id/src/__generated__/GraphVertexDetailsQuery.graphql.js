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
export type GraphVertexDetailsQueryVariables = {|
  id: string
|};
export type GraphVertexDetailsQueryResponse = {|
  +vertex: ?{|
    +id: string,
    +type: string,
    +fields: $ReadOnlyArray<{|
      +name: string,
      +value: string,
    |}>,
  |}
|};
export type GraphVertexDetailsQuery = {|
  variables: GraphVertexDetailsQueryVariables,
  response: GraphVertexDetailsQueryResponse,
|};
*/


/*
query GraphVertexDetailsQuery(
  $id: ID!
) {
  vertex(id: $id) {
    id
    type
    fields {
      name
      value
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
v1 = [
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
      {
        "alias": null,
        "args": null,
        "kind": "ScalarField",
        "name": "type",
        "storageKey": null
      },
      {
        "alias": null,
        "args": null,
        "concreteType": "Field",
        "kind": "LinkedField",
        "name": "fields",
        "plural": true,
        "selections": [
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "name",
            "storageKey": null
          },
          {
            "alias": null,
            "args": null,
            "kind": "ScalarField",
            "name": "value",
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
    "name": "GraphVertexDetailsQuery",
    "selections": (v1/*: any*/),
    "type": "Query",
    "abstractKey": null
  },
  "kind": "Request",
  "operation": {
    "argumentDefinitions": (v0/*: any*/),
    "kind": "Operation",
    "name": "GraphVertexDetailsQuery",
    "selections": (v1/*: any*/)
  },
  "params": {
    "cacheID": "3ab5c4ce13ed255d1183d7fab15a5feb",
    "id": null,
    "metadata": {},
    "name": "GraphVertexDetailsQuery",
    "operationKind": "query",
    "text": "query GraphVertexDetailsQuery(\n  $id: ID!\n) {\n  vertex(id: $id) {\n    id\n    type\n    fields {\n      name\n      value\n    }\n  }\n}\n"
  }
};
})();
// prettier-ignore
(node/*: any*/).hash = 'b810a4963a64c2695e2cb81fbd02be9e';

module.exports = node;
