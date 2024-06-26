package models

import (
	commentsv1 "github.com/KBcHMFollower/blog_posts_service/api/protos/gen/comments"
)

type Comment struct {
	Id        string `json:"id"`
	PostId    string `json:"post_id"`
	UserId    string `json:"user_id"`
	Content   string `json:"content"`
	Likes     int32  `json:"likes"`
	CreatedAt string `json:"created_at"`
}

func ConvertCommentFromProto(c *commentsv1.Comment) *Comment {
	return &Comment{
		Id:        c.GetId(),
		UserId:    c.GetUserId(),
		PostId:    c.GetPostId(),
		Content:   c.GetContent(),
		Likes:     c.GetLikes(),
		CreatedAt: c.GetCreatedAt().String(),
	}
}

func CommentsArrayFromProto(comments []*commentsv1.Comment) []*Comment {
	result := make([]*Comment, 0, len(comments))
	for _, comment := range comments {
		result = append(result, ConvertCommentFromProto(comment))
	}

	return result
}
