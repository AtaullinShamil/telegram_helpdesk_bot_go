package requst_processor

type Request struct {
	UserId      int64
	Department  string
	Tittle      string
	Discription string
	Status      RequestStatus
}

type RequestStatus struct {
	IsDepartment  bool
	IsTittle      bool
	IsDiscription bool
}

var Admins map[string]int64 = map[string]int64{
	"Support": 0,
	"IT":      0,
	"Billing": 0,
}
