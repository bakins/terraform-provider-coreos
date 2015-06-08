package heroku

import (
	"fmt"
	"testing"

	"github.com/cyberdelia/heroku-go/v3"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccHerokuApp_Basic(t *testing.T) {
	var app heroku.App

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHerokuAppDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckHerokuAppConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHerokuAppExists("heroku_app.foobar", &app),
					testAccCheckHerokuAppAttributes(&app),
					resource.TestCheckResourceAttr(
						"heroku_app.foobar", "name", "terraform-test-app"),
					resource.TestCheckResourceAttr(
						"heroku_app.foobar", "config_vars.0.FOO", "bar"),
				),
			},
		},
	})
}

func TestAccHerokuApp_NameChange(t *testing.T) {
	var app heroku.App

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHerokuAppDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckHerokuAppConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHerokuAppExists("heroku_app.foobar", &app),
					testAccCheckHerokuAppAttributes(&app),
					resource.TestCheckResourceAttr(
						"heroku_app.foobar", "name", "terraform-test-app"),
					resource.TestCheckResourceAttr(
						"heroku_app.foobar", "config_vars.0.FOO", "bar"),
				),
			},
			resource.TestStep{
				Config: testAccCheckHerokuAppConfig_updated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHerokuAppExists("heroku_app.foobar", &app),
					testAccCheckHerokuAppAttributesUpdated(&app),
					resource.TestCheckResourceAttr(
						"heroku_app.foobar", "name", "terraform-test-renamed"),
					resource.TestCheckResourceAttr(
						"heroku_app.foobar", "config_vars.0.FOO", "bing"),
					resource.TestCheckResourceAttr(
						"heroku_app.foobar", "config_vars.0.BAZ", "bar"),
				),
			},
		},
	})
}

func TestAccHerokuApp_NukeVars(t *testing.T) {
	var app heroku.App

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckHerokuAppDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckHerokuAppConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHerokuAppExists("heroku_app.foobar", &app),
					testAccCheckHerokuAppAttributes(&app),
					resource.TestCheckResourceAttr(
						"heroku_app.foobar", "name", "terraform-test-app"),
					resource.TestCheckResourceAttr(
						"heroku_app.foobar", "config_vars.0.FOO", "bar"),
				),
			},
			resource.TestStep{
				Config: testAccCheckHerokuAppConfig_no_vars,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckHerokuAppExists("heroku_app.foobar", &app),
					testAccCheckHerokuAppAttributesNoVars(&app),
					resource.TestCheckResourceAttr(
						"heroku_app.foobar", "name", "terraform-test-app"),
					resource.TestCheckResourceAttr(
						"heroku_app.foobar", "config_vars.0.FOO", ""),
				),
			},
		},
	})
}

func testAccCheckHerokuAppDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*heroku.Service)

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "heroku_app" {
			continue
		}

		_, err := client.AppInfo(rs.Primary.ID)

		if err == nil {
			return fmt.Errorf("App still exists")
		}
	}

	return nil
}

func testAccCheckHerokuAppAttributes(app *heroku.App) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*heroku.Service)

		if app.Region.Name != "us" {
			return fmt.Errorf("Bad region: %s", app.Region.Name)
		}

		if app.Stack.Name != "cedar-14" {
			return fmt.Errorf("Bad stack: %s", app.Stack.Name)
		}

		if app.Name != "terraform-test-app" {
			return fmt.Errorf("Bad name: %s", app.Name)
		}

		vars, err := client.ConfigVarInfo(app.Name)
		if err != nil {
			return err
		}

		if vars["FOO"] != "bar" {
			return fmt.Errorf("Bad config vars: %v", vars)
		}

		return nil
	}
}

func testAccCheckHerokuAppAttributesUpdated(app *heroku.App) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*heroku.Service)

		if app.Name != "terraform-test-renamed" {
			return fmt.Errorf("Bad name: %s", app.Name)
		}

		vars, err := client.ConfigVarInfo(app.Name)
		if err != nil {
			return err
		}

		// Make sure we kept the old one
		if vars["FOO"] != "bing" {
			return fmt.Errorf("Bad config vars: %v", vars)
		}

		if vars["BAZ"] != "bar" {
			return fmt.Errorf("Bad config vars: %v", vars)
		}

		return nil

	}
}

func testAccCheckHerokuAppAttributesNoVars(app *heroku.App) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*heroku.Service)

		if app.Name != "terraform-test-app" {
			return fmt.Errorf("Bad name: %s", app.Name)
		}

		vars, err := client.ConfigVarInfo(app.Name)
		if err != nil {
			return err
		}

		if len(vars) != 0 {
			return fmt.Errorf("vars exist: %v", vars)
		}

		return nil
	}
}

func testAccCheckHerokuAppExists(n string, app *heroku.App) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("Not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("No App Name is set")
		}

		client := testAccProvider.Meta().(*heroku.Service)

		foundApp, err := client.AppInfo(rs.Primary.ID)

		if err != nil {
			return err
		}

		if foundApp.Name != rs.Primary.ID {
			return fmt.Errorf("App not found")
		}

		*app = *foundApp

		return nil
	}
}

const testAccCheckHerokuAppConfig_basic = `
resource "heroku_app" "foobar" {
	name = "terraform-test-app"
	region = "us"

	config_vars {
		FOO = "bar"
	}
}`

const testAccCheckHerokuAppConfig_updated = `
resource "heroku_app" "foobar" {
	name = "terraform-test-renamed"
	region = "us"

	config_vars {
		FOO = "bing"
		BAZ = "bar"
	}
}`

const testAccCheckHerokuAppConfig_no_vars = `
resource "heroku_app" "foobar" {
	name = "terraform-test-app"
	region = "us"
}`
