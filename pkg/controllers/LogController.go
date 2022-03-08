package controllers

// import (
// 	"net/http"

// 	"gitlab.com/verdant-sso/pkg/util/database"
// 	"gitlab.com/verdant-sso/pkg/util/network"

// 	"github.com/go-chi/chi"
// )

// func NewLogRoute() *chi.Mux {

// 	route := chi.NewRouter()
// 	//user
// 	route.Post("/write", WriteLog)
// 	route.Get("/write", GetLog)

// 	return route
// }

// func WriteLog(w http.ResponseWriter, r *http.Request) {
// 	model := database.AuditLog{}
// 	jsonErr := network.ReadJSONData(r, &model)

// 	if jsonErr != nil {
// 		network.ResponseJSON(w, true, http.StatusBadRequest, jsonErr)
// 		return
// 	}

// 	getID := network.GetID{}
// 	database.GetInstance().
// 		Table("audit_logs").
// 		Order("audit_id DESC").
// 		Limit("1").
// 		Find(&getID)

// 	model.AuditID = getID.AuditID
// 	model.AuditID++

// 	err := database.GetInstance().
// 		Create(&model).
// 		Error

// 	if err != nil {
// 		network.ResponseJSON(w, true, http.StatusInternalServerError, ErrDbTransFail)
// 	} else {
// 		network.ResponseJSON(w, false, http.StatusCreated, "Log inserted successfully.")
// 	}

// }

// func GetLog(w http.ResponseWriter, r *http.Request) {
// 	// model := []database.AuditLog{}

// }
