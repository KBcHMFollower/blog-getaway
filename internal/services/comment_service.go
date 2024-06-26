package services

import (
	"context"
	commentsv1 "github.com/KBcHMFollower/blog_posts_service/api/protos/gen/comments"
	"log/slog"
	postdepend "test-plate/internal/dependencies/postservice"
)

type CommentService struct {
	log          *slog.Logger
	commentsGRpc *postdepend.PostsGrpcClient
}

func NewCommentService(postGRpc *postdepend.PostsGrpcClient, log *slog.Logger) *CommentService {
	return &CommentService{
		log:          log,
		commentsGRpc: postGRpc,
	}
}

func (s *CommentService) GetComment(ctx context.Context, dto *commentsv1.GetCommentRequest) (*commentsv1.GetCommentResponse, error) {
	op := "CommentService.GetComments"
	log := s.log.With(
		slog.String("op", op))

	res, err := s.commentsGRpc.CommentsApi.GetComment(ctx, dto)
	if err != nil {
		log.Error("can`t get comment from grpc:", err)
	}

	return res, nil
}

func (s *CommentService) CreateComment(ctx context.Context, dto *commentsv1.CreateCommentRequest) (*commentsv1.CreateCommentResponse, error) {
	op := "CommentService.CreateComment"
	log := s.log.With(
		slog.String("op", op))

	res, err := s.commentsGRpc.CommentsApi.CreateComment(ctx, dto)
	if err != nil {
		log.Error("can`t create comment:", err)
	}

	return res, nil
}

func (s *CommentService) UpdateComment(ctx context.Context, dto *commentsv1.UpdateCommentRequest) (*commentsv1.UpdateCommentResponse, error) {
	op := "CommentService.UpdateComment"
	log := s.log.With(
		slog.String("op", op))

	res, err := s.commentsGRpc.CommentsApi.UpdateComment(ctx, dto)
	if err != nil {
		log.Error("can`t update comment:", err)
	}

	return res, nil
}

func (s *CommentService) DeleteComment(ctx context.Context, dto *commentsv1.DeleteCommentRequest) (*commentsv1.DeleteCommentResponse, error) {
	op := "CommentService.DeleteComment"
	log := s.log.With(
		slog.String("op", op))

	res, err := s.commentsGRpc.CommentsApi.DeleteComment(ctx, dto)
	if err != nil {
		log.Error("can`t delete comment:", err)
	}

	return res, nil
}

func (s *CommentService) GetPostComments(ctx context.Context, dto *commentsv1.GetPostCommentsRequest) (*commentsv1.GetPostCommentsResponse, error) {
	op := "CommentService.GetUserPosts"
	log := s.log.With(
		slog.String("op", op))

	res, err := s.commentsGRpc.CommentsApi.GetPostComments(ctx, dto)
	if err != nil {
		log.Error("can`t get post comments:", err)
	}

	return res, nil
}
