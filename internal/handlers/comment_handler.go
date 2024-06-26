package handlers

import (
	"encoding/json"
	commentsv1 "github.com/KBcHMFollower/blog_posts_service/api/protos/gen/comments"
	"github.com/go-chi/chi/v5"
	"github.com/unrolled/render"
	"net/http"
	"strconv"
	"test-plate/internal/domain/models"
	"test-plate/internal/services"
)

type CommentDeleteResponse struct {
	IsDeleted bool `json:"is_deleted"`
}

type CommentsGetResponse struct {
	TotalCount int32             `json:"total_count"`
	Comments   []*models.Comment `json:"comments"`
}

type UpdateCommentItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type UpdateCommentBody struct {
	UpdateItems []UpdatePostItem `json:"update_items"`
}

type CreateCommentBody struct {
	Content string `json:"content"`
	UserId  string `json:"user_id"`
}

type CommentHandler struct {
	commService *services.CommentService
	rnd         *render.Render
}

func NewCommentHandler(commService *services.CommentService, rnd *render.Render) *CommentHandler {
	return &CommentHandler{commService, rnd}
}

func (ch *CommentHandler) GetComment(w http.ResponseWriter, r *http.Request) {
	commentId := chi.URLParam(r, "commentId")
	if commentId == "" {
		ch.rnd.Text(w, http.StatusBadRequest, "Unable to parse commentId")
		return
	}

	res, err := ch.commService.GetComment(r.Context(), &commentsv1.GetCommentRequest{Id: commentId})
	if err != nil {
		ch.rnd.Text(w, http.StatusInternalServerError, err.Error())
		return
	}

	ch.rnd.JSON(w, http.StatusOK, models.ConvertCommentFromProto(res.GetComments()))
}

func (ch *CommentHandler) DeleteComment(w http.ResponseWriter, r *http.Request) {
	commentId := chi.URLParam(r, "commentId")
	if commentId == "" {
		ch.rnd.Text(w, http.StatusBadRequest, "Unable to parse commentId")
		return
	}

	ok, err := ch.commService.DeleteComment(r.Context(), &commentsv1.DeleteCommentRequest{Id: commentId})
	if err != nil {
		ch.rnd.Text(w, http.StatusInternalServerError, err.Error())
		return
	}

	ch.rnd.JSON(w, http.StatusOK, CommentDeleteResponse{IsDeleted: ok.IsDeleted})
}

func (ch *CommentHandler) GetPostComments(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")
	if postId == "" {
		ch.rnd.Text(w, http.StatusBadRequest, "Unable to parse postId")
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

	res, err := ch.commService.GetPostComments(r.Context(), &commentsv1.GetPostCommentsRequest{PostId: postId, Page: int32(page), Size: int32(size)})
	if err != nil {
		ch.rnd.Text(w, http.StatusInternalServerError, err.Error())
		return
	}

	ch.rnd.JSON(w, http.StatusOK, CommentsGetResponse{
		TotalCount: res.TotalCount,
		Comments:   models.CommentsArrayFromProto(res.Comments),
	})
}

func (ch *CommentHandler) CreateComment(w http.ResponseWriter, r *http.Request) {
	postId := chi.URLParam(r, "postId")
	if postId == "" {
		ch.rnd.Text(w, http.StatusBadRequest, "Unable to parse postId")
		return
	}

	var createData CreateCommentBody

	if err := json.NewDecoder(r.Body).Decode(&createData); err != nil {
		ch.rnd.Text(w, http.StatusBadRequest, "Unable to parse body")
		return
	}

	res, err := ch.commService.CreateComment(r.Context(), &commentsv1.CreateCommentRequest{
		UserId:  createData.UserId,
		Content: createData.Content,
		PostId:  postId,
	})
	if err != nil {
		ch.rnd.Text(w, http.StatusInternalServerError, err.Error())
		return
	}

	ch.rnd.JSON(w, http.StatusOK, models.ConvertCommentFromProto(res.Comment))
}

func (ch *CommentHandler) UpdateComment(w http.ResponseWriter, r *http.Request) {
	commId := chi.URLParam(r, "commentId")
	if commId == "" {
		ch.rnd.Text(w, http.StatusBadRequest, "Unable to parse commId")
		return
	}

	var updateBody UpdateCommentBody

	if err := json.NewDecoder(r.Body).Decode(&updateBody); err != nil {
		ch.rnd.Text(w, http.StatusBadRequest, "Unable to parse body")
		return
	}

	var protoUpdateItems = make([]*commentsv1.CommUpdateItem, 0, len(updateBody.UpdateItems))
	for _, item := range updateBody.UpdateItems {
		protoUpdateItems = append(protoUpdateItems, &commentsv1.CommUpdateItem{
			Name:  item.Name,
			Value: item.Value,
		})
	}

	res, err := ch.commService.UpdateComment(r.Context(), &commentsv1.UpdateCommentRequest{
		Id:         commId,
		UpdateData: protoUpdateItems,
	})
	if err != nil {
		ch.rnd.Text(w, http.StatusInternalServerError, err.Error())
		return
	}

	ch.rnd.JSON(w, http.StatusOK, models.ConvertCommentFromProto(res.Comment))
}
