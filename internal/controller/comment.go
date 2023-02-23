package controller

import (
	"context"
	"errors"
	"github.com/YOJIA-yukino/simple-douyin-backend/api"
	"github.com/YOJIA-yukino/simple-douyin-backend/internal/service"
	"github.com/YOJIA-yukino/simple-douyin-backend/internal/utils/constants"
	"github.com/YOJIA-yukino/simple-douyin-backend/internal/utils/jwt"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"strconv"
)

type CommentListResponse struct {
	api.Response
	CommentList []api.Comment `json:"comment_list,omitempty"`
}

type CommentActionResponse struct {
	api.Response
	Comment api.Comment `json:"comment,omitempty"`
}

// CommentAction 视频评论接口
func CommentAction(c context.Context, ctx *app.RequestContext) {
	var err error
	var text string
	var commentId int64

	loginUserId, err := jwt.GetUserId(c, ctx)
	if err != nil {
		ctx.JSON(consts.StatusOK, api.Response{
			StatusCode: int32(api.TokenInvalidErr),
			StatusMsg:  api.ErrorCodeToMsg[api.TokenInvalidErr],
		})
		return
	}
	videoIdStr := ctx.Query("video_id")
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		ctx.JSON(consts.StatusOK, api.Response{
			StatusCode: int32(api.InputFormatCheckErr),
			StatusMsg:  api.ErrorCodeToMsg[api.InputFormatCheckErr],
		})
		return
	}
	actionTypeStr := ctx.Query("action_type")
	actionType, err := strconv.ParseInt(actionTypeStr, 10, 64)
	if err != nil {
		ctx.JSON(consts.StatusOK, api.Response{
			StatusCode: int32(api.InputFormatCheckErr),
			StatusMsg:  api.ErrorCodeToMsg[api.InputFormatCheckErr],
		})
		return
	}

	// 对不同actionType 获取对应参数
	switch actionType {
	case api.PushComment:
		text = ctx.Query("comment_text")
		if text == "" {
			ctx.JSON(consts.StatusOK, api.Response{
				StatusCode: int32(api.InputFormatCheckErr),
				StatusMsg:  api.ErrorCodeToMsg[api.InputFormatCheckErr],
			})
			return
		}
		commentId, err = CommentActionPush(loginUserId, videoId, text)

	case api.DeleteComment:
		commentIdStr := ctx.Query("comment_id")
		commentId, err := strconv.ParseInt(commentIdStr, 10, 64)
		if err != nil {
			ctx.JSON(consts.StatusOK, api.Response{
				StatusCode: int32(api.InputFormatCheckErr),
				StatusMsg:  api.ErrorCodeToMsg[api.InputFormatCheckErr],
			})
			return
		}
		err = CommentActionDelete(loginUserId, videoId, commentId)
	default:
		ctx.JSON(consts.StatusOK, api.Response{
			StatusCode: int32(api.UnKnownActionType),
			StatusMsg:  api.ErrorCodeToMsg[api.UnKnownActionType],
		})
		return
	}
	if err != nil {
		if errors.Is(constants.RecordNotExistErr, err) {
			ctx.JSON(consts.StatusOK, api.Response{
				StatusCode: int32(api.RecordNotExistErr),
				StatusMsg:  api.ErrorCodeToMsg[api.RecordNotExistErr],
			})
		} else if errors.Is(constants.InnerDataBaseErr, err) {
			ctx.JSON(consts.StatusOK, api.Response{
				StatusCode: int32(api.InnerDataBaseErr),
				StatusMsg:  api.ErrorCodeToMsg[api.InnerDataBaseErr],
			})
		}
		return
	}

	// 返回json response
	// userInfo, err := service.GetUserServiceInstance().GetUserByUserId(loginUserId)
	switch actionType {
	case api.PushComment:
		if commentId == -1 {
			ctx.JSON(consts.StatusOK, api.Response{
				StatusCode: int32(api.InnerDataBaseErr),
				StatusMsg:  api.ErrorCodeToMsg[api.InnerDataBaseErr] + " & commentId Init error",
			})
		}
		comment, err := service.GetCommentServiceInstance().GetCommentByCommentId(commentId)
		if err != nil {
			ctx.JSON(consts.StatusOK, api.Response{
				StatusCode: int32(api.InnerDataBaseErr),
				StatusMsg:  api.ErrorCodeToMsg[api.InnerDataBaseErr],
			})
		}
		ctx.JSON(consts.StatusOK, api.CommentActionResponse{
			StatusCode: 0,
			Comment:    *comment,
		})
	case api.DeleteComment:
		ctx.JSON(consts.StatusOK, api.Response{
			StatusCode: 0,
		})
	}

}

func CommentActionPush(loginUserId, videoId int64, text string) (int64, error) {
	var commentId int64
	commentId, err := service.GetCommentServiceInstance().CommentInfoPush(loginUserId, videoId, text)
	return commentId, err
}

func CommentActionDelete(loginUserId, videoId, commentId int64) error {
	err := service.GetCommentServiceInstance().CommentInfoDelete(loginUserId, videoId, commentId)
	return err
}

// 获取一个视频的所有评论，按发布时间倒序
func CommentList(c context.Context, ctx *app.RequestContext) {
	//todo
	var err error
	_, err = jwt.GetUserId(c, ctx) // loginUserId not used yet
	if err != nil {
		ctx.JSON(consts.StatusOK, api.Response{
			StatusCode: int32(api.TokenInvalidErr),
			StatusMsg:  api.ErrorCodeToMsg[api.TokenInvalidErr],
		})
		return
	}
	videoIdStr := ctx.Query("video_id")
	videoId, err := strconv.ParseInt(videoIdStr, 10, 64)
	if err != nil {
		ctx.JSON(consts.StatusOK, api.Response{
			StatusCode: int32(api.InputFormatCheckErr),
			StatusMsg:  api.ErrorCodeToMsg[api.InputFormatCheckErr],
		})
		return
	}

	commentList, err := service.GetCommentServiceInstance().GetCommentList(videoId)
	if err != nil {
		ctx.JSON(consts.StatusOK, api.Response{
			StatusCode: int32(api.InnerDataBaseErr),
			StatusMsg:  api.ErrorCodeToMsg[api.InnerDataBaseErr],
		})
		return
	}
	ctx.JSON(consts.StatusOK, api.CommentListResponse{
		Response: api.Response{
			StatusCode: 0,
		},
		CommentList: *commentList,
	})
}
