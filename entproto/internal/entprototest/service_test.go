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

func (suite *AdapterTestSuite) TestServiceGeneration() {
	fd, err := suite.adapter.GetFileDescriptor("BlogPost")
	suite.Require().NoError(err)

	svc := fd.FindService("entpb.BlogPostService")
	suite.NotNil(svc)

	getMeth := svc.FindMethodByName("Get")
	suite.Require().NotNil(getMeth)
	suite.EqualValues("GetBlogPostRequest", getMeth.GetInputType().GetName())
	suite.EqualValues("BlogPost", getMeth.GetOutputType().GetName())

	createMeth := svc.FindMethodByName("Create")
	suite.Require().NotNil(createMeth)
	suite.EqualValues("CreateBlogPostRequest", createMeth.GetInputType().GetName())
	suite.EqualValues("BlogPost", createMeth.GetOutputType().GetName())

	deleteMeth := svc.FindMethodByName("Delete")
	suite.Require().NotNil(deleteMeth)
	suite.EqualValues("DeleteBlogPostRequest", deleteMeth.GetInputType().GetName())
	suite.EqualValues("google.protobuf.Empty", deleteMeth.GetOutputType().GetFullyQualifiedName())

	updateMeth := svc.FindMethodByName("Update")
	suite.Require().NotNil(updateMeth)
	suite.EqualValues("UpdateBlogPostRequest", updateMeth.GetInputType().GetName())
	suite.EqualValues("BlogPost", updateMeth.GetOutputType().GetName())
}
