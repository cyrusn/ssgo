package route

import (
	"database/sql"
	"net/http"

	"ssgo-server/model/dbmanager"
	modelAuth "ssgo-server/model/auth"
	"ssgo-server/model/student"
	"ssgo-server/model/subject"
	"ssgo-server/model/signature"

	"ssgo-server/route/handler/auth"
	signatureHandler "ssgo-server/route/handler/signature"
	studentHandler "ssgo-server/route/handler/student"
	subjectHandler "ssgo-server/route/handler/subject"
	"ssgo-server/route/handler/superadmin"

	jwtAuth "github.com/cyrusn/goJWTAuthHelper"
)

// Env contains stores for providing values to handlers dynamically
type Env struct {
	Manager *dbmanager.Manager
	Secret  *jwtAuth.Secret
}

// Route stores information of a route in mux
type Route struct {
	Path    string
	Methods []string
	Scopes  []string
	Auth    bool
	Handler func(http.ResponseWriter, *http.Request)
}

// Helper to return 503 if no active cohort is configured for standard users
func (env *Env) withActiveDB(h func(db *sql.DB, w http.ResponseWriter, r *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		db, _, err := env.Manager.GetActiveDB()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusServiceUnavailable)
			_, _ = w.Write([]byte(`{"error": "no_active_cohort", "message": "System is currently offline. Superadmin must configure an active cohort."}`))
			return
		}
		h(db, w, r)
	}
}

// Routes return routes by given env
func (env *Env) Routes() []Route {
	saHandler := &superadmin.Handler{
		Manager: env.Manager,
		Secret:  env.Secret,
	}

	return []Route{
		// ==========================================
		// PUBLIC GLOBAL ENDPOINTS
		// ==========================================
		{
			Path:    "/config",
			Methods: []string{"GET"},
			Scopes:  []string{},
			Auth:    false,
			Handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				active, err := env.Manager.Master().GetActiveCohort()
				if err != nil {
					w.WriteHeader(http.StatusServiceUnavailable)
					_, _ = w.Write([]byte(`{"error": "no_active_cohort", "message": "System is currently offline. Superadmin must configure an active cohort."}`))
					return
				}
				_, _ = w.Write([]byte(active.Config))
			},
		},

		// ==========================================
		// SYSTEM CONFIGURATION ENDPOINTS (Global scope, no active cohort required)
		// ==========================================
		{
			Path:    "/config/login",
			Methods: []string{"POST"},
			Scopes:  []string{},
			Auth:    false,
			Handler: saHandler.Login(),
		},
		{
			Path:    "/config/system-defaults",
			Methods: []string{"GET"},
			Scopes:  []string{"SUPERADMIN"},
			Auth:    true,
			Handler: saHandler.GetSystemDefaults(),
		},
		{
			Path:    "/config/cohorts",
			Methods: []string{"GET"},
			Scopes:  []string{"SUPERADMIN"},
			Auth:    true,
			Handler: saHandler.ListCohorts(),
		},
		{
			Path:    "/config/cohorts",
			Methods: []string{"POST"},
			Scopes:  []string{"SUPERADMIN"},
			Auth:    true,
			Handler: saHandler.CreateCohort(),
		},
		{
			Path:    "/config/cohort/{id}/active",
			Methods: []string{"PUT"},
			Scopes:  []string{"SUPERADMIN"},
			Auth:    true,
			Handler: saHandler.SetCohortActive(),
		},
		{
			Path:    "/config/cohort/{id}/deactivate",
			Methods: []string{"PUT"},
			Scopes:  []string{"SUPERADMIN"},
			Auth:    true,
			Handler: saHandler.DeactivateCohort(),
		},
		{
			Path:    "/config/cohort/{id}/config",
			Methods: []string{"PUT"},
			Scopes:  []string{"SUPERADMIN"},
			Auth:    true,
			Handler: saHandler.UpdateCohortConfig(),
		},
		{
			Path:    "/config/cohort/{id}",
			Methods: []string{"DELETE"},
			Scopes:  []string{"SUPERADMIN"},
			Auth:    true,
			Handler: saHandler.DeleteCohort(),
		},
		{
			Path:    "/config/cohort/{id}/import/{type}",
			Methods: []string{"POST"},
			Scopes:  []string{"SUPERADMIN"},
			Auth:    true,
			Handler: saHandler.ImportData(),
		},
		{
			Path:    "/config/cohort/{id}/user/{userAlias}",
			Methods: []string{"GET"},
			Scopes:  []string{"SUPERADMIN"},
			Auth:    true,
			Handler: saHandler.GetUser(),
		},
		{
			Path:    "/config/cohort/{id}/user/{userAlias}",
			Methods: []string{"PUT"},
			Scopes:  []string{"SUPERADMIN"},
			Auth:    true,
			Handler: saHandler.UpdateUser(),
		},
		{
			Path:    "/config/cohort/{id}/user/{userAlias}",
			Methods: []string{"DELETE"},
			Scopes:  []string{"SUPERADMIN"},
			Auth:    true,
			Handler: saHandler.DeleteUser(),
		},
		{
			Path:    "/config/manage/admins",
			Methods: []string{"GET"},
			Scopes:  []string{"SUPERADMIN"},
			Auth:    true,
			Handler: saHandler.ListSuperadmins(),
		},
		{
			Path:    "/config/manage/admins",
			Methods: []string{"POST"},
			Scopes:  []string{"SUPERADMIN"},
			Auth:    true,
			Handler: saHandler.CreateSuperadmin(),
		},
		{
			Path:    "/config/manage/root/reset",
			Methods: []string{"POST"},
			Scopes:  []string{"SUPERADMIN"},
			Auth:    true,
			Handler: saHandler.ResetRootPassword(),
		},
		{
			Path:    "/config/manage/admins/{username}/password",
			Methods: []string{"PUT"},
			Scopes:  []string{"SUPERADMIN"},
			Auth:    true,
			Handler: saHandler.UpdateSuperadminPassword(),
		},
		{
			Path:    "/config/manage/admins/{username}",
			Methods: []string{"DELETE"},
			Scopes:  []string{"SUPERADMIN"},
			Auth:    true,
			Handler: saHandler.DeleteSuperadmin(),
		},

		// ==========================================
		// STANDARD ENDPOINTS (Routed to active cohort database)
		// ==========================================
		{
			// Login request for standard cohort users
			Path:    "/auth/login",
			Methods: []string{"POST"},
			Scopes:  []string{},
			Auth:    false,
			Handler: env.withActiveDB(func(db *sql.DB, w http.ResponseWriter, r *http.Request) {
				authStore := &modelAuth.DB{DB: db, Secret: env.Secret}
				auth.LoginHandler(authStore)(w, r)
			}),
		},
		{
			Path:    "/auth/refresh/{jwtKeyName}",
			Methods: []string{"GET"},
			Scopes:  []string{},
			Auth:    true,
			Handler: env.withActiveDB(func(db *sql.DB, w http.ResponseWriter, r *http.Request) {
				authStore := &modelAuth.DB{DB: db, Secret: env.Secret}
				auth.RefreshHandler(authStore)(w, r)
			}),
		},
		{
			// Get all students' status
			Path:    "/students",
			Methods: []string{"GET"},
			Scopes:  []string{"TEACHER", "ADMIN", "SUPERADMIN"},
			Auth:    true,
			Handler: env.withActiveDB(func(db *sql.DB, w http.ResponseWriter, r *http.Request) {
				store := &student.DB{DB: db}
				studentHandler.ListHandler(store)(w, r)
			}),
		},
		{
			Path:    "/student/{userAlias}",
			Methods: []string{"GET"},
			Scopes:  []string{"STUDENT", "SUPERADMIN"},
			Auth:    true,
			Handler: env.withActiveDB(func(db *sql.DB, w http.ResponseWriter, r *http.Request) {
				store := &student.DB{DB: db}
				studentHandler.GetHandler(store)(w, r)
			}),
		},
		{
			// Update student's priorities
			Path:    "/student/{userAlias}/priorities",
			Methods: []string{"PUT"},
			Scopes:  []string{"STUDENT", "SUPERADMIN"},
			Auth:    true,
			Handler: env.withActiveDB(func(db *sql.DB, w http.ResponseWriter, r *http.Request) {
				store := &student.DB{DB: db}
				studentHandler.UpdatePrioritiesHandler(store)(w, r)
			}),
		},
		{
			// Get all student's signatures
			Path:    "/signatures",
			Methods: []string{"GET"},
			Scopes:  []string{"TEACHER", "ADMIN", "SUPERADMIN"},
			Auth:    true,
			Handler: env.withActiveDB(func(db *sql.DB, w http.ResponseWriter, r *http.Request) {
				store := &signature.DB{DB: db}
				signatureHandler.ListHandler(store)(w, r)
			}),
		},
		{
			// Update student's signature address
			Path:    "/signature/{userAlias}",
			Methods: []string{"PUT"},
			Scopes:  []string{"STUDENT", "SUPERADMIN"},
			Auth:    true,
			Handler: env.withActiveDB(func(db *sql.DB, w http.ResponseWriter, r *http.Request) {
				store := &signature.DB{DB: db}
				signatureHandler.UpdateAddressHandler(store)(w, r)
			}),
		},
		{
			// Get signature
			Path:    "/signature/{userAlias}",
			Methods: []string{"GET"},
			Scopes:  []string{"STUDENT", "TEACHER", "ADMIN", "SUPERADMIN"},
			Auth:    true,
			Handler: env.withActiveDB(func(db *sql.DB, w http.ResponseWriter, r *http.Request) {
				store := &signature.DB{DB: db}
				signatureHandler.GetHandler(store)(w, r)
			}),
		},
		{
			// Set student's isSigned value
			Path:    "/signature/{userAlias}/issigned/{bool}",
			Methods: []string{"PUT"},
			Scopes:  []string{"STUDENT", "ADMIN", "SUPERADMIN"},
			Auth:    true,
			Handler: env.withActiveDB(func(db *sql.DB, w http.ResponseWriter, r *http.Request) {
				store := &signature.DB{DB: db}
				signatureHandler.UpdateIsSignedHandler(store)(w, r)
			}),
		},
		{
			// Set student's isConfirmed value
			Path:    "/student/{userAlias}/isconfirmed/{bool}",
			Methods: []string{"PUT"},
			Scopes:  []string{"STUDENT", "ADMIN", "SUPERADMIN"},
			Auth:    true,
			Handler: env.withActiveDB(func(db *sql.DB, w http.ResponseWriter, r *http.Request) {
				store := &student.DB{DB: db}
				studentHandler.IsConfirmHandler(store)(w, r)
			}),
		},
		{
			// Set student's isX3 value
			Path:    "/student/{userAlias}/isx3/{bool}",
			Methods: []string{"PUT"},
			Scopes:  []string{"ADMIN", "SUPERADMIN"},
			Auth:    true,
			Handler: env.withActiveDB(func(db *sql.DB, w http.ResponseWriter, r *http.Request) {
				store := &student.DB{DB: db}
				studentHandler.IsX3Handler(store)(w, r)
			}),
		},
		{
			// Set student's rank value
			Path:    "/students/rank",
			Methods: []string{"PUT"},
			Scopes:  []string{"ADMIN", "SUPERADMIN"},
			Auth:    true,
			Handler: env.withActiveDB(func(db *sql.DB, w http.ResponseWriter, r *http.Request) {
				store := &student.DB{DB: db}
				studentHandler.UpdateRankHandler(store)(w, r)
			}),
		},
		{
			// List all subjects information
			Path:    "/subjects",
			Methods: []string{"GET"},
			Scopes:  []string{"ADMIN", "SUPERADMIN"},
			Auth:    true,
			Handler: env.withActiveDB(func(db *sql.DB, w http.ResponseWriter, r *http.Request) {
				store := &subject.DB{DB: db}
				subjectHandler.ListHandler(store)(w, r)
			}),
		},
		{
			// Update subject's capacity
			Path:    "/subject/{subjectCode}/capacity/{capacity}",
			Methods: []string{"PUT"},
			Scopes:  []string{"ADMIN", "SUPERADMIN"},
			Auth:    true,
			Handler: env.withActiveDB(func(db *sql.DB, w http.ResponseWriter, r *http.Request) {
				store := &subject.DB{DB: db}
				subjectHandler.UpdateCapacityHandler(store)(w, r)
			}),
		},
	}
}
