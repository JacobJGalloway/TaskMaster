package main

import (
	"log"
	"net/http"
	"database/sql"
	_ "modernc.org/sqlite"
	"TaskMaster-api/handler"
	"TaskMaster-api/service"
	"TaskMaster-api/repository"
)

const schema = `
CREATE TABLE IF NOT EXISTS users (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      name TEXT NOT NULL
  );

  CREATE TABLE IF NOT EXISTS study_tasks (
      id INTEGER PRIMARY KEY AUTOINCREMENT,
      user_id INTEGER NOT NULL REFERENCES users(id),
      title TEXT NOT NULL,
      duration_minutes INTEGER NOT NULL,
      period_start TEXT NOT NULL,
      period_end TEXT NOT NULL,
      completed INTEGER NOT NULL DEFAULT 0,
      completed_at TEXT
  );`

func withCORS(next http.Handler) http.Handler {      
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:5173")          
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main(){
	db, err := sql.Open("sqlite", "./db.sqlite")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	_, err = db.Exec(schema)
	if err != nil {
		log.Fatalf("❌ Failed to initialize schema: %v", err)
	}

	log.Println("Database initialized and schema applied successfully!")

	db.Ping()

	uerrRepo := repository.NewSQLiteUserRepository(db)
	userService := service.NewUserService(uerrRepo)
	userHandler := handler.NewUserHandler(userService)

	taskRepo := repository.NewSQLiteTaskRepository(db)
	taskService := service.NewTaskService(taskRepo)
	taskHandler := handler.NewTaskHandler(taskService)

	mux := http.NewServeMux()
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})
	mux.HandleFunc("POST /users", userHandler.Create)
	mux.HandleFunc("GET /users", userHandler.GetUsers)
	mux.HandleFunc("GET /users/{id}", userHandler.GetUserByID)
	mux.HandleFunc("GET /users/{id}/tasks", taskHandler.ListActiveByUser)
  	mux.HandleFunc("POST /tasks", taskHandler.Create)
	mux.HandleFunc("POST /tasks/{id}/complete", taskHandler.MarkComplete)

	log.Fatal(http.ListenAndServe(":8080", withCORS(mux)))

}