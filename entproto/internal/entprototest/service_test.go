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
	// Test default method generation
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

	listMeth := svc.FindMethodByName("List")
	suite.Require().NotNil(listMeth)
	suite.EqualValues("ListBlogPostRequest", listMeth.GetInputType().GetName())
	suite.EqualValues("ListBlogPostResponse", listMeth.GetOutputType().GetName())

	batchCreateMeth := svc.FindMethodByName("BatchCreate")
	suite.Require().NotNil(batchCreateMeth)
	suite.EqualValues("BatchCreateBlogPostsRequest", batchCreateMeth.GetInputType().GetName())
	suite.EqualValues("BatchCreateBlogPostsResponse", batchCreateMeth.GetOutputType().GetName())

	// Test all method generation
	fd, err = suite.adapter.GetFileDescriptor("AllMethodsService")
	suite.Require().NoError(err)

	svc = fd.FindService("entpb.AllMethodsServiceService")
	suite.NotNil(svc)

	getMeth = svc.FindMethodByName("Get")
	suite.Require().NotNil(getMeth)
	suite.EqualValues("GetAllMethodsServiceRequest", getMeth.GetInputType().GetName())
	suite.EqualValues("AllMethodsService", getMeth.GetOutputType().GetName())

	createMeth = svc.FindMethodByName("Create")
	suite.Require().NotNil(createMeth)
	suite.EqualValues("CreateAllMethodsServiceRequest", createMeth.GetInputType().GetName())
	suite.EqualValues("AllMethodsService", createMeth.GetOutputType().GetName())

	deleteMeth = svc.FindMethodByName("Delete")
	suite.Require().NotNil(deleteMeth)
	suite.EqualValues("DeleteAllMethodsServiceRequest", deleteMeth.GetInputType().GetName())
	suite.EqualValues("google.protobuf.Empty", deleteMeth.GetOutputType().GetFullyQualifiedName())

	updateMeth = svc.FindMethodByName("Update")
	suite.Require().NotNil(updateMeth)
	suite.EqualValues("UpdateAllMethodsServiceRequest", updateMeth.GetInputType().GetName())
	suite.EqualValues("AllMethodsService", updateMeth.GetOutputType().GetName())

	listMeth = svc.FindMethodByName("List")
	suite.Require().NotNil(listMeth)
	suite.EqualValues("ListAllMethodsServiceRequest", listMeth.GetInputType().GetName())
	suite.EqualValues("ListAllMethodsServiceResponse", listMeth.GetOutputType().GetName())

	batchCreateMeth = svc.FindMethodByName("BatchCreate")
	suite.Require().NotNil(batchCreateMeth)
	suite.EqualValues("BatchCreateAllMethodsServicesRequest", batchCreateMeth.GetInputType().GetName())
	suite.EqualValues("BatchCreateAllMethodsServicesResponse", batchCreateMeth.GetOutputType().GetName())

	// Test single method generation
	fd, err = suite.adapter.GetFileDescriptor("OneMethodService")
	suite.Require().NoError(err)

	svc = fd.FindService("entpb.OneMethodServiceService")
	suite.NotNil(svc)

	batchCreateMeth = svc.FindMethodByName("BatchCreate")
	suite.Require().NotNil(batchCreateMeth)
	suite.EqualValues("BatchCreateOneMethodServicesRequest", batchCreateMeth.GetInputType().GetName())
	suite.EqualValues("BatchCreateOneMethodServicesResponse", batchCreateMeth.GetOutputType().GetName())

	getMeth = svc.FindMethodByName("Get")
	suite.Require().Nil(getMeth)

	createMeth = svc.FindMethodByName("Create")
	suite.Require().Nil(createMeth)

	deleteMeth = svc.FindMethodByName("Delete")
	suite.Require().Nil(deleteMeth)

	updateMeth = svc.FindMethodByName("Update")
	suite.Require().Nil(updateMeth)

	listMeth = svc.FindMethodByName("List")
	suite.Require().Nil(listMeth)

	// Test two method generation
	fd, err = suite.adapter.GetFileDescriptor("TwoMethodService")
	suite.Require().NoError(err)

	svc = fd.FindService("entpb.TwoMethodServiceService")
	suite.NotNil(svc)

	getMeth = svc.FindMethodByName("Get")
	suite.Require().NotNil(getMeth)
	suite.EqualValues("GetTwoMethodServiceRequest", getMeth.GetInputType().GetName())
	suite.EqualValues("TwoMethodService", getMeth.GetOutputType().GetName())

	createMeth = svc.FindMethodByName("Create")
	suite.Require().NotNil(createMeth)
	suite.EqualValues("CreateTwoMethodServiceRequest", createMeth.GetInputType().GetName())
	suite.EqualValues("TwoMethodService", createMeth.GetOutputType().GetName())

	deleteMeth = svc.FindMethodByName("Delete")
	suite.Require().Nil(deleteMeth)

	updateMeth = svc.FindMethodByName("Update")
	suite.Require().Nil(updateMeth)

	listMeth = svc.FindMethodByName("List")
	suite.Require().Nil(listMeth)

	batchCreateMeth = svc.FindMethodByName("BatchCreate")
	suite.Require().Nil(batchCreateMeth)

	// Test message with id generation
	fd, err = suite.adapter.GetFileDescriptor("MessageWithID")
	suite.Require().NoError(err)

	svc = fd.FindService("entpb.MessageWithIDService")
	suite.NotNil(svc)

	getMeth = svc.FindMethodByName("Get")
	suite.Require().NotNil(getMeth)
	suite.EqualValues("GetMessageWithIDRequest", getMeth.GetInputType().GetName())
	suite.EqualValues("MessageWithID", getMeth.GetOutputType().GetName())

	createMeth = svc.FindMethodByName("Create")
	suite.Require().NotNil(createMeth)
	suite.EqualValues("CreateMessageWithIDRequest", createMeth.GetInputType().GetName())
	suite.EqualValues("MessageWithID", createMeth.GetOutputType().GetName())

	deleteMeth = svc.FindMethodByName("Delete")
	suite.Require().NotNil(deleteMeth)
	suite.EqualValues("DeleteMessageWithIDRequest", deleteMeth.GetInputType().GetName())
	suite.EqualValues("google.protobuf.Empty", deleteMeth.GetOutputType().GetFullyQualifiedName())

	updateMeth = svc.FindMethodByName("Update")
	suite.Require().NotNil(updateMeth)
	suite.EqualValues("UpdateMessageWithIDRequest", updateMeth.GetInputType().GetName())
	suite.EqualValues("MessageWithID", updateMeth.GetOutputType().GetName())

	listMeth = svc.FindMethodByName("List")
	suite.Require().NotNil(listMeth)
	suite.EqualValues("ListMessageWithIDRequest", listMeth.GetInputType().GetName())
	suite.EqualValues("ListMessageWithIDResponse", listMeth.GetOutputType().GetName())

	batchCreateMeth = svc.FindMethodByName("BatchCreate")
	suite.Require().NotNil(batchCreateMeth)
	suite.EqualValues("BatchCreateMessageWithIDsRequest", batchCreateMeth.GetInputType().GetName())
	suite.EqualValues("BatchCreateMessageWithIDsResponse", batchCreateMeth.GetOutputType().GetName())

}
