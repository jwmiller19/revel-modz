package forum

import (
	"errors"
	"time"

	"github.com/jinzhu/gorm"
)

// database tables
type ForumTopic struct {
	// gorm fields
	Id        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	TopicId   int64 `sql:"not null;unique"`
	TopicName string
}

type ForumTopicSubscriber struct {
	// gorm fields
	Id        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	// keys
	TopicId   int64
	UserEmail string
}

type ForumMessage struct {
	// gorm fields
	Id        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	// keys
	TopicId   int64
	MessageId int64

	// message details
	AuthorName  string
	MessageBody string
}

type ForumTopicTag struct {
	// gorm fields
	Id        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	TopicId  int64
	TopicTag string
}

type ForumTopicStats struct {
	// gorm fields
	Id        int64
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time

	TopicId int64

	NumMessages int64
	NumViews    int64
}

func AddTables(db *gorm.DB) error {
	var err error
	err = db.AutoMigrate(ForumTopic{}).Error
	if err != nil {
		return err
	}
	err = db.AutoMigrate(ForumTopicSubscriber{}).Error
	if err != nil {
		return err
	}
	err = db.AutoMigrate(ForumMessage{}).Error
	if err != nil {
		return err
	}
	err = db.AutoMigrate(ForumTopicTag{}).Error
	if err != nil {
		return err
	}
	err = db.AutoMigrate(ForumTopicStats{}).Error
	if err != nil {
		return err
	}
	return err
}

func DropTables(db *gorm.DB) error {
	var err error
	err = db.DropTable(ForumTopic{}).Error
	if err != nil {
		return err
	}
	err = db.DropTable(ForumTopicSubscriber{}).Error
	if err != nil {
		return err
	}
	err = db.DropTable(ForumMessage{}).Error
	if err != nil {
		return err
	}
	err = db.DropTable(ForumTopicTag{}).Error
	if err != nil {
		return err
	}
	err = db.DropTable(ForumTopicStats{}).Error
	if err != nil {
		return err
	}
	return nil
}

func FillTables(db *gorm.DB) error {
	for _, t := range fakeTopics {
		err := db.Save(t).Error
		if err != nil {
			return err
		}
		stats := &ForumTopicStats{
			TopicId:     t.TopicId,
			NumMessages: 0,
			NumViews:    0,
		}
		err = db.Save(stats).Error
		if err != nil {
			return err
		}
	}
	for _, m := range fakeMessages {
		err := AddNewMessage(db, m.AuthorName, m.MessageBody, m.TopicId)
		if err != nil {
			return err
		}
	}

	return nil
}
func TestTables(db *gorm.DB) error {
	return errors.New("TODO")
}

var fakeTopics = []ForumTopic{
	ForumTopic{TopicId: 1, TopicName: "Topic1"},
	ForumTopic{TopicId: 2, TopicName: "Topic2"},
	ForumTopic{TopicId: 3, TopicName: "Topic3"},
	ForumTopic{TopicId: 4, TopicName: "Topic4"},
	ForumTopic{TopicId: 5, TopicName: "Topic5"},
	ForumTopic{TopicId: 6, TopicName: "Topic6"},
}

var fakeMessages = []ForumMessage{
	ForumMessage{TopicId: 1, MessageId: 1, AuthorName: "user1", MessageBody: "t1 msg1 body"},
	ForumMessage{TopicId: 1, MessageId: 2, AuthorName: "user2", MessageBody: "t1 msg2 body"},
	ForumMessage{TopicId: 1, MessageId: 3, AuthorName: "user3", MessageBody: "t1 msg3 body"},

	ForumMessage{TopicId: 2, MessageId: 1, AuthorName: "user2", MessageBody: "t2 msg1 body"},
	ForumMessage{TopicId: 2, MessageId: 2, AuthorName: "user3", MessageBody: "t2 msg2 body"},
	ForumMessage{TopicId: 2, MessageId: 3, AuthorName: "user4", MessageBody: "t2 msg3 body"},

	ForumMessage{TopicId: 3, MessageId: 1, AuthorName: "user3", MessageBody: "t3 msg1 body"},
	ForumMessage{TopicId: 3, MessageId: 2, AuthorName: "user4", MessageBody: "t3 msg2 body"},
	ForumMessage{TopicId: 3, MessageId: 3, AuthorName: "user5", MessageBody: "t3 msg3 body"},

	ForumMessage{TopicId: 4, MessageId: 1, AuthorName: "user4", MessageBody: "t4 msg1 body"},
	ForumMessage{TopicId: 4, MessageId: 2, AuthorName: "user5", MessageBody: "t4 msg2 body"},
	ForumMessage{TopicId: 4, MessageId: 3, AuthorName: "user6", MessageBody: "t4 msg3 body"},

	ForumMessage{TopicId: 5, MessageId: 1, AuthorName: "user5", MessageBody: "t5 msg1 body"},
	ForumMessage{TopicId: 5, MessageId: 2, AuthorName: "user6", MessageBody: "t5 msg2 body"},
	ForumMessage{TopicId: 5, MessageId: 3, AuthorName: "user7", MessageBody: "t5 msg3 body"},

	ForumMessage{TopicId: 6, MessageId: 1, AuthorName: "user6", MessageBody: "t6 msg1 body"},
	ForumMessage{TopicId: 6, MessageId: 2, AuthorName: "user7", MessageBody: "t6 msg2 body"},
	ForumMessage{TopicId: 6, MessageId: 3, AuthorName: "user8", MessageBody: "t6 msg3 body"},
}
