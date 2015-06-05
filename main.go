package main

import (
	"github.com/bakins/terraform-provider-coreos/coreos"
	"github.com/hashicorp/terraform/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: coreos.Provider,
	})
}
