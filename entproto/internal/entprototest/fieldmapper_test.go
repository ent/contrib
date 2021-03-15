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
	mp, err := suite.adapter.FieldMap("User")
	suite.Require().NoError(err)
	suite.Require().NotNil(mp)

	userName, ok := mp["user_name"]
	suite.Require().True(ok)
	suite.Assert().False(userName.IsIDField)
	suite.Assert().False(userName.IsEdgeField)
	suite.Assert().NotNil(userName.EntField)
	suite.Assert().NotNil(userName.PbFieldDescriptor)
	suite.Assert().EqualValues("UserName", userName.PbStructField())

	id, ok := mp["id"]
	suite.Require().True(ok)
	suite.Assert().True(id.IsIDField)
	suite.Assert().False(id.IsEdgeField)
	suite.Assert().EqualValues("Id", id.PbStructField())

	blogPosts, ok := mp["blog_posts"]
	suite.Require().True(ok)
	suite.Assert().EqualValues("BlogPosts", blogPosts.PbStructField())
	suite.Assert().False(blogPosts.IsIDField)
	suite.Assert().True(blogPosts.IsEdgeField)
}
