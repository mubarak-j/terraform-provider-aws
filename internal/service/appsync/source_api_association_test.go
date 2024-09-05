// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package appsync_test

import (
	"context"
	"errors"
	"fmt"
	"github.com/YakDriver/regexache"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/appsync"
	"github.com/aws/aws-sdk-go-v2/service/appsync/types"
	sdkacctest "github.com/hashicorp/terraform-plugin-testing/helper/acctest"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
	"github.com/hashicorp/terraform-provider-aws/internal/acctest"
	"github.com/hashicorp/terraform-provider-aws/internal/conns"
	"github.com/hashicorp/terraform-provider-aws/internal/create"
	tfappsync "github.com/hashicorp/terraform-provider-aws/internal/service/appsync"
	"github.com/hashicorp/terraform-provider-aws/internal/tfresource"
	"github.com/hashicorp/terraform-provider-aws/names"
)

func testAccAppSyncSourceApiAssociation_basic(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var sourceapiassociation types.SourceApiAssociation
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_appsync_source_api_association.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.AppSyncEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.AppSyncServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSourceApiAssociationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSourceApiAssociationConfig_basic(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourceApiAssociationExists(ctx, resourceName, &sourceapiassociation),
					resource.TestCheckResourceAttr(resourceName, names.AttrDescription, rName),
					acctest.MatchResourceAttrRegionalARN(resourceName, names.AttrARN, "appsync", regexache.MustCompile(`apis/.+/sourceApiAssociations/.+`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAppSyncSourceApiAssociation_update(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var sourceapiassociation types.SourceApiAssociation
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_appsync_source_api_association.test"
	updateDesc := rName + "2"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.AppSyncEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.AppSyncServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSourceApiAssociationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSourceApiAssociationConfig_basic(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourceApiAssociationExists(ctx, resourceName, &sourceapiassociation),
					resource.TestCheckResourceAttr(resourceName, names.AttrDescription, rName),
					acctest.MatchResourceAttrRegionalARN(resourceName, names.AttrARN, "appsync", regexache.MustCompile(`apis/.+/sourceApiAssociations/.+`)),
				),
			},
			{
				Config: testAccSourceApiAssociationConfig_basic(rName, updateDesc),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourceApiAssociationExists(ctx, resourceName, &sourceapiassociation),
					resource.TestCheckResourceAttr(resourceName, names.AttrDescription, updateDesc),
					acctest.MatchResourceAttrRegionalARN(resourceName, names.AttrARN, "appsync", regexache.MustCompile(`apis/.+/sourceApiAssociations/.+`)),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAppSyncSourceApiAssociation_disappears(t *testing.T) {
	ctx := acctest.Context(t)
	if testing.Short() {
		t.Skip("skipping long-running test in short mode")
	}

	var sourceapiassociation types.SourceApiAssociation
	rName := sdkacctest.RandomWithPrefix(acctest.ResourcePrefix)
	resourceName := "aws_appsync_source_api_association.test"

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acctest.PreCheck(ctx, t)
			acctest.PreCheckPartitionHasService(t, names.AppSyncEndpointID)
			testAccPreCheck(ctx, t)
		},
		ErrorCheck:               acctest.ErrorCheck(t, names.AppSyncServiceID),
		ProtoV5ProviderFactories: acctest.ProtoV5ProviderFactories,
		CheckDestroy:             testAccCheckSourceApiAssociationDestroy(ctx),
		Steps: []resource.TestStep{
			{
				Config: testAccSourceApiAssociationConfig_basic(rName, rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckSourceApiAssociationExists(ctx, resourceName, &sourceapiassociation),
					acctest.CheckFrameworkResourceDisappears(ctx, acctest.Provider, tfappsync.ResourceSourceApiAssociation, resourceName),
				),
				ExpectNonEmptyPlan: true,
			},
		},
	})
}

func testAccCheckSourceApiAssociationDestroy(ctx context.Context) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		conn := acctest.Provider.Meta().(*conns.AWSClient).AppSyncClient(ctx)

		for _, rs := range s.RootModule().Resources {
			if rs.Type != "aws_appsync_source_api_association" {
				continue
			}

			_, err := tfappsync.FindSourceApiAssociationByTwoPartKey(ctx, conn, rs.Primary.Attributes["association_id"], rs.Primary.Attributes["merged_api_id"])

			if tfresource.NotFound(err) {
				continue
			}

			if err != nil {
				return err
			}

			return create.Error(names.AppSync, create.ErrActionCheckingDestroyed, tfappsync.ResNameSourceApiAssociation, rs.Primary.ID, errors.New("not destroyed"))
		}

		return nil
	}
}

func testAccCheckSourceApiAssociationExists(ctx context.Context, name string, sourceapiassociation *types.SourceApiAssociation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return create.Error(names.AppSync, create.ErrActionCheckingExistence, tfappsync.ResNameSourceApiAssociation, name, errors.New("not found"))
		}

		if rs.Primary.ID == "" {
			return create.Error(names.AppSync, create.ErrActionCheckingExistence, tfappsync.ResNameSourceApiAssociation, name, errors.New("not set"))
		}

		conn := acctest.Provider.Meta().(*conns.AWSClient).AppSyncClient(ctx)
		resp, err := tfappsync.FindSourceApiAssociationByTwoPartKey(ctx, conn, rs.Primary.Attributes["association_id"], rs.Primary.Attributes["merged_api_id"])

		if err != nil {
			return create.Error(names.AppSync, create.ErrActionCheckingExistence, tfappsync.ResNameSourceApiAssociation, rs.Primary.ID, err)
		}

		*sourceapiassociation = *resp

		return nil
	}
}

func testAccPreCheck(ctx context.Context, t *testing.T) {
	conn := acctest.Provider.Meta().(*conns.AWSClient).AppSyncClient(ctx)

	input := &appsync.ListGraphqlApisInput{}
	_, err := conn.ListGraphqlApis(ctx, input)

	if acctest.PreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}
	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccCheckSourceApiAssociationNotRecreated(before, after *types.SourceApiAssociation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if before, after := aws.ToString(before.AssociationId), aws.ToString(after.AssociationId); before != after {
			return create.Error(names.AppSync, create.ErrActionCheckingNotRecreated, tfappsync.ResNameSourceApiAssociation, before, errors.New("recreated"))
		}

		return nil
	}
}

func testAccSourceApiAssociationConfig_basic(rName, description string) string {
	return fmt.Sprintf(`
resource "aws_iam_role" "test" {
  assume_role_policy = data.aws_iam_policy_document.test.json
  name_prefix        = %[1]q
}

data "aws_caller_identity" "current" {}

data "aws_partition" "current" {}

data "aws_region" "current" {}

data "aws_iam_policy_document" "test" {
  statement {
    actions = ["sts:AssumeRole"]
    principals {
      identifiers = ["appsync.amazonaws.com"]
      type        = "Service"
    }
    condition {
      test     = "StringEquals"
      values   = [data.aws_caller_identity.current.account_id]
      variable = "aws:SourceAccount"
    }

    condition {
      test     = "ArnLike"
      values   = ["arn:${data.aws_partition.current.partition}:appsync:${data.aws_region.current.name}:${data.aws_caller_identity.current.account_id}::apis/*"]
      variable = "aws:SourceArn"
    }
  }
}

resource "aws_appsync_graphql_api" "merged" {
  authentication_type           = "API_KEY"
  name                          = %[1]q
  api_type                      = "MERGED"
  merged_api_execution_role_arn = aws_iam_role.test.arn
}

resource "aws_appsync_graphql_api" "source" {
  authentication_type           = "API_KEY"
  name                          = %[1]q
  schema = <<EOF
schema {
    query: Query
}
type Query {
  test: Int
}
EOF
}

resource "aws_appsync_source_api_association" "test" {
  description   = %[2]q
  merged_api_id = aws_appsync_graphql_api.merged.id
  source_api_id = aws_appsync_graphql_api.source.id
}
`, rName, description)
}
