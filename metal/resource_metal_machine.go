package metal

import (
	"time"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	metalgo "github.com/metal-pod/metal-go"
)

func resourceMetalMachine() *schema.Resource {
	return &schema.Resource{
		Create: resourceMetalMachineCreate,
		Read:   resourceMetalMachineRead,
		Delete: resourceMetalMachineDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"hostname": {
				Type:     schema.TypeString,
				Required: true,
			},

			"name": {
				Type:     schema.TypeString,
				Required: true,
			},

			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},

			"image": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"partition": {
				Type:     schema.TypeString,
				Required: true,
			},
			"size": {
				Type:     schema.TypeString,
				Required: true,
			},

			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"metal_password": {
				Type:      schema.TypeString,
				Computed:  true,
				Sensitive: true,
			},

			"access_public_ipv4": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"access_private_ipv4": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"network": {
				Type:     schema.TypeList,
				Computed: true,
			},

			"created": {
				Type:     schema.TypeString,
				Computed: true,
			},

			"user_data": {
				Type:      schema.TypeString,
				Optional:  true,
				ForceNew:  true,
				Sensitive: true,
			},

			"tags": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
			"ssh_key_ids": {
				Type:     schema.TypeList,
				Computed: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func resourceMetalMachineCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*metalgo.Driver)

	createRequest := &metalgo.MachineCreateRequest{
		Hostname:    d.Get("hostname").(string),
		Name:        d.Get("name").(string),
		Partition:   d.Get("partition").(string),
		Description: d.Get("description").(string),
		UserData:    d.Get("userdata").(string),
		Size:        d.Get("size").(string),
		Project:     d.Get("project").(string),
		Tenant:      d.Get("tenant").(string),
		Image:       d.Get("image").(string),
	}

	tags := d.Get("tags.#").(int)
	if tags > 0 {
		createRequest.Tags = convertStringArr(d.Get("tags").([]interface{}))
	}
	sshPublicKeys := d.Get("ssh_key_ids.#").(int)
	if sshPublicKeys > 0 {
		createRequest.SSHPublicKeys = convertStringArr(d.Get("ssh_key_ids").([]interface{}))
	}

	newMachine, err := client.MachineCreate(createRequest)
	if err != nil {
		return err
	}

	d.SetId(*newMachine.Machine.ID)

	// Wait for the machine so we can get the networking attributes that show up after a while.
	_, err = waitForMachineAttribute(d, "Phone Home", "state", meta)
	if err != nil {
		return err
	}

	return resourceMetalMachineRead(d, meta)
}

func resourceMetalMachineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*metalgo.Driver)

	machine, err := client.MachineGet(d.Id())
	if err != nil {
		return err
	}
	m := machine.Machine
	alloc := m.Allocation
	d.Set("hostname", alloc.Hostname)
	d.Set("name", alloc.Name)
	d.Set("partition", m.Partition)
	d.Set("description", alloc.Description)
	d.Set("userdata", alloc.UserData)
	d.Set("size", m.Size)
	d.Set("project", alloc.Project)
	d.Set("tenant", alloc.Tenant)
	d.Set("image", alloc.Image)

	// FIXME tags are not populated
	// d.Set("tags", alloc.Tags)
	d.Set("ssh_key_ids", alloc.SSHPubKeys)
	d.Set("metal_password", alloc.ConsolePassword)

	// FIXME wrong
	d.Set("access_public_ipv4", alloc.Networks[0].Ips)
	d.Set("access_private_ipv4", alloc.Networks[0].Ips)
	return nil
}

func resourceMetalMachineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*metalgo.Driver)

	if _, err := client.MachineDelete(d.Id()); err != nil {
		return err
	}

	return nil
}

func waitForMachineAttribute(d *schema.ResourceData, target string, attribute string, meta interface{}) (interface{}, error) {
	stateConf := &resource.StateChangeConf{
		Target:     []string{target},
		Refresh:    newMachineStateRefreshFunc(d, attribute, meta),
		Timeout:    60 * time.Minute,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}
	return stateConf.WaitForState()
}

func newMachineStateRefreshFunc(d *schema.ResourceData, attribute string, meta interface{}) resource.StateRefreshFunc {
	client := meta.(*metalgo.Driver)

	return func() (interface{}, string, error) {
		if err := resourceMetalMachineRead(d, meta); err != nil {
			return nil, "", err
		}

		if attr, ok := d.GetOk(attribute); ok {
			machine, err := client.MachineGet(d.Id())
			if err != nil {
				return nil, "", err
			}
			return &machine, attr.(string), nil
		}

		return nil, "", nil
	}
}
