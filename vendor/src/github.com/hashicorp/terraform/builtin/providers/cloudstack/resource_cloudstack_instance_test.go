package cloudstack

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/xanzy/go-cloudstack/cloudstack"
)

func TestAccCloudStackInstance_basic(t *testing.T) {
	var instance cloudstack.VirtualMachine

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCloudStackInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackInstanceExists(
						"cloudstack_instance.foobar", &instance),
					testAccCheckCloudStackInstanceAttributes(&instance),
					resource.TestCheckResourceAttr(
						"cloudstack_instance.foobar",
						"user_data",
						"0cf3dcdc356ec8369494cb3991985ecd5296cdd5"),
					resource.TestCheckResourceAttr(
						"cloudstack_instance.foobar",
						"keypair",
						CLOUDSTACK_SSH_KEYPAIR),
				),
			},
		},
	})
}

func TestAccCloudStackInstance_update(t *testing.T) {
	var instance cloudstack.VirtualMachine

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCloudStackInstance_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackInstanceExists(
						"cloudstack_instance.foobar", &instance),
					testAccCheckCloudStackInstanceAttributes(&instance),
					resource.TestCheckResourceAttr(
						"cloudstack_instance.foobar",
						"user_data",
						"0cf3dcdc356ec8369494cb3991985ecd5296cdd5"),
					resource.TestCheckResourceAttr(
						"cloudstack_instance.foobar",
						"keypair",
						CLOUDSTACK_SSH_KEYPAIR),
				),
			},

			resource.TestStep{
				Config: testAccCloudStackInstance_renameAndResize,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackInstanceExists(
						"cloudstack_instance.foobar", &instance),
					testAccCheckCloudStackInstanceRenamedAndResized(&instance),
					resource.TestCheckResourceAttr(
						"cloudstack_instance.foobar", "display_name", "terraform-updated"),
					resource.TestCheckResourceAttr(
						"cloudstack_instance.foobar", "service_offering", CLOUDSTACK_SERVICE_OFFERING_2),
				),
			},
		},
	})
}

func TestAccCloudStackInstance_fixedIP(t *testing.T) {
	var instance cloudstack.VirtualMachine

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckCloudStackInstanceDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCloudStackInstance_fixedIP,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckCloudStackInstanceExists(
						"cloudstack_instance.foobar", &instance),
					resource.TestCheckResourceAttr(
						"cloudstack_instance.foobar", "ipaddress", CLOUDSTACK_NETWORK_1_IPADDRESS),
				),
			},
		},
	})
}

func testAccCheckCloudStackInstanceExists(
	n string, instance *cloudstack.VirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)
		vm, _, err := cs.VirtualMachine.GetVirtualMachineByID(rs.Primary.ID)

		if err != nil {
			return err
		}

		if vm.Id != rs.Primary.ID {
			return fmt.Errorf("Instance not found")
		}

		*instance = *vm

		return nil
	}
}

func testAccCheckCloudStackInstanceAttributes(
	instance *cloudstack.VirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if instance.Name != "terraform-test" {
			return fmt.Errorf("Bad name: %s", instance.Name)
		}

		if instance.Displayname != "terraform" {
			return fmt.Errorf("Bad display name: %s", instance.Displayname)
		}

		if instance.Serviceofferingname != CLOUDSTACK_SERVICE_OFFERING_1 {
			return fmt.Errorf("Bad service offering: %s", instance.Serviceofferingname)
		}

		if instance.Templatename != CLOUDSTACK_TEMPLATE {
			return fmt.Errorf("Bad template: %s", instance.Templatename)
		}

		if instance.Nic[0].Networkname != CLOUDSTACK_NETWORK_1 {
			return fmt.Errorf("Bad network: %s", instance.Nic[0].Networkname)
		}

		return nil
	}
}

func testAccCheckCloudStackInstanceRenamedAndResized(
	instance *cloudstack.VirtualMachine) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if instance.Displayname != "terraform-updated" {
			return fmt.Errorf("Bad display name: %s", instance.Displayname)
		}

		if instance.Serviceofferingname != CLOUDSTACK_SERVICE_OFFERING_2 {
			return fmt.Errorf("Bad service offering: %s", instance.Serviceofferingname)
		}

		return nil
	}
}

func testAccCheckCloudStackInstanceDestroy(s *terraform.State) error {
	cs := testAccProvider.Meta().(*cloudstack.CloudStackClient)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "cloudstack_instance" {
			continue
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No instance ID is set")
		}

		p := cs.VirtualMachine.NewDestroyVirtualMachineParams(rs.Primary.ID)
		_, err := cs.VirtualMachine.DestroyVirtualMachine(p)

		if err != nil {
			return fmt.Errorf(
				"Error deleting instance (%s): %s",
				rs.Primary.ID, err)
		}
	}

	return nil
}

var testAccCloudStackInstance_basic = fmt.Sprintf(`
resource "cloudstack_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform"
  service_offering= "%s"
  network = "%s"
  template = "%s"
  zone = "%s"
  keypair = "%s"
  user_data = "foobar\nfoo\nbar"
  expunge = true
}`,
	CLOUDSTACK_SERVICE_OFFERING_1,
	CLOUDSTACK_NETWORK_1,
	CLOUDSTACK_TEMPLATE,
	CLOUDSTACK_ZONE,
	CLOUDSTACK_SSH_KEYPAIR)

var testAccCloudStackInstance_renameAndResize = fmt.Sprintf(`
resource "cloudstack_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform-updated"
  service_offering= "%s"
  network = "%s"
  template = "%s"
  zone = "%s"
  keypair = "%s"
  user_data = "foobar\nfoo\nbar"
  expunge = true
}`,
	CLOUDSTACK_SERVICE_OFFERING_2,
	CLOUDSTACK_NETWORK_1,
	CLOUDSTACK_TEMPLATE,
	CLOUDSTACK_ZONE,
	CLOUDSTACK_SSH_KEYPAIR)

var testAccCloudStackInstance_fixedIP = fmt.Sprintf(`
resource "cloudstack_instance" "foobar" {
  name = "terraform-test"
  display_name = "terraform"
  service_offering= "%s"
  network = "%s"
  ipaddress = "%s"
  template = "%s"
  zone = "%s"
  expunge = true
}`,
	CLOUDSTACK_SERVICE_OFFERING_1,
	CLOUDSTACK_NETWORK_1,
	CLOUDSTACK_NETWORK_1_IPADDRESS,
	CLOUDSTACK_TEMPLATE,
	CLOUDSTACK_ZONE)
