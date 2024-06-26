package postdepend

import (
	commentsv1 "github.com/KBcHMFollower/blog_posts_service/api/protos/gen/comments"
	postsv1 "github.com/KBcHMFollower/blog_posts_service/api/protos/gen/posts"
	"google.golang.org/grpc"
)

type PostsGrpcClient struct {
	PostsApi    postsv1.PostServiceClient
	CommentsApi commentsv1.CommentServiceClient
}

func NewPostsClient(addr string) (*PostsGrpcClient, error) {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		return nil, err
	}

	return &PostsGrpcClient{
		PostsApi:    postsv1.NewPostServiceClient(conn),
		CommentsApi: commentsv1.NewCommentServiceClient(conn),
	}, nil
}
