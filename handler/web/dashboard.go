package web

import (
	"a21hc3NpZ25tZW50/client"
	"embed"
	"fmt"
	"log"
	"net/http"
	"path"
	"text/template"
)

type DashboardWeb interface {
	Dashboard(w http.ResponseWriter, r *http.Request)
}

type dashboardWeb struct {
	categoryClient client.CategoryClient
	embed          embed.FS
}

func NewDashboardWeb(catClient client.CategoryClient, embed embed.FS) *dashboardWeb {
	return &dashboardWeb{catClient, embed}
}

func (d *dashboardWeb) Dashboard(w http.ResponseWriter, r *http.Request) {
	userId := r.Context().Value("id")

	categories, err := d.categoryClient.GetCategories(userId.(string))
	if err != nil {
		log.Println("error get cat: ", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var dataTemplate = map[string]interface{}{
		"categories": categories,
	}

	var funcMap = template.FuncMap{
		"categoryInc": func(catId int) int {
			return catId + 1
		},
		"categoryDec": func(catId int) int {
			return catId - 1
		},
		"write": func() string {
			mainTemplate := `
			<!-- cat -->
            <div
                class="h-fit flex flex-col flex-1 items-center bg-gradient-to-b from-%s-400 to-transparent pb-12 rounded-md mx-4">
                <div class="w-full h-fit flex flex-row flex-1 items-center pb-8 rounded-md mx-4">
                    <h2 class="flex-1 text-xl font-semibold px-4 text-black">%s</h2>
                    <a href="/task/add?category=%d"> <button class="text-3xl">+</button> </a>
                    <div class="w-2"></div>
                    <a href="/category/delete?category_id=%d"> <button class="text-3xl">x</button> </a>
                    <div class="w-4"></div>
                </div>
                <!-- content -->
                <div class="flex flex-col items-center w-full">
					%s
                </div>
            </div>
			`
			taskTemplate := `
			<div class="flex flex-col h-98 w-full px-4 shadow-lg mb-8">
                        <div class="flex flex-row items-center bg-%s-500 rounded-t-md text-white p-3">
                            <h1 class="text-2xl font-bold flex-1">%s</h1>
							<a href="/task/update?task_id=%d">
                                <button type="button" class="px-2">
                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                                        stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                                        <path stroke-linecap="round" stroke-linejoin="round"
                                            d="M16.862 4.487l1.687-1.688a1.875 1.875 0 112.652 2.652L10.582 16.07a4.5 4.5 0 01-1.897 1.13L6 18l.8-2.685a4.5 4.5 0 011.13-1.897l8.932-8.931zm0 0L19.5 7.125M18 14v4.75A2.25 2.25 0 0115.75 21H5.25A2.25 2.25 0 013 18.75V8.25A2.25 2.25 0 015.25 6H10" />
                                    </svg>
                                </button>
                            </a>

                            <a href="/task/delete?task_id=%d">
                                <button type="button">
                                    <svg xmlns="http://www.w3.org/2000/svg" fill="none" viewBox="0 0 24 24"
                                        stroke-width="1.5" stroke="currentColor" class="w-6 h-6">
                                        <path stroke-linecap="round" stroke-linejoin="round"
                                            d="M14.74 9l-.346 9m-4.788 0L9.26 9m9.968-3.21c.342.052.682.107 1.022.166m-1.022-.165L18.16 19.673a2.25 2.25 0 01-2.244 2.077H8.084a2.25 2.25 0 01-2.244-2.077L4.772 5.79m14.456 0a48.108 48.108 0 00-3.478-.397m-12 .562c.34-.059.68-.114 1.022-.165m0 0a48.11 48.11 0 013.478-.397m7.5 0v-.916c0-1.18-.91-2.164-2.09-2.201a51.964 51.964 0 00-3.32 0c-1.18.037-2.09 1.022-2.09 2.201v.916m7.5 0a48.667 48.667 0 00-7.5 0" />
                                    </svg>
                                </button>
                            </a>
                        </div>
                        <div class="flex bg-zinc-100 rounded-b-md text-black p-4">
                            <p>%s</p>
                        </div>
                    </div>`
			addCategoryData := `
			<div
                class="h-fit flex flex-col flex-0 items-center">
                <div class="w-full h-fit flex flex-row flex-1 items-center pb-8 rounded-md mx-4">
					<a href="/category/add">
						<button type="button"
							class="bg-blue-600 hover:bg-blue-900 py-1.5 px-4 rounded-lg">Add Category</button>
					</a>
                </div>
            </div>
			`

			resData := ""
			for _, category := range categories {
				tasksData := ""
				mainColor := "violet"
				switch category.Type {
				case "Todo":
					mainColor = "orange"
				case "In Progress":
					mainColor = "yellow"
				case "Done":
					mainColor = "green"
				case "Backlog":
					mainColor = "blue"
				}
				// Add task
				taskColor := "red"
				idx := 0
				for _, task := range category.Tasks {
					if idx&1 == 0 {
						taskColor = "blue"
					} else {
						taskColor = "red"
					}
					tasksData += fmt.Sprintf(taskTemplate, taskColor, task.Title, task.ID, task.ID, task.Description)
					idx++
				}
				resData += fmt.Sprintf(mainTemplate, mainColor, category.Type, category.ID, category.ID, tasksData)
			}

			// add Category

			return resData + addCategoryData
		},
	}

	// ignore this
	_ = dataTemplate
	_ = funcMap
	//

	var filepath = path.Join("views", "main", "dashboard.html")
	var header = path.Join("views", "general", "header.html")

	// var tmpl = template.Must(template.ParseFS(d.embed, filepath, header)).Funcs(funcMap)
	var tmpl = template.Must(template.New("dashboard.html").Funcs(funcMap).ParseFS(d.embed, filepath, header))

	// append data
	dataTemplate["user"] = map[string]string{
		"name": "Maria Nearl",
	}

	err = tmpl.Execute(w, dataTemplate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// TODO: answer here
}
