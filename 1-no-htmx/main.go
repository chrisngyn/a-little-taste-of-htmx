package main

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	todohtmx "todo-htmx"
)

func main() {
	repo := todohtmx.NewRepository()
	// seed data
	repo.Add("Buy milk")
	id := repo.Add("Buy eggs")
	repo.SetStatus(id, true)
	repo.Add("Buy bread")

	tmpl := template.Must(template.ParseFiles("1-no-htmx/index.html"))

	http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		tasks := repo.All()
		total, completed := repo.Count()

		err := tmpl.Execute(writer, struct {
			Tasks     []todohtmx.Todo
			Total     int
			Completed int
		}{
			Tasks:     tasks,
			Total:     total,
			Completed: completed,
		})
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/add", func(writer http.ResponseWriter, request *http.Request) {
		title := request.FormValue("title")
		repo.Add(title)
		http.Redirect(writer, request, "/", http.StatusFound)
	})

	http.HandleFunc("/delete/", func(writer http.ResponseWriter, request *http.Request) {
		pathSegments := strings.Split(request.URL.Path, "/")

		if len(pathSegments) != 3 {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		id, err := strconv.Atoi(pathSegments[2])
		if err != nil {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		repo.Delete(id)
		http.Redirect(writer, request, "/", http.StatusFound)
	})

	http.HandleFunc("/set-status", func(writer http.ResponseWriter, request *http.Request) {
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

	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}
