package controllers

import (
	"github.com/revel/revel"

	"github.com/iassic/revel-modz/modules/forum"
	"github.com/iassic/revel-modz/sample/app/routes"
)

func (c App) Forum(msg_pos, count int) revel.Result {
	if count == 0 {
		count = 20
	}
	topics, err := forum.GetTopicList(c.Txn, msg_pos, count)
	if err != nil {
		revel.ERROR.Println("Getting forum topic list: ", err)
	}

	return c.Render(topics)
}

func (c App) ForumTopic(topic_id, msg_id int) revel.Result {
	revel.INFO.Println("Forum: ", topic_id, msg_id)
	if msg_id != 0 {
		// enable the scroll to message
	}

	detail, err := forum.GetAllMessagesByTopicId(c.Txn, int64(topic_id))
	if err != nil {
		revel.ERROR.Println("Getting forum topic detail: ", topic_id, err)
	}

	return c.Render(detail)
}

func (c App) ForumMessage(topic_id, msg_id int) revel.Result {
	return c.Redirect(routes.App.ForumTopic(topic_id, msg_id))
}

func (c User) ForumPost(author, subject, content string, tags []string) revel.Result {
	revel.INFO.Println("Forum POST: ", author, subject, content, tags)
	return c.Render()
}