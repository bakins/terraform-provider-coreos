package google

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/googleapi"
)

func resourceComputeForwardingRule() *schema.Resource {
	return &schema.Resource{
		Create: resourceComputeForwardingRuleCreate,
		Read:   resourceComputeForwardingRuleRead,
		Delete: resourceComputeForwardingRuleDelete,
		Update: resourceComputeForwardingRuleUpdate,

		Schema: map[string]*schema.Schema{
			"ip_address": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"ip_protocol": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
				Computed: true,
			},

			"description": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"name": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},

			"port_range": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},

			"self_link": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},

			"target": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: false,
			},
		},
	}
}

func resourceComputeForwardingRuleCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	frule := &compute.ForwardingRule{
		IPAddress:   d.Get("ip_address").(string),
		IPProtocol:  d.Get("ip_protocol").(string),
		Description: d.Get("description").(string),
		Name:        d.Get("name").(string),
		PortRange:   d.Get("port_range").(string),
		Target:      d.Get("target").(string),
	}

	log.Printf("[DEBUG] ForwardingRule insert request: %#v", frule)
	op, err := config.clientCompute.ForwardingRules.Insert(
		config.Project, config.Region, frule).Do()
	if err != nil {
		return fmt.Errorf("Error creating ForwardingRule: %s", err)
	}

	// It probably maybe worked, so store the ID now
	d.SetId(frule.Name)

	// Wait for the operation to complete
	w := &OperationWaiter{
		Service: config.clientCompute,
		Op:      op,
		Region:  config.Region,
		Project: config.Project,
		Type:    OperationWaitRegion,
	}
	state := w.Conf()
	state.Timeout = 2 * time.Minute
	state.MinTimeout = 1 * time.Second
	opRaw, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for ForwardingRule to create: %s", err)
	}
	op = opRaw.(*compute.Operation)
	if op.Error != nil {
		// The resource didn't actually create
		d.SetId("")

		// Return the error
		return OperationError(*op.Error)
	}

	return resourceComputeForwardingRuleRead(d, meta)
}

func resourceComputeForwardingRuleUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	d.Partial(true)

	if d.HasChange("target") {
		target_name := d.Get("target").(string)
		target_ref := &compute.TargetReference{Target: target_name}
		op, err := config.clientCompute.ForwardingRules.SetTarget(
			config.Project, config.Region, d.Id(), target_ref).Do()
		if err != nil {
			return fmt.Errorf("Error updating target: %s", err)
		}

		// Wait for the operation to complete
		w := &OperationWaiter{
			Service: config.clientCompute,
			Op:      op,
			Region:  config.Region,
			Project: config.Project,
			Type:    OperationWaitRegion,
		}
		state := w.Conf()
		state.Timeout = 2 * time.Minute
		state.MinTimeout = 1 * time.Second
		opRaw, err := state.WaitForState()
		if err != nil {
			return fmt.Errorf("Error waiting for ForwardingRule to update target: %s", err)
		}
		op = opRaw.(*compute.Operation)
		if op.Error != nil {
			// The resource didn't actually create
			d.SetId("")

			// Return the error
			return OperationError(*op.Error)
		}
		d.SetPartial("target")
	}

	d.Partial(false)

	return resourceComputeForwardingRuleRead(d, meta)
}

func resourceComputeForwardingRuleRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	frule, err := config.clientCompute.ForwardingRules.Get(
		config.Project, config.Region, d.Id()).Do()
	if err != nil {
		if gerr, ok := err.(*googleapi.Error); ok && gerr.Code == 404 {
			// The resource doesn't exist anymore
			d.SetId("")

			return nil
		}

		return fmt.Errorf("Error reading ForwardingRule: %s", err)
	}

	d.Set("ip_address", frule.IPAddress)
	d.Set("ip_protocol", frule.IPProtocol)
	d.Set("self_link", frule.SelfLink)

	return nil
}

func resourceComputeForwardingRuleDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	// Delete the ForwardingRule
	log.Printf("[DEBUG] ForwardingRule delete request")
	op, err := config.clientCompute.ForwardingRules.Delete(
		config.Project, config.Region, d.Id()).Do()
	if err != nil {
		return fmt.Errorf("Error deleting ForwardingRule: %s", err)
	}

	// Wait for the operation to complete
	w := &OperationWaiter{
		Service: config.clientCompute,
		Op:      op,
		Region:  config.Region,
		Project: config.Project,
		Type:    OperationWaitRegion,
	}
	state := w.Conf()
	state.Timeout = 2 * time.Minute
	state.MinTimeout = 1 * time.Second
	opRaw, err := state.WaitForState()
	if err != nil {
		return fmt.Errorf("Error waiting for ForwardingRule to delete: %s", err)
	}
	op = opRaw.(*compute.Operation)
	if op.Error != nil {
		// Return the error
		return OperationError(*op.Error)
	}

	d.SetId("")
	return nil
}
