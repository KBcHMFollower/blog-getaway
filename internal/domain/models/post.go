package models

import postsv1 "github.com/KBcHMFollower/blog_posts_service/api/protos/gen/posts"

type Post struct {
	Id            string `json:"id"`
	UserId        string `json:"user_id"`
	Title         string `json:"title"`
	TextContent   string `json:"text_content"`
	ImagesContent string `json:"images_content"`
	Likes         int32  `json:"likes"`
	CreatedAt     string `json:"created_at"`
}

func ConvertPostFromProto(p *postsv1.Post) *Post {
	return &Post{
		Id:            p.GetId(),
		UserId:        p.GetUserId(),
		Title:         p.GetTitle(),
		TextContent:   p.GetTextContent(),
		ImagesContent: p.GetImagesContent(),
		Likes:         p.GetLikes(),
		CreatedAt:     p.GetCreatedAt().String(),
	}
}

func PostArrayFromProto(protoPosts []*postsv1.Post) []*Post {
	var posts = make([]*Post, 0, len(protoPosts))
	for _, protoPost := range protoPosts {
		posts = append(posts, ConvertPostFromProto(protoPost))
	}

	return posts
}
