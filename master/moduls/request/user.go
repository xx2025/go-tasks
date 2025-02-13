package request

type UserSave struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Password string `json:"password"`
	RoleId   int    `json:"roleId"`
	Status   int    `json:"status"`
	Avatar   string `json:"avatar"`
}

type UserSaveRequest struct {
	Id       int    `json:"id"`
	Nickname string `json:"nickname"`
	RoleId   int    `json:"roleId"`
	Status   int    `json:"status"`
	Password string `json:"password"`
}
