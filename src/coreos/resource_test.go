package coreos

import "github.com/hashicorp/terraform/terraform"

var testProviders = map[string]terraform.ResourceProvider{
	"coreos": Provider(),
}

/*
func TestCoreOSAMI(t *testing.T) {
	r.Test(t, r.TestCase{
			Providers: testProviders,
			Steps: []r.TestStep{
				r.TestStep{
					Config: `
resource "template_file" "t0" {
	filename = "mock"
	vars = ` + tt.vars + `
}
output "rendered" {
    value = "${template_file.t0.rendered}"
}
`,
					Check: func(s *terraform.State) error {
						got := s.RootModule().Outputs["rendered"]
						if tt.want != got {
							return fmt.Errorf("template:\n%s\nvars:\n%s\ngot:\n%s\nwant:\n%s\n", tt.template, tt.vars, got, tt.want)
						}
						return nil
					},
				},
			},
		})
	}
}
*/
