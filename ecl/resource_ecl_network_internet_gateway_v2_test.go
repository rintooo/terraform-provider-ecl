package ecl

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/nttcom/terraform-provider-ecl/ecl/testhelper/mock"

	"github.com/nttcom/eclcloud/ecl/network/v2/internet_gateways"
)

func TestAccNetworkV2InternetGatewayBasic(t *testing.T) {
	var internet_gateway internet_gateways.InternetGateway

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckInternetGateway(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkV2InternetGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkV2InternetGatewayBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkV2InternetGatewayExists("ecl_network_internet_gateway_v2.internet_gateway_1", &internet_gateway),
				),
			},
			resource.TestStep{
				Config: testAccNetworkV2InternetGatewayUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ecl_network_internet_gateway_v2.internet_gateway_1", "name", stringMaxLength),
					resource.TestCheckResourceAttr(
						"ecl_network_internet_gateway_v2.internet_gateway_1", "description", stringMaxLength),
				),
			},
			resource.TestStep{
				Config: testAccNetworkV2InternetGatewayUpdate2,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ecl_network_internet_gateway_v2.internet_gateway_1", "name", ""),
					resource.TestCheckResourceAttr(
						"ecl_network_internet_gateway_v2.internet_gateway_1", "description", ""),
				),
			},
		},
	})
}

func TestMockedAccNetworkV2InternetGatewayBasic(t *testing.T) {
	var internet_gateway internet_gateways.InternetGateway

	mc := mock.NewMockController()
	defer mc.TerminateMockControllerSafety()

	postKeystone := fmt.Sprintf(fakeKeystonePostTmpl, mc.Endpoint())
	mc.Register(t, "keystone", "/v3/auth/tokens", postKeystone)
	mc.Register(t, "internet_gateway", "/v2.0/internet_gateways", testMockNetworkV2InternetGatewayPost)
	mc.Register(t, "internet_gateway", "/v2.0/internet_gateways/", testMockNetworkV2InternetGatewayGetBasic)
	mc.Register(t, "internet_gateway", "/v2.0/internet_gateways/", testMockNetworkV2InternetGatewayGetPendingCreate)
	mc.Register(t, "internet_gateway", "/v2.0/internet_gateways/", testMockNetworkV2InternetGatewayGetPendingUpdate)
	mc.Register(t, "internet_gateway", "/v2.0/internet_gateways/", testMockNetworkV2InternetGatewayGetPendingDelete)
	mc.Register(t, "internet_gateway", "/v2.0/internet_gateways/", testMockNetworkV2InternetGatewayGetUpdated)
	mc.Register(t, "internet_gateway", "/v2.0/internet_gateways/", testMockNetworkV2InternetGatewayGetDeleted)
	mc.Register(t, "internet_gateway", "/v2.0/internet_gateways/", testMockNetworkV2InternetGatewayPut)
	mc.Register(t, "internet_gateway", "/v2.0/internet_gateways/", testMockNetworkV2InternetGatewayDelete)

	mc.StartServer(t)

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheckInternetGateway(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckNetworkV2InternetGatewayDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccNetworkV2InternetGatewayBasic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckNetworkV2InternetGatewayExists("ecl_network_internet_gateway_v2.internet_gateway_1", &internet_gateway),
				),
			},
			resource.TestStep{
				Config: testAccNetworkV2InternetGatewayUpdate,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"ecl_network_internet_gateway_v2.internet_gateway_1", "description", "test_internet_gateway2"),
				),
			},
		},
	})
}

func testAccCheckNetworkV2InternetGatewayDestroy(s *terraform.State) error {
	config := testAccProvider.Meta().(*Config)
	networkClient, err := config.networkV2Client(OS_REGION_NAME)
	if err != nil {
		return fmt.Errorf("Error creating ECL network client: %s", err)
	}

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "ecl_network_internet_gateway_v2" {
			continue
		}

		_, err := internet_gateways.Get(networkClient, rs.Primary.ID).Extract()
		if err == nil {
			return fmt.Errorf("Internet gateway still exists")
		}
	}

	return nil
}

func testAccCheckNetworkV2InternetGatewayExists(n string, internet_gateway *internet_gateways.InternetGateway) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]
		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No ID is set")
		}

		config := testAccProvider.Meta().(*Config)
		networkClient, err := config.networkV2Client(OS_REGION_NAME)
		if err != nil {
			return fmt.Errorf("Error creating ECL network client: %s", err)
		}

		found, err := internet_gateways.Get(networkClient, rs.Primary.ID).Extract()
		if err != nil {
			return err
		}

		if found.ID != rs.Primary.ID {
			return fmt.Errorf("Internet gateway not found")
		}

		*internet_gateway = *found

		return nil
	}
}

var testAccNetworkV2InternetGatewayBasic = fmt.Sprintf(`
resource "ecl_network_internet_gateway_v2" "internet_gateway_1" {
    name = "Terraform_Test_Internet_Gateway_01"
    description = "test_internet_gateway"
    internet_service_id = "%s"
    qos_option_id = "%s"
}
`,
	OS_INTERNET_SERVICE_ID,
	OS_QOS_OPTION_ID_10M)

var testAccNetworkV2InternetGatewayUpdate = fmt.Sprintf(`
resource "ecl_network_internet_gateway_v2" "internet_gateway_1" {
    name = "%s",
    description = "%s",
    internet_service_id = "%s"
    qos_option_id = "%s"
}
`,
	stringMaxLength,
	stringMaxLength,
	OS_INTERNET_SERVICE_ID,
	OS_QOS_OPTION_ID_100M)

var testAccNetworkV2InternetGatewayUpdate2 = fmt.Sprintf(`
resource "ecl_network_internet_gateway_v2" "internet_gateway_1" {
    name = "",
    description = "",
    internet_service_id = "%s"
    qos_option_id = "%s"
}
`,
	OS_INTERNET_SERVICE_ID,
	OS_QOS_OPTION_ID_10M)

var testMockNetworkV2InternetGatewayPost = `
request:
    method: POST
response:
    code: 201
    body: >
        {
            "internet_gateway": {
                "description": "test_internet_gateway",
                "id": "3e71cf00-ddb5-4eb5-9ed0-ed4c481f6d61",
                "internet_service_id": "5536154d-9a00-4b11-81fb-b185c9111d90",
                "name": "Terraform_Test_Internet_Gateway_01",
                "qos_option_id": "e497bbc3-1127-4490-a51d-93582c40ab40",
                "status": "PENDING_CREATE",
                "tenant_id": "01234567890123456789abcdefabcdef"
            }
        }
newStatus: Created
`

var testMockNetworkV2InternetGatewayGetBasic = `
request:
    method: GET
response:
    code: 200
    body: >
        {
            "internet_gateway": {
                "description": "test_internet_gateway",
                "id": "3e71cf00-ddb5-4eb5-9ed0-ed4c481f6d61",
                "internet_service_id": "5536154d-9a00-4b11-81fb-b185c9111d90",
                "name": "Terraform_Test_Internet_Gateway_01",
                "qos_option_id": "e497bbc3-1127-4490-a51d-93582c40ab40",
                "status": "ACTIVE",
                "tenant_id": "01234567890123456789abcdefabcdef"
            }
        }
expectedStatus:
    - Created
counter:
    min: 4
`

var testMockNetworkV2InternetGatewayGetPendingCreate = `
request:
    method: GET
response:
    code: 200
    body: >
        {
            "internet_gateway": {
                "description": "test_internet_gateway",
                "id": "3e71cf00-ddb5-4eb5-9ed0-ed4c481f6d61",
                "internet_service_id": "5536154d-9a00-4b11-81fb-b185c9111d90",
                "name": "Terraform_Test_Internet_Gateway_01",
                "qos_option_id": "e497bbc3-1127-4490-a51d-93582c40ab40",
                "status": "PENDING_CREATE",
                "tenant_id": "01234567890123456789abcdefabcdef"
            }
        }
expectedStatus:
    - Created
counter:
    max: 3
`

var testMockNetworkV2InternetGatewayGetUpdated = `
request:
    method: GET
response:
    code: 200
    body: >
        {
            "internet_gateway": {
                "description": "test_internet_gateway2",
                "id": "3e71cf00-ddb5-4eb5-9ed0-ed4c481f6d61",
                "internet_service_id": "5536154d-9a00-4b11-81fb-b185c9111d90",
                "name": "Terraform_Test_Internet_Gateway_01",
                "qos_option_id": "e497bbc3-1127-4490-a51d-93582c40ab40",
                "status": "ACTIVE",
                "tenant_id": "01234567890123456789abcdefabcdef"
            }
        }
expectedStatus:
    - Updated
counter:
    min: 4
`

var testMockNetworkV2InternetGatewayGetPendingUpdate = `
request:
    method: GET
response:
    code: 200
    body: >
        {
            "internet_gateway": {
                "description": "test_internet_gateway2",
                "id": "3e71cf00-ddb5-4eb5-9ed0-ed4c481f6d61",
                "internet_service_id": "5536154d-9a00-4b11-81fb-b185c9111d90",
                "name": "Terraform_Test_Internet_Gateway_01",
                "qos_option_id": "e497bbc3-1127-4490-a51d-93582c40ab40",
                "status": "PENDING_UPDATE",
                "tenant_id": "01234567890123456789abcdefabcdef"
            }
        }
expectedStatus:
    - Updated
counter:
    max: 3
`

var testMockNetworkV2InternetGatewayGetDeleted = `
request:
    method: GET
response:
    code: 404
expectedStatus:
    - Deleted
counter:
    min: 4
`

var testMockNetworkV2InternetGatewayGetPendingDelete = `
request:
    method: GET
response:
    code: 200
    body: >
        {
            "internet_gateway": {
                "description": "test_internet_gateway2",
                "id": "3e71cf00-ddb5-4eb5-9ed0-ed4c481f6d61",
                "internet_service_id": "5536154d-9a00-4b11-81fb-b185c9111d90",
                "name": "Terraform_Test_Internet_Gateway_01",
                "qos_option_id": "e497bbc3-1127-4490-a51d-93582c40ab40",
                "status": "PENDING_DELETE",
                "tenant_id": "01234567890123456789abcdefabcdef"
            }
        }
expectedStatus:
    - Deleted
counter:
    max: 3
`

var testMockNetworkV2InternetGatewayPut = `
request:
    method: PUT
response:
    code: 200
    body: >
        {
            "internet_gateway": {
                "description": "test_internet_gateway2",
                "id": "3e71cf00-ddb5-4eb5-9ed0-ed4c481f6d61",
                "internet_service_id": "5536154d-9a00-4b11-81fb-b185c9111d90",
                "name": "Terraform_Test_Internet_Gateway_01",
                "qos_option_id": "e497bbc3-1127-4490-a51d-93582c40ab40",
                "status": "PENDING_UPDATE",
                "tenant_id": "dcb2d589c0c646d0bad45c0cf9f90cf1"
            }
        }
expectedStatus:
    - Created
newStatus: Updated
`

var testMockNetworkV2InternetGatewayDelete = `
request:
    method: DELETE
response:
    code: 204
expectedStatus:
    - Created
    - Updated
newStatus: Deleted
`
