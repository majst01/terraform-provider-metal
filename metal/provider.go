package metal

import (
	"time"

	"github.com/hashicorp/terraform/helper/mutexkv"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
)

var metalMutexKV = mutexkv.NewMutexKV()

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"auth_token": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("METAL_AUTH_TOKEN", nil),
				Description: "The API auth key for API operations.",
			},
		},
		DataSourcesMap: map[string]*schema.Resource{
			"metal_operating_system": dataSourceOperatingSystem(),
		},

		ResourcesMap: map[string]*schema.Resource{
			"metal_machine": resourceMetalMachine(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	config := Config{
		hmac: d.Get("auth_token").(string),
	}
	return config.Client()
}

var resourceDefaultTimeouts = &schema.ResourceTimeout{
	Create:  schema.DefaultTimeout(60 * time.Minute),
	Update:  schema.DefaultTimeout(60 * time.Minute),
	Delete:  schema.DefaultTimeout(60 * time.Minute),
	Default: schema.DefaultTimeout(60 * time.Minute),
}
