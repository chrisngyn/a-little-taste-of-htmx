package todo_htmx

type Todo struct {
	ID        int
	Title     string
	Completed bool
}

type Repository struct {
	todos []Todo
}

func NewRepository() *Repository {
	return &Repository{
		todos: []Todo{},
	}
}

func (r *Repository) Add(title string) int {
	id := len(r.todos) + 1
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
	for i, todo := range r.todos {
		if todo.ID == id {
			r.todos = append(r.todos[:i], r.todos[i+1:]...)
		}
	}
}

func (r *Repository) Count() (total int, completed int) {
	for _, todo := range r.todos {
		if todo.Completed {
			completed++
		}
	}
	return len(r.todos), completed
}
