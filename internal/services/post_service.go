package services

import (
	"context"
	postsv1 "github.com/KBcHMFollower/blog_posts_service/api/protos/gen/posts"
	"log/slog"
	dependencies "test-plate/internal/dependencies/postservice"
)

type PostService struct {
	postsGRpc *dependencies.PostsGrpcClient
	log       *slog.Logger
}

func NewPostService(postsGRpc *dependencies.PostsGrpcClient, log *slog.Logger) *PostService {
	return &PostService{postsGRpc, log}
}

func (p *PostService) GetPost(ctx context.Context, dto *postsv1.GetPostRequest) (*postsv1.GetPostResponse, error) {
	op := "PostService.GetPosts"
	log := p.log.With(
		slog.String("op", op))

	res, err := p.postsGRpc.PostsApi.GetPost(ctx, dto)
	if err != nil {
		log.Error("can`t get post from grpc:", err)
	}

	return res, nil
}

func (p *PostService) CreatePost(ctx context.Context, dto *postsv1.CreatePostRequest) (*postsv1.CreatePostResponse, error) {
	op := "PostService.CreatePost"
	log := p.log.With(
		slog.String("op", op))

	res, err := p.postsGRpc.PostsApi.CreatePost(ctx, dto)
	if err != nil {
		log.Error("can`t create post:", err)
	}

	return res, nil
}

func (p *PostService) UpdatePost(ctx context.Context, dto *postsv1.UpdatePostRequest) (*postsv1.UpdatePostResponse, error) {
	op := "PostService.UpdatePost"
	log := p.log.With(
		slog.String("op", op))

	res, err := p.postsGRpc.PostsApi.UpdatePost(ctx, dto)
	if err != nil {
		log.Error("can`t update post:", err)
	}

	return res, nil
}

func (p *PostService) DeletePost(ctx context.Context, dto *postsv1.DeletePostRequest) (*postsv1.DeletePostResponse, error) {
	op := "PostService.DeletePost"
	log := p.log.With(
		slog.String("op", op))

	res, err := p.postsGRpc.PostsApi.DeletePost(ctx, dto)
	if err != nil {
		log.Error("can`t delete post:", err)
	}

	return res, nil
}

func (p *PostService) GetUserPosts(ctx context.Context, dto *postsv1.GetUserPostsRequest) (*postsv1.GetUserPostsResponse, error) {
	op := "PostService.GetUserPosts"
	log := p.log.With(
		slog.String("op", op))

	res, err := p.postsGRpc.PostsApi.GetUserPosts(ctx, dto)
	if err != nil {
		log.Error("can`t get user posts:", err)
	}

	return res, nil
}
