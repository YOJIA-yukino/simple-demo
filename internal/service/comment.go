package service

import (
	"fmt"
	"github.com/YOJIA-yukino/simple-douyin-backend/api"
	"github.com/YOJIA-yukino/simple-douyin-backend/internal/dao"
	"github.com/YOJIA-yukino/simple-douyin-backend/internal/model"
	"github.com/YOJIA-yukino/simple-douyin-backend/internal/utils/constants"
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
func (*commentService) CommentInfoPush(userId, videoId int64, text string) (int64, error) {
	var err error
	var commentId int64
	commentId, err = dao.GetCommentDaoInstance().Add(userId, videoId, text)
	if err != nil {
		return -1, err
	}
	err = dao.GetCommentDaoInstance().AddCommentCount(videoId)
	if err != nil {
		return -1, err
	}
	return commentId, nil
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

func (*commentService) GetCommentByCommentId(commentId int64) (*api.Comment, error) {
	// var commentApi api.Comment
	comment, err := dao.GetCommentDaoInstance().GetCommentByCommentId(commentId)
	if err != nil {
		return nil, err
	}
	commentApi, err := newComment(comment)
	if err != nil {
		return nil, err
	}
	return commentApi, nil
}

// 获取一条视频的评论列表 后续改成发送数据而不是指针?：可不用同一个地址空间
func (*commentService) GetCommentList(videoId int64) (*[]api.Comment, error) {
	comments, err := dao.GetCommentDaoInstance().GetCommentList(videoId)
	if err != nil {
		return nil, err
	}
	commentList, err := newCommentList(comments)
	if err != nil {
		return nil, err
	}
	return commentList, nil // 返回slice的引用
}

// 转换到API层的Comment
func newComment(comment *model.Comment) (*api.Comment, error) {
	userInfo, err := GetUserServiceInstance().GetUserByUserId(comment.UserID)
	if err != nil {
		return nil, constants.InnerDataBaseErr
	}
	dateStr := fmt.Sprintf("%s-%d", comment.CreatedAt.Month().String(), comment.CreatedAt.Day())
	commentApi := api.Comment{
		Id: int64(comment.ID),
		User: api.User{
			Id:            userInfo.UserID,
			Name:          userInfo.UserName,
			FollowCount:   userInfo.FollowCount,
			FollowerCount: userInfo.FollowerCount,
			IsFollow:      false,
		},
		Content:    comment.Content,
		CreateDate: dateStr,
	}
	return &commentApi, nil
}

// 转换到API层的Comment列表
func newCommentList(comments []*model.Comment) (*[]api.Comment, error) {
	commentList := make([]api.Comment, len(comments))

	for i, v := range comments {
		userInfo, err := GetUserServiceInstance().GetUserByUserId(v.UserID)
		if err != nil {
			return nil, constants.InnerDataBaseErr
		}
		dateStr := fmt.Sprintf("%s-%d", v.CreatedAt.Month().String(), v.CreatedAt.Day())

		commentList[i] = api.Comment{
			Id: int64(v.ID),
			User: api.User{
				Id:            userInfo.UserID,
				Name:          userInfo.UserName,
				FollowCount:   userInfo.FollowCount,
				FollowerCount: userInfo.FollowerCount,
				IsFollow:      false,
			},
			Content:    v.Content,
			CreateDate: dateStr,
		}
	}
	return &commentList, nil
}
