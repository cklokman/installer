// ----------------------------------------------------------------------------
//
//     ***     AUTO GENERATED CODE    ***    AUTO GENERATED CODE     ***
//
// ----------------------------------------------------------------------------
//
//     This file is automatically generated by Magic Modules and manual
//     changes will be clobbered when the file is regenerated.
//
//     Please read more about how to change this file in
//     .github/CONTRIBUTING.md.
//
// ----------------------------------------------------------------------------

package google

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceSpannerInstance() *schema.Resource {
	return &schema.Resource{
		Create: resourceSpannerInstanceCreate,
		Read:   resourceSpannerInstanceRead,
		Update: resourceSpannerInstanceUpdate,
		Delete: resourceSpannerInstanceDelete,

		Importer: &schema.ResourceImporter{
			State: resourceSpannerInstanceImport,
		},

		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(4 * time.Minute),
			Update: schema.DefaultTimeout(4 * time.Minute),
			Delete: schema.DefaultTimeout(4 * time.Minute),
		},

		Schema: map[string]*schema.Schema{
			"config": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: compareSelfLinkOrResourceName,
				Description: `The name of the instance's configuration (similar but not
quite the same as a region) which defines defines the geographic placement and
replication of your databases in this instance. It determines where your data
is stored. Values are typically of the form 'regional-europe-west1' , 'us-central' etc.
In order to obtain a valid list please consult the
[Configuration section of the docs](https://cloud.google.com/spanner/docs/instances).`,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
				Description: `The descriptive name for this instance as it appears in UIs. Must be
unique per project and between 4 and 30 characters in length.`,
			},
			"name": {
				Type:         schema.TypeString,
				Computed:     true,
				Optional:     true,
				ForceNew:     true,
				ValidateFunc: validateRegexp(`^[a-z][-a-z0-9]*[a-z0-9]$`),
				Description: `A unique identifier for the instance, which cannot be changed after
the instance is created. The name must be between 6 and 30 characters
in length.


If not provided, a random string starting with 'tf-' will be selected.`,
			},
			"labels": {
				Type:     schema.TypeMap,
				Optional: true,
				Description: `An object containing a list of "key": value pairs.
Example: { "name": "wrench", "mass": "1.3kg", "count": "3" }.`,
				Elem: &schema.Schema{Type: schema.TypeString},
			},
			"num_nodes": {
				Type:        schema.TypeInt,
				Optional:    true,
				Description: `The number of nodes allocated to this instance.`,
				Default:     1,
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: `Instance status: 'CREATING' or 'READY'.`,
			},
			"project": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
		},
	}
}

func resourceSpannerInstanceCreate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	obj := make(map[string]interface{})
	nameProp, err := expandSpannerInstanceName(d.Get("name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("name"); !isEmptyValue(reflect.ValueOf(nameProp)) && (ok || !reflect.DeepEqual(v, nameProp)) {
		obj["name"] = nameProp
	}
	configProp, err := expandSpannerInstanceConfig(d.Get("config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("config"); !isEmptyValue(reflect.ValueOf(configProp)) && (ok || !reflect.DeepEqual(v, configProp)) {
		obj["config"] = configProp
	}
	displayNameProp, err := expandSpannerInstanceDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(displayNameProp)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	nodeCountProp, err := expandSpannerInstanceNumNodes(d.Get("num_nodes"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("num_nodes"); !isEmptyValue(reflect.ValueOf(nodeCountProp)) && (ok || !reflect.DeepEqual(v, nodeCountProp)) {
		obj["nodeCount"] = nodeCountProp
	}
	labelsProp, err := expandSpannerInstanceLabels(d.Get("labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(labelsProp)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	obj, err = resourceSpannerInstanceEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{SpannerBasePath}}projects/{{project}}/instances")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Creating new Instance: %#v", obj)
	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	res, err := sendRequestWithTimeout(config, "POST", project, url, obj, d.Timeout(schema.TimeoutCreate))
	if err != nil {
		return fmt.Errorf("Error creating Instance: %s", err)
	}

	// Store the ID now
	id, err := replaceVars(d, config, "{{project}}/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	// Use the resource in the operation response to populate
	// identity fields and d.Id() before read
	var opRes map[string]interface{}
	err = spannerOperationWaitTimeWithResponse(
		config, res, &opRes, project, "Creating Instance",
		d.Timeout(schema.TimeoutCreate))
	if err != nil {
		// The resource didn't actually create
		d.SetId("")
		return fmt.Errorf("Error waiting to create Instance: %s", err)
	}

	opRes, err = resourceSpannerInstanceDecoder(d, meta, opRes)
	if err != nil {
		return fmt.Errorf("Error decoding response from operation: %s", err)
	}
	if opRes == nil {
		return fmt.Errorf("Error decoding response from operation, could not find object")
	}

	if err := d.Set("name", flattenSpannerInstanceName(opRes["name"], d, config)); err != nil {
		return err
	}

	// This may have caused the ID to update - update it if so.
	id, err = replaceVars(d, config, "{{project}}/{{name}}")
	if err != nil {
		return fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	log.Printf("[DEBUG] Finished creating Instance %q: %#v", d.Id(), res)

	// This is useful if the resource in question doesn't have a perfectly consistent API
	// That is, the Operation for Create might return before the Get operation shows the
	// completed state of the resource.
	time.Sleep(5 * time.Second)

	return resourceSpannerInstanceRead(d, meta)
}

func resourceSpannerInstanceRead(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	url, err := replaceVars(d, config, "{{SpannerBasePath}}projects/{{project}}/instances/{{name}}")
	if err != nil {
		return err
	}

	project, err := getProject(d, config)
	if err != nil {
		return err
	}
	res, err := sendRequest(config, "GET", project, url, nil)
	if err != nil {
		return handleNotFoundError(err, d, fmt.Sprintf("SpannerInstance %q", d.Id()))
	}

	res, err = resourceSpannerInstanceDecoder(d, meta, res)
	if err != nil {
		return err
	}

	if res == nil {
		// Decoding the object has resulted in it being gone. It may be marked deleted
		log.Printf("[DEBUG] Removing SpannerInstance because it no longer exists.")
		d.SetId("")
		return nil
	}

	if err := d.Set("project", project); err != nil {
		return fmt.Errorf("Error reading Instance: %s", err)
	}

	if err := d.Set("name", flattenSpannerInstanceName(res["name"], d, config)); err != nil {
		return fmt.Errorf("Error reading Instance: %s", err)
	}
	if err := d.Set("config", flattenSpannerInstanceConfig(res["config"], d, config)); err != nil {
		return fmt.Errorf("Error reading Instance: %s", err)
	}
	if err := d.Set("display_name", flattenSpannerInstanceDisplayName(res["displayName"], d, config)); err != nil {
		return fmt.Errorf("Error reading Instance: %s", err)
	}
	if err := d.Set("num_nodes", flattenSpannerInstanceNumNodes(res["nodeCount"], d, config)); err != nil {
		return fmt.Errorf("Error reading Instance: %s", err)
	}
	if err := d.Set("labels", flattenSpannerInstanceLabels(res["labels"], d, config)); err != nil {
		return fmt.Errorf("Error reading Instance: %s", err)
	}
	if err := d.Set("state", flattenSpannerInstanceState(res["state"], d, config)); err != nil {
		return fmt.Errorf("Error reading Instance: %s", err)
	}

	return nil
}

func resourceSpannerInstanceUpdate(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	obj := make(map[string]interface{})
	configProp, err := expandSpannerInstanceConfig(d.Get("config"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("config"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, configProp)) {
		obj["config"] = configProp
	}
	displayNameProp, err := expandSpannerInstanceDisplayName(d.Get("display_name"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("display_name"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, displayNameProp)) {
		obj["displayName"] = displayNameProp
	}
	nodeCountProp, err := expandSpannerInstanceNumNodes(d.Get("num_nodes"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("num_nodes"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, nodeCountProp)) {
		obj["nodeCount"] = nodeCountProp
	}
	labelsProp, err := expandSpannerInstanceLabels(d.Get("labels"), d, config)
	if err != nil {
		return err
	} else if v, ok := d.GetOkExists("labels"); !isEmptyValue(reflect.ValueOf(v)) && (ok || !reflect.DeepEqual(v, labelsProp)) {
		obj["labels"] = labelsProp
	}

	obj, err = resourceSpannerInstanceUpdateEncoder(d, meta, obj)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{SpannerBasePath}}projects/{{project}}/instances/{{name}}")
	if err != nil {
		return err
	}

	log.Printf("[DEBUG] Updating Instance %q: %#v", d.Id(), obj)
	res, err := sendRequestWithTimeout(config, "PATCH", project, url, obj, d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return fmt.Errorf("Error updating Instance %q: %s", d.Id(), err)
	}

	err = spannerOperationWaitTime(
		config, res, project, "Updating Instance",
		d.Timeout(schema.TimeoutUpdate))

	if err != nil {
		return err
	}

	return resourceSpannerInstanceRead(d, meta)
}

func resourceSpannerInstanceDelete(d *schema.ResourceData, meta interface{}) error {
	config := meta.(*Config)

	project, err := getProject(d, config)
	if err != nil {
		return err
	}

	url, err := replaceVars(d, config, "{{SpannerBasePath}}projects/{{project}}/instances/{{name}}")
	if err != nil {
		return err
	}

	var obj map[string]interface{}
	log.Printf("[DEBUG] Deleting Instance %q", d.Id())

	res, err := sendRequestWithTimeout(config, "DELETE", project, url, obj, d.Timeout(schema.TimeoutDelete))
	if err != nil {
		return handleNotFoundError(err, d, "Instance")
	}

	log.Printf("[DEBUG] Finished deleting Instance %q: %#v", d.Id(), res)
	return nil
}

func resourceSpannerInstanceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	config := meta.(*Config)
	if err := parseImportId([]string{
		"projects/(?P<project>[^/]+)/instances/(?P<name>[^/]+)",
		"(?P<project>[^/]+)/(?P<name>[^/]+)",
		"(?P<name>[^/]+)",
	}, d, config); err != nil {
		return nil, err
	}

	// Replace import id for the resource id
	id, err := replaceVars(d, config, "{{project}}/{{name}}")
	if err != nil {
		return nil, fmt.Errorf("Error constructing id: %s", err)
	}
	d.SetId(id)

	return []*schema.ResourceData{d}, nil
}

func flattenSpannerInstanceName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}
	return NameFromSelfLinkStateFunc(v)
}

func flattenSpannerInstanceConfig(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	if v == nil {
		return v
	}
	return ConvertSelfLinkToV1(v.(string))
}

func flattenSpannerInstanceDisplayName(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenSpannerInstanceNumNodes(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	// Handles the string fixed64 format
	if strVal, ok := v.(string); ok {
		if intVal, err := strconv.ParseInt(strVal, 10, 64); err == nil {
			return intVal
		}
	}

	// number values are represented as float64
	if floatVal, ok := v.(float64); ok {
		intVal := int(floatVal)
		return intVal
	}

	return v // let terraform core handle it otherwise
}

func flattenSpannerInstanceLabels(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func flattenSpannerInstanceState(v interface{}, d *schema.ResourceData, config *Config) interface{} {
	return v
}

func expandSpannerInstanceName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandSpannerInstanceConfig(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	r := regexp.MustCompile("projects/(.+)/instanceConfigs/(.+)")
	if r.MatchString(v.(string)) {
		return v.(string), nil
	}

	project, err := getProject(d, config)
	if err != nil {
		return nil, err
	}

	return fmt.Sprintf("projects/%s/instanceConfigs/%s", project, v.(string)), nil
}

func expandSpannerInstanceDisplayName(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandSpannerInstanceNumNodes(v interface{}, d TerraformResourceData, config *Config) (interface{}, error) {
	return v, nil
}

func expandSpannerInstanceLabels(v interface{}, d TerraformResourceData, config *Config) (map[string]string, error) {
	if v == nil {
		return map[string]string{}, nil
	}
	m := make(map[string]string)
	for k, val := range v.(map[string]interface{}) {
		m[k] = val.(string)
	}
	return m, nil
}

func resourceSpannerInstanceEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	newObj := make(map[string]interface{})
	newObj["instance"] = obj
	if obj["name"] == nil {
		d.Set("name", resource.PrefixedUniqueId("tfgen-spanid-")[:30])
		newObj["instanceId"] = d.Get("name").(string)
	} else {
		newObj["instanceId"] = obj["name"]
	}
	delete(obj, "name")
	return newObj, nil
}

func resourceSpannerInstanceUpdateEncoder(d *schema.ResourceData, meta interface{}, obj map[string]interface{}) (map[string]interface{}, error) {
	project, err := getProject(d, meta.(*Config))
	if err != nil {
		return nil, err
	}
	obj["name"] = fmt.Sprintf("projects/%s/instances/%s", project, obj["name"])
	newObj := make(map[string]interface{})
	newObj["instance"] = obj
	updateMask := make([]string, 0)
	if d.HasChange("num_nodes") {
		updateMask = append(updateMask, "nodeCount")
	}
	if d.HasChange("display_name") {
		updateMask = append(updateMask, "displayName")
	}
	if d.HasChange("labels") {
		updateMask = append(updateMask, "labels")
	}
	newObj["fieldMask"] = strings.Join(updateMask, ",")
	return newObj, nil
}

func resourceSpannerInstanceDecoder(d *schema.ResourceData, meta interface{}, res map[string]interface{}) (map[string]interface{}, error) {
	config := meta.(*Config)
	d.SetId(res["name"].(string))
	if err := parseImportId([]string{"projects/(?P<project>[^/]+)/instances/(?P<name>[^/]+)"}, d, config); err != nil {
		return nil, err
	}
	res["project"] = d.Get("project").(string)
	res["name"] = d.Get("name").(string)
	id, err := replaceVars(d, config, "{{project}}/{{name}}")
	if err != nil {
		return nil, err
	}
	d.SetId(id)
	return res, nil
}
