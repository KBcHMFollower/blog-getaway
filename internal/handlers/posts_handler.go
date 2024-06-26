package handlers

import (
	"encoding/json"
	postsv1 "github.com/KBcHMFollower/blog_posts_service/api/protos/gen/posts"
	"github.com/go-chi/chi/v5"
	"github.com/unrolled/render"
	"net/http"
	"strconv"
	"test-plate/internal/domain/models"
	"test-plate/internal/services"
)

type PostDeleteResponse struct {
	IsDeleted bool `json:"is_deleted"`
}

type PostsGetResponse struct {
	TotalCount int32          `json:"total_count"`
	Posts      []*models.Post `json:"posts"`
}

type UpdatePostItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type UpdatePostBody struct {
	UpdateItems []UpdatePostItem `json:"update_items"`
}

type CreatePostBody struct {
	Title         string  `json:"title"`
	TextContent   string  `json:"text_content"`
	imagesContent *string `json:"images_content"`
}

type PostHandler struct {
	postService *services.PostService
	rnd         *render.Render
}

func NewPostHandler(postService *services.PostService, rnd *render.Render) *PostHandler {
	return &PostHandler{postService, rnd}
}

func (ph *PostHandler) GetPost(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")
	if postId == "" {
		ph.rnd.Text(w, http.StatusBadRequest, "Unable to parse postId")
		return
	}

	post, err := ph.postService.GetPost(r.Context(), &postsv1.GetPostRequest{Id: postId})
	if err != nil {
		ph.rnd.Text(w, http.StatusInternalServerError, err.Error())
		return
	}

	ph.rnd.JSON(w, http.StatusOK, models.ConvertPostFromProto(post.GetPosts()))
}

func (ph *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")
	if postId == "" {
		ph.rnd.Text(w, http.StatusBadRequest, "Unable to parse postId")
		return
	}

	ok, err := ph.postService.DeletePost(r.Context(), &postsv1.DeletePostRequest{Id: postId})
	if err != nil {
		ph.rnd.Text(w, http.StatusInternalServerError, err.Error())
		return
	}

	ph.rnd.JSON(w, http.StatusOK, PostDeleteResponse{IsDeleted: ok.IsDeleted})
}

func (ph *PostHandler) GetUserPosts(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	if userId == "" {
		ph.rnd.Text(w, http.StatusBadRequest, "Unable to parse userId")
		return
	}

	query := r.URL.Query()

	pageStr := query.Get("page")
	page, err := strconv.Atoi(pageStr)
	if err != nil {
		page = 1
	}

	sizeStr := query.Get("size")
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		size = 10
	}

	res, err := ph.postService.GetUserPosts(r.Context(), &postsv1.GetUserPostsRequest{UserId: userId, Page: int32(page), Size: int32(size)})
	if err != nil {
		ph.rnd.Text(w, http.StatusInternalServerError, err.Error())
		return
	}

	ph.rnd.JSON(w, http.StatusOK, PostsGetResponse{
		TotalCount: res.TotalCount,
		Posts:      models.PostArrayFromProto(res.Posts),
	})
}

func (ph *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	if userId == "" {
		ph.rnd.Text(w, http.StatusBadRequest, "Unable to parse userId")
		return
	}

	var createData CreatePostBody

	if err := json.NewDecoder(r.Body).Decode(&createData); err != nil {
		ph.rnd.Text(w, http.StatusBadRequest, "Unable to parse body")
		return
	}

	res, err := ph.postService.CreatePost(r.Context(), &postsv1.CreatePostRequest{
		UserId:        userId,
		Title:         createData.Title,
		TextContent:   createData.TextContent,
		ImagesContent: createData.imagesContent,
	})
	if err != nil {
		ph.rnd.Text(w, http.StatusInternalServerError, err.Error())
		return
	}

	ph.rnd.JSON(w, http.StatusOK, models.ConvertPostFromProto(res.Post))
}

func (ph *PostHandler) UpdatePost(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")
	if postId == "" {
		ph.rnd.Text(w, http.StatusBadRequest, "Unable to parse postId")
		return
	}

	var updateBody UpdatePostBody

	if err := json.NewDecoder(r.Body).Decode(&updateBody); err != nil {
		ph.rnd.Text(w, http.StatusBadRequest, "Unable to parse body")
		return
	}

	var protoUpdateItems = make([]*postsv1.PostUpdateItem, 0, len(updateBody.UpdateItems))
	for _, item := range updateBody.UpdateItems {
		protoUpdateItems = append(protoUpdateItems, &postsv1.PostUpdateItem{
			Name:  item.Name,
			Value: item.Value,
		})
	}

	res, err := ph.postService.UpdatePost(r.Context(), &postsv1.UpdatePostRequest{
		Id:         postId,
		UpdateData: protoUpdateItems,
	})
	if err != nil {
		ph.rnd.Text(w, http.StatusInternalServerError, err.Error())
		return
	}

	ph.rnd.JSON(w, http.StatusOK, models.ConvertPostFromProto(res.Post))
}
