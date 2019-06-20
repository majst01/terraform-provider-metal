package metal

import (
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	metalgo "github.com/metal-pod/metal-go"
)

func dataSourceOperatingSystem() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMetalOperatingSystemRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"distro": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"slug": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMetalOperatingSystemRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*metalgo.Driver)

	slug, slugOK := d.GetOk("slug")

	if !slugOK {
		return fmt.Errorf("slug must be assigned")
	}

	log.Println("[DEBUG] ******")
	log.Println("[DEBUG] params", slug)
	log.Println("[DEBUG] ******")

	image, err := client.ImageGet(slug.(string))
	if err != nil {
		return err
	}

	d.Set("name", image.Image.Name)
	d.Set("slug", image.Image.ID)
	d.SetId(*image.Image.ID)
	return nil
}
