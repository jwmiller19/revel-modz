package forum

import (
	"github.com/jinzhu/gorm"
	"github.com/revel/revel"
)

func getNextTopicId(db *gorm.DB) (int64, error) {
	var ft ForumTopic
	err := db.Where(&ForumTopic{}).Order("topic_id desc").First(&ft).Error
	if err == gorm.RecordNotFound {
		return 1, nil
	}
	if err != nil {
		revel.TRACE.Println(err)
		return 0, err
	}
	return ft.TopicId + 1, nil
}

func getNextMsgIdByTopicId(db *gorm.DB, tId int64) (int64, error) {
	var fm ForumMessage
	err := db.Where(&ForumMessage{TopicId: tId}).Order("message_id desc").First(&fm).Error
	if err == gorm.RecordNotFound {
		return 0, nil
	}
	if err != nil {
		revel.TRACE.Println(err)
		return 0, err
	}
	return fm.MessageId + 1, nil
}

func getStatsByTopicId(db *gorm.DB, tId int64) (*ForumTopicStats, error) {
	var stats ForumTopicStats
	err := db.Where(&ForumTopicStats{TopicId: tId}).First(&stats).Error
	if err == gorm.RecordNotFound {
		return nil, nil
	}
	if err != nil {
		revel.TRACE.Println(err)
		return nil, err
	}
	return &stats, nil
}

func getSubscriberByTopicId(db *gorm.DB, tId int64) ([]ForumTopicSubscriber, error) {
	var subscriber []ForumTopicSubscriber
	err := db.Where(&ForumTopicSubscriber{TopicId: tId}).Find(subscriber).Error
	if err == gorm.RecordNotFound {
		return nil, nil
	}
	if err != nil {
		revel.TRACE.Println(err)
		return nil, err
	}
	return subscriber, nil
}

func checkSubscriberByTopicId(db *gorm.DB, tId int64, email string) (bool, error) {
	var subscriber ForumTopicSubscriber
	err := db.Where(&ForumTopicSubscriber{TopicId: tId, UserEmail: email}).First(&subscriber).Error
	if err == gorm.RecordNotFound {
		return false, nil
	}
	if err != nil {
		revel.TRACE.Println(err)
		return false, err
	}
	return true, nil
}
