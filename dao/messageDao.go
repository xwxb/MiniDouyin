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
	return err == nil, nil
}

// returns the list of messages from fromUserId to toUserId after tm(Unix time)
func GetRecentMessageListByUserId(tm, fromUserId, toUserId int64) ([]Message, error) {
	condi := "from_user_id = ? AND to_user_id = ? AND create_time > ?"
	var msgList []Message
	err := Db.Where(condi, fromUserId, toUserId, tm).Order("create_time ASC").Find(&msgList).Error
	return msgList, err
}
