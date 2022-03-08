package controllers

import (
	"bytes"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"strconv"
	"text/template"
	"time"

	"gitlab.com/verdant-sso/pkg/util/database"
	"gitlab.com/verdant-sso/pkg/util/network"
	"gitlab.com/verdant-sso/pkg/util/session/jwt"
	"gitlab.com/verdant-sso/pkg/util/validation"
	"golang.org/x/crypto/bcrypt"

	"github.com/go-chi/chi"
)

type (
	UserLoginForm struct {
		Email    string `json:"email" `
		Password string `json:"password" `
	}

	RegisterNewUser struct {
		FirstName    string
		LastName     string
		Email        string
		Username     string
		UserPassword string `json:"password" `
		PhoneNo      string
		RoleID       int `json:"role_id" `
	}

	GetAuthDetail struct {
		UserID       string
		FirstName    string
		LastName     string
		Email        string
		RoleName     string
		RoleID       string
		UserPassword string
	}
)

func NewAuthRoute() *chi.Mux {
	router := chi.NewRouter()

	//normal staff
	router.Post("/login", UserLogin)
	router.Put("/logout/{email}", UserLogout)
	router.Post("/forgot-password", ForgotPassword)
	router.Post("/set-password", ResetPassword)

	//token check
	router.Get("/tokencheck/{email}", TokenCheck)

	return router
}

func UserLogin(w http.ResponseWriter, r *http.Request) {

	loginmodel := UserLoginForm{}
	validatemodel := GetAuthDetail{}
	if err := network.ReadJSONData(r, &loginmodel); err != nil {
		log.Println(err)
		network.ResponseJSON(w, true, http.StatusBadRequest, ErrInvalidParam)
	} else {
		// if CheckLock(loginmodel.EmployerID, "CHECK") {
		// 	network.ResponseJSON(w, true, ErrAuthSuspended)
		// } else {
		// if !CheckSession(loginmodel.Email) {
		if validation.ValidateForm(&loginmodel) {

			//check empid adak x
			if err := database.GetInstance().
				Table("users").
				Select("users.*, roles.role_name").
				Joins("JOIN roles ON roles.role_id = users.role_id").
				Where("users.email = ?", loginmodel.Email).
				Find(&validatemodel).
				Error; err != nil {
				log.Println(err)
				network.ResponseJSON(w, true, http.StatusInternalServerError, err.Error())

			} else {
				// aclvalue := string(validatemodel.AccessAvai)
				// aclvalue = aclvalue[1 : len(aclvalue)-1]
				// finalacl := strings.Split(aclvalue, ",")
				// log.Println(aclvalue[1])
				// token := ""
				// errfail := truex
				//ni cara cari pasal namanya , NEED REVIEW , TEMP SETUP

				//try match password
				log.Println(loginmodel)
				log.Println(validatemodel)
				err := bcrypt.CompareHashAndPassword([]byte(validatemodel.UserPassword),
					[]byte(loginmodel.Password))

				if err != nil {
					log.Println("WRONG_PW", err)
					network.ResponseJSON(w, true, http.StatusUnauthorized, "Wrong password")
					return
				} else {

					//if password is correct , generate jwt and push to redis
					token, _ := jwt.GenNewToken(
						validatemodel.UserID,
						validatemodel.FirstName,
						validatemodel.Email,
						validatemodel.RoleID,
					)

					//update auth status

					// if err := database.GetInstance().
					// 	Table("employees").
					// 	Where("emp_id = ? ", loginmodel.EmployerID).
					// 	Update("auth_status", "1").Error; err != nil {
					// 	log.Println("ERR_UPDATE_LOGSTATUS: ", err)
					// } else {
					//push tu redis
					key := "vl:verdant:login:" + loginmodel.Email
					err = database.GetRedisInstance().
						Set(key,
							token,
							time.Hour*72).Err()

					if err != nil {

						log.Println(err)
						network.ResponseJSON(w, true, http.StatusInternalServerError, err)
						return
					} else {

						hj, _ := database.GetRedisInstance().
							Get(key).
							Result()
						log.Println("FROM REDIS: ", hj)

					}

					network.ResponseJSON(w, false, http.StatusOK, map[string]interface{}{
						"jwt_token": token,
					})
				}
			}
		} else {

		}
		// } else {
		// 	network.ResponseJSON(w, true, http.StatusConflict, "Multiple logins are not supported. Please log out from previous session.")
		// }

	}
}

func UserLogout(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")

	log.Println("LOGOUT: ", email)
	if email != "" {
		if CheckSession(email) {
			if err := database.GetInstance().
				Table("users").
				Where("email = ? ", email).
				Update("auth_status", 0).Error; err != nil {
				log.Println("ERR_UPDATE_LOGSTATUS: ", err)
				network.ResponseJSON(w, true, http.StatusInternalServerError, "Unable to logout")
			} else {
				key := "vl:verdant:login:" + email
				log.Println(key)
				database.GetRedisInstance().Del(key).Err()
				network.ResponseJSON(w, false, http.StatusOK, "Logout OK")
			}
		} else {
			network.ResponseJSON(w, true, http.StatusInternalServerError, "Key not found")
		}
	} else {
		network.ResponseJSON(w, true, http.StatusBadRequest, "failed logout")
	}
}

func TokenCheck(w http.ResponseWriter, r *http.Request) {
	email := chi.URLParam(r, "email")

	if email != "" {
		staffkey, _ := database.GetRedisInstance().Get(email).Result()
		if staffkey != "" {
			network.ResponseJSON(w, false, http.StatusOK, "OK")
			database.GetRedisInstance().
				Set(email,
					staffkey,
					time.Second*10)

		} else {
			database.GetRedisInstance().Del(email)
			network.ResponseJSON(w, true, http.StatusInternalServerError, "NOT OK")

		}
	} else {
		network.ResponseJSON(w, true, http.StatusBadRequest, "Email is required to perform token check.")
	}

}

func CheckSession(email string) bool {
	key := "vl:verdant:login:" + email
	if email != "" {
		staffkey, _ := database.GetRedisInstance().Get(key).Result()

		if staffkey != "" {
			return true
		} else {
			return false
		}
	}
	return false
}

func ForgotPassword(w http.ResponseWriter, r *http.Request) {
	req := network.ForgotPasswordReq{}
	err := network.ReadJSONData(r, &req)
	if err != nil {
		network.ResponseJSON(w, true, http.StatusBadRequest, ErrInvalidParam)
		return
	}

	min := 10000000
	max := 99999999
	token := "vl:verdant:pw:" + strconv.Itoa(rand.Intn(max-min)+min)
	log.Println(token)

	err = database.GetRedisInstance().
		Set(token,
			req.Email,
			time.Minute*20).Err()
	if err != nil {
		network.ResponseJSON(w, true, http.StatusInternalServerError, "Unable to access redis.")
		return
	}

	reset_link := "http://verdant.vectolabs.com/#/auth/set-password" + token

	// Sender data.
	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASSWD")

	// Receiver email address.
	to := []string{
		req.Email,
	}

	// smtp server configuration.
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// Message.
	log.Println(reset_link)
	//message := []byte(reset_link)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, smtpHost)

	path, err := os.Getwd()
	if err != nil {
		log.Println(err)
	}
	fmt.Println(path)
	templateFile := fmt.Sprintf("%s/cmd/template.html", path)
	fmt.Println(templateFile)
	emailBody, _ := template.ParseFiles(templateFile)

	var body bytes.Buffer

	mimeHeaders := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	body.Write([]byte(fmt.Sprintf("Subject: Password Reset Link \n%s\n\n", mimeHeaders)))

	err = emailBody.Execute(&body, struct {
		Name string
		Url  string
	}{
		Name: req.Email,
		Url:  reset_link,
	})

	log.Println(err)

	// Sending email.
	err = smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, body.Bytes())
	if err != nil {
		log.Println(err)
		network.ResponseJSON(w, true, http.StatusInternalServerError, "Email failed to send.")
		return
	}
	network.ResponseJSON(w, false, http.StatusOK, "Email Sent Successfully!")
}

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	req := network.ResetPasswordReq{}
	err := network.ReadJSONData(r, &req)
	if err != nil {
		network.ResponseJSON(w, true, http.StatusBadRequest, ErrInvalidParam)
		return
	}

	email, err := database.GetRedisInstance().Get("vl:verdant:pw:" + req.Token).Result()
	if err != nil {
		network.ResponseJSON(w, true, http.StatusForbidden, "Token Expired.")
		return
	}

	hpwd, err := bcrypt.GenerateFromPassword([]byte(req.Password),
		bcrypt.DefaultCost)
	if err != nil {
		network.ResponseJSON(w, true, http.StatusInternalServerError, "Password encryption error.")
		return
	}
	usertable := database.User{}
	err = database.GetInstance().Model(&usertable).Where("email=?", email).Update("user_password", string(hpwd)).Error
	if err != nil {
		network.ResponseJSON(w, true, http.StatusInternalServerError, err)
		return
	}

	err = database.GetRedisInstance().Del("vl:verdant:pw:" + req.Token).Err()
	if err != nil {
		log.Println("Error deleting key vl:verdant:pw:", req.Token, "from Redis! Error: ", err)
	}
	network.ResponseJSON(w, true, http.StatusOK, "Password reset successfully.")
}

// func CheckLock(email, act string) bool {
// 	userModel := &database.User{}
// 	if act == "CHECK" {
// 		log.Println("CHECK")
// 		database.GetInstance().
// 			Table("employees").
// 			Select("status_id").
// 			Where("emp_id = ?", email).
// 			Scan(&userModel)

// 	} else if act == "LOCK" {
// 		log.Println("SALAHH")
// 		//shitf string
// 		tempArr := []string{}
// 		email := len(email)
// 		for i := 1; i <= email; {
// 			tempArr = append(tempArr, string(email[len(email)-(i)]))
// 			i++
// 		}

// 		invertedID := strings.Join(tempArr, "")

// 		//check count

// 		// hj, _ := database.GetRedisInstance().Get(invertedID).Result()
// 		if err := database.GetInstance().
// 			Table("employees").
// 			Select("auth_trial").
// 			Where("emp_id = ? ", email).
// 			Find(userModel).Error; err != nil {
// 			log.Println("ERR_UPDATA_GETCOUNT:", err)
// 		}

// 		count, _ := strconv.Atoi(userModel.AuthTrial)
// 		log.Println("FROM REDIS LOCK: ", count)
// 		if count >= 2 {
// 			userModel.StatusID = 2
// 			//suspended acc
// 			database.GetInstance().
// 				Table("employees").
// 				Where("emp_id = ?", email).
// 				Update(&userModel)

// 			go UnlockAccount(email, invertedID)
// 			return true
// 		} else if count == 0 {
// 			if err := database.GetInstance().
// 				Table("employees").
// 				Where("emp_id = ? ", email).
// 				Update("auth_trial", 1).Error; err != nil {
// 				log.Println("ERR_UPDATA_TRYCOUNT:", err)
// 			}
// 			// database.GetRedisInstance().Set(invertedID, 1, time.Hour*24).Err()
// 		} else {
// 			count++
// 			if err := database.GetInstance().
// 				Table("employees").
// 				Where("emp_id = ? ", email).
// 				Update("auth_trial", count).Error; err != nil {
// 				log.Println("ERR_UPDATA_TRYCOUNT:", err)
// 			}
// 			// database.GetRedisInstance().Set(invertedID, count, time.Hour*24).Err()
// 		}
// 	}

// 	return false

// }

// func UnlockAccount(email, inverID string) {

// 	log.Println("UNLOCK_ACC_COUNTDOWN")
// 	time.Sleep(time.Second * 10)

// 	userModel := &database.Employee{}

// 	userModel.StatusID = 1
// 	userModel.AuthTrial = "0"
// 	//suspended acc
// 	if err := database.GetInstance().
// 		Table("employees").
// 		Where("emp_id = ? ", email).
// 		Update(userModel).Error; err != nil {
// 		log.Println("ERR_UPDATA_TRYCOUNT:", err)
// 	}

// 	database.GetRedisInstance().Del(inverID).Result()
// }
