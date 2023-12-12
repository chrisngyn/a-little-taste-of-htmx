package main

import (
	"embed"
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	todohtmx "todo-htmx"
)

//go:embed *.html
var templateFS embed.FS

func main() {
	repo := todohtmx.NewRepository()
	// seed data
	repo.Add("Buy milk")
	id := repo.Add("Buy eggs")
	repo.SetStatus(id, true)
	repo.Add("Buy bread")

	t := template.Must(template.ParseFS(templateFS, "*.html"))

	r := chi.NewRouter()

	r.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		tasks := repo.All()
		count := repo.Count()

		err := t.ExecuteTemplate(writer, "index.html", struct {
			Tasks []todohtmx.Todo
			Count todohtmx.Count
		}{
			Tasks: tasks,
			Count: count,
		})
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	})

	r.Post("/tasks", func(writer http.ResponseWriter, request *http.Request) {
		title := request.FormValue("title")
		repo.Add(title)
		http.Redirect(writer, request, "/", http.StatusFound)
	})

	r.Put("/tasks", func(writer http.ResponseWriter, request *http.Request) {
		err := request.ParseForm()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		allTasks := repo.All()
		for _, task := range allTasks {
			repo.SetStatus(task.ID, false)
		}

		doneTaskIDs := request.Form["task"]
		for _, idStr := range doneTaskIDs {
			id, err := strconv.Atoi(idStr)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusBadRequest)
				return
			}

			repo.SetStatus(id, true)
		}

		http.Redirect(writer, request, "/", http.StatusFound)

	})

	r.Delete("/tasks/{id}", func(writer http.ResponseWriter, request *http.Request) {
		idStr := chi.URLParam(request, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		repo.Delete(id)
		http.Redirect(writer, request, "/", http.StatusFound)
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}

}
