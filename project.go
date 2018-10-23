package todoist

import "encoding/json"
import "strconv"
import "net/http"
import "fmt"

const ProjectsUrl string = DefaultRestUrl + "/projects"

type Project struct {
	Id           string `json:"id"`
	Name         string `json:"name"`
	Order        int    `json:"order"`
	Indent       int    `json:"indent"`
	CommentCount int    `json:"comment_count"`
}

type jProject struct {
	Id           int    `json:"id"`
	Name         string `json:"name"`
	Order        int    `json:"order"`
	Indent       int    `json:"indent"`
	CommentCount int    `json:"comment_count"`
}

func (p *Project) UnmarshalJSON(b []byte) error {
	var jp jProject

	err := json.Unmarshal(b, &jp)
	if err != nil {
		return err
	}

	p.Id = strconv.Itoa(jp.Id)
	p.Name = jp.Name
	p.Order = jp.Order
	p.Indent = jp.Indent
	p.CommentCount = jp.CommentCount

	return nil
}

func (c *Client) GetProjects() ([]Project, error) {
	req, err := http.NewRequest("GET", ProjectsUrl, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var projects []Project

	err = json.Unmarshal(res, &projects)
	if err != nil {
		return nil, err
	}

	return projects, nil
}

func (c *Client) GetProject(id string) (*Project, error) {
	var project Project

	req, err := http.NewRequest("GET", fmt.Sprintf(ProjectsUrl+"/%s", id), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &project)
	if err != nil {
		return nil, err
	}

	return &project, nil
}

func (c *Client) GetProjectByName(name string) (*Project, error) {
	projects, err := c.GetProjects()
	if err != nil {
		return nil, err
	}

	for _, project := range projects {
		if project.Name == name {
			return &project, nil
		}
	}

	return nil, fmt.Errorf("Project with name %s not found", name)
}
