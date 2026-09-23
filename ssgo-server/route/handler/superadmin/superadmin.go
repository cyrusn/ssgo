package superadmin

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"
	"ssgo-server/model/auth"
	"ssgo-server/model/dbmanager"
	"ssgo-server/model/master"
	"ssgo-server/model/student"
	"ssgo-server/model/subject"
	"ssgo-server/model/signature"

	helper "github.com/cyrusn/goHTTPHelper"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	Manager *dbmanager.Manager
	Secret  SecretGenerator
}

type SecretGenerator interface {
	GenerateToken(claims jwt.Claims) (string, error)
}

func (h *Handler) Login() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		form := new(struct { Username string `json:"username"`; Password string `json:"password"` })
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, form)
		err := h.Manager.Master().Authenticate(form.Username, form.Password)
		if err != nil { helper.PrintError(w, err, http.StatusUnauthorized); return }
		claims := auth.Claims{ UserAlias: form.Username, Name: "系統管理員", Cname: "系統管理員", Role: "SUPERADMIN", StandardClaims: jwt.StandardClaims{ ExpiresAt: expiresAfter30Min() } }
		token, _ := h.Secret.GenerateToken(claims)
		w.Write([]byte(token))
	}
}

func (h *Handler) ResetRootPassword() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := h.Manager.Master().ResetRootPasswordToDefault(); err != nil { helper.PrintError(w, err, http.StatusInternalServerError); return }
		w.WriteHeader(http.StatusOK); w.Write([]byte(`{"status":"success"}`))
	}
}

func getEnvSafe(key string) string {
	val := strings.TrimSpace(os.Getenv(key))
	for len(val) >= 2 && ((val[0] == '"' && val[len(val)-1] == '"') || (val[0] == '\'' && val[len(val)-1] == '\'')) {
		val = strings.TrimSpace(val[1 : len(val)-1])
	}
	return val
}

func jsonEscape(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	s = strings.ReplaceAll(s, "\r\n", "\\n")
	s = strings.ReplaceAll(s, "\n", "\\n")
	return s
}

func (h *Handler) GetSystemDefaults() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defaults := map[string]string{
			"schoolName":            getEnvSafe("DEFAULT_SCHOOL_NAME"),
			"schoolWebsite":         getEnvSafe("DEFAULT_SCHOOL_WEBSITE"),
			"systemTitle":           getEnvSafe("DEFAULT_SYSTEM_TITLE"),
			"mockSystemTitle":       getEnvSafe("DEFAULT_MOCK_SYSTEM_TITLE"),
			"introMarkdown":         getEnvSafe("DEFAULT_INTRODUCTION_MAKRDOWN"),
			"mockIntroMarkdown":     getEnvSafe("DEFAULT_MOCK_INTRODUCTION_MAKRDOWN"),
			"instrMarkdown":         getEnvSafe("DEFAULT_INSTRUCTION_MARKDOWN"),
			"notAcceptMarkdown":     getEnvSafe("DEFAULT_NOT_ACCEPT_MARKDOWN"),
			"mockNotAcceptMarkdown": getEnvSafe("DEFAULT_MOCK_NOT_ACCEPT_MARKDOWN"),
			"confirmMarkdown":       getEnvSafe("DEFAULT_CONFIRM_MARKDOWN"),
			"subjectsJSON":          getEnvSafe("DEFAULT_SUBJECTS_JSON"),
			"combinationsJSON":      getEnvSafe("DEFAULT_COMBINATIONS_JSON"),
		}
		helper.PrintJSON(w, defaults)
	}
}

func (h *Handler) CreateCohort() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		form := new(struct { ID string `json:"id"`; Name string `json:"name"`; Config string `json:"config"` })
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, form)
		if form.ID == "" || form.Name == "" { helper.PrintError(w, fmt.Errorf("ID/Name required"), http.StatusBadRequest); return }

		configJSON := form.Config
		if configJSON == "" || configJSON == "{}" || configJSON == "null" {
			isMock := strings.Contains(strings.ToLower(form.ID), "mock")
			intro := getEnvSafe("DEFAULT_INTRODUCTION_MAKRDOWN")
			title := getEnvSafe("DEFAULT_SYSTEM_TITLE")
			notAccept := getEnvSafe("DEFAULT_NOT_ACCEPT_MARKDOWN")
			if isMock {
				intro = getEnvSafe("DEFAULT_MOCK_INTRODUCTION_MAKRDOWN")
				title = getEnvSafe("DEFAULT_MOCK_SYSTEM_TITLE")
				notAccept = getEnvSafe("DEFAULT_MOCK_NOT_ACCEPT_MARKDOWN")
			}
			subjects := getEnvSafe("DEFAULT_SUBJECTS_JSON"); if subjects == "" { subjects = "[]" }
			combos := getEnvSafe("DEFAULT_COMBINATIONS_JSON"); if combos == "" { combos = "[]" }

			configJSON = fmt.Sprintf(`{
				"schoolName": "%s", "systemTitle": "%s", "schoolWebsite": "%s",
				"isMock": %v, "isAcceptResponses": true,
				"introMarkdown": "%s", "instrMarkdown": "%s", "notAcceptMarkdown": "%s", "confirmMarkdown": "%s",
				"subjectsJSON": "%s", "combinationsJSON": "%s"
			}`,
				jsonEscape(getEnvSafe("DEFAULT_SCHOOL_NAME")), jsonEscape(title), jsonEscape(getEnvSafe("DEFAULT_SCHOOL_WEBSITE")),
				isMock, jsonEscape(intro), jsonEscape(getEnvSafe("DEFAULT_INSTRUCTION_MARKDOWN")),
				jsonEscape(notAccept), jsonEscape(getEnvSafe("DEFAULT_CONFIRM_MARKDOWN")),
				jsonEscape(subjects), jsonEscape(combos),
			)
		}
		cohort := &master.Cohort{ID: form.ID, Name: form.Name, DBFileName: fmt.Sprintf("cohort_%s.sqlite", form.ID), IsActive: false, Config: configJSON}
		h.Manager.CreateCohortDB(cohort); helper.PrintJSON(w, cohort)
	}
}

func (h *Handler) ListCohorts() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		list, _ := h.Manager.Master().ListCohorts(); helper.PrintJSON(w, list)
	}
}

func (h *Handler) SetCohortActive() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]; h.Manager.Master().SetCohortActive(id); w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) DeactivateCohort() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]; h.Manager.Master().DeactivateCohort(id); w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) UpdateCohortConfig() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]; body, _ := ioutil.ReadAll(r.Body)
		var f struct { Name string `json:"name"`; Config string `json:"config"` }
		json.Unmarshal(body, &f)
		if f.Name != "" { h.Manager.Master().UpdateCohortName(id, f.Name) }
		h.Manager.Master().UpdateCohortConfig(id, f.Config); w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) DeleteCohort() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]; h.Manager.DeleteCohortDB(id); w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) ListSuperadmins() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		list, _ := h.Manager.Master().ListSuperadmins(); helper.PrintJSON(w, list)
	}
}

func (h *Handler) CreateSuperadmin() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var f struct { Username string `json:"username"`; Password string `json:"password"` }
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &f)
		h.Manager.Master().CreateSuperadmin(f.Username, f.Password); w.WriteHeader(http.StatusCreated)
	}
}

func (h *Handler) UpdateSuperadminPassword() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		username := mux.Vars(r)["username"]; var f struct { Password string `json:"password"` }
		body, _ := ioutil.ReadAll(r.Body)
		json.Unmarshal(body, &f)
		h.Manager.Master().UpdateSuperadminPassword(username, f.Password); w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) DeleteSuperadmin() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		username := mux.Vars(r)["username"]; h.Manager.Master().DeleteSuperadmin(username); w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) ImportData() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		cohortID := mux.Vars(r)["id"]
		dataType := mux.Vars(r)["type"]
		body, _ := ioutil.ReadAll(r.Body)
		db, err := h.Manager.GetCohortDB(cohortID)
		if err != nil {
			helper.PrintError(w, err, http.StatusBadRequest)
			return
		}

		if dataType == "student" {
			type StudentImport struct {
				UserAlias string `json:"userAlias"`
				Name      string `json:"name"`
				Cname     string `json:"cname"`
				Password  string `json:"password"`
				ClassCode string `json:"classCode"`
				ClassNo   int    `json:"classNo"`
			}
			var imports []StudentImport
			if err := json.Unmarshal(body, &imports); err != nil {
				helper.PrintError(w, err, http.StatusBadRequest)
				return
			}
			dbS := &student.DB{DB: db}
			dbA := &auth.DB{DB: db}
			dbSig := &signature.DB{DB: db}
			for _, imp := range imports {
				// 1. Insert Credential FIRST (satisfy FK)
				err := dbA.Insert(&auth.Credential{
					UserAlias: imp.UserAlias,
					Password:  imp.Password,
					Role:      "STUDENT",
					Name:      imp.Name,
					CName:     imp.Cname,
				})
				if err != nil {
					if strings.Contains(err.Error(), "UNIQUE constraint failed") {
						helper.PrintError(w, fmt.Errorf("匯入失敗：部分帳號已存在。"), http.StatusConflict)
					} else {
						helper.PrintError(w, err, http.StatusInternalServerError)
					}
					return
				}
				// 2. Insert Student record
				err = dbS.Insert(&student.Student{
					UserAlias: imp.UserAlias,
					ClassCode: imp.ClassCode,
					ClassNo:   imp.ClassNo,
				})
				if err != nil {
					helper.PrintError(w, err, http.StatusInternalServerError)
					return
				}
				// 3. Init Signature
				dbSig.Insert(imp.UserAlias)
			}
		} else if dataType == "teacher" || dataType == "admin" {
			var credentials []auth.Credential
			json.Unmarshal(body, &credentials)
			dbA := &auth.DB{DB: db}
			role := "TEACHER"
			if dataType == "admin" {
				role = "ADMIN"
			}
			for _, c := range credentials {
				c.Role = role
				err := dbA.Insert(&c)
				if err != nil {
					if strings.Contains(err.Error(), "UNIQUE constraint failed") {
						helper.PrintError(w, fmt.Errorf("匯入失敗：部分帳號已存在。"), http.StatusConflict)
					} else {
						helper.PrintError(w, err, http.StatusInternalServerError)
					}
					return
				}
			}
		} else if dataType == "subject" {
			var codes []string
			json.Unmarshal(body, &codes)
			dbSub := &subject.DB{DB: db}
			for _, code := range codes {
				dbSub.Insert(&subject.Subject{Code: code})
			}
		}
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) GetUser() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		id := mux.Vars(r)["id"]
		alias := mux.Vars(r)["userAlias"]
		db, err := h.Manager.GetCohortDB(id)
		if err != nil {
			helper.PrintError(w, err, http.StatusBadRequest)
			return
		}
		var u struct {
			UserAlias string `json:"userAlias"`
			Name      string `json:"name"`
			Cname     string `json:"cname"`
			Role      string `json:"role"`
			ClassCode string `json:"classCode"`
			ClassNo   int    `json:"classNo"`
		}
		err = db.QueryRow(`SELECT c.userAlias, c.name, c.cname, c.role, COALESCE(s.classCode,''), COALESCE(s.classNo,0) FROM Credential c LEFT JOIN Student s ON c.userAlias=s.userAlias WHERE c.userAlias=?`, alias).Scan(&u.UserAlias, &u.Name, &u.Cname, &u.Role, &u.ClassCode, &u.ClassNo)
		if err != nil {
			fmt.Printf("GetUser Error: %v\n", err)
		}
		helper.PrintJSON(w, u)
	}
}

func (h *Handler) UpdateUser() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := mux.Vars(r)["userAlias"]
		db, _ := h.Manager.GetCohortDB(mux.Vars(r)["id"])
		body, _ := ioutil.ReadAll(r.Body)
		var f struct {
			Name      string `json:"name"`
			Cname     string `json:"cname"`
			Password  string `json:"password"`
			ClassCode string `json:"classCode"`
			ClassNo   int    `json:"classNo"`
			Role      string `json:"role"`
		}
		json.Unmarshal(body, &f)

		// Fetch current role
		var currentRole string
		db.QueryRow("SELECT role FROM Credential WHERE userAlias=?", alias).Scan(&currentRole)

		// Role restriction logic
		targetRole := f.Role
		if targetRole == "" {
			targetRole = currentRole
		}

		// Only allow change between TEACHER and ADMIN
		if currentRole == "STUDENT" {
			targetRole = "STUDENT"
		} else if targetRole == "STUDENT" {
			targetRole = currentRole // Don't allow changing teacher/admin to student
		}

		if f.Password != "" {
			hp, _ := bcrypt.GenerateFromPassword([]byte(f.Password), bcrypt.DefaultCost)
			db.Exec("UPDATE Credential SET name=?, cname=?, password=?, role=? WHERE userAlias=?", f.Name, f.Cname, hp, targetRole, alias)
		} else {
			db.Exec("UPDATE Credential SET name=?, cname=?, role=? WHERE userAlias=?", f.Name, f.Cname, targetRole, alias)
		}

		db.Exec("UPDATE Student SET classCode=?, classNo=? WHERE userAlias=?", f.ClassCode, f.ClassNo, alias)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handler) DeleteUser() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := mux.Vars(r)["userAlias"]; db, _ := h.Manager.GetCohortDB(mux.Vars(r)["id"])
		db.Exec("DELETE FROM Student WHERE userAlias=?", alias); db.Exec("DELETE FROM Signature WHERE userAlias=?", alias); db.Exec("DELETE FROM Credential WHERE userAlias=?", alias); w.WriteHeader(http.StatusOK)
	}
}

func expiresAfter30Min() int64 { return time.Now().Add(time.Minute * 30).Unix() }
