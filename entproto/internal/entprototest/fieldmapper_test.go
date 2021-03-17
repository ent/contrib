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

	status, ok := mp["status"]
	require.True(ok)
	assert.EqualValues("Status", status.PbStructField())
	assert.True(status.IsEnumFIeld)

	for _, en := range mp.Edges() {
		assert.True(en.IsEdgeField, "expected .Edges() to return only Edge fields")
	}
	for _, en := range mp.Enums() {
		assert.True(en.IsEnumFIeld, "expected .Enums() to return only enum fields")
	}
}
