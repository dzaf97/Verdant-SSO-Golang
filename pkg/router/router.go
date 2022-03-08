package router

import (
	"log"
	"net/http"
	"os"

	"gitlab.com/verdant-sso/pkg/util/session/jwt"

	"gitlab.com/verdant-sso/pkg/util/database"

	"gitlab.com/verdant-sso/pkg/controllers"
	"gitlab.com/verdant-sso/pkg/util/cron"

	"github.com/go-chi/chi"
	"github.com/joho/godotenv"
)

type RouterHandler struct {
	mux  *chi.Mux
	mux1 *chi.Mux
	jwtt *jwt.JWTStruct
}

func NewRouter() *RouterHandler {

	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	mux := chi.NewRouter()
	mux1 := chi.NewRouter()

	//MQTT
	// topic - /verdant/canbus/dvc001/tx
	// mqhost := os.Getenv("MQ_HOST1")
	// mqtt.NewMQTTClient(mqhost)
	// mqtt.GetMqttClient().Subscribe("/verdant/canbus/+/tx", 0, nil)

	//redis
	err, _ := database.NewRedis()

	// aws
	err, _ = database.NewAwsSession()

	//Auth
	auth := jwt.NewAuth()

	//cron
	err, _ = cron.NewCron()

	//db
	err, _ = database.NewDB()
	// database.NewMongoObj()

	if err != nil {
		log.Println("ERR_ROUTER_INIT: ", err)
	} else {
		return &RouterHandler{
			mux:  mux,
			mux1: mux1,
			jwtt: auth,
		}
	}

	return nil
}

func (a *RouterHandler) Run() {

	//Rerun cron
	// cron.ReRunCron()

	//init jwt

	a.mux.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("\n\n"))
		w.Write([]byte("Scheduler Service\n"))

	})

	//main logiun
	a.mux.Route("/auth", func(r chi.Router) {

		//changes for teting
		// cors := cors.New(cors.Options{
		// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		// 	AllowedOrigins: []string{"*"},
		// 	// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// 	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		// 	ExposedHeaders:   []string{"Link"},
		// 	AllowCredentials: true,
		// 	MaxAge:           300, // Maximum value not ignored by any of major browsers
		// })

		// r.Use(cors.Handler)

		// login
		r.Mount("/", controllers.NewAuthRoute())

	})

	//main route managemetn
	a.mux.Route("/api/v1", func(r chi.Router) {

		// cors := cors.New(cors.Options{
		// 	// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
		// 	AllowedOrigins: []string{"*"},
		// 	// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		// 	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// 	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		// 	ExposedHeaders:   []string{"Link"},
		// 	AllowCredentials: true,
		// 	MaxAge:           300, // Maximum value not ignored by any of major browsers
		// })

		// r.Use(cors.Handler)

		// jwt check middleware
		// r.Use(a.jwtt.AuthMiddleware)

		//management
		r.Mount("/management", controllers.ManagementRoute())

		//r.Mount("/log", controllers.NewLogRoute())

	})

	// //INTERNAL ROUTE
	// a.mux1.Route("/api/internal", func(r chi.Router) {
	// 	r.Mount("/tenant", controllers.NewRegRoute())
	// 	r.Mount("/manage", controllers.NewManageRoute())
	// })

	port1 := os.Getenv("OUT_PORT")
	port2 := os.Getenv("IN_PORT")
	//start api server
	log.Println("START PUBLIC PORT: ", port1)
	go http.ListenAndServe(":"+port1, a.mux)
	log.Println("START INTERNAL PORT: ", port2)
	http.ListenAndServe(":"+port2, a.mux1)
}
