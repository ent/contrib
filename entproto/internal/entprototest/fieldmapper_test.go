// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package entprototest

func (suite *AdapterTestSuite) TestFieldMap() {
	require := suite.Require()
	assert := suite.Assert()

	mp, err := suite.adapter.FieldMap("User")
	require.NoError(err)
	require.NotNil(mp)

	userName, ok := mp["user_name"]
	require.True(ok)
	assert.False(userName.IsIDField)
	assert.False(userName.IsEdgeField)
	assert.NotNil(userName.EntField)
	assert.NotNil(userName.PbFieldDescriptor)
	assert.EqualValues("UserName", userName.PbStructField())

	profilePic, ok := mp["profile_pic"]
	require.True(ok)
	assert.True(profilePic.IsEdgeField)

	id, ok := mp["id"]
	require.True(ok)
	assert.True(id.IsIDField)
	assert.False(id.IsEdgeField)
	assert.EqualValues("Id", id.PbStructField())

	blogPosts, ok := mp["blog_posts"]
	require.True(ok)
	assert.EqualValues("BlogPosts", blogPosts.PbStructField())
	assert.False(blogPosts.IsIDField)
	assert.True(blogPosts.IsEdgeField)
	assert.EqualValues("Id", blogPosts.EdgeIDPbStructField())
	assert.EqualValues("id", blogPosts.EdgeIDPbStructFieldDesc().GetName())

	status, ok := mp["status"]
	require.True(ok)
	assert.EqualValues("Status", status.PbStructField())
	assert.True(status.IsEnumField)

	for _, en := range mp.Edges() {
		assert.True(en.IsEdgeField, "expected .Edges() to return only Edge fields")
		assert.NotNil(en.EntEdge)
	}
	for _, en := range mp.Enums() {
		assert.True(en.IsEnumField, "expected .Enums() to return only enum fields")
	}
}

func (suite *AdapterTestSuite) TestExternalId() {
	require := suite.Require()
	assert := suite.Assert()

	mp, err := suite.adapter.FieldMap("BlogPost")
	require.NoError(err)
	require.NotNil(mp)
	eid, ok := mp["external_id"]
	require.True(ok)
	assert.EqualValues("ExternalId", eid.PbStructField())
}

func (suite *AdapterTestSuite) TestReferenced() {
	require := suite.Require()

	mp, err := suite.adapter.FieldMap("BlogPost")
	require.NoError(err)
	require.NotNil(mp)
	cats, ok := mp["categories"]
	require.True(ok)
	require.NotNil(cats)
	require.EqualValues(cats.ReferencedPbType.GetName(), "Category")
	auth, ok := mp["author"]
	require.True(ok)
	require.NotNil(auth)
	require.EqualValues(auth.ReferencedPbType.GetName(), "User")
}

func (suite *AdapterTestSuite) TestNoBackref() {
	require := suite.Require()
	mp, err := suite.adapter.FieldMap("NoBackref")
	require.NoError(err)
	require.Equal("Id", mp.Edges()[0].EdgeIDPbStructField())
}
