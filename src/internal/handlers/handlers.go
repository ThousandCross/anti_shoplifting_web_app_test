package handlers

import (
	"anti-shoplifting-webapp/internal/config"
	"anti-shoplifting-webapp/internal/forms"
	"anti-shoplifting-webapp/internal/helpers"
	"anti-shoplifting-webapp/internal/models"
	"anti-shoplifting-webapp/internal/render"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var prefectures []models.Prefectures

// for save jwt token
//var cookies []*http.Cookie
var cookie *http.Cookie

func InitFormData() {

	method := "GET"
	url := "https://anti-shoplifting-dev.cf/api/search/prefectures"
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic("Error")
	}

	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		panic("Error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll err=%s", err.Error())
	}
	json.Unmarshal(body, &prefectures)

	//fmt.Printf("%-v", prefectures)
}

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the handler for the home page
// func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
// 	remoteIP := r.RemoteAddr
// 	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

// 	render.RenderTemplate(w, "home.page.tmpl", &models.TemplateData{})
// }

// About is the handler for the about page
// func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
// 	// perform some logic
// 	stringMap := make(map[string]string)
// 	stringMap["test"] = "Hello, again"

// 	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
// 	stringMap["remote_ip"] = remoteIP

// 	// send data to the template
// 	render.RenderTemplate(w, "about.page.tmpl", &models.TemplateData{
// 		StringMap: stringMap,
// 	})
// }

// PostFcmToken handles request for availability and sends JSON response
func (m *Repository) PostFcmToken(w http.ResponseWriter, r *http.Request) {
	// login api call
	type Request struct {
		FcmToken string `json:"fcm_token"`
	}

	// save as session
	token := r.Form.Get("fcm_token")
	m.App.Session.Put(r.Context(), "fcm_token", token)

	request := &Request{
		FcmToken: token,
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	method := "POST"
	url := "https://anti-shoplifting-dev.cf/api/user/register/fcm"
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonRequest))
	if err != nil {
		// redirect to signin
		m.App.ErrorLog.Println("Http NewRequest Error")
		m.App.Session.Put(r.Context(), "error", "Http NewRequest Error")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		// redirect to signin
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot Regist FcmToken")
		m.App.Session.Put(r.Context(), "error", "Cannot Regist FcmToken")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot retrieve value from response")
		m.App.Session.Put(r.Context(), "error", "Cannot retrieve value from response")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	type jsonResponse struct {
		Result  string `json:"result"`
		Message string `json:"message"`
	}

	var jsonresponse jsonResponse
	json.Unmarshal(body, &jsonresponse)

	out, err := json.MarshalIndent(jsonresponse, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Dashboard is the handler for the dashboard page
func (m *Repository) Dashboard(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "dashboard.page.tmpl", &models.TemplateData{})
}

// Blacklists is the handler for the blacklists page
func (m *Repository) Blacklists(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Blacklists called!!")

	// login api call
	type Request struct {
		From string `json:"from"`
	}

	request := &Request{
		From: "web",
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	method := "POST"
	url := "https://anti-shoplifting-dev.cf/api/user/blacklists"
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonRequest))
	if err != nil {
		// redirect to signin
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Http NewRequest Error")
		m.App.Session.Put(r.Context(), "error", "Http NewRequest Error")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// cookie(jwt token)をリクエストパラメータに追加
	req.AddCookie(cookie)

	req.Header.Set("Content-Type", "application/json")
	// cookieが取得できるように以下を追加
	req.Header.Set("Access-Control-Allow-Credentials", "true")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		// force to logout
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot get blacklists")
		m.App.Session.Put(r.Context(), "error", "Cannot get blacklists")
		// ブラックリストが取得できない場合は強制的にサインアウト
		http.Redirect(w, r, "/signout", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("%+v\n", resp)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// force to logout
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot retrieve value from response")
		m.App.Session.Put(r.Context(), "error", "Cannot retrieve value from response")
		// インシデントが取得できない場合は強制的にサインアウト
		http.Redirect(w, r, "/signout", http.StatusTemporaryRedirect)
		return
	}

	//fmt.Printf("body: %v\n", body)

	var blacklists []models.Blacklists
	json.Unmarshal(body, &blacklists)
	//fmt.Printf("incidents: %v\n", incidents)

	data := make(map[string]interface{})
	data["blacklists"] = blacklists

	render.RenderTemplate(w, r, "blacklists.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// RegisterBlacklist handles blacklist registeration
func (m *Repository) BlacklistRegister(w http.ResponseWriter, r *http.Request) {
	fmt.Println("RegisterBlacklist called!!")

	// login api call
	type Request struct {
		From       string `json:"from"`
		IncidentId string `json:"incident_id"`
		Name       string `json:"name"`
	}

	// not save as session
	incidentId := r.Form.Get("incident_id")
	name := r.Form.Get("name")
	//m.App.Session.Put(r.Context(), "incident_id", incidentId)
	//incidentId, _ := strconv.ParseUint(incidentIdStr, 10, 64)
	request := &Request{
		From:       "web",
		IncidentId: incidentId,
		Name:       name,
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	method := "POST"
	url := "https://anti-shoplifting-dev.cf/api/user/blacklists/register/"
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonRequest))
	if err != nil {
		// redirect to signin
		m.App.ErrorLog.Println("Http NewRequest Error")
		m.App.Session.Put(r.Context(), "error", "Http NewRequest Error")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// cookie(jwt token)をリクエストパラメータに追加
	req.AddCookie(cookie)

	req.Header.Set("Content-Type", "application/json")
	// cookieが取得できるように以下を追加
	req.Header.Set("Access-Control-Allow-Credentials", "true")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		// redirect to signin
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot Regist a record of Blacklist")
		m.App.Session.Put(r.Context(), "error", "Cannot Regist a record of Blacklist")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot retrieve value from response")
		m.App.Session.Put(r.Context(), "error", "Cannot retrieve value from response")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	type jsonResponse struct {
		Result  string `json:"result"`
		Message string `json:"message"`
	}

	var jsonresponse jsonResponse
	json.Unmarshal(body, &jsonresponse)

	out, err := json.MarshalIndent(jsonresponse, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Incidents is the handler for the incidents page
func (m *Repository) Incidents(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Incidents called!!")

	// login api call
	type Request struct {
		From string `json:"from"`
	}

	request := &Request{
		From: "web",
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	method := "POST"
	url := "https://anti-shoplifting-dev.cf/api/user/incidents"
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonRequest))
	if err != nil {
		// redirect to signin
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Http NewRequest Error")
		m.App.Session.Put(r.Context(), "error", "Http NewRequest Error")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// cookie(jwt token)をリクエストパラメータに追加
	req.AddCookie(cookie)

	req.Header.Set("Content-Type", "application/json")
	// cookieが取得できるように以下を追加
	req.Header.Set("Access-Control-Allow-Credentials", "true")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		// force to logout
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot get incidents")
		m.App.Session.Put(r.Context(), "error", "Cannot get incidents")
		// インシデントが取得できない場合は強制的にサインアウト
		http.Redirect(w, r, "/signout", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("%+v\n", resp)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		// force to logout
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot retrieve value from response")
		m.App.Session.Put(r.Context(), "error", "Cannot retrieve value from response")
		// インシデントが取得できない場合は強制的にサインアウト
		http.Redirect(w, r, "/signout", http.StatusTemporaryRedirect)
		return
	}

	//fmt.Printf("body: %v\n", body)

	//var incidents []models.Incidents
	//json.Unmarshal(body, &incidents)
	//fmt.Printf("incidents: %v\n", incidents)

	data := make(map[string]interface{})
	//data["incidents"] = incidents

	// pass to javascript
	data["incidents"] = string(body)

	render.RenderTemplate(w, r, "incidents.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// Incidents is the handler for the incidents page
func (m *Repository) IncidentsByGlobalId(w http.ResponseWriter, r *http.Request) {
	fmt.Println("IncidentsByGlobalId called!!")

	// login api call
	type Request struct {
		From     string `json:"from"`
		GlobalId string `json:"global_id"`
	}

	// not save as session
	globalId := r.Form.Get("global_id")
	fmt.Printf("%-v", globalId)

	request := &Request{
		From:     "web",
		GlobalId: globalId,
	}

	fmt.Printf("%-v", request)

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	method := "POST"
	url := "https://anti-shoplifting-dev.cf/api/user/incidents/global_id"
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonRequest))
	if err != nil {
		// redirect to signin
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Http NewRequest Error")
		m.App.Session.Put(r.Context(), "error", "Http NewRequest Error")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// cookie(jwt token)をリクエストパラメータに追加
	req.AddCookie(cookie)

	req.Header.Set("Content-Type", "application/json")
	// cookieが取得できるように以下を追加
	req.Header.Set("Access-Control-Allow-Credentials", "true")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		// redirect to signin
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot retrieve records of Global Id")
		m.App.Session.Put(r.Context(), "error", "Cannot retrieve records of Global Id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot retrieve value from response")
		m.App.Session.Put(r.Context(), "error", "Cannot retrieve value from response")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var incidents []models.Incidents
	json.Unmarshal(body, &incidents)

	out, err := json.MarshalIndent(incidents, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// IncidentsSalesItem is the handler to get all the sales items reayed to the incident
func (m *Repository) IncidentsSalesItem(w http.ResponseWriter, r *http.Request) {
	fmt.Println("IncidentsSalesItem called!!")

	// login api call
	type Request struct {
		From       string `json:"from"`
		IncidentId string `json:"incident_id"`
	}

	// not save as session
	incidentId := r.Form.Get("incident_id")
	fmt.Printf("%-v", incidentId)

	request := &Request{
		From:       "web",
		IncidentId: incidentId,
	}

	fmt.Printf("%-v", request)

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	method := "POST"
	url := "https://anti-shoplifting-dev.cf/api/user/incidents/sales_items"
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonRequest))
	if err != nil {
		// redirect to signin
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Http NewRequest Error")
		m.App.Session.Put(r.Context(), "error", "Http NewRequest Error")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// cookie(jwt token)をリクエストパラメータに追加
	req.AddCookie(cookie)

	req.Header.Set("Content-Type", "application/json")
	// cookieが取得できるように以下を追加
	req.Header.Set("Access-Control-Allow-Credentials", "true")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		// redirect to signin
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot retrieve records of Global Id")
		m.App.Session.Put(r.Context(), "error", "Cannot retrieve records of Global Id")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot retrieve value from response")
		m.App.Session.Put(r.Context(), "error", "Cannot retrieve value from response")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	var IncidentSalesItem []models.IncidentSalesItem
	json.Unmarshal(body, &IncidentSalesItem)

	out, err := json.MarshalIndent(IncidentSalesItem, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Settings is the handler for the settings page
func (m *Repository) Settings(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "settings.page.tmpl", &models.TemplateData{})
}

// Signin is the handler for the signin page
func (m *Repository) Signin(w http.ResponseWriter, r *http.Request) {
	var emptySignin models.Signin
	data := make(map[string]interface{})

	signin, ok := m.App.Session.Get(r.Context(), "signin").(models.Signin)
	if !ok {
		data["signin"] = emptySignin
	} else {
		data["signin"] = signin
	}

	render.RenderTemplate(w, r, "signin.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostSignin is the handler for the signin page
func (m *Repository) PostSignin(w http.ResponseWriter, r *http.Request) {
	// m.App.Session.RenewToken以前に実行してfcm_tokenを保持
	fcmToken := m.App.Session.GetString(r.Context(), "fcm_token")
	if len(fcmToken) == 0 {
		// redirect to signin
		m.App.ErrorLog.Println("Cannot Login")
		m.App.Session.Put(r.Context(), "error", "Can't Login")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// セッション初期化
	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	remember_me := false
	if r.Form.Get("remember_me") == "1" {
		remember_me = true
	}

	signin := models.Signin{
		CompanyCd:  r.Form.Get("company_code"),
		StoreCd:    r.Form.Get("store_code"),
		Password:   r.Form.Get("password"),
		RememberMe: remember_me,
	}

	form := forms.New(r.PostForm)
	// validation (move to api server inthe future)
	form.Required("company_code", "store_code", "password")

	// login api call
	type Request struct {
		From      string `json:"from"`
		CompanyCd string `json:"company_cd"`
		StoreCd   string `json:"store_cd"`
		Password  string `json:"password"`
		FcmToken  string `json:"fcm_token"`
	}

	request := &Request{
		From:      "web",
		CompanyCd: r.Form.Get("company_code"),
		StoreCd:   r.Form.Get("store_code"),
		Password:  r.Form.Get("password"),
		FcmToken:  fcmToken,
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	method := "POST"
	url := "https://anti-shoplifting-dev.cf/api/user/login"
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonRequest))
	if err != nil {
		// redirect to signin
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Http NewRequest Error")
		m.App.Session.Put(r.Context(), "error", "Http NewRequest Error")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	// cookieが取得できるように以下を追加
	req.Header.Set("Access-Control-Allow-Credentials", "true")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		// redirect to signin
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot Login")
		m.App.Session.Put(r.Context(), "error", "Can't Login")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("%+v\n", resp)
	defer resp.Body.Close()

	// save jwt token to gloabal cookie and session
	for _, Cookie := range resp.Cookies() {
		fmt.Println("Found a cookie named:", Cookie.Name)
		if Cookie.Name == "jwt" {
			fmt.Println("save jwt token to gloabal cookie and session!!")
			m.App.Session.Put(r.Context(), "jwt", Cookie.Value)
			cookie = Cookie
		}
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot retrieve value from response")
		m.App.Session.Put(r.Context(), "error", "Cannot retrieve value from response")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// struct for api response
	type jsonResponse struct {
		Result  string `json:"result"`
		Message string `json:"message"`
		Jwt     string `json:"jwt"`
	}

	var jsonresponse jsonResponse
	json.Unmarshal(body, &jsonresponse)

	// redirect signin.tmpol if api returning message is "ng"
	if jsonresponse.Result == "ng" {
		m.App.ErrorLog.Println("Cannot Login due to response ng from API")
		m.App.Session.Put(r.Context(), "error", "Cannot Login due to response ng from API")
		data := make(map[string]interface{})
		data["signin"] = signin
		render.RenderTemplate(w, r, "signin.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// redirect signin.tmpol if form parameter is invalid
	if !form.Valid() {
		data := make(map[string]interface{})
		data["signin"] = signin
		render.RenderTemplate(w, r, "signin.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	if r.Form.Get("remember_me") == "1" {
		// remember_me check on 時の処理
		m.App.Session.Put(r.Context(), "signin", signin)
	} else {
		// remember_me check off 時の処理
		m.App.Session.Remove(r.Context(), "signin")
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// Signout sign user out
func (m *Repository) Signout(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Signout called!!")

	// login api call
	type Request struct {
		From string `json:"from"`
	}

	request := &Request{
		From: "web",
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	method := "POST"
	url := "https://anti-shoplifting-dev.cf/api/user/logout"
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonRequest))
	if err != nil {
		// redirect to signin
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Http NewRequest Error")
		m.App.Session.Put(r.Context(), "error", "Http NewRequest Error")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// cookie(jwt token)をリクエストパラメータに追加
	req.AddCookie(cookie)

	req.Header.Set("Content-Type", "application/json")
	// cookieが取得できるように以下を追加
	req.Header.Set("Access-Control-Allow-Credentials", "true")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		// redirect to signin
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot Logout")
		m.App.Session.Put(r.Context(), "error", "Can't Logout")
		// ログアウトできない場合はリダイレクトさせなくて良い
		//http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("%+v\n", resp)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot retrieve value from response")
		m.App.Session.Put(r.Context(), "error", "Cannot retrieve value from response")
		// ログアウトできない場合はリダイレクトさせなくて良い
		//http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// struct for api response
	type jsonResponse struct {
		Result  string `json:"result"`
		Message string `json:"message"`
	}

	var jsonresponse jsonResponse
	json.Unmarshal(body, &jsonresponse)

	// redirect signin.tmpol if api returning message is "ng"
	if jsonresponse.Result == "ng" {
		m.App.ErrorLog.Println("Cannot Login due to response ng from API")
		m.App.Session.Put(r.Context(), "error", "Cannot Login due to response ng from API")
		// ログアウトできない場合はリダイレクトさせなくて良い
		//http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// フロント側のセッション削除処理
	//
	companyCd := ""
	storeCd := ""
	password := ""
	rememberMe := false

	// サインイン時にremeber meにチェックON時のみ取得可能
	signin, ok := m.App.Session.Get(r.Context(), "signin").(models.Signin)
	if ok {
		companyCd = signin.CompanyCd
		storeCd = signin.StoreCd
		password = signin.Password
		rememberMe = signin.RememberMe
	}

	// セッション情報全削除
	_ = m.App.Session.Destroy(r.Context())
	_ = m.App.Session.RenewToken(r.Context())

	// 次回ログイン用にセッション情報保存
	if ok {
		signinNext := models.Signin{
			CompanyCd:  companyCd,
			StoreCd:    storeCd,
			Password:   password,
			RememberMe: rememberMe,
		}
		m.App.Session.Put(r.Context(), "signin", signinNext)
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Signup is the handler for the signup page
func (m *Repository) Signup(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "signup.page.tmpl", &models.TemplateData{})
}

// SignupCompany is the handler for the signup-company page
func (m *Repository) SignupCompany(w http.ResponseWriter, r *http.Request) {
	var emptyCompanyBasicInfo models.CompanyBasicInfo
	data := make(map[string]interface{})
	data["basicinfo"] = emptyCompanyBasicInfo
	data["prefectures"] = prefectures

	render.RenderTemplate(w, r, "company-registration-input-basic.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// SignupCompanyBasicInfo is the handler for basic information input of signup-company page
func (m *Repository) PostSignupCompanyBasicInfo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	prefecture_id, _ := strconv.Atoi(r.Form.Get("prefecture"))
	prefecture := ""
	if prefecture_id > 0 {
		prefecture = prefectures[prefecture_id-1].Name
	}

	basicinfo := models.CompanyBasicInfo{
		CompanyName:                  r.Form.Get("company_name"),
		RepresentativeFamilyName:     r.Form.Get("representative_family_name"),
		RepresentativeFirstName:      r.Form.Get("representative_first_name"),
		RepresentativeFamilyNameKana: r.Form.Get("representative_family_name_kana"),
		RepresentativeFirstNameKana:  r.Form.Get("representative_first_name_kana"),
		Zipcode:                      r.Form.Get("zipcode"),
		Prefecture:                   prefecture,
		PrefectureId:                 uint(prefecture_id),
		City:                         r.Form.Get("city"),
		Street:                       r.Form.Get("street"),
		Building:                     r.Form.Get("building"),
		Tel:                          r.Form.Get("tel"),
		Mail:                         r.Form.Get("mail"),
	}

	form := forms.New(r.PostForm)

	//fmt.Printf("%-v", form)

	// form.Has("company_name", r)
	form.Required(
		"company_name",
		"representative_family_name",
		"representative_first_name",
		"representative_family_name_kana",
		"representative_first_name_kana",
		"zipcode",
		"prefecture",
		"city",
		"street",
		//"building",
		"tel",
		"mail",
	)

	form.MaxLength("company_name", 50, r)
	form.MaxLength("representative_family_name", 30, r)
	form.MaxLength("representative_first_name", 30, r)
	form.MaxLength("representative_family_name_kana", 30, r)
	form.MaxLength("representative_first_name_kana", 30, r)
	form.MaxLength("street", 50, r)
	form.MinLength("tel", 10, r)
	form.MaxLength("tel", 11, r)
	form.IsEmail("mail")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["basicinfo"] = basicinfo
		data["prefectures"] = prefectures
		render.RenderTemplate(w, r, "company-registration-input-basic.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "company_basicinfo", basicinfo)
	http.Redirect(w, r, "/signup/company/basicinfo/contd", http.StatusSeeOther)
}

// SignupCompanyBasicInfoContd is the handler for basic information continued input of signup-company page
func (m *Repository) SignupCompanyBasicInfoContd(w http.ResponseWriter, r *http.Request) {
	// how to retriev session data!!!
	// reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	// if !ok {
	// 	log.Println("Cannot get item from session")
	// 	m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }
	// 削除するタイミングは再検討すること
	//m.App.Session.Remove(r.Context(), "reservation")

	// data := make(map[string]interface{})
	// data["reservation"] = reservation

	// render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
	// 	Data: data,
	// })
	var emptyCompanyBasicInfoContd models.CompanyBasicInfoContd
	data := make(map[string]interface{})
	data["basicinfocontd"] = emptyCompanyBasicInfoContd

	render.RenderTemplate(w, r, "company-registration-input-basic-contd.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostSignupCompanyBasicInfoContd is the handler for posting basic information continued input of signup-company page
func (m *Repository) PostSignupCompanyBasicInfoContd(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	basicinfocontd := models.CompanyBasicInfoContd{
		ManagerFamilyName:     r.Form.Get("manager_family_name"),
		ManagerFirstName:      r.Form.Get("manager_first_name"),
		ManagerFamilyNameKana: r.Form.Get("manager_family_name_kana"),
		ManagerFirstNameKana:  r.Form.Get("manager_first_name_kana"),
		ManagerTel:            r.Form.Get("manager_tel"),
		ManagerMail:           r.Form.Get("manager_mail"),
	}

	form := forms.New(r.PostForm)

	//form.Has("manager_family_name", r)
	form.Required(
		"manager_family_name",
		"manager_first_name",
		"manager_family_name_kana",
		"manager_first_name_kana",
		"manager_tel",
		"manager_mail",
	)

	form.MaxLength("manager_family_name", 30, r)
	form.MaxLength("manager_first_name", 30, r)
	form.MaxLength("manager_family_name_kana", 30, r)
	form.MaxLength("manager_first_name_kana", 30, r)
	form.MinLength("manager_tel", 10, r)
	form.MaxLength("manager_tel", 11, r)
	form.IsEmail("manager_mail")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["basicinfocontd"] = basicinfocontd
		render.RenderTemplate(w, r, "company-registration-input-basic-contd.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "company_basicinfocontd", basicinfocontd)
	http.Redirect(w, r, "/signup/company/payment", http.StatusSeeOther)
}

// SignupCompanyPayment is the handler for payment input of signup-company page
func (m *Repository) SignupCompanyPayment(w http.ResponseWriter, r *http.Request) {
	// how to retriev session data!!!
	// reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	// if !ok {
	// 	log.Println("Cannot get item from session")
	// 	m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }
	// 削除するタイミングは再検討すること
	//m.App.Session.Remove(r.Context(), "reservation")

	// data := make(map[string]interface{})
	// data["reservation"] = reservation

	// render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
	// 	Data: data,
	// })
	var emptyCompanyPayment models.CompanyPayment
	data := make(map[string]interface{})
	data["payment"] = emptyCompanyPayment

	render.RenderTemplate(w, r, "company-registration-input-payment.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// SignupCompanyPayment is the handler for payment input of signup-company page
func (m *Repository) PostSignupCompanyPayment(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	payment := models.CompanyPayment{
		CardNo:                   r.Form.Get("card_no"),
		CardHolderFamilyNameKana: r.Form.Get("card_holder_family_name_kana"),
		CardHolderFirstNameKana:  r.Form.Get("card_holder_first_name_kana"),
		CardMonth:                r.Form.Get("card_month"),
		CardYear:                 r.Form.Get("card_year"),
		SecurityCd:               r.Form.Get("security_cd"),
	}

	form := forms.New(r.PostForm)

	//form.Has("card_no", r)
	form.Required(
		"card_no",
		"card_holder_family_name_kana",
		"card_holder_first_name_kana",
		"card_month",
		"card_year",
		"security_cd",
	)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["payment"] = payment
		render.RenderTemplate(w, r, "company-registration-input-payment.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "company_payment", payment)
	http.Redirect(w, r, "/signup/company/confirm", http.StatusSeeOther)
}

// SignupCompanyConfirm is the handler for confirm of signup-company page
func (m *Repository) SignupCompanyConfirm(w http.ResponseWriter, r *http.Request) {
	// get basicinfo
	basicinfo, ok := m.App.Session.Get(r.Context(), "company_basicinfo").(models.CompanyBasicInfo)
	if !ok {
		m.App.ErrorLog.Println("Cannot get basicinfo from session")
		m.App.Session.Put(r.Context(), "error", "Can't get basicinfo from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// get basicinfocontd
	basicinfocontd, ok := m.App.Session.Get(r.Context(), "company_basicinfocontd").(models.CompanyBasicInfoContd)
	if !ok {
		m.App.ErrorLog.Println("Cannot get basicinfocontd from session")
		m.App.Session.Put(r.Context(), "error", "Can't get basicinfocontd from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// get payment
	payment, ok := m.App.Session.Get(r.Context(), "company_payment").(models.CompanyPayment)
	if !ok {
		m.App.ErrorLog.Println("Cannot get payment from session")
		m.App.Session.Put(r.Context(), "error", "Can't get payment from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// have to cast prefecture id from int to string
	prefectureId := strconv.Itoa(int(basicinfo.PrefectureId))
	companySignup := models.CompanySignupConfirm{
		From:                         "web",
		CompanyName:                  basicinfo.CompanyName,
		RepresentativeFamilyName:     basicinfo.RepresentativeFamilyName,
		RepresentativeFirstName:      basicinfo.RepresentativeFirstName,
		RepresentativeFamilyNameKana: basicinfo.RepresentativeFamilyNameKana,
		RepresentativeFirstNameKana:  basicinfo.RepresentativeFirstNameKana,
		Zipcode:                      basicinfo.Zipcode,
		Prefecture:                   basicinfo.Prefecture,
		PrefectureId:                 prefectureId,
		City:                         basicinfo.City,
		Street:                       basicinfo.Street,
		Building:                     basicinfo.Building,
		Tel:                          basicinfo.Tel,
		Mail:                         basicinfo.Mail,
		ManagerFamilyName:            basicinfocontd.ManagerFamilyName,
		ManagerFirstName:             basicinfocontd.ManagerFirstName,
		ManagerFamilyNameKana:        basicinfocontd.ManagerFamilyNameKana,
		ManagerFirstNameKana:         basicinfocontd.ManagerFirstNameKana,
		ManagerTel:                   basicinfocontd.ManagerTel,
		ManagerMail:                  basicinfocontd.ManagerMail,
		CardNo:                       payment.CardNo,
		CardHolderFamilyNameKana:     payment.CardHolderFamilyNameKana,
		CardHolderFirstNameKana:      payment.CardHolderFirstNameKana,
		CardMonth:                    payment.CardMonth,
		CardYear:                     payment.CardYear,
		SecurityCd:                   payment.SecurityCd,
	}
	m.App.Session.Put(r.Context(), "company_signup_confirm", companySignup)
	data := make(map[string]interface{})
	data["confirm"] = companySignup

	render.RenderTemplate(w, r, "company-registration-confirm.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// PostSignupCompanyConfirm is the handler for confirm of signup-company page
func (m *Repository) PostSignupCompanyConfirm(w http.ResponseWriter, r *http.Request) {
	// get company_signup
	request, ok := m.App.Session.Get(r.Context(), "company_signup_confirm").(models.CompanySignupConfirm)
	if !ok {
		m.App.ErrorLog.Println("Cannot get company signup confirm from session")
		m.App.Session.Put(r.Context(), "error", "Can't get company signup confirm from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	method := "POST"
	url := "https://anti-shoplifting-dev.cf/api/user/register/company"
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonRequest))
	if err != nil {
		// redirect to signin
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Http NewRequest Error")
		m.App.Session.Put(r.Context(), "error", "Http NewRequest Error")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		// redirect to signin
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot Regist Company")
		m.App.Session.Put(r.Context(), "error", "Cannot Regist Company")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("%+v\n", resp)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot retrieve value from response")
		m.App.Session.Put(r.Context(), "error", "Cannot retrieve value from response")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// struct for api response
	type jsonResponse struct {
		Result  string `json:"result"`
		Message string `json:"message"`
		//CompanyKey        string `json:"company_key"`
		//CompanyManagerKey string `json:"company_manager_key"`
	}

	var jsonresponse jsonResponse
	json.Unmarshal(body, &jsonresponse)

	// redirect signin.tmpol if api returning message is "ng"
	if jsonresponse.Result == "ng" {
		m.App.ErrorLog.Println("Error occuered while company registration!")
		m.App.Session.Put(r.Context(), "error", "Error occuered while company registration!")

		// session delete and redirect to signup
		m.App.Session.Remove(r.Context(), "company_basicinfo")
		m.App.Session.Remove(r.Context(), "company_basicinfocontd")
		m.App.Session.Remove(r.Context(), "company_payment")
		m.App.Session.Remove(r.Context(), "company_signup_confirm")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//m.App.Session.Put(r.Context(), "company_key", jsonresponse.CompanyKey)
	//m.App.Session.Put(r.Context(), "company_manager_key", jsonresponse.CompanyManagerKey)
	http.Redirect(w, r, "/signup/company/complete", http.StatusSeeOther)
}

// SignupCompanyComplete is the handler for complete of signup-company page
func (m *Repository) SignupCompanyComplete(w http.ResponseWriter, r *http.Request) {
	// get session before being deleted
	// companyKey := m.App.Session.GetString(r.Context(), "company_key")
	// companyManagerKey := m.App.Session.GetString(r.Context(), "company_manager_key")

	// if len(companyKey) == 0 || len(companyManagerKey) == 0 {
	// 	m.App.ErrorLog.Println("Cannot get company key information from session")
	// 	m.App.Session.Put(r.Context(), "error", "Can't get company key information from session")
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	// session delete and redirect to signup
	m.App.Session.Remove(r.Context(), "company_basicinfo")
	m.App.Session.Remove(r.Context(), "company_basicinfocontd")
	m.App.Session.Remove(r.Context(), "company_payment")
	m.App.Session.Remove(r.Context(), "company_signup_confirm")
	// m.App.Session.Remove(r.Context(), "company_key")
	// m.App.Session.Remove(r.Context(), "company_manager_key")

	// stringMap := make(map[string]string)
	// stringMap["company_key"] = companyKey
	// stringMap["company_manager_key"] = companyManagerKey

	render.RenderTemplate(w, r, "company-registration-success-temporary.page.tmpl", &models.TemplateData{
		//Form:      forms.New(nil),
		//StringMap: stringMap,
	})
}

// SignupCompanyVeryfyEmail is the handler for email verification and show company cd to user
func (m *Repository) SignupCompanyVeryfyEmail(w http.ResponseWriter, r *http.Request) {

	key1 := r.URL.Query().Get("key1")
	key2 := r.URL.Query().Get("key2")

	method := "GET"
	url := "https://anti-shoplifting-dev.cf/api/user/register/company/verify?key1=" + key1 + "&key2=" + key2
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic("Error")
	}

	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		panic("Error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll err=%s", err.Error())
	}

	// struct for api response
	type jsonResponse struct {
		Result    string `json:"result"`
		Message   string `json:"message"`
		CompanyCd string `json:"company_cd"`
	}

	var jsonresponse jsonResponse
	json.Unmarshal(body, &jsonresponse)

	// redirect signin.tmpol if api returning message is "ng"
	if jsonresponse.Result == "ng" {
		m.App.ErrorLog.Println("Error occuered while compnay email verification!")
		m.App.Session.Put(r.Context(), "error", "Error occuered while compnay email verification!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	stringMap := make(map[string]string)
	stringMap["company_cd"] = jsonresponse.CompanyCd

	render.RenderTemplate(w, r, "company-registration-success.page.tmpl", &models.TemplateData{
		//Form:      forms.New(nil),
		StringMap: stringMap,
	})
}

// SignupStore is the handler for the signup-store page
func (m *Repository) SignupStore(w http.ResponseWriter, r *http.Request) {
	var emptyStoreBasicInfo models.StoreBasicInfo
	data := make(map[string]interface{})
	data["basicinfo"] = emptyStoreBasicInfo
	data["prefectures"] = prefectures

	render.RenderTemplate(w, r, "store-registration-input-basic.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// SignupStoreBasicInfo is the handler for basic information input of signup-company page
func (m *Repository) SignupStoreBasicInfo(w http.ResponseWriter, r *http.Request) {
	var emptyStoreBasicInfo models.StoreBasicInfo
	data := make(map[string]interface{})
	data["basicinfo"] = emptyStoreBasicInfo
	data["prefectures"] = prefectures

	render.RenderTemplate(w, r, "store-registration-input-basic.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostSignupStoreBasicInfo is the handler for basic information post input of signup-store page
func (m *Repository) PostSignupStoreBasicInfo(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	prefecture_id, _ := strconv.Atoi(r.Form.Get("prefecture"))
	prefecture := ""
	if prefecture_id > 0 {
		prefecture = prefectures[prefecture_id-1].Name
	}
	basicinfo := models.StoreBasicInfo{
		StoreName:    r.Form.Get("store_name"),
		CompanyKey:   r.Form.Get("company_key"),
		CompanyCd:    r.Form.Get("company_cd"),
		Zipcode:      r.Form.Get("zipcode"),
		Prefecture:   prefecture,
		PrefectureId: uint(prefecture_id),
		City:         r.Form.Get("city"),
		Street:       r.Form.Get("street"),
		Building:     r.Form.Get("building"),
	}

	form := forms.New(r.PostForm)

	//form.Has("store_name", r)
	form.Required(
		"store_name",
		"company_key",
		"company_cd",
		"zipcode",
		"prefecture",
		"city",
		"street",
		//"building",
	)

	form.MaxLength("street", 50, r)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["basicinfo"] = basicinfo
		data["prefectures"] = prefectures
		render.RenderTemplate(w, r, "store-registration-input-basic.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "store_basicinfo", basicinfo)
	http.Redirect(w, r, "/signup/store/basicinfo/contd", http.StatusSeeOther)
}

// SignupStoreBasicInfoContd is the handler for basic information continued input of signup-store page
func (m *Repository) SignupStoreBasicInfoContd(w http.ResponseWriter, r *http.Request) {
	var emptyStoreBasicInfoContd models.StoreBasicInfoContd
	data := make(map[string]interface{})
	data["basicinfocontd"] = emptyStoreBasicInfoContd

	render.RenderTemplate(w, r, "store-registration-input-basic-contd.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostSignupStoreBasicInfoContd is the handler for basic information continued input post of signup-store page
func (m *Repository) PostSignupStoreBasicInfoContd(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
		return
	}

	basicinfocontd := models.StoreBasicInfoContd{
		ManagerFamilyName:     r.Form.Get("manager_family_name"),
		ManagerFirstName:      r.Form.Get("manager_first_name"),
		ManagerFamilyNameKana: r.Form.Get("manager_family_name_kana"),
		ManagerFirstNameKana:  r.Form.Get("manager_first_name_kana"),
		ManagerTel:            r.Form.Get("manager_tel"),
		ManagerMail:           r.Form.Get("manager_mail"),
	}

	form := forms.New(r.PostForm)

	//form.Has("manager_family_name", r)
	form.Required(
		"manager_family_name",
		"manager_first_name",
		"manager_family_name_kana",
		"manager_first_name_kana",
		"manager_tel",
		"manager_mail",
	)

	form.MaxLength("manager_family_name", 30, r)
	form.MaxLength("manager_first_name", 30, r)
	form.MaxLength("manager_family_name_kana", 30, r)
	form.MaxLength("manager_first_name_kana", 30, r)
	form.MinLength("manager_tel", 10, r)
	form.MaxLength("manager_tel", 11, r)
	form.IsEmail("manager_mail")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["basicinfocontd"] = basicinfocontd
		render.RenderTemplate(w, r, "store-registration-input-basic-contd.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "store_basicinfocontd", basicinfocontd)
	http.Redirect(w, r, "/signup/store/password", http.StatusSeeOther)
}

// SignupStorePassword is the handler for password input of signup-store page
func (m *Repository) SignupStorePassword(w http.ResponseWriter, r *http.Request) {
	var emptyStorePassword models.StorePassword
	data := make(map[string]interface{})
	data["password"] = emptyStorePassword

	render.RenderTemplate(w, r, "store-registration-input-password.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
		Data: data,
	})
}

// PostSignupStorePassword is the handler for password input post of signup-store page
func (m *Repository) PostSignupStorePassword(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	password := models.StorePassword{
		Password:        r.Form.Get("password"),
		PasswordConfirm: r.Form.Get("password_confirm"),
	}

	form := forms.New(r.PostForm)

	//form.Has("password", r)
	form.Required(
		"password",
		"password_confirm",
	)

	if !form.Valid() {
		data := make(map[string]interface{})
		data["password"] = password
		render.RenderTemplate(w, r, "store-registration-input-password.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	m.App.Session.Put(r.Context(), "store_password", password)
	http.Redirect(w, r, "/signup/store/confirm", http.StatusSeeOther)
}

// SignupStoreConfirm is the handler for confirm of signup-store page
func (m *Repository) SignupStoreConfirm(w http.ResponseWriter, r *http.Request) {
	// get basicinfo
	basicinfo, ok := m.App.Session.Get(r.Context(), "store_basicinfo").(models.StoreBasicInfo)
	if !ok {
		m.App.ErrorLog.Println("Cannot get basicinfo from session")
		m.App.Session.Put(r.Context(), "error", "Can't get basicinfo from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// get basicinfocontd
	basicinfocontd, ok := m.App.Session.Get(r.Context(), "store_basicinfocontd").(models.StoreBasicInfoContd)
	if !ok {
		m.App.ErrorLog.Println("Cannot get basicinfocontd from session")
		m.App.Session.Put(r.Context(), "error", "Can't get basicinfocontd from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// get password
	password, ok := m.App.Session.Get(r.Context(), "store_password").(models.StorePassword)
	if !ok {
		m.App.ErrorLog.Println("Cannot get password from session")
		m.App.Session.Put(r.Context(), "error", "Can't get password from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// have to cast prefecture id from int to string
	prefectureId := strconv.Itoa(int(basicinfo.PrefectureId))
	storeSignup := models.StoreSignupConfirm{
		From:                  "web",
		CompanyKey:            basicinfo.CompanyKey,
		CompanyCd:             basicinfo.CompanyCd,
		StoreName:             basicinfo.StoreName,
		Zipcode:               basicinfo.Zipcode,
		Prefecture:            basicinfo.Prefecture,
		PrefectureId:          prefectureId, // cast from int to string
		City:                  basicinfo.City,
		Street:                basicinfo.Street,
		Building:              basicinfo.Building,
		ManagerFamilyName:     basicinfocontd.ManagerFamilyName,
		ManagerFirstName:      basicinfocontd.ManagerFirstName,
		ManagerFamilyNameKana: basicinfocontd.ManagerFamilyNameKana,
		ManagerFirstNameKana:  basicinfocontd.ManagerFirstNameKana,
		ManagerTel:            basicinfocontd.ManagerTel,
		ManagerMail:           basicinfocontd.ManagerMail,
		Password:              password.Password,
		PasswordConfirm:       password.PasswordConfirm,
	}

	m.App.Session.Put(r.Context(), "store_signup_confirm", storeSignup)
	data := make(map[string]interface{})
	data["confirm"] = storeSignup

	render.RenderTemplate(w, r, "store-registration-confirm.page.tmpl", &models.TemplateData{
		Data: data,
	})
}

// PostSignupStoreConfirm is the handler for confirm of signup-store page
func (m *Repository) PostSignupStoreConfirm(w http.ResponseWriter, r *http.Request) {
	// get store_signup_confirm
	request, ok := m.App.Session.Get(r.Context(), "store_signup_confirm").(models.StoreSignupConfirm)
	if !ok {
		m.App.ErrorLog.Println("Cannot get store signup confirm from session")
		m.App.Session.Put(r.Context(), "error", "Can't get store signup confirm from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	method := "POST"
	url := "https://anti-shoplifting-dev.cf/api/user/register/store"
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonRequest))
	if err != nil {
		// redirect to signin
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Http NewRequest Error")
		m.App.Session.Put(r.Context(), "error", "Http NewRequest Error")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		// redirect to signin
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot Regist Store")
		m.App.Session.Put(r.Context(), "error", "Cannot Regist Store")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("%+v\n", resp)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot retrieve value from response")
		m.App.Session.Put(r.Context(), "error", "Cannot retrieve value from response")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// struct for api response
	type jsonResponse struct {
		Result  string `json:"result"`
		Message string `json:"message"`
		// StoreKey string `json:"store_key"`
	}

	var jsonresponse jsonResponse
	json.Unmarshal(body, &jsonresponse)

	// redirect signin.tmpol if api returning message is "ng"
	if jsonresponse.Result == "ng" {
		m.App.ErrorLog.Println("Error occuered while store registration!")
		m.App.Session.Put(r.Context(), "error", "Error occuered while store registration!")

		// session delete and redirect to signup
		m.App.Session.Remove(r.Context(), "store_basicinfo")
		m.App.Session.Remove(r.Context(), "store_basicinfocontd")
		m.App.Session.Remove(r.Context(), "store_password")
		m.App.Session.Remove(r.Context(), "store_signup_confirm")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	//m.App.Session.Put(r.Context(), "store_key", jsonresponse.StoreKey)
	http.Redirect(w, r, "/signup/store/complete", http.StatusSeeOther)
}

// SignupCompanyComplete is the handler for complete of signup-store page
func (m *Repository) SignupStoreComplete(w http.ResponseWriter, r *http.Request) {
	// get session before being deleted
	//storeKey := m.App.Session.GetString(r.Context(), "store_key")
	// if len(storeKey) == 0 {
	// 	m.App.ErrorLog.Println("Cannot get store key information from session")
	// 	m.App.Session.Put(r.Context(), "error", "Can't get store key information from session")
	// 	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	// 	return
	// }

	// session delete and redirect to signup
	m.App.Session.Remove(r.Context(), "store_basicinfo")
	m.App.Session.Remove(r.Context(), "store_basicinfocontd")
	m.App.Session.Remove(r.Context(), "store_password")
	m.App.Session.Remove(r.Context(), "store_signup_confirm")
	//m.App.Session.Remove(r.Context(), "store_key")

	// stringMap := make(map[string]string)
	// stringMap["store_key"] = storeKey

	render.RenderTemplate(w, r, "store-registration-success-temporary.page.tmpl", &models.TemplateData{
		//Form:      forms.New(nil),
		//StringMap: stringMap,
	})
}

// SignupStoreVeryfyEmail is the handler for email verification and show company cd to user
func (m *Repository) SignupStoreVeryfyEmail(w http.ResponseWriter, r *http.Request) {

	key1 := r.URL.Query().Get("key1")
	key2 := r.URL.Query().Get("key2")

	method := "GET"
	url := "https://anti-shoplifting-dev.cf/api/user/register/store/verify?key1=" + key1 + "&key2=" + key2
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic("Error")
	}

	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		panic("Error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll err=%s", err.Error())
	}

	// struct for api response
	type jsonResponse struct {
		Result  string `json:"result"`
		Message string `json:"message"`
	}

	var jsonresponse jsonResponse
	json.Unmarshal(body, &jsonresponse)

	// redirect signin.tmpol if api returning message is "ng"
	if jsonresponse.Result == "ng" {
		m.App.ErrorLog.Println("Error occuered while compnay email verification!")
		m.App.Session.Put(r.Context(), "error", "Error occuered while compnay email verification!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	render.RenderTemplate(w, r, "store-registration-success.page.tmpl", &models.TemplateData{
		//Form:      forms.New(nil),
		//StringMap: stringMap,
	})
}

// SignupStoreApproveNewStore is the handler for new store approval
func (m *Repository) SignupStoreApproveNewStore(w http.ResponseWriter, r *http.Request) {

	key1 := r.URL.Query().Get("key1")
	key2 := r.URL.Query().Get("key2")
	key3 := r.URL.Query().Get("key3")

	method := "GET"
	url := "https://anti-shoplifting-dev.cf/api/user/register/store/approve?key1=" + key1 + "&key2=" + key2 + "&key3=" + key3
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		panic("Error")
	}

	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		panic("Error")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("ioutil.ReadAll err=%s", err.Error())
	}

	// struct for api response
	type jsonResponse struct {
		Result  string `json:"result"`
		Message string `json:"message"`
		StoreCd string `json:"store_cd"`
	}

	var jsonresponse jsonResponse
	json.Unmarshal(body, &jsonresponse)

	// redirect signin.tmpol if api returning message is "ng"
	if jsonresponse.Result == "ng" {
		m.App.ErrorLog.Println("Error occuered while compnay email verification!")
		m.App.Session.Put(r.Context(), "error", "Error occuered while compnay email verification!")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	stringMap := make(map[string]string)
	stringMap["store_cd"] = jsonresponse.StoreCd

	render.RenderTemplate(w, r, "store-registration-approve-new-store.page.tmpl", &models.TemplateData{
		//Form:      forms.New(nil),
		StringMap: stringMap,
	})
}

//  SignupStoreResetPassword is the handler for resseting password
func (m *Repository) SignupStoreResetPassword(w http.ResponseWriter, r *http.Request) {
	var emptyStoreResetPassword models.StoreResetPassword
	data := make(map[string]interface{})
	data["resetpassword"] = emptyStoreResetPassword

	// get result of resetting password from session
	reset := m.App.Session.GetString(r.Context(), "reset")
	if len(reset) == 0 {
		reset = ""
	} else {
		m.App.Session.Remove(r.Context(), "reset")
		// if resetting password was succeeded, delete jwt token not to login anymore.
		m.App.Session.Remove(r.Context(), "jwt")
	}
	stringMap := make(map[string]string)
	stringMap["reset"] = reset

	render.RenderTemplate(w, r, "store-registration-reset-password.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

//  SignupStoreResetPassword is the handler for resseting password
func (m *Repository) PostSignupStoreResetPassword(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	resetPassword := models.StoreResetPassword{
		CompanyCd:          r.Form.Get("company_code"),
		StoreCd:            r.Form.Get("store_code"),
		StoreKey:           r.Form.Get("store_key"),
		OldPassword:        r.Form.Get("old_password"),
		NewPassword:        r.Form.Get("new_password"),
		NewPasswordConfirm: r.Form.Get("new_password_confirm"),
	}

	form := forms.New(r.PostForm)
	// validation (move to api server inthe future)
	form.Required("company_code", "store_code", "store_key", "old_password", "new_password", "new_password_confirm")

	// reset password api call
	type Request struct {
		From               string `json:"from"`
		CompanyCd          string `json:"company_cd"`
		StoreCd            string `json:"store_cd"`
		StoreKey           string `json:"store_key"`
		OldPassword        string `json:"old_password"`
		NewPassword        string `json:"new_password"`
		NewPasswordConfirm string `json:"new_password_confirm"`
	}

	request := &Request{
		From:               "web",
		CompanyCd:          r.Form.Get("company_code"),
		StoreCd:            r.Form.Get("store_code"),
		StoreKey:           r.Form.Get("store_key"),
		OldPassword:        r.Form.Get("old_password"),
		NewPassword:        r.Form.Get("new_password"),
		NewPasswordConfirm: r.Form.Get("new_password_confirm"),
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	method := "POST"
	url := "https://anti-shoplifting-dev.cf/api/user/register/store/reset-password"
	req, err := http.NewRequest(method, url, bytes.NewBuffer(jsonRequest))
	if err != nil {
		// redirect to signin
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Http NewRequest Error")
		m.App.Session.Put(r.Context(), "error", "Http NewRequest Error")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		// redirect to signin
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot Reset Password")
		m.App.Session.Put(r.Context(), "error", "Can't Reset Password")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	fmt.Printf("%+v\n", resp)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		helpers.ServerError(w, err)
		m.App.ErrorLog.Println("Cannot retrieve value from response")
		m.App.Session.Put(r.Context(), "error", "Cannot retrieve value from response")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// struct for api response
	type jsonResponse struct {
		Result  string `json:"result"`
		Message string `json:"message"`
	}

	var jsonresponse jsonResponse
	json.Unmarshal(body, &jsonresponse)

	// redirect same screen doe to resetting password failed (NG)
	if jsonresponse.Result == "ng" {
		fmt.Printf("%-v", jsonresponse.Message)
		m.App.ErrorLog.Println("Cannot reset password due to response ng from API")
		m.App.Session.Put(r.Context(), "error", "Cannot reset password due to response ng from API")
		data := make(map[string]interface{})
		data["resetpassword"] = resetPassword
		render.RenderTemplate(w, r, "/signup/store/reset-password", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// redirect signin.tmpol if form parameter is invalid
	if !form.Valid() {
		data := make(map[string]interface{})
		data["resetpassword"] = resetPassword
		render.RenderTemplate(w, r, "/signup/store/reset-password", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// if resetting password waa succeeded, save result to session
	m.App.Session.Put(r.Context(), "reset", "succeeded")
	http.Redirect(w, r, "/signup/store/reset-password", http.StatusSeeOther)
}
