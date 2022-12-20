package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type CategoryAPI interface {
	GetCategory(w http.ResponseWriter, r *http.Request)
	CreateNewCategory(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
	GetCategoryWithTasks(w http.ResponseWriter, r *http.Request)
}

type categoryAPI struct {
	categoryService service.CategoryService
}

func NewCategoryAPI(categoryService service.CategoryService) *categoryAPI {
	return &categoryAPI{categoryService}
}

func (c *categoryAPI) GetCategory(w http.ResponseWriter, r *http.Request) {
	userIDT := r.Context().Value("id")

	if userIDT == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	userID, err := strconv.Atoi(userIDT.(string))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	categories, err := c.categoryService.GetCategories(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("[categoryAPI][GetCategory]" + err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}
	json.NewEncoder(w).Encode(categories)
	// TODO: answer here
}

func (c *categoryAPI) CreateNewCategory(w http.ResponseWriter, r *http.Request) {
	var category entity.CategoryRequest

	err := json.NewDecoder(r.Body).Decode(&category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	if category.Type == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid category request"))
		return
	}

	userIDT := r.Context().Value("id").(string)

	if userIDT == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	userID, _ := strconv.Atoi(userIDT)
	categoryType := entity.Category{UserID: userID, Type: category.Type}

	categoryNew, err := c.categoryService.StoreCategory(r.Context(), &categoryType)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("[categoryAPI][CreateNewCategory]" + err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	res := map[string]interface{}{}
	res["user_id"] = categoryNew.UserID
	res["category_id"] = categoryNew.ID
	res["message"] = "success create new category"

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)

	// TODO: answer here
}

func (c *categoryAPI) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	userIDT := r.Context().Value("id").(string)

	if userIDT == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	userID, _ := strconv.Atoi(userIDT)
	categoryIDT := r.URL.Query().Get("category_id")
	if categoryIDT == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("empty category id"))
		return
	}
	categoryID, _ := strconv.Atoi(categoryIDT)

	err := c.categoryService.DeleteCategory(r.Context(), categoryID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("[categoryAPI][DeleteCategory]" + err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	res := map[string]interface{}{}
	res["user_id"] = userID
	res["category_id"] = categoryID
	res["message"] = "success delete category"

	json.NewEncoder(w).Encode(res)
	// TODO: answer here
}

func (c *categoryAPI) GetCategoryWithTasks(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println("get category task", err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	categories, err := c.categoryService.GetCategoriesWithTasks(r.Context(), int(idLogin))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("internal server error"))
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(categories)

}
