package controllers

import (
	"log"
	"net/http"

	"gitlab.com/verdant-sso/pkg/util/database"
	"gitlab.com/verdant-sso/pkg/util/session/jwt"
	"golang.org/x/crypto/bcrypt"

	"gitlab.com/verdant-sso/pkg/util/network"

	"github.com/go-chi/chi"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
)

func ManagementRoute() *chi.Mux {
	route := chi.NewRouter()

	// users
	route.Get("/users", ListUser)
	route.Get("/users/{userID}", GetUser)
	route.Post("/users", CreateUser)
	route.Put("/users/{userID}", UpdateUser)
	route.Delete("/users/{userID}", DeleteUser)

	//roles
	route.Post("/roles", CreateRoles)
	route.Get("/roles", GetRoles)
	route.Put("/roles/{roleid}", UpdateRoles)
	route.Delete("/roles/{roleid}", DeleteRoles)

	return route
}

// ############################ USERS ############################ //

func ListUser(w http.ResponseWriter, r *http.Request) {
	resp := []network.ListUser{}
	err := database.GetInstance().
		Select("username, role_name, users.created_at as date_registered").
		Table("users").
		Joins("JOIN roles ON users.role_id = roles.role_id").
		Scan(&resp).Error
	if err != nil {
		network.ResponseJSON(w, false, http.StatusInternalServerError, "Error")
		return
	}
	network.ResponseJSON(w, false, http.StatusOK, resp)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	userID := chi.URLParam(r, "userID")
	log.Println(userID)
	// resp := network

}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	empmodel := database.User{}
	regmodel := RegisterNewUser{}
	var err1 error
	if err := network.ReadJSONData(r, &regmodel); err != nil {
		log.Println(err)
		network.ResponseJSON(w, true, http.StatusBadRequest, ErrInvalidParam)
		return
	} else {
		//check email
		if emailFound := database.GetInstance().
			Select("email").
			Where("email = ?", regmodel.Email).
			Find(&database.User{}).
			RecordNotFound(); emailFound {

			//hash password
			hpwd, err := bcrypt.GenerateFromPassword([]byte(regmodel.UserPassword),
				bcrypt.DefaultCost)
			if err != nil {
				log.Println("ERR_CR8USER_ENCPASSWD: ", err)
			}

			copier.Copy(&empmodel, &regmodel)
			empmodel.UserPassword = string(hpwd)
			empmodel.RoleID = 3

			err1 = database.DoTransaction(func(tx *gorm.DB) error {
				if err := tx.Table("users").
					Create(&empmodel).
					Error; err != nil {
					return err

				}
				return nil
			})

			if err1 != nil {

				network.ResponseJSON(w, true, http.StatusInternalServerError, ErrDbInsertFail)
			} else {
				network.ResponseJSON(w, false, http.StatusCreated, RegSuccess)
			}
		} else {
			network.ResponseJSON(w, true, http.StatusConflict, "Email not available")
		}
	}

}

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	//userid := chi.URLParam(r, "userID")

	// //if has passwd change , will encrypt
	// if formmodel.EmpPassword != "" {
	// 	hpwd, err := bcrypt.GenerateFromPassword([]byte(formmodel.EmpPassword), bcrypt.DefaultCost)
	// 	if err != nil {
	// 		log.Println("ERR_CR8USER_ENCPASSWD: ", err)
	// 	}

	// 	usermodel.UserPassword = string(hpwd)
	// }

}

func DeleteUser(w http.ResponseWriter, r *http.Request) {
	userid := chi.URLParam(r, "userID")
	data, _ := jwt.Get(r.Context())
	deluserid := data.Get("userid").MustInt64()
	println(deluserid)

	//check if user ade x
	usercheck := database.GetInstance().
		Where("user_id = ?", userid).
		Find(&database.User{}).
		RecordNotFound()

	if !usercheck {
		if err := database.DoTransaction(func(tx *gorm.DB) error {
			if err := tx.Table("users").
				Unscoped().
				Where("user_id = ?", userid).
				Delete(database.User{}).
				Error; err != nil {
				return err
			}

			return nil
		}); err != nil {
			log.Println("ERR_DELEMP: ", err)
			network.ResponseJSON(w, true, http.StatusInternalServerError, ErrDbTransFail)
		} else {
			network.ResponseJSON(w, false, http.StatusOK, "User successfully deleted.")
		}
	} else {
		network.ResponseJSON(w, true, http.StatusBadRequest, ErrInvalidParam)
	}

}

// ############################ ROLE ############################ //

func GetRoles(w http.ResponseWriter, r *http.Request) {
	response := []network.GetRole{}
	rolesmodel := []database.Role{}

	if err := database.GetInstance().
		Find(&rolesmodel).
		Scan(&response).
		Error; err != nil {

		log.Println(err)
		network.ResponseJSON(w, true, http.StatusInternalServerError, ErrDbInsertFail)

	} else {
		network.ResponseJSON(w, false, http.StatusOK, response)
	}
}

func CreateRoles(w http.ResponseWriter, r *http.Request) {
	body := network.AddRole{}
	dbmodel := database.Role{}
	err := network.ReadJSONData(r, &body)
	if err != nil {
		network.ResponseJSON(w, true, http.StatusBadRequest, "Invalid parameter")
		return
	}

	database.GetInstance().Last(&dbmodel)
	dbmodel.RoleID++
	copier.Copy(&dbmodel, &body)

	err = database.GetInstance().Create(&dbmodel).Error
	if err != nil {
		network.ResponseJSON(w, true, http.StatusInternalServerError, "An error occured when performing data transaction.")
		log.Println(err.Error())
		return
	}
	network.ResponseJSON(w, false, http.StatusCreated, "Role added successfully.")

}

func UpdateRoles(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "roleid")
	uptmodel := network.AddRole{}
	dbmodel := database.Role{}
	err := network.ReadJSONData(r, &uptmodel)
	if err != nil {
		network.ResponseJSON(w, true, http.StatusBadRequest, "Invalid parameter")
		return
	}
	tx := database.GetInstance().Model(&dbmodel).Where("role_id = ?", id).
		Updates(database.Role{RoleName: uptmodel.RoleName})
	if tx.Error != nil {
		network.ResponseJSON(w, true, http.StatusInternalServerError, "An error occured when performing data transaction.")
		return
	}
	if tx.RowsAffected <= 0 {
		network.ResponseJSON(w, true, http.StatusBadRequest, "Nothing is updated.")
		return
	}
	network.ResponseJSON(w, false, http.StatusOK, "Role updated successfully")
}

func DeleteRoles(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "roleid")
	dbmodel := database.Role{}
	tx := database.GetInstance().Where("role_id = ?", id).Delete(&dbmodel)

	if tx.Error != nil {
		network.ResponseJSON(w, true, http.StatusInternalServerError, tx.Error.Error())
		return
	}

	if tx.RowsAffected <= 0 {
		network.ResponseJSON(w, true, http.StatusNotFound, "Requested ID for delete operation not found.")
		return
	}
	network.ResponseJSON(w, false, http.StatusOK, "Role successfully deleted")
}
