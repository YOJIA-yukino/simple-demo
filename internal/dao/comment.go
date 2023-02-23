package dao

import (
	"errors"
	"github.com/YOJIA-yukino/simple-douyin-backend/internal/model"
	"github.com/YOJIA-yukino/simple-douyin-backend/internal/utils/constants"
	"gorm.io/gorm"
	"sync"
	"time"
)

type commentDao struct{}

var (
	commentDaoInstance *commentDao
	commentOnce        sync.Once
)

// 获取一个Dao层与Comment操作有关的Instance
func GetCommentDaoInstance() *commentDao {
	dataBaseInitialization()
	commentOnce.Do(func() {
		commentDaoInstance = &commentDao{}
	})
	return commentDaoInstance
}

// 向数据库添加一条评论
func (*commentDao) Add(userId, videoId int64, text string) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var err error
		var comment model.Comment
		comment.CreatedAt = time.Now()
		comment.UserID = userId
		comment.VideoID = videoId
		comment.Content = text
		comment.LikeCount = 0
		comment.TeaseCount = 0

		if err = tx.Create(&comment).Error; err != nil {
			return constants.InnerDataBaseErr
		}
		return nil
	})
}

// 数据库中删除一条评论
func (*commentDao) Del(userId, videoId, commentId int64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var err error
		var comment model.Comment

		err = tx.Where("comment_id = ? And video_id = ? And user_id = ?", commentId, videoId, userId).
			First(&comment).Error
		if errors.Is(gorm.ErrRecordNotFound, err) {
			return constants.RecordNotExistErr
		} else if err != nil {
			return constants.InnerDataBaseErr
		}

		if err = tx.Delete(&comment).Error; err != nil {
			return constants.InnerDataBaseErr
		}
		return nil
	})
}

// 视频的评论加一
func (*commentDao) AddCommentCount(videoId int64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var video model.Video
		if err := db.Where("video_id = ?", videoId).First(&video).Error; err != nil {
			if errors.Is(gorm.ErrRecordNotFound, err) {
				return constants.RecordNotExistErr
			} else {
				return constants.InnerDataBaseErr
			}
		}
		if err := tx.Model(&model.Video{}).
			Where("video_id = ?", videoId).
			Update("comment_count", video.CommentCount+1).Error; err != nil {
			return constants.InnerDataBaseErr
		}
		return nil
	})
}

// 视频的评论减一
func (*commentDao) SubCommentCount(videoId int64) error {
	return db.Transaction(func(tx *gorm.DB) error {
		var video model.Video
		if err := db.Where("video_id = ?", videoId).First(&video).Error; err != nil {
			if errors.Is(gorm.ErrRecordNotFound, err) {
				return constants.RecordNotExistErr
			} else {
				return constants.InnerDataBaseErr
			}
		}
		if err := tx.Model(&model.Video{}).
			Where("video_id = ?", videoId).
			Update("comment_count", video.CommentCount-1).Error; err != nil {
			return constants.InnerDataBaseErr
		}
		return nil
	})
}

// 从数据库中获取视频的评论列表 按时间降序排序
func (*commentDao) getCommentList(videoId int64) ([]*model.Comment, error) {
	comments := make([]*model.Comment, 0)
	err := db.Where("video_id = ?", videoId).Order("created_at desc").Find(&comments).Error
	if err != nil {
		return nil, constants.InnerDataBaseErr
	}
	return comments, nil
}
