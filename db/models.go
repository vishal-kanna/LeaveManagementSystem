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
