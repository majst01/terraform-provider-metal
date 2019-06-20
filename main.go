package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/terraform-providers/terraform-provider-metal/metal"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: metal.Provider})
}
