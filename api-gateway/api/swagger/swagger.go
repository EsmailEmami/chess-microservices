package swagger

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
)

func SetupSwagger(r *mux.Router) {

	r.HandleFunc("/swagger/{filename}.json", func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		filename := vars["filename"]

		directory := viper.GetString("app.docs_directory")
		bts, err := os.ReadFile(directory + "/" + filename + ".json")

		if err != nil {
			http.Error(w, "File not found", http.StatusNotFound)
			return
		}

		// Set the content type header
		w.Header().Set("Content-Type", "application/json")

		// Write the JSON content to the response
		w.Write(bts)

	})

	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("user.json"),
	))

}
