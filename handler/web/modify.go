package web

import (
	"a21hc3NpZ25tZW50/client"
	"embed"
	"net/http"
	"path"
	"text/template"
)

type ModifyWeb interface {
	AddTask(w http.ResponseWriter, r *http.Request)
	AddTaskProcess(w http.ResponseWriter, r *http.Request)
	AddCategory(w http.ResponseWriter, r *http.Request)
	AddCategoryProcess(w http.ResponseWriter, r *http.Request)

	UpdateTask(w http.ResponseWriter, r *http.Request)
	UpdateTaskProcess(w http.ResponseWriter, r *http.Request)

	DeleteTask(w http.ResponseWriter, r *http.Request)
	DeleteCategory(w http.ResponseWriter, r *http.Request)
}

type modifyWeb struct {
	taskClient     client.TaskClient
	categoryClient client.CategoryClient
	embed          embed.FS
}

func NewModifyWeb(tC client.TaskClient, cC client.CategoryClient, embed embed.FS) *modifyWeb {
	return &modifyWeb{tC, cC, embed}
}

func (a *modifyWeb) AddTask(w http.ResponseWriter, r *http.Request) {
	catId := r.URL.Query().Get("category")

	// ignore this
	_ = catId
	//
	var filepath = path.Join("views", "main", "add-task.html")
	var header = path.Join("views", "general", "header.html")

	var tmpl = template.Must(template.ParseFS(a.embed, filepath, header))

	err := tmpl.Execute(w, map[string]string{"category": catId})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: answer here
}

func (a *modifyWeb) AddTaskProcess(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	title := r.FormValue("title")
	description := r.FormValue("description")
	category := r.URL.Query().Get("category")

	respCode, err := a.taskClient.CreateTask(title, description, category, userId.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ignore this
	_ = respCode
	//
	switch respCode {
	case http.StatusCreated:
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	case http.StatusOK:
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	default:
		http.Redirect(w, r, "/task/add", http.StatusSeeOther)
	}

	// TODO: answer here
}

func (a *modifyWeb) AddCategory(w http.ResponseWriter, r *http.Request) {
	var filepath = path.Join("views", "main", "add-category.html")
	var header = path.Join("views", "general", "header.html")

	var tmpl = template.Must(template.ParseFS(a.embed, filepath, header))

	err := tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (a *modifyWeb) AddCategoryProcess(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	category := r.FormValue("type")

	respCode, err := a.categoryClient.AddCategories(category, userId.(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ignore this
	_ = respCode
	//
	switch respCode {
	case http.StatusCreated:
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	default:
		http.Redirect(w, r, "/category/add", http.StatusSeeOther)
	}

	// TODO: answer here
}

func (a *modifyWeb) UpdateTask(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Query().Get("task_id")

	task, err := a.taskClient.GetTaskById(taskId, r.Context().Value("id").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// ignore this
	_ = task
	var filepath = path.Join("views", "main", "update-task.html")
	var header = path.Join("views", "general", "header.html")

	var tmpl = template.Must(template.ParseFS(a.embed, filepath, header))

	data := map[string]interface{}{
		"category_id": task.CategoryID,
		"task_id":     task.ID,
		"title":       task.Title,
		"description": task.Description,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: answer here
}

func (a *modifyWeb) UpdateTaskProcess(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Query().Get("task_id")
	categoryId := r.URL.Query().Get("category_id")

	if categoryId == "" {
		title := r.FormValue("title")
		description := r.FormValue("description")

		respCode, err := a.taskClient.UpdateTask(taskId, title, description, r.Context().Value("id").(string))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if respCode == 200 {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		} else {
			http.Redirect(w, r, "/task/update?task_id="+taskId, http.StatusSeeOther)
		}
	} else {
		_, err := a.taskClient.UpdateCategoryTask(taskId, categoryId, r.Context().Value("id").(string))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}

}

func (a *modifyWeb) DeleteTask(w http.ResponseWriter, r *http.Request) {
	taskId := r.URL.Query().Get("task_id")

	respCode, err := a.taskClient.DeleteTask(taskId, r.Context().Value("id").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch respCode {
	case http.StatusCreated:
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	case http.StatusOK:
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	default:
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}

	// TODO: answer here
}

func (a *modifyWeb) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	categoryId := r.URL.Query().Get("category_id")

	respCode, err := a.categoryClient.DeleteCategory(categoryId, r.Context().Value("id").(string))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	switch respCode {
	case http.StatusCreated:
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	case http.StatusOK:
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	default:
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
	}

	// TODO: answer here
}
