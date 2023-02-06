package db

type Student struct {
	Name string `json:"name"`
	Id   string `json:"id"`
}
type LeaveRequset struct {
	Id            string `json:"id"  `
	Reason        string `json:"reason" `
	Date_of_leave string `json:"dateofleave"`
	Status        string `json:"status"`
}
type Approve struct {
	Id            string `json:"id"`
	Date_of_leave string `json:"dateofleave"`
}
type Admin struct {
	Username string `json:"username"`
	Passwd   int    `json:"passwd"`
}
type AdminToken struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}
type StudentToken struct {
	Id          string `json:"id"`
	Studentoken string `json:"studentoken"`
}
type StudentCredentials struct {
	Name   string `json:"name"`
	Id     string `json:"id"`
	Passwd string `json:"passwd"`
}
type StudentLoginCredentials struct {
	Id     string `json:"id"`
	Passwd string `json:"passwd"`
}
type StudentLogout struct {
	Id string `json:"id"`
}
