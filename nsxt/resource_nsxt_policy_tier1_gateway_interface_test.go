/* Copyright © 2019 VMware, Inc. All Rights Reserved.
   SPDX-License-Identifier: MPL-2.0 */

package nsxt

import (
	"fmt"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	gm_locale_services "github.com/vmware/vsphere-automation-sdk-go/services/nsxt-gm/global_infra/tier_1s/locale_services"
	"github.com/vmware/vsphere-automation-sdk-go/services/nsxt/infra/tier_1s/locale_services"
)

var nsxtPolicyTier1GatewayName = "test"

func TestAccResourceNsxtPolicyTier1GatewayInterface_basic(t *testing.T) {
	name := "test-nsx-policy-tier1-interface-basic"
	updatedName := fmt.Sprintf("%s-update", name)
	mtu := "1500"
	updatedMtu := "1800"
	subnet := "1.1.12.2/24"
	updatedSubnet := "1.2.12.2/24"
	testResourceName := "nsxt_policy_tier1_gateway_interface.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier1InterfaceCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier1InterfaceTemplate(name, subnet, mtu),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1InterfaceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "mtu", mtu),
					resource.TestCheckResourceAttr(testResourceName, "subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.0", subnet),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "segment_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "locale_service_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyTier1InterfaceTemplate(updatedName, updatedSubnet, updatedMtu),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1InterfaceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "mtu", updatedMtu),
					resource.TestCheckResourceAttr(testResourceName, "subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.0", updatedSubnet),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "segment_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "locale_service_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyTier1InterfaceThinTemplate(name, subnet),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1InterfaceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttr(testResourceName, "mtu", "0"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.0", subnet),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "segment_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "locale_service_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTier1GatewayInterface_withID(t *testing.T) {
	name := "test-nsx-policy-tier1-interface-id"
	updatedName := fmt.Sprintf("%s-update", name)
	subnet := "1.1.12.2/24"
	// Update to 2 addresses
	ipv6Subnet := "4003::12/64"
	updatedSubnet := fmt.Sprintf("%s\",\"%s", subnet, ipv6Subnet)
	testResourceName := "nsxt_policy_tier1_gateway_interface.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier1InterfaceCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier1InterfaceTemplateWithID(name, subnet),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1InterfaceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "nsx_id", "test"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.0", subnet),
					resource.TestCheckResourceAttr(testResourceName, "urpf_mode", "NONE"),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "segment_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "locale_service_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyTier1InterfaceTemplateWithID(updatedName, updatedSubnet),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1InterfaceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "nsx_id", "test"),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.0", subnet),
					resource.TestCheckResourceAttr(testResourceName, "subnets.1", ipv6Subnet),
					resource.TestCheckResourceAttr(testResourceName, "urpf_mode", "NONE"),
					resource.TestCheckResourceAttrSet(testResourceName, "segment_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "locale_service_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTier1GatewayInterface_withSite(t *testing.T) {
	name := "test-nsx-policy-tier1-interface-site"
	updatedName := fmt.Sprintf("%s-update", name)
	mtu := "1500"
	updatedMtu := "1800"
	subnet := "1.1.12.2/24"
	updatedSubnet := "1.2.12.2/24"
	testResourceName := "nsxt_policy_tier1_gateway_interface.test"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccOnlyGlobalManager(t); testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier1InterfaceCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier1InterfaceTemplate(name, subnet, mtu),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1InterfaceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "mtu", mtu),
					resource.TestCheckResourceAttr(testResourceName, "subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.0", subnet),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "segment_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "site_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "locale_service_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyTier1InterfaceTemplate(updatedName, updatedSubnet, updatedMtu),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1InterfaceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "mtu", updatedMtu),
					resource.TestCheckResourceAttr(testResourceName, "subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.0", updatedSubnet),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "1"),
					resource.TestCheckResourceAttrSet(testResourceName, "segment_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "site_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "locale_service_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyTier1InterfaceThinTemplate(name, subnet),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1InterfaceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "description", ""),
					resource.TestCheckResourceAttr(testResourceName, "mtu", "0"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.0", subnet),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "segment_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "site_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "locale_service_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func TestAccResourceNsxtPolicyTier1GatewayInterface_withIPv6(t *testing.T) {
	name := "test-nsx-policy-tier1-interface-ipv6"
	updatedName := fmt.Sprintf("%s-update", name)
	subnet := "1.1.12.2/24"
	// Update to 2 addresses
	ipv6Subnet := "4003::12/64"
	updatedSubnet := fmt.Sprintf("%s\",\"%s", subnet, ipv6Subnet)
	testResourceName := "nsxt_policy_tier1_gateway_interface.test"
	profilePath := testAccAdjustPolicyInfraConfig("/infra/ipv6-ndra-profiles/default")
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t); testAccNSXVersion(t, "3.0.0") },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier1InterfaceCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier1InterfaceTemplateWithIPv6(name, subnet),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1InterfaceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", name),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.#", "1"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.0", subnet),
					resource.TestCheckResourceAttr(testResourceName, "urpf_mode", "NONE"),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_ndra_profile_path", profilePath),
					resource.TestCheckResourceAttr(testResourceName, "tag.#", "0"),
					resource.TestCheckResourceAttrSet(testResourceName, "segment_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "locale_service_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
			{
				Config: testAccNsxtPolicyTier1InterfaceTemplateWithIPv6(updatedName, updatedSubnet),
				Check: resource.ComposeTestCheckFunc(
					testAccNsxtPolicyTier1InterfaceExists(testResourceName),
					resource.TestCheckResourceAttr(testResourceName, "display_name", updatedName),
					resource.TestCheckResourceAttr(testResourceName, "description", "Acceptance Test"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.#", "2"),
					resource.TestCheckResourceAttr(testResourceName, "subnets.0", subnet),
					resource.TestCheckResourceAttr(testResourceName, "subnets.1", ipv6Subnet),
					resource.TestCheckResourceAttr(testResourceName, "urpf_mode", "NONE"),
					resource.TestCheckResourceAttr(testResourceName, "ipv6_ndra_profile_path", profilePath),
					resource.TestCheckResourceAttrSet(testResourceName, "segment_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "gateway_path"),
					resource.TestCheckResourceAttrSet(testResourceName, "path"),
					resource.TestCheckResourceAttrSet(testResourceName, "nsx_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "locale_service_id"),
					resource.TestCheckResourceAttrSet(testResourceName, "revision"),
				),
			},
		},
	})
}

func testAccNSXPolicyTier1InterfaceImporterGetID(s *terraform.State) (string, error) {
	testResourceName := "nsxt_policy_tier1_gateway_interface.test"
	rs, ok := s.RootModule().Resources[testResourceName]
	if !ok {
		return "", fmt.Errorf("NSX Policy Tier1 Interface resource %s not found in resources", testResourceName)
	}
	resourceID := rs.Primary.ID
	localeServiceID := rs.Primary.Attributes["locale_service_id"]
	if resourceID == "" {
		return "", fmt.Errorf("NSX Policy Tier1 Interface resource ID not set in resources ")
	}
	gwPath := rs.Primary.Attributes["gateway_path"]
	if gwPath == "" {
		return "", fmt.Errorf("NSX Policy Interface Gateway Policy Path not set in resources ")
	}
	_, gwID := parseGatewayPolicyPath(gwPath)
	return fmt.Sprintf("%s/%s/%s", gwID, localeServiceID, resourceID), nil
}

func TestAccResourceNsxtPolicyTier1GatewayInterface_importBasic(t *testing.T) {
	name := "test-nsx-policy-tier1-interface-import"
	testResourceName := "nsxt_policy_tier1_gateway_interface.test"
	subnet := "1.1.12.2/24"

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		CheckDestroy: func(state *terraform.State) error {
			return testAccNsxtPolicyTier1InterfaceCheckDestroy(state, name)
		},
		Steps: []resource.TestStep{
			{
				Config: testAccNsxtPolicyTier1InterfaceThinTemplate(name, subnet),
			},
			{
				ResourceName:      testResourceName,
				ImportState:       true,
				ImportStateVerify: true,
				ImportStateIdFunc: testAccNSXPolicyTier1InterfaceImporterGetID,
			},
		},
	})
}

func testAccNsxtPolicyTier1InterfaceExists(resourceName string) resource.TestCheckFunc {
	return func(state *terraform.State) error {

		connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))

		rs, ok := state.RootModule().Resources[resourceName]
		if !ok {
			return fmt.Errorf("Policy Tier1 Interface resource %s not found in resources", resourceName)
		}

		resourceID := rs.Primary.ID
		localeServiceID := rs.Primary.Attributes["locale_service_id"]
		gwID := getPolicyIDFromPath(rs.Primary.Attributes["gateway_path"])

		if resourceID == "" {
			return fmt.Errorf("Policy Tier1 Interface resource ID not set in resources")
		}

		var err error
		if testAccIsGlobalManager() {
			nsxClient := gm_locale_services.NewDefaultInterfacesClient(connector)
			_, err = nsxClient.Get(gwID, localeServiceID, resourceID)
		} else {
			nsxClient := locale_services.NewDefaultInterfacesClient(connector)
			_, err = nsxClient.Get(gwID, localeServiceID, resourceID)
		}

		if err != nil {
			return fmt.Errorf("Error while retrieving policy Tier1 Interface ID %s. Error: %v", resourceID, err)
		}

		return nil
	}
}

func testAccNsxtPolicyTier1InterfaceCheckDestroy(state *terraform.State, displayName string) error {
	connector := getPolicyConnector(testAccProvider.Meta().(nsxtClients))
	var nsxClient interface{}
	var err error
	if testAccIsGlobalManager() {
		nsxClient = gm_locale_services.NewDefaultInterfacesClient(connector)

	} else {
		nsxClient = locale_services.NewDefaultInterfacesClient(connector)
	}
	for _, rs := range state.RootModule().Resources {

		if rs.Type != "nsxt_policy_tier1_gateway_interface" {
			continue
		}

		resourceID := rs.Primary.Attributes["id"]
		localeServicesID := rs.Primary.Attributes["locale_service_id"]
		gwID := getPolicyIDFromPath(rs.Primary.Attributes["gateway_path"])
		if testAccIsGlobalManager() {
			_, err = nsxClient.(*gm_locale_services.DefaultInterfacesClient).Get(gwID, localeServicesID, resourceID)
		} else {
			_, err = nsxClient.(*locale_services.DefaultInterfacesClient).Get(gwID, localeServicesID, resourceID)
		}
		if err == nil {
			return fmt.Errorf("Policy Tier1 Interface %s still exists", displayName)
		}
	}
	return nil
}

func testAccNsxtPolicyTier1InterfaceTemplate(name string, subnet string, mtu string) string {
	return testAccNsxtPolicyGatewayInterfaceDeps("11") + fmt.Sprintf(`
resource "nsxt_policy_tier1_gateway" "test" {
  display_name      = "%s"
  %s
}

resource "nsxt_policy_tier1_gateway_interface" "test" {
  display_name = "%s"
  description  = "Acceptance Test"
  mtu          = %s
  gateway_path = nsxt_policy_tier1_gateway.test.path
  segment_path = nsxt_policy_vlan_segment.test.path
  subnets      = ["%s"]
  %s

  tag {
    scope = "scope1"
    tag   = "tag1"
  }
}`, nsxtPolicyTier1GatewayName, testAccNsxtPolicyTier0EdgeClusterTemplate(), name, mtu, subnet, testAccNsxtPolicyTier0InterfaceSiteTemplate()) +
		testAccNextPolicyTier1InterfaceRealizationTemplate()
}

func testAccNsxtPolicyTier1InterfaceThinTemplate(name string, subnet string) string {
	return testAccNsxtPolicyGatewayInterfaceDeps("11") + fmt.Sprintf(`
resource "nsxt_policy_tier1_gateway" "test" {
  display_name      = "%s"
  %s
}

resource "nsxt_policy_tier1_gateway_interface" "test" {
  display_name = "%s"
  gateway_path = nsxt_policy_tier1_gateway.test.path
  segment_path = nsxt_policy_vlan_segment.test.path
  subnets      = ["%s"]
  %s
}`, nsxtPolicyTier1GatewayName, testAccNsxtPolicyTier0EdgeClusterTemplate(), name, subnet, testAccNsxtPolicyTier0InterfaceSiteTemplate()) +
		testAccNextPolicyTier1InterfaceRealizationTemplate()
}

func testAccNsxtPolicyTier1InterfaceTemplateWithID(name string, subnet string) string {
	return testAccNsxtPolicyGatewayInterfaceDeps("11") + fmt.Sprintf(`
resource "nsxt_policy_tier1_gateway" "test" {
  display_name      = "%s"
  %s
}

resource "nsxt_policy_tier1_gateway_interface" "test" {
  nsx_id                 = "test"
  display_name           = "%s"
  description            = "Acceptance Test"
  gateway_path           = nsxt_policy_tier1_gateway.test.path
  segment_path           = nsxt_policy_vlan_segment.test.path
  subnets                = ["%s"]
  urpf_mode              = "NONE"
  %s
}`, nsxtPolicyTier1GatewayName, testAccNsxtPolicyTier0EdgeClusterTemplate(), name, subnet, testAccNsxtPolicyTier0InterfaceSiteTemplate()) +
		testAccNextPolicyTier1InterfaceRealizationTemplate()
}

func testAccNextPolicyTier1InterfaceRealizationTemplate() string {
	return strings.Replace(testAccNsxtPolicyTier0InterfaceRealizationTemplate(), "tier0", "tier1", -1)
}

func testAccNsxtPolicyTier1InterfaceTemplateWithIPv6(name string, subnet string) string {
	return testAccNsxtPolicyGatewayInterfaceDeps("11") + fmt.Sprintf(`
data "nsxt_policy_ipv6_ndra_profile" "default" {
  display_name = "default"
}

resource "nsxt_policy_tier1_gateway" "test" {
  display_name      = "%s"
  %s
}

resource "nsxt_policy_tier1_gateway_interface" "test" {
  display_name           = "%s"
  description            = "Acceptance Test"
  gateway_path           = nsxt_policy_tier1_gateway.test.path
  segment_path           = nsxt_policy_vlan_segment.test.path
  subnets                = ["%s"]
  ipv6_ndra_profile_path = data.nsxt_policy_ipv6_ndra_profile.default.path
  urpf_mode              = "NONE"
  %s
}`, nsxtPolicyTier1GatewayName, testAccNsxtPolicyTier0EdgeClusterTemplate(), name, subnet, testAccNsxtPolicyTier0InterfaceSiteTemplate()) +
		testAccNextPolicyTier1InterfaceRealizationTemplate()
}
