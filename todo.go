package todo_htmx

import (
	"sync"
)

type Todo struct {
	ID        int
	Title     string
	Completed bool
}

type Count struct {
	Total     int
	Completed int
}

type Repository struct {
	todos []Todo

	mu        sync.Mutex
	idCounter int
}

func NewRepository() *Repository {
	return &Repository{
		todos: []Todo{},
	}
}

func (r *Repository) Add(title string) int {
	id := r.generateID()
	r.todos = append(r.todos, Todo{
		ID:        id,
		Title:     title,
		Completed: false,
	})
	return id
}

func (r *Repository) All() []Todo {
	return r.todos
}

func (r *Repository) SetStatus(id int, completed bool) {
	for i, todo := range r.todos {
		if todo.ID == id {
			r.todos[i].Completed = completed
		}
	}
}

func (r *Repository) Delete(id int) {
	newTodos := make([]Todo, 0, len(r.todos)-1)
	for _, todo := range r.todos {
		if todo.ID != id {
			newTodos = append(newTodos, todo)
		}
	}
	r.todos = newTodos
}

func (r *Repository) Count() Count {
	var completed int
	for _, todo := range r.todos {
		if todo.Completed {
			completed++
		}
	}
	return Count{
		Total:     len(r.todos),
		Completed: completed,
	}
}

func (r *Repository) generateID() int {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.idCounter++
	return r.idCounter
}
