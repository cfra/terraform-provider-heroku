package heroku

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	heroku "github.com/heroku/heroku-go/v5"
)

func resourceHerokuPipeline() *schema.Resource {
	return &schema.Resource{
		Create: resourceHerokuPipelineCreate,
		Update: resourceHerokuPipelineUpdate,
		Read:   resourceHerokuPipelineRead,
		Delete: resourceHerokuPipelineDelete,

		Importer: &schema.ResourceImporter{
			State: resourceHerokuPipelineImport,
		},

		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceHerokuPipelineImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	client := meta.(*Config).Api

	p, err := client.PipelineInfo(context.TODO(), d.Id())
	if err != nil {
		return nil, err
	}

	d.Set("name", p.Name)

	return []*schema.ResourceData{d}, nil
}

func resourceHerokuPipelineCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).Api

	opts := heroku.PipelineCreateOpts{
		Name: d.Get("name").(string),
	}

	log.Printf("[DEBUG] Pipeline create configuration: %#v", opts)

	p, err := client.PipelineCreate(context.TODO(), opts)
	if err != nil {
		return fmt.Errorf("Error creating pipeline: %s", err)
	}

	d.SetId(p.ID)
	d.Set("name", p.Name)

	log.Printf("[INFO] Pipeline ID: %s", d.Id())

	return resourceHerokuPipelineUpdate(d, meta)
}

func resourceHerokuPipelineUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).Api

	if d.HasChange("name") {
		name := d.Get("name").(string)
		opts := heroku.PipelineUpdateOpts{
			Name: &name,
		}

		_, err := client.PipelineUpdate(context.TODO(), d.Id(), opts)
		if err != nil {
			return err
		}
	}

	return resourceHerokuPipelineRead(d, meta)
}

func resourceHerokuPipelineDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).Api

	log.Printf("[INFO] Deleting pipeline: %s", d.Id())

	_, err := client.PipelineDelete(context.TODO(), d.Id())
	if err != nil {
		return fmt.Errorf("Error deleting pipeline: %s", err)
	}

	return nil
}

func resourceHerokuPipelineRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*Config).Api

	p, err := client.PipelineInfo(context.TODO(), d.Id())
	if err != nil {
		return fmt.Errorf("Error retrieving pipeline: %s", err)
	}

	d.Set("name", p.Name)

	return nil
}
