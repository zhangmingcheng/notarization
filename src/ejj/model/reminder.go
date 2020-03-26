package model

//Reminder 催告记录
type Reminder struct {
	Type				string				`json:"Type"`				//类别，类似表名，区分大小写
	DebtID				string				`json:"DebtID"`				//债项编号
	ReminderHashList	[]string			`json:"ReminderHashList"`	//催告记录的哈希列表，包含短信记录、邮件记录、打电话录音、执行证书、附件的哈希
}

func (rr Reminder) GetType() string {
	return rr.Type
}

func (rr Reminder) GetID() string {
	return rr.DebtID
}

func (rr *Reminder) AppendString(strs ...string) error {
	rr.ReminderHashList = append(rr.ReminderHashList, strs...)
	return nil
}

