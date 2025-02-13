package request

type NodeSave struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Url    string `json:"url"`
	Status int    `json:"status"`
}

type ProjectSave struct {
	Id       int    `json:"id"`
	Name     string `json:"name"`
	Describe string `json:"describe"`
}

type TaskSave struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Spec      string `json:"spec"`
	Status    int    `json:"status"`
	ProjectId int    `json:"projectId"`
	NodeId    int    `json:"nodeId"`
	Describe  string `json:"describe"`
}

type ProcessSave struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Status     int    `json:"status"`
	ProjectId  int    `json:"projectId"`
	NodeId     int    `json:"nodeId"`
	MaxRetries int    `json:"maxRetries"`
	Describe   string `json:"describe"`
}
