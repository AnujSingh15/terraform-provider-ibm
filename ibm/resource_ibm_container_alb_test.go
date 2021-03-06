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

	v1 "github.com/IBM-Cloud/bluemix-go/api/container/containerv1"
	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccIBMContainerALB_Basic(t *testing.T) {
	clusterName := fmt.Sprintf("terraform_%d", acctest.RandIntRange(10, 100))

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIBMContainerALBDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckIBMContainerALBBasic(clusterName, true),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_container_alb.alb", "enable", "true"),
				),
			},
			{
				Config: testAccCheckIBMContainerALBBasic(clusterName, false),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ibm_container_alb.alb", "enable", "false"),
				),
			},
		},
	})
}

func testAccCheckIBMContainerALBDestroy(s *terraform.State) error {
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ibm_container_alb" {
			continue
		}

		albID := rs.Primary.ID
		targetEnv := v1.ClusterTargetHeader{
			Region: "us-south",
		}

		csClient, err := testAccProvider.Meta().(ClientSession).ContainerAPI()
		if err != nil {
			return err
		}
		albAPI := csClient.Albs()
		_, err = albAPI.GetALB(albID, targetEnv)

		if err == nil {
			return fmt.Errorf("Instance still exists: %s", rs.Primary.ID)
		} else if !strings.Contains(err.Error(), "404") {
			return fmt.Errorf("Error checking if instance (%s) has been destroyed: %s", rs.Primary.ID, err)
		}
	}
	return nil
}

func testAccCheckIBMContainerALBBasic(clusterName string, enable bool) string {
	return fmt.Sprintf(`
resource "ibm_container_cluster" "testacc_cluster" {
  name       = "%s"
  datacenter = "%s"
  default_pool_size = 1
  machine_type    = "%s"
  hardware        = "shared"
  public_vlan_id  = "%s"
  private_vlan_id = "%s"
  timeouts {
    create = "720m"
	update = "720m"
  }
}

resource "ibm_container_alb" "alb" {
  alb_id = ibm_container_cluster.testacc_cluster.albs[0].id
  enable = "%t"
}`, clusterName, datacenter, machineType, publicVlanID, privateVlanID, enable)
}
