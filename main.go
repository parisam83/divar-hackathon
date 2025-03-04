package main

import (
	"Realestate-POI/handlers"
	"Realestate-POI/services"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

// type OauthResourceType string

// const (
// 	POST_ADDON_CREATE OauthResourceType = "POST_ADDON_CREATE"
// 	USER_PHONE        OauthResourceType = "USER_PHONE"
// 	OFFLINE_ACCESS    OauthResourceType = "offline_access"
// )

// type Scope struct {
// 	resourceType OauthResourceType
// 	resourceID   string
// }

// func addonOauth(w http.ResponseWriter, r *http.Request) {
// 	post_token := r.URL.Query().Get("post_token")

// 	if post_token == "" {
// 		http.Error(w, "post_token is required", http.StatusBadRequest)
// 		return
// 	}
// 	callback_url := r.URL.Query().Get("return_url") //the adress the user will be redirected after oauth
// 	if callback_url == "" {
// 		http.Error(w, "return_url is required", http.StatusBadRequest)
// 		return
// 	}

// 	// TODO
// 	// create a post with it's token in databse

// 	oauthScopes := []Scope{
// 		{resourceType: POST_ADDON_CREATE, resourceID: post_token},
// 		{resourceType: USER_PHONE},
// 	}
// 	var scopes []string

// 	for _, scope := range oauthScopes {
// 		if scope.resourceID != "" {
// 			scopes = append(scopes, fmt.Sprintf("%s.%s", scope.resourceType, scope.resourceID))
// 		} else {
// 			scopes = append(scopes, string(scope.resourceType))
// 		}
// 	}

// 	conf := &oauth2.Config{
// 		ClientID:     os.Getenv("KENAR_APP_SLUG"),
// 		ClientSecret: os.Getenv("KENAR_OAUTH_SECRET"),
// 		Scopes:       scopes,
// 		Endpoint: oauth2.Endpoint{
// 			AuthURL: "https://api.divar.ir/oauth2/auth",
// 			// TokenURL: "https://oryx-meet-elf.ngrok-free.app/oauth/callback",
// 		},
// 	}
// 	state := uuid.New().String()
// 	url := conf.AuthCodeURL(state, oauth2.AccessTypeOffline)
// 	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)
// 	http.Redirect(w, r, url, http.StatusFound)

//create a specifc session for user
// }
// func oauthCallback(w http.ResponseWriter, r *http.Request) {
// 	log.Printf(r.URL.String())
// 	authorizationCode := r.URL.Query().Get("code")
// 	requestState := r.URL.Query().Get("state")
// 	pic_t := r.URL.Query().Get("scope")
// 	if authorizationCode == "" {
// 		http.Error(w, "code is required", http.StatusBadRequest)
// 		return
// 	}
// 	if requestState == "" {
// 		http.Error(w, "state is required", http.StatusBadRequest)
// 		return
// 	}
// 	// fmt.Printf("code: %s, state: %s", authorizationCode, requestState)
// 	conf := &oauth2.Config{
// 		ClientID:     os.Getenv("KENAR_APP_SLUG"),
// 		ClientSecret: os.Getenv("KENAR_OAUTH_SECRET"),
// 		Endpoint: oauth2.Endpoint{
// 			AuthURL:  "https://api.divar.ir/oauth2/auth",
// 			TokenURL: "https://api.divar.ir/oauth2/token",
// 		},
// 	}
// 	token, err := conf.Exchange(r.Context(), authorizationCode)
// 	if err != nil {
// 		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	fmt.Fprintf(w, "AccessToken: %s\n", token.AccessToken)
// w.Write([]byte(fmt.Sprintf(token.AccessToken)))
// w.Write([]byte("salammmmmmmmmmmmm" + time.Now().Format(time.RFC3339) + "\n"))

// w.Write([]byte(fmt.Sprintf(pic_t) + "\n"))
// re := regexp.MustCompile(`POST_ADDON_CREATE\.([^\s]+)`)
// match := re.FindStringSubmatch(pic_t)
// w.Write([]byte(fmt.Sprintf(match[1]) + "\n"))

// url := fmt.Sprintf("https://api.divar.ir/v1/open-platform/finder/post/%s", match[1])
// req, _ := http.NewRequest("GET", url, nil)
// req.Header.Add("X-Api-Key", "eyJhbGciOiJSUzI1NiIsImtpZCI6InByaXZhdGVfa2V5XzIiLCJ0eXAiOiJKV1QifQ.eyJhcHBfc2x1ZyI6InBsYW5ldC1yaXBwbGUtbGVnZW5kIiwiYXVkIjoic2VydmljZXByb3ZpZGVycyIsImV4cCI6MTc0NTc0OTYyNywianRpIjoiM2FlMzM4NDQtZjQyYy0xMWVmLTlhMTktZmFkODI3M2I1OGM1IiwiaWF0IjoxNzQwNTY1NjI3LCJpc3MiOiJkaXZhciIsInN1YiI6ImFwaWtleSJ9.dml7gkBuE26fXELKUznkOx1WePJ1qZXJyq6i50ZbEgGmmaiNFTlIXTvQSZ_OTjj2sJay9T2iUuNa8uh2tlTFnxtEJSIJsWblzga2_uD8m3RWf76yzBznJmCia3fRkEt8dVbekMzqBg3seDppMzJctuJaVFE0Zhctbm9GFaY2ee1ikxhk65AVLjry6UbEv263Bsk7uQolS49MT7nx0Ij9kMmTrcXfUxECEoj_yFJADsInLkzzNVQNKycfOJdP7D0jDsnOPugYIET9AHZqS0X2KGD_6nz1ugb1QJo-8g0yn22NMTcw_RIvePdLWiStQhENcl5Rf6j2jOUem27JOiG4Dw")

// res, _ := http.DefaultClient.Do(req)

// defer res.Body.Close()
// body, _ := io.ReadAll(res.Body)

// var jsonData map[string]interface{}
// err = json.Unmarshal(body, &jsonData)
// if err != nil {
// 	// return -1
// }
// //fmt.Println(string(body))

// fmt.Println(jsonData["data"].(map[string]interface{})["latitude"])
// fmt.Println(jsonData["data"].(map[string]interface{})["longitude"])

// url := "https://oryx-meet-elf.ngrok-free.app/getlatlong"
// http.Redirect(w, r, url, http.StatusFound)
//redirect to the

// }
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// var configFilePath string
	// var conf *config
	// var err error
	// viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	// viper.AutomaticEnv()

	// // check if config file is not provided
	// if configFilePath != "" {
	// 	viper.SetConfigFile(configFilePath)
	// 	if err := viper.ReadInConfig(); err != nil {
	// 		logrus.Fatalln("failed to read config file.")
	// 	}
	// }
	// conf, err = loadConfig()
	// if err != nil {
	// 	log.Fatal("failed to load configurations: %s", err)
	// }

	oauthService := services.NewOAuthService()
	oauthHandler := handlers.NewOAuthHandler(oauthService)

	r := mux.NewRouter()

	r.HandleFunc("/hello", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello, World!"))
	})
	r.HandleFunc("/addon/oauth", oauthHandler.AddonOauth)
	r.HandleFunc("/oauth/callback", oauthHandler.OauthCallback)
	port := os.Getenv("PORT")
	log.Printf("Server started on port %s", port)
	http.ListenAndServe(":"+port, r)
}
