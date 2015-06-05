package coreos

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
)

type (
	ami struct {
		Name string `json:"name"`
		PV   string `json:"pv"`
		HVM  string `json:"hvm"`
	}

	amiInfo struct {
		AMIs []ami `json:"amis"`
	}
)

func resourceCoreOSAMI() *schema.Resource {
	return &schema.Resource{
		Create: Create,
		Delete: Delete,
		Exists: Exists,
		Read:   Read,

		Schema: map[string]*schema.Schema{
			"region": &schema.Schema{
				Type:        schema.TypeString,
				Description: "AWS region",
				Default:     "us-west-2",
				Optional:    true,
				ForceNew:    true,
			},
			"channel": &schema.Schema{
				Type:        schema.TypeString,
				Description: "CoreOS update channel",
				Default:     "stable",
				Optional:    true,
				ForceNew:    true,
			},
			"type": &schema.Schema{
				Type:        schema.TypeString,
				Description: "virtualization type",
				Default:     "pv",
				Optional:    true,
				ForceNew:    true,
			},
			"ami": &schema.Schema{
				Type:        schema.TypeString,
				Computed:    true,
				Description: "ami",
			},
		},
	}
}

func Create(d *schema.ResourceData, meta interface{}) error {
	log.Println("[INFO] calling create")
	ami, err := getAMI(d)
	if err != nil {
		return err
	}
	d.Set("ami", ami)
	d.SetId(getID(d))
	return nil
}

func Delete(d *schema.ResourceData, meta interface{}) error {
	log.Println("[INFO] calling delete")
	d.SetId("")
	return nil
}

func Exists(d *schema.ResourceData, meta interface{}) (bool, error) {
	log.Println("[INFO] calling exists")
	return getID(d) == d.Id(), nil
}

func Read(d *schema.ResourceData, meta interface{}) error {
	log.Println("[INFO] calling read")
	ami, err := getAMI(d)
	if err != nil {
		return err
	}
	d.Set("ami", ami)
	d.SetId(getID(d))
	return nil
}

func getAMI(d *schema.ResourceData) (string, error) {
	url := fmt.Sprintf("http://%s.release.core-os.net/amd64-usr/current/coreos_production_ami_all.json", d.Get("channel").(string))
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var data amiInfo
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", err
	}

	t := d.Get("type").(string)
	r := d.Get("region").(string)

	for _, a := range data.AMIs {
		if a.Name == r {
			switch t {
			case "pv":
				return a.PV, nil
			case "hvm":
				return a.HVM, nil
			default:
				return "", fmt.Errorf("invalid type: %s", t)
			}
		}
	}
	return "", fmt.Errorf("no ami found")
}

func getID(d *schema.ResourceData) string {
	channel := d.Get("channel").(string)
	r := d.Get("region").(string)
	t := d.Get("type").(string)

	return strings.Join([]string{channel, r, t}, ":")
}
