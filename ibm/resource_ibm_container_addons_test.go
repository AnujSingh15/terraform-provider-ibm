/* IBM Confidential
*  Object Code Only Source Materials
*  5747-SM3
*  (c) Copyright IBM Corp. 2017,2021
*
*  The source code for this program is not published or otherwise divested
*  of its trade secrets, irrespective of what has been deposited with the
*  U.S. Copyright Office. */

package ibm

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
)

func TestAccIBMContainerAddOns_Basic(t *testing.T) {
	name := fmt.Sprintf("tf-cluster-addon-%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMContainerAddOnsDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMContainerAddOnsBasic(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_container_addons.addons", "addons.#", "2"),
				),
			},
			{
				Config: testAccCheckIBMContainerAddOnsUpdate(name),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_container_addons.addons", "addons.#", "1"),
				),
			},
			{
				ResourceName:      "ibm_container_addons.addons",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccCheckIBMContainerAddOnsDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_container_addons" {
			continue
		}
		targetEnv := v1.ClusterTargetHeader{
			Region: "us-south",
		}
		csClient, err := testAccProvider.Meta().(ClientSession).ContainerAPI()
		if err != nil {
			return err
		}
		cluster := rs.Primary.ID
		addOnAPI := csClient.AddOns()
		_, err = addOnAPI.GetAddons(cluster, targetEnv)
		if err == nil {
			return fmt.Errorf("AddOns still exists: %s", rs.Primary.ID)
		} else if !strings.Contains(err.Error(), "404") {
			return fmt.Errorf("Error checking if AddOns (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}
	return nil
}

func testAccCheckIBMContainerAddOnsBasic(name string) string {
	return fmt.Sprintf(`
	provider "ibm"{
		region = "eu-de"
	}
	resource "ibm_is_vpc" "vpc" {
		name = "%[1]s"
	}
	resource "ibm_is_subnet" "subnet" {
		name                     = "%[1]s"
		vpc                      = ibm_is_vpc.vpc.id
		zone                     = "eu-de-1"
		total_ipv4_address_count = 256
	}
	resource "ibm_container_vpc_cluster" "cluster" {
		name              = "%[1]s"
		vpc_id            = ibm_is_vpc.vpc.id
		flavor            = "cx2.2x4"
		worker_count      = 1
		wait_till         = "OneWorkerNodeReady"
		zones {
			subnet_id = ibm_is_subnet.subnet.id
			name      = "eu-de-1"
		}
	}
	resource "ibm_container_addons" "addons" {
		cluster = ibm_container_vpc_cluster.cluster.id
		addons {
			name    = "vpc-block-csi-driver"
		}
		addons {
			name    = "cluster-autoscaler"
		}
}`, name)
}
func testAccCheckIBMContainerAddOnsUpdate(name string) string {
	return fmt.Sprintf(`
	provider "ibm"{
		region = "eu-de"
	}
	resource "ibm_is_vpc" "vpc" {
		name = "%[1]s"
	}
	resource "ibm_is_subnet" "subnet" {
		name                     = "%[1]s"
		vpc                      = ibm_is_vpc.vpc.id
		zone                     = "eu-de-1"
		total_ipv4_address_count = 256
	}
	resource "ibm_container_vpc_cluster" "cluster" {
		name              = "%[1]s"
		vpc_id            = ibm_is_vpc.vpc.id
		flavor            = "cx2.2x4"
		worker_count      = 1
		wait_till         = "OneWorkerNodeReady"
		zones {
			subnet_id = ibm_is_subnet.subnet.id
			name      = "eu-de-1"
		}
	}
	resource "ibm_container_addons" "addons" {
		cluster = ibm_container_vpc_cluster.cluster.id
		addons {
			name    = "cluster-autoscaler"
		}
}`, name)
}
