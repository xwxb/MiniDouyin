package dao

type Message struct {
	Id         int64  `gorm:"primary_key;AUTO_INCREMENT" json:"id,omitempty"`
	ToUserId   int64  `gorm:"to_user_id" json:"to_user_id,omitempty"`
	FromUserId int64  `gorm:"from_user_id" json:"from_user_id,omitempty"`
	Content    string `gorm:"content" json:"content,omitempty"`
	CreateTime int64  `gorm:"create_time" json:"create_time,omitempty"`
}

func (Message) TableName() string {
	return "message"
}

// SendMessage(msg) sends a message, and returns if it is successful
func SendMessage(message *Message) (bool, error) {
	err := Db.Model(&Message{}).Create(&message).Error
	return err == nil, err
}

// returns the list of messages between u1 and u2 after tm(Unix time)
func GetRecentMessageListByUserId(tm, u1, u2 int64) ([]Message, error) {
	condi := "from_user_id = ? AND to_user_id = ? AND create_time > ?"
	order := "create_time ASC"
	var msgList []Message
	err := Db.Where(condi, u1, u2, tm).Or(condi, u2, u1, tm).Order(order).Find(&msgList).Error
	return msgList, err
}
