package ecl

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/nttcom/eclcloud/ecl/objectstorage/v1/containers"
)

func TestAccObjectStorageV1Container_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckSwift(t)
		},
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckObjectStorageV1ContainerDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccObjectStorageV1Container_basic,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ecl_objectstorage_container_v1.container_1", "name", "container_1"),
					resource.TestCheckResourceAttr(
						"ecl_objectstorage_container_v1.container_1", "content_type", "application/json"),
				),
			},
			resource.TestStep{
				Config: testAccObjectStorageV1Container_update,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ecl_objectstorage_container_v1.container_1", "content_type", "text/plain"),
				),
			},
		},
	})
}

func testAccCheckObjectStorageV1ContainerDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	objectStorageClient, err := config.objectStorageV1Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating ECL object storage client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecl_objectstorage_container_v1" {
			continue
		}

		_, err := containers.Get(objectStorageClient, rs.Primary.ID, nil).Extract()
		if err == nil {
			return fmt.Errorf("Container still exists")
		}
	}

	return nil
}

const testAccObjectStorageV1Container_basic = `
resource "ecl_objectstorage_container_v1" "container_1" {
  name = "container_1"
  metadata {
    test = "true"
  }
  content_type = "application/json"
}
`

const testAccObjectStorageV1Container_update = `
resource "ecl_objectstorage_container_v1" "container_1" {
  name = "container_1"
  metadata {
    test = "true"
  }
  content_type = "text/plain"
}
`
