package main

import (
	"auth_service/internal/database"
	"auth_service/internal/handlers"
	"auth_service/internal/repository"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Can't initializa zap logger: %v", err)
	}
	defer logger.Sync() // Flushes buffer
	sugar := logger.Sugar()

	sugar.Infow("starting up", "version", "v1.0.0", "port", 8080)

	db, err := database.Initialize(sugar)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	err = database.SetupSchema(db)
	if err != nil {
		log.Fatal(err)
	}

	// Обработчик для отдачи Swagger UI
	http.Handle("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/openapi.yaml"), // Путь к файлу спецификации API
	))
	http.Handle("/openapi.yaml", http.FileServer(http.Dir("/root")))
	//http.Handle("/openapi.yaml", http.FileServer(http.Dir("../../openapi/")))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/", http.StatusMovedPermanently)
		sugar.Infow("User redirrected to /swagger/")
	})

	userRepository := repository.NewUserRepository(db)

	http.HandleFunc("/register", handlers.RegisterUser(userRepository, sugar))

	http.HandleFunc("/login", handlers.LoginHandler(userRepository, sugar))

	log.Println("Starting Auth Service on port 8080...")
	//if err := http.ListenAndServeTLS(":8443", "/root/cert.pem", "/root/key.pem", nil); err != nil {
	//	log.Fatal(err)
	//}
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}

}
