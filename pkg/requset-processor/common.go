package requst_processor

type Request struct {
	UserId      int64
	ChatId      int64
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
