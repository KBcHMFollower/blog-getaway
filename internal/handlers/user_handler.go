package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	usersv1 "github.com/KBcHMFollower/blog_user_service/api/protos/gen/users"
	"github.com/go-chi/chi/v5"
	"github.com/unrolled/render"
	"io"
	"net/http"
	"strconv"
	"test-plate/internal/domain/models"
	"test-plate/internal/services"
)

type UserDeleteResponse struct {
	IsDeleted bool `json:"is_deleted"`
}

type UserSubscribeResponse struct {
	Ok bool `json:"ok"`
}

type UsersGetResponse struct {
	Users      []*models.User `json:"users"`
	TotalCount int32          `json:"total_count"`
}

type AvatarsResponse struct {
	Avatar     string `json:"avatar"`
	AvatarMini string `json:"avatar_mini"`
}

type UpdateUserBody struct {
	UpdateItems []*UpdateUserItem `json:"update_items"`
}

type UpdateUserItem struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type UsersHandler struct {
	userService *services.UserService
	rnd         *render.Render
	baseUrl     string
}

func NewUsersHandler(userService *services.UserService, rnd *render.Render, baseUrl string) *UsersHandler {
	return &UsersHandler{userService, rnd, baseUrl}
}

func (uh *UsersHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	if userId == "" {
		uh.rnd.Text(w, http.StatusBadRequest, "Unable to parse form")
		return
	}

	user, err := uh.userService.GetUser(r.Context(), userId)
	if err != nil {
		uh.rnd.Text(w, http.StatusInternalServerError, err.Error())
		return
	}

	uh.rnd.JSON(w, http.StatusOK, models.ConvertUserFromProto(user))
}

func (uh *UsersHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	if userId == "" {
		uh.rnd.Text(w, http.StatusBadRequest, "Unable to parse form")
		return
	}

	res, err := uh.userService.DeleteUser(r.Context(), userId)
	if err != nil {
		uh.rnd.Text(w, http.StatusInternalServerError, err.Error())
		return
	}

	uh.rnd.JSON(w, http.StatusOK, &UserDeleteResponse{IsDeleted: res})
}

func (uh *UsersHandler) GetSubscribers(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	if userId == "" {
		uh.rnd.Text(w, http.StatusBadRequest, "Unable to parse form")
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

	res, err := uh.userService.GetSubscribers(r.Context(), &usersv1.GetSubscribersDTO{
		BloggerId: userId,
		Page:      int32(page),
		Size:      int32(size),
	})
	if err != nil {
		uh.rnd.Text(w, http.StatusInternalServerError, err.Error())
		return
	}

	uh.rnd.JSON(w, http.StatusOK, &UsersGetResponse{
		TotalCount: res.TotalCount,
		Users:      models.UsersArrayFromProto(res.Subscribers),
	})
}

func (uh *UsersHandler) GetSubscriptions(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	if userId == "" {
		uh.rnd.Text(w, http.StatusBadRequest, "Unable to parse form")
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

	res, err := uh.userService.GetSubscriptions(r.Context(), &usersv1.GetSubscriptionsDTO{
		SubscriberId: userId,
		Page:         int32(page),
		Size:         int32(size),
	})
	if err != nil {
		uh.rnd.Text(w, http.StatusInternalServerError, err.Error())
		return
	}

	uh.rnd.JSON(w, http.StatusOK, &UsersGetResponse{
		TotalCount: res.TotalCount,
		Users:      models.UsersArrayFromProto(res.Subscriptions),
	})
}

func (uh *UsersHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	if userId == "" {
		uh.rnd.Text(w, http.StatusBadRequest, "Unable to parse form")
		return
	}

	var updateData UpdateUserBody

	if err := json.NewDecoder(r.Body).Decode(&updateData); err != nil {
		uh.rnd.Text(w, http.StatusBadRequest, "Unable to parse body")
		return
	}

	var updateFieldsProto []*usersv1.UpdateUserItem

	for _, updateItem := range updateData.UpdateItems {
		updateFieldsProto = append(updateFieldsProto, &usersv1.UpdateUserItem{
			Name:  updateItem.Name,
			Value: updateItem.Value,
		})
	}

	res, err := uh.userService.UpdateUser(r.Context(), &usersv1.UpdateUserDTO{
		Id:         userId,
		UpdateData: updateFieldsProto,
	})
	if err != nil {
		uh.rnd.Text(w, http.StatusInternalServerError, err.Error())
		return
	}

	uh.rnd.JSON(w, http.StatusOK, models.ConvertUserFromProto(res.User))
}

func (uh *UsersHandler) UploadAvatar(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	if userId == "" {
		uh.rnd.Text(w, http.StatusBadRequest, "Unable to parse form")
		return
	}

	r.ParseMultipartForm(10 << 20)

	file, _, err := r.FormFile("image")
	if err != nil {
		uh.rnd.Text(w, http.StatusBadRequest, "Unable to parse form-file")
		return
	}
	defer file.Close()

	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, file)
	if err != nil {
		uh.rnd.Text(w, http.StatusBadRequest, "Unable to copy image")
		return
	}

	res, err := uh.userService.UploadAvatar(r.Context(), &usersv1.UploadAvatarDTO{
		UserId:  userId,
		Image:   buffer.Bytes(),
		BaseUrl: uh.baseUrl + "/api/images/",
	})
	if err != nil {
		uh.rnd.Text(w, http.StatusInternalServerError, err.Error())
		return
	}

	uh.rnd.JSON(w, http.StatusOK, AvatarsResponse{
		Avatar:     res.AvatarUrl,
		AvatarMini: res.AvatarMiniUrl,
	})
}

func (uh *UsersHandler) GetAvatar(w http.ResponseWriter, r *http.Request) {
	imgName := chi.URLParam(r, "imgName")
	if imgName == "" {
		uh.rnd.Text(w, http.StatusBadRequest, "Unable to parse form")
		return
	}

	fmt.Println(imgName)

	res, err := uh.userService.GetAvatar(r.Context(), imgName+".jpeg")
	if err != nil {
		uh.rnd.Text(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.Header().Set("Content-Type", "image/png")

	if _, err := w.Write(res.GetImage()); err != nil {
		uh.rnd.Text(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (uh *UsersHandler) Subscribe(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	if userId == "" {
		uh.rnd.Text(w, http.StatusBadRequest, "Unable to parse userId")
		return
	}

	bloggerId := r.URL.Query().Get("bloggerId")
	if bloggerId == "" {
		uh.rnd.Text(w, http.StatusBadRequest, "Unable to parse bloggerId")
		return
	}

	ok, err := uh.userService.Subscribe(r.Context(), &usersv1.SubscribeDTO{
		SubscriberId: userId,
		BloggerId:    bloggerId,
	})
	if err != nil {
		uh.rnd.Text(w, http.StatusInternalServerError, err.Error())
		return
	}

	uh.rnd.JSON(w, http.StatusOK, UserSubscribeResponse{
		Ok: ok,
	})
}

func (uh *UsersHandler) Unsubscribe(w http.ResponseWriter, r *http.Request) {
	userId := chi.URLParam(r, "userId")
	if userId == "" {
		uh.rnd.Text(w, http.StatusBadRequest, "Unable to parse form")
		return
	}

	bloggerId := r.URL.Query().Get("bloggerId")
	if bloggerId == "" {
		uh.rnd.Text(w, http.StatusBadRequest, "Unable to parse bloggerId")
		return
	}

	ok, err := uh.userService.Unsubscribe(r.Context(), &usersv1.SubscribeDTO{
		SubscriberId: userId,
		BloggerId:    bloggerId,
	})
	if err != nil {
		uh.rnd.Text(w, http.StatusInternalServerError, err.Error())
		return
	}

	uh.rnd.JSON(w, http.StatusOK, UserSubscribeResponse{
		Ok: ok,
	})
}
