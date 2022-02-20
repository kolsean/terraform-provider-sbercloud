package cbr

import (
	"fmt"
	"testing"

	"github.com/chnsz/golangsdk/openstack/cbr/v3/vaults"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/config"
	"github.com/huaweicloud/terraform-provider-huaweicloud/huaweicloud/services/cbr"
	"github.com/sbercloud-terraform/terraform-provider-sbercloud/sbercloud/acceptance"
)

func getVaultResourceFunc(conf *config.Config, state *terraform.ResourceState) (interface{}, error) {
	c, err := conf.CbrV3Client(acceptance.SBC_REGION_NAME)
	if err != nil {
		return nil, fmt.Errorf("error creating SberCloud CBR client: %s", err)
	}
	return vaults.Get(c, state.Primary.ID).Extract()
}

func TestAccCBRV3Vault_BasicServer(t *testing.T) {
	var vault vaults.Vault
	randName := acceptance.RandomAccResourceName()
	resourceName := "sbercloud_cbr_vault.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vault,
		getVaultResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCBRV3Vault_serverBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "app_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "200"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.SBC_ENTERPRISE_PROJECT_ID),
				),
			},
			{
				Config: testAccCBRV3Vault_serverUpdate(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "app_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeServer),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "300"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.SBC_ENTERPRISE_PROJECT_ID),
					acceptance.TestCheckResourceAttrWithVariable(resourceName, "policy_id",
						"${sbercloud_cbr_policy.test.id}"),
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

func TestAccCBRV3Vault_BasicVolume(t *testing.T) {
	var vault vaults.Vault
	randName := acceptance.RandomAccResourceName()
	resourceName := "sbercloud_cbr_vault.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vault,
		getVaultResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCBRV3Vault_volumeBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeDisk),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "50"),
					resource.TestCheckResourceAttr(resourceName, "auto_expand", "false"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.SBC_ENTERPRISE_PROJECT_ID),
					resource.TestCheckResourceAttr(resourceName, "resources.0.includes.#", "2"),
				),
			},
			{
				Config: testAccCBRV3Vault_volumeUpdate(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeDisk),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "100"),
					resource.TestCheckResourceAttr(resourceName, "auto_expand", "true"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.SBC_ENTERPRISE_PROJECT_ID),
					resource.TestCheckResourceAttr(resourceName, "resources.0.includes.#", "2"),
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

func TestAccCBRV3Vault_BasicTurbo(t *testing.T) {
	var vault vaults.Vault
	randName := acceptance.RandomAccResourceName()
	resourceName := "sbercloud_cbr_vault.test"

	rc := acceptance.InitResourceCheck(
		resourceName,
		&vault,
		getVaultResourceFunc,
	)

	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			acceptance.TestAccPreCheck(t)
			acceptance.TestAccPreCheckEpsID(t)
		},
		ProviderFactories: acceptance.TestAccProviderFactories,
		CheckDestroy:      rc.CheckResourceDestroy(),
		Steps: []resource.TestStep{
			{
				Config: testAccCBRV3Vault_turboBasic(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "800"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.SBC_ENTERPRISE_PROJECT_ID),
				),
			},
			{
				Config: testAccCBRV3Vault_turboUpdate(randName),
				Check: resource.ComposeTestCheckFunc(
					rc.CheckResourceExists(),
					resource.TestCheckResourceAttr(resourceName, "name", randName+"-update"),
					resource.TestCheckResourceAttr(resourceName, "consistent_level", "crash_consistent"),
					resource.TestCheckResourceAttr(resourceName, "type", cbr.VaultTypeTurbo),
					resource.TestCheckResourceAttr(resourceName, "protection_type", "backup"),
					resource.TestCheckResourceAttr(resourceName, "size", "1000"),
					resource.TestCheckResourceAttrSet(resourceName, "policy_id"),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "enterprise_project_id", acceptance.SBC_ENTERPRISE_PROJECT_ID),
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

func testAccCBRV3Vault_policy(rName string) string {
	return fmt.Sprintf(`
resource "sbercloud_cbr_policy" "test" {
  name        = "%s"
  type        = "backup"
  time_period = 20

  backup_cycle {
    days            = "MO,TU"
    execution_times = ["06:00", "18:00"]
  }
}
`, rName)
}

func testAccCBRV3VaultBasicConfiguration(rName string) string {
	return fmt.Sprintf(`
variable "volume_configuration" {
  type = list(object({
    volume_type = string
    size        = number
    device_type = string
  }))
  default = [
    {volume_type = "SSD", size = 100, device_type = "VBD"},
    {volume_type = "SSD", size = 100, device_type = "SCSI"},
  ]
}

data "sbercloud_availability_zones" "test" {}

data "sbercloud_compute_flavors" "test" {
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  performance_type  = "normal"
  cpu_core_count    = 2
  memory_size       = 4
}

data "sbercloud_images_image" "test" {
  name        = "Ubuntu 18.04 server 64bit"
  most_recent = true
}

resource "sbercloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/20"
}

resource "sbercloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/24"
  vpc_id     = sbercloud_vpc.test.id
  gateway_ip = "192.168.0.1"
}

resource "sbercloud_networking_secgroup" "test" {
  name                 = "%s"
  delete_default_rules = true
}

resource "sbercloud_compute_keypair" "test" {
  name = "%s"
  lifecycle {
    ignore_changes = [
      public_key,
    ]
  }
}

resource "sbercloud_compute_instance" "test" {
  availability_zone = data.sbercloud_availability_zones.test.names[0]
  name              = "%s"
  image_id          = data.sbercloud_images_image.test.id
  flavor_id         = data.sbercloud_compute_flavors.test.ids[0]
  key_pair          = sbercloud_compute_keypair.test.name

  system_disk_type = "SSD"
  system_disk_size = 50

  security_group_ids = [
    sbercloud_networking_secgroup.test.id
  ]

  network {
    uuid = sbercloud_vpc_subnet.test.id
  }
}

resource "sbercloud_evs_volume" "test" {
  count = length(var.volume_configuration)

  availability_zone = data.sbercloud_availability_zones.test.names[0]
  volume_type       = var.volume_configuration[count.index].volume_type
  name              = "%s_${tostring(count.index)}"
  size              = var.volume_configuration[count.index].size
  device_type       = var.volume_configuration[count.index].device_type
}

resource "sbercloud_compute_volume_attach" "test" {
  count = length(sbercloud_evs_volume.test)

  instance_id = sbercloud_compute_instance.test.id
  volume_id   = sbercloud_evs_volume.test[count.index].id
}`, rName, rName, rName, rName, rName, rName)
}

func testAccCBRV3Vault_serverBasic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_cbr_vault" "test" {
  name                  = "%s"
  type                  = "server"
  consistent_level      = "app_consistent"
  protection_type       = "backup"
  size                  = 200
  enterprise_project_id = "%s"

  resources {
    server_id = sbercloud_compute_instance.test.id
  }
}
`, testAccCBRV3VaultBasicConfiguration(rName),
		rName, acceptance.SBC_ENTERPRISE_PROJECT_ID)
}

func testAccCBRV3Vault_serverUpdate(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "sbercloud_cbr_vault" "test" {
  name                  = "%s-update"
  type                  = "server"
  consistent_level      = "app_consistent"
  protection_type       = "backup"
  size                  = 300
  enterprise_project_id = "%s"
  policy_id             = sbercloud_cbr_policy.test.id

  resources {
    server_id = sbercloud_compute_instance.test.id
 }
}
`, testAccCBRV3VaultBasicConfiguration(rName), testAccCBRV3Vault_policy(rName),
		rName, acceptance.SBC_ENTERPRISE_PROJECT_ID)
}

func testAccCBRV3Vault_volumeBasic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_cbr_vault" "test" {
  name                  = "%s"
  consistent_level      = "crash_consistent"
  type                  = "disk"
  protection_type       = "backup"
  size                  = 50
  enterprise_project_id = "%s"

  resources {
    includes = sbercloud_compute_volume_attach.test[*].volume_id
  }
}
`, testAccCBRV3VaultBasicConfiguration(rName),
		rName, acceptance.SBC_ENTERPRISE_PROJECT_ID)
}

func testAccCBRV3Vault_volumeUpdate(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "sbercloud_cbr_vault" "test" {
  name                  = "%s-update"
  consistent_level      = "crash_consistent"
  type                  = "disk"
  protection_type       = "backup"
  size                  = 100
  auto_expand           = true
  enterprise_project_id = "%s"
  policy_id             = sbercloud_cbr_policy.test.id

  resources {
    includes = sbercloud_compute_volume_attach.test[*].volume_id
  }
}
`, testAccCBRV3VaultBasicConfiguration(rName),
		testAccCBRV3Vault_policy(rName), rName, acceptance.SBC_ENTERPRISE_PROJECT_ID)
}

//Vaults of type 'turbo'
func testAccCBRV3Vault_turboBase(rName string) string {
	return fmt.Sprintf(`
data "sbercloud_availability_zones" "test" {}

resource "sbercloud_vpc" "test" {
  name = "%s"
  cidr = "192.168.0.0/20"
}

resource "sbercloud_vpc_subnet" "test" {
  name       = "%s"
  cidr       = "192.168.0.0/22"
  gateway_ip = "192.168.0.1"
  vpc_id     = sbercloud_vpc.test.id
}

resource "sbercloud_networking_secgroup" "test" {
  name = "%s"
}

resource "sbercloud_sfs_turbo" "test1" {
  name              = "%s-1"
  size              = 500
  share_proto       = "NFS"
  vpc_id            = sbercloud_vpc.test.id
  subnet_id         = sbercloud_vpc_subnet.test.id
  security_group_id = sbercloud_networking_secgroup.test.id
  availability_zone = data.sbercloud_availability_zones.test.names[0]
}`, rName, rName, rName, rName)
}

func testAccCBRV3Vault_turboBasic(rName string) string {
	return fmt.Sprintf(`
%s

resource "sbercloud_cbr_vault" "test" {
  name                  = "%s"
  consistent_level      = "crash_consistent"
  type                  = "turbo"
  protection_type       = "backup"
  size                  = 800
  enterprise_project_id = "%s"

  resources {
    includes = [
      sbercloud_sfs_turbo.test1.id
    ]
  }
}
`, testAccCBRV3Vault_turboBase(rName), rName, acceptance.SBC_ENTERPRISE_PROJECT_ID)
}

func testAccCBRV3Vault_turboUpdate(rName string) string {
	return fmt.Sprintf(`
%s

%s

resource "sbercloud_sfs_turbo" "test2" {
  name              = "%s-2"
  size              = 500
  share_proto       = "NFS"
  vpc_id            = sbercloud_vpc.test.id
  subnet_id         = sbercloud_vpc_subnet.test.id
  security_group_id = sbercloud_networking_secgroup.test.id
  availability_zone = data.sbercloud_availability_zones.test.names[0]
}

resource "sbercloud_cbr_vault" "test" {
  name                  = "%s-update"
  consistent_level      = "crash_consistent"
  type                  = "turbo"
  protection_type       = "backup"
  size                  = 1000
  enterprise_project_id = "%s"
  policy_id             = sbercloud_cbr_policy.test.id

  resources {
    includes = [
      sbercloud_sfs_turbo.test1.id,
      sbercloud_sfs_turbo.test2.id
    ]
  }
}
`, testAccCBRV3Vault_turboBase(rName), testAccCBRV3Vault_policy(rName), rName, rName,
		acceptance.SBC_ENTERPRISE_PROJECT_ID)
}
