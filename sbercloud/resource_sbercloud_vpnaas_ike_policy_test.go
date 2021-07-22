package sbercloud

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
	"github.com/huaweicloud/golangsdk"
	"github.com/huaweicloud/golangsdk/openstack/networking/v2/extensions/vpnaas/ikepolicies"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
)

func TestAccVpnIKEPolicyV2_basic(t *testing.T) {
	var policy ikepolicies.Policy
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIKEPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyV2_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists(
						"sbercloud_vpnaas_ike_policy.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("sbercloud_vpnaas_ike_policy.policy_1", "name", &policy.Name),
					resource.TestCheckResourceAttrPtr("sbercloud_vpnaas_ike_policy.policy_1", "description", &policy.Description),
					resource.TestCheckResourceAttrPtr("sbercloud_vpnaas_ike_policy.policy_1", "tenant_id", &policy.TenantID),
				),
			},
			{
				Config: testAccIKEPolicyV2_Update,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists(
						"sbercloud_vpnaas_ike_policy.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("sbercloud_vpnaas_ike_policy.policy_1", "name", &policy.Name),
				),
			},
		},
	})
}

func TestAccVpnIKEPolicyV2_withLifetime(t *testing.T) {
	var policy ikepolicies.Policy
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckIKEPolicyV2Destroy,
		Steps: []resource.TestStep{
			{
				Config: testAccIKEPolicyV2_withLifetime,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists(
						"sbercloud_vpnaas_ike_policy.policy_1", &policy),
					resource.TestCheckResourceAttrPtr("sbercloud_vpnaas_ike_policy.policy_1", "auth_algorithm", &policy.AuthAlgorithm),
					resource.TestCheckResourceAttrPtr("sbercloud_vpnaas_ike_policy.policy_1", "pfs", &policy.PFS),
				),
			},
			{
				Config: testAccIKEPolicyV2_withLifetimeUpdate,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckIKEPolicyV2Exists(
						"sbercloud_vpnaas_ike_policy.policy_1", &policy),
				),
			},
		},
	})
}

func testAccCheckIKEPolicyV2Destroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*config.Config)
	networkingClient, err := config.NetworkingV2Client(SBC_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating SberCloud networking client: %s", err)
	}
	for _, rs := range s.RootModule().Resources {
		if rs.Type != "sbercloud_vpnaas_ike_policy" {
			continue
		}
		_, err = ikepolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("IKE policy (%s) still exists.", rs.Primary.ID)
		}
		if _, ok := err.(golangsdk.ErrDefault404); !ok {
			return err
		}
	}
	return nil
}

func testAccCheckIKEPolicyV2Exists(n string, policy *ikepolicies.Policy) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*config.Config)
		networkingClient, err := config.NetworkingV2Client(SBC_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating SberCloud networking client: %s", err)
		}

		found, err := ikepolicies.Get(networkingClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}
		*policy = *found

		return nil
	}
}

const testAccIKEPolicyV2_basic = `
resource "sbercloud_vpnaas_ike_policy" "policy_1" {
}
`

const testAccIKEPolicyV2_Update = `
resource "sbercloud_vpnaas_ike_policy" "policy_1" {
	name = "updatedname"
}
`

const testAccIKEPolicyV2_withLifetime = `
resource "sbercloud_vpnaas_ike_policy" "policy_1" {
	auth_algorithm = "sha2-256"
	pfs = "group14"
	lifetime {
		units = "seconds"
		value = 1200
	}
}
`

const testAccIKEPolicyV2_withLifetimeUpdate = `
resource "sbercloud_vpnaas_ike_policy" "policy_1" {
	auth_algorithm = "sha2-256"
	pfs = "group14"
	lifetime {
		units = "seconds"
		value = 1400
	}
}
`
