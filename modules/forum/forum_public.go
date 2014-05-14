package forum

import (
	"errors"
	"fmt"
	"time"

	"github.com/jinzhu/gorm"
)

// not a database table, for sending over the wire
// Used in the list of topics
type TopicBrief struct {
	TopicId   int64
	TopicName string
	TopicTags []string

	NumViews    int64
	NumMessages int64

	OpenedAt  time.Time
	OpenedBy  string
	LastMsgAt time.Time
	LastMsgBy string
}

type MessageWire struct {
	MessageId   int64
	MessageTime time.Time
	MessageAuth string
	MessageBody string
}

type TopicDetail struct {
	Brief    TopicBrief
	Messages []MessageWire
}

func AddNewTopic(db *gorm.DB, author, subject, body string) error {
	var err error
	tId, _ := getNextTopicId(db)
	if err != nil {
		return err
	}
	mId, _ := getNextMsgIdByTopicId(db, tId)
	if err != nil {
		return err
	}

	topic := &ForumTopic{
		TopicId:   tId,
		TopicName: subject,
	}

	// create subscriber struct
	subscriber := &ForumTopicSubscriber{
		TopicId:   tId,
		UserEmail: author,
	}

	msg := &ForumMessage{
		TopicId:     tId,
		MessageId:   mId,
		AuthorName:  author,
		MessageBody: body,
	}
	stats := &ForumTopicStats{
		TopicId:     tId,
		NumMessages: 1,
		NumViews:    0,
	}

	err = db.Save(topic).Error
	if err != nil {
		return err
	}
	// save the subsrcibers
	err = db.Save(subscriber).Error
	if err != nil {
		return err
	}
	err = db.Save(msg).Error
	if err != nil {
		return err
	}
	err = db.Save(stats).Error
	if err != nil {
		return err
	}
	return nil
}

func AddNewMessage(db *gorm.DB, author, message string, topicId int64) error {
	var err error
	mId, err := getNextMsgIdByTopicId(db, topicId)
	if err != nil {
		return err
	}
	stats, err := getStatsByTopicId(db, topicId)
	if err != nil {
		return err
	}

	msg := &ForumMessage{
		TopicId:     topicId,
		MessageId:   mId,
		AuthorName:  author,
		MessageBody: message,
	}

	// check if already subscribed to email. if not,
	// create subscriber struct

	stats.NumMessages += 1

	// save the subscriber

	err = db.Save(msg).Error
	if err != nil {
		return err
	}
	err = db.Save(stats).Error
	if err != nil {
		return err
	}
	return nil
}

func AddTopicTag(db *gorm.DB, tag string, topicId int64) error {
	var err error
	ftag := &ForumTopicTag{
		TopicId:  topicId,
		TopicTag: tag,
	}

	err = db.Save(ftag).Error
	if err != nil {
		return err
	}
	return nil
}

func GetTopicList(db *gorm.DB, start, count int) ([]TopicBrief, error) {
	var topics []ForumTopic

	err := db.Order("updated_at desc").Offset(start).Limit(count).Find(&topics).Error
	if err != nil {
		return nil, errors.New("ForumTopic Table " + fmt.Sprint(err))
	}

	fmt.Println("#Topics: ", len(topics))

	briefs := make([]TopicBrief, len(topics))
	for i, topic := range topics {
		var (
			firstMsg ForumMessage
			lastMsg  ForumMessage
			tags     []ForumTopicTag
			stats    ForumTopicStats
			err      error
		)

		id := topic.TopicId

		// get current info
		err = db.Where(ForumMessage{TopicId: id}).Order("created_at").First(&firstMsg).Last(&lastMsg).Error
		if err != nil {
			return nil, errors.New("ForumMessage Table " + fmt.Sprint(err))
		}
		err = db.Where(ForumTopicTag{TopicId: id}).Find(&tags).Error
		if err != nil && err != gorm.RecordNotFound {
			return nil, errors.New("ForumTopicTag Table " + fmt.Sprint(err))
		}
		err = db.Where(ForumTopic{TopicId: id}).First(&stats).Error
		if err != nil && err != gorm.RecordNotFound {
			return nil, errors.New("ForumTopicStats Table " + fmt.Sprint(err))
		}

		// build detail
		mtags := make([]string, len(tags))
		for i, t := range tags {
			mtags[i] = t.TopicTag
		}

		briefs[i] = TopicBrief{
			TopicId:     topic.TopicId,
			TopicName:   topic.TopicName,
			TopicTags:   mtags,
			NumViews:    stats.NumViews,
			NumMessages: stats.NumMessages,
			OpenedAt:    firstMsg.CreatedAt,
			OpenedBy:    firstMsg.AuthorName,
			LastMsgAt:   lastMsg.CreatedAt,
			LastMsgBy:   lastMsg.AuthorName,
		}

	}

	return briefs, nil
}

func GetAllMessagesByTopicId(db *gorm.DB, id int64) (*TopicDetail, error) {
	var (
		topic    ForumTopic
		messages []ForumMessage
		tags     []ForumTopicTag
		stats    ForumTopicStats
		err      error
	)

	// get current info
	err = db.Where(ForumTopic{TopicId: id}).First(&topic).Error
	if err != nil {
		return nil, errors.New("ForumTopic Table " + fmt.Sprint(err))
	}
	err = db.Where(ForumMessage{TopicId: id}).Order("created_at").Find(&messages).Error
	if err != nil {
		return nil, errors.New("ForumMessage Table " + fmt.Sprint(err))
	}
	err = db.Where(ForumTopicTag{TopicId: id}).Find(&tags).Error
	if err != nil && err != gorm.RecordNotFound {
		return nil, errors.New("ForumTopicTag Table " + fmt.Sprint(err))
	}
	err = db.Where(ForumTopic{TopicId: id}).First(&stats).Error
	if err != nil && err != gorm.RecordNotFound {
		return nil, errors.New("ForumTopicStats Table " + fmt.Sprint(err))
	}

	lastMsgPos := len(messages) - 1

	// build detail
	mtags := make([]string, len(tags))
	for i, t := range tags {
		mtags[i] = t.TopicTag
	}

	detail := new(TopicDetail)
	detail.Brief = TopicBrief{
		TopicId:     topic.TopicId,
		TopicName:   topic.TopicName,
		TopicTags:   mtags,
		NumViews:    stats.NumViews,
		NumMessages: stats.NumMessages,
		OpenedAt:    messages[0].CreatedAt,
		OpenedBy:    messages[0].AuthorName,
		LastMsgAt:   messages[lastMsgPos].CreatedAt,
		LastMsgBy:   messages[lastMsgPos].AuthorName,
	}
	detail.Messages = make([]MessageWire, len(messages))
	for i, m := range messages {
		detail.Messages[i] = MessageWire{
			MessageId:   m.MessageId,
			MessageTime: m.CreatedAt,
			MessageAuth: m.AuthorName,
			MessageBody: m.MessageBody,
		}
	}

	// update stats
	stats.NumViews += 1
	err = db.Save(&stats).Error
	if err != nil {
		return nil, err
	}

	return detail, nil
}

// add unsubscribe function
