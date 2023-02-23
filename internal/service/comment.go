package service

import (
	"github.com/YOJIA-yukino/simple-douyin-backend/internal/dao"
	"sync"
)

type commentService struct{}

var (
	commentServiceInstance *commentService
	commentOnce            sync.Once
)

// 获取一个commentService实例
func GetCommentServiceInstance() *commentService {
	commentOnce.Do(func() {
		commentServiceInstance = &commentService{}
	})
	return commentServiceInstance
}

// 添加评论
func (*commentService) CommentInfoPush(userId, videoId int64, text string) error {
	var err error

	err = dao.GetCommentDaoInstance().Add(userId, videoId, text)
	if err != nil {
		return err
	}
	err = dao.GetCommentDaoInstance().AddCommentCount(videoId)
	if err != nil {
		return err
	}
	return nil
}

// 删除评论
func (*commentService) CommentInfoDelete(userId, videoId, commentId int64) error {
	var err error

	err = dao.GetCommentDaoInstance().Del(userId, videoId, commentId)
	if err != nil {
		return err
	}
	err = dao.GetCommentDaoInstance().SubCommentCount(videoId)
	if err != nil {
		return nil
	}
	return nil
}

func (*commentService) getCommentList(videoid int64) error {

}
