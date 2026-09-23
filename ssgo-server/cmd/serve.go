package cmd

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"ssgo-server/model/auth"
	"ssgo-server/model/dbmanager"
	"ssgo-server/route"

	helper "github.com/cyrusn/goHTTPHelper"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start Subject Selection System Backend Server",
	Run: func(cmd *cobra.Command, args []string) {
		auth.UpdateLifeTime(lifeTime)
		checkPathExist(staticFolderLocation)
		
		// Use the directory of dsn as the data directory where master.sqlite and cohort DB files reside
		dataDir := filepath.Dir(dsn)
		mgr, err := dbmanager.New(dataDir)
		if err != nil {
			log.Fatal(err)
		}
		defer mgr.Close()

		env := route.Env{
			Manager: mgr,
			Secret:  &secret,
		}

		Serve(&env)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.PersistentFlags().StringVarP(
		&port,
		"port",
		"p",
		DEFAULT_PORT,
		"port value",
	)
	serveCmd.PersistentFlags().StringVarP(
		&staticFolderLocation,
		"static",
		"s",
		STATIC_FOLDER_LOCATION,
		"location of static folder for serving",
	)
	serveCmd.PersistentFlags().Int64VarP(
		&lifeTime,
		"time",
		"t",
		DEFAULT_LIFE_TIME,
		"update the life time (minutes) of jwt",
	)

	viper.BindPFlags(serveCmd.PersistentFlags())
}

// Serve serve the routers
func Serve(env *route.Env) {
	r := mux.NewRouter()
	
	// Create a dedicated subrouter for all API endpoints to handle /api prefix
	apiRouter := r.PathPrefix("/api").Subrouter()

	routes := env.Routes()
	for _, ro := range routes {
		handler := http.HandlerFunc(ro.Handler)

		// pass Access to handler first
		if len(ro.Scopes) != 0 {
			handler = secret.Access(ro.Scopes, handler).(http.HandlerFunc)
		}

		// then pass Authenticate at last
		if ro.Auth {
			handler = secret.Authenticate(handler).(http.HandlerFunc)
		}

		apiRouter.
			Methods(ro.Methods...).
			Path(ro.Path).
			HandlerFunc(handler)
	}

	serveStaticFolder(r, staticFolderLocation)

	srv := &http.Server{
		Handler: helper.Logger(r),
		Addr:    "0.0.0.0" + port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Bind to a port and pass our router in
	fmt.Println("Available on http://0.0.0.0" + port)
	log.Fatal(srv.ListenAndServe())
}

func serveStaticFolder(r *mux.Router, staticFolderLocation string) {
	staticFolder := http.Dir(staticFolderLocation)

	// Strip the /ss/ prefix used in production deployment subpath
	r.PathPrefix("/ss/").Handler(
		http.StripPrefix("/ss/", http.FileServer(staticFolder)),
	)

	// Serve static files at root / for local testing
	r.PathPrefix("/").Handler(
		http.FileServer(staticFolder),
	)
}
