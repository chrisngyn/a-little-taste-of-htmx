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
	todohtmx.Seed(repo)

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
		if title != "" {
			id := repo.Add(title)
			if err := t.ExecuteTemplate(writer, "task.html", repo.Get(id)); err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
		}
	})

	r.Put("/tasks/{id}", func(writer http.ResponseWriter, request *http.Request) {
		idStr := chi.URLParam(request, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}
		err = request.ParseForm()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusBadRequest)
			return
		}

		completed := false
		values := request.Form["task"]
		if len(values) > 0 {
			completed = values[0] == "on"
		}

		repo.SetStatus(id, completed)

		if err := t.ExecuteTemplate(writer, "task.html", repo.Get(id)); err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	})

	r.Delete("/tasks/{id}", func(writer http.ResponseWriter, request *http.Request) {
		idStr := chi.URLParam(request, "id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		repo.Delete(id)
	})

	if err := http.ListenAndServe(":8080", r); err != nil {
		panic(err)
	}

}
