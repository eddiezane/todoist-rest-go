package todoist

import "net/http"
import "encoding/json"
import "fmt"
import "bytes"

const TasksUrl string = DefaultUrl + "/tasks"

type Task struct {
	Id           string `json:"id,Number"`
	ProjectId    int    `json:"project_id"`
	Content      string `json:"content"`
	Completed    bool   `json:"completed"`
	LabelIds     []int  `json:"label_ids"`
	Order        int    `json:"order"`
	Indent       int    `json:"indent"`
	Priority     int    `json:"priority"`
	Due          Due    `json:"due"`
	Url          string `json:"url"`
	CommentCount int    `json:"comment_count"`
}

type jTask struct {
	Task
	Id int `json:"id"`
}

type NewTask struct {
	Content     string `json:"content"` // required
	ProjectId   int    `json:"project_id,omitempty"`
	Order       int    `json:"order,omitempty"`
	LabelIds    []int  `json:"label_ids,omitempty"`
	Priority    int    `json:"priority,omitempty"`
	DueString   string `json:"due_string,omitempty"`
	DueDate     string `json:"due_date,omitempty"`
	DueDatetime string `json:"due_datetime,omitempty"`
	DueLang     string `json:"due_lang,omitempty"`
}

type Due struct {
	String   string `json:"string"` // required
	Date     string `json:"date"`   // required
	Datetime string `json:"datetime"`
	Timezone string `json:"timezone"`
}

func (t *Task) UnmarshalJSON(b []byte) error {
	var task jTask

	err := json.Unmarshal(b, &task)
	if err != nil {
		return err
	}

	*t = task.Task

	return nil
}

func (c *Client) GetTasks() ([]Task, error) {
	req, err := http.NewRequest("GET", TasksUrl, nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var tasks []Task

	err = json.Unmarshal(res, &tasks)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}

func (c *Client) GetTask(id string) (*Task, error) {
	var task Task

	req, err := http.NewRequest("GET", fmt.Sprintf(TasksUrl+"/%s", id), nil)
	if err != nil {
		return nil, err
	}

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(res, &task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (c *Client) CreateTask(t *NewTask) (*Task, error) {
	j, err := json.Marshal(t)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", TasksUrl, bytes.NewBuffer(j))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	res, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	var task Task
	err = json.Unmarshal(res, &task)
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (c *Client) CloseTask(id string) error {
	req, err := http.NewRequest("POST", TasksUrl+"/"+id+"/close", nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) ReopenTask(id string) error {
	req, err := http.NewRequest("POST", TasksUrl+"/"+id+"/reopen", nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) DeleteTask(id string) error {
	req, err := http.NewRequest("DELETE", TasksUrl+"/"+id, nil)
	if err != nil {
		return err
	}

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) UpdateTask(t *Task) error {
	j, err := json.Marshal(t)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", TasksUrl+"/"+t.Id, bytes.NewBuffer(j))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	_, err = c.doRequest(req)
	if err != nil {
		return err
	}

	return nil
}
