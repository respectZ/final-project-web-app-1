package api

import (
	"a21hc3NpZ25tZW50/entity"
	"a21hc3NpZ25tZW50/service"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
)

type TaskAPI interface {
	GetTask(w http.ResponseWriter, r *http.Request)
	CreateNewTask(w http.ResponseWriter, r *http.Request)
	UpdateTask(w http.ResponseWriter, r *http.Request)
	DeleteTask(w http.ResponseWriter, r *http.Request)
	UpdateTaskCategory(w http.ResponseWriter, r *http.Request)
}

type taskAPI struct {
	taskService service.TaskService
}

func NewTaskAPI(taskService service.TaskService) *taskAPI {
	return &taskAPI{taskService}
}

func (t *taskAPI) GetTask(w http.ResponseWriter, r *http.Request) {
	taskID, err := strconv.Atoi(r.URL.Query().Get("task_id"))
	if err != nil {
		taskID = -1
	}

	userIDT := r.Context().Value("id").(string)

	if userIDT == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	userID, _ := strconv.Atoi(userIDT)
	tasks := []entity.Task{}
	switch taskID {
	case -1:
		tasks, err = t.taskService.GetTasks(r.Context(), userID)
	default:
		task, _ := t.taskService.GetTaskByID(r.Context(), taskID)
		tasks = append(tasks, task)
	}
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("[taskAPI][GetTask]" + err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}
	switch taskID {
	case -1:
		json.NewEncoder(w).Encode(tasks)
	default:
		json.NewEncoder(w).Encode(tasks[0])
	}
	// TODO: answer here
}

func (t *taskAPI) CreateNewTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}

	if task.Title == "" || task.Description == "" || task.CategoryID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid task request"))
		return
	}

	userIDT := r.Context().Value("id").(string)

	if userIDT == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	userID, _ := strconv.Atoi(userIDT)
	taskT := entity.Task{Title: task.Title, Description: task.Description, CategoryID: task.CategoryID, UserID: userID}

	taskNew, err := t.taskService.StoreTask(r.Context(), &taskT)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("[taskAPI][CreateNewTask]" + err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	res := map[string]interface{}{}
	res["user_id"] = taskNew.UserID
	res["task_id"] = taskNew.ID
	res["message"] = "success create new task"

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(res)

	// TODO: answer here
}

func (t *taskAPI) DeleteTask(w http.ResponseWriter, r *http.Request) {
	userIDT := r.Context().Value("id").(string)

	if userIDT == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	userID, _ := strconv.Atoi(userIDT)
	taskIDT := r.URL.Query().Get("task_id")
	if taskIDT == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("empty task id"))
		return
	}
	taskID, _ := strconv.Atoi(taskIDT)

	err := t.taskService.DeleteTask(r.Context(), taskID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("[categoryAPI][DeleteCategory]" + err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	res := map[string]interface{}{}
	res["user_id"] = userID
	res["task_id"] = taskID
	res["message"] = "success delete task"

	json.NewEncoder(w).Encode(res)
	// TODO: answer here
}

func (t *taskAPI) UpdateTask(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	userIDT := r.Context().Value("id").(string)

	if userIDT == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	userID, _ := strconv.Atoi(userIDT)
	taskT := entity.Task{ID: task.ID, Title: task.Title, Description: task.Description, CategoryID: task.CategoryID, UserID: userID}

	taskNew, err := t.taskService.UpdateTask(r.Context(), &taskT)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println("[taskAPI][UpdateTask]" + err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	res := map[string]interface{}{}
	res["user_id"] = taskNew.UserID
	res["task_id"] = taskNew.ID
	res["message"] = "success update task"

	json.NewEncoder(w).Encode(res)

	// TODO: answer here
}

func (t *taskAPI) UpdateTaskCategory(w http.ResponseWriter, r *http.Request) {
	var task entity.TaskCategoryRequest

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid decode json"))
		return
	}

	userId := r.Context().Value("id")

	idLogin, err := strconv.Atoi(userId.(string))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("invalid user id"))
		return
	}

	var updateTask = entity.Task{
		ID:         task.ID,
		CategoryID: task.CategoryID,
		UserID:     int(idLogin),
	}

	_, err = t.taskService.UpdateTask(r.Context(), &updateTask)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Println(err.Error())
		json.NewEncoder(w).Encode(entity.NewErrorResponse("error internal server"))
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"user_id": userId,
		"task_id": task.ID,
		"message": "success update task category",
	})
}
