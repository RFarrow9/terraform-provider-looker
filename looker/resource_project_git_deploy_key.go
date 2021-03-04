package looker

import (
	"strings"

	"github.com/billtrust/looker-go-sdk/client/project"

	apiclient "github.com/billtrust/looker-go-sdk/client"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceProjectGitDeployKey() *schema.Resource {
	return &schema.Resource{
		Create: resourceProjectGitDeployKeyCreate,
		Read:   resourceProjectGitDeployKeyRead,
		Delete: resourceProjectGitDeployKeyDelete,
		Exists: resourceProjectGitDeployKeyExists,
		Importer: &schema.ResourceImporter{
			State: resourceProjectGitDeployKeyImport,
		},

		Schema: map[string]*schema.Schema{
			"project_id": &schema.Schema{
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"ssh_deploy_key": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceProjectGitDeployKeyCreate(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerAPI30Reference)

	err := updateSession(client, "dev")
	if err != nil {
		return err
	}

	projectID := d.Get("project_id").(string)
	params := project.NewCreateProjectGitDeployKeyParams()
	params.ProjectID = projectID

	_, err = client.Project.CreateProjectGitDeployKey(params)
	if err != nil {
		return err
	}

	d.SetId(projectID)

	return resourceProjectGitDeployKeyRead(d, m)
}

func resourceProjectGitDeployKeyRead(d *schema.ResourceData, m interface{}) error {
	client := m.(*apiclient.LookerAPI30Reference)

	err := updateSession(client, "dev")
	if err != nil {
		return err
	}

	projectID := d.Id()

	params := project.NewProjectGitDeployKeyParams()
	params.ProjectID = projectID

	result, err := client.Project.ProjectGitDeployKey(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("project_id", projectID)

	// the payload is a string with 3 values separated by spaces.  The first index contains "ssh-rsa", the second index includes the key, the third index contains the project id
	// the project id doesn't appear to be part of the key though since when adding it to github, it is ignored it seems
	sshKey := strings.Fields(result.Payload)
	d.Set("ssh_deploy_key", sshKey[0]+" "+sshKey[1])

	return nil
}

func resourceProjectGitDeployKeyDelete(d *schema.ResourceData, m interface{}) error {
	// TODO There is no way to delete a git deploy key, possibly put this into the project resource (but there is no way to delete project either)
	return nil
}

func resourceProjectGitDeployKeyExists(d *schema.ResourceData, m interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	client := m.(*apiclient.LookerAPI30Reference)

	// TODO Not sure if we should always set session to "dev" instead of "production" when checking if it exists? will dev always show all dev+prod projects?
	err := updateSession(client, "dev")
	if err != nil {
		return false, err
	}

	params := project.NewProjectGitDeployKeyParams()
	params.ProjectID = d.Id()

	_, err = client.Project.ProjectGitDeployKey(params)
	if err != nil {
		if strings.Contains(err.Error(), "Not found") {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func resourceProjectGitDeployKeyImport(d *schema.ResourceData, m interface{}) ([]*schema.ResourceData, error) {
	if err := resourceProjectGitDeployKeyRead(d, m); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
