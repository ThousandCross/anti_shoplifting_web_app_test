package models

import "time"

type Blacklists struct {
	RowNumber            uint      `json:"row_number"`
	Id                   uint      `json:"incident_id"`
	UserId               uint      `json:"user_id"`
	GlobalId             uint      `json:"global_id"`
	Name                 string    `json:"name"`
	Img                  string    `json:"img"`
	LatestVisitDatetime  time.Time `json:"latest_visit_datetime"`
	NumberOfVisits       uint      `json:"number_of_visits"`
	RegistrationDatetime time.Time `json:"registration_datetime"`
}

type Incidents struct {
	RowNumber      uint            `json:"row_number"`
	Id             uint            `json:"incident_id"`
	IncidentTypeId uint            `json:"incident_type_id"`
	UserId         uint            `json:"user_id"`
	GlobalId       uint            `json:"global_id"`
	StartDatetime  time.Time       `json:"start_datetime"`
	EndDatetime    time.Time       `json:"end_datetime"`
	Name           string          `json:"incident_type_name"`
	IncidentVideo  []IncidentVideo `json:"videos"`
}

type IncidentVideo struct {
	CameraId      uint      `json:"camera_id"`
	SerialNo      string    `json:"serial_no"`
	StartDatetime time.Time `json:"start_datetime"`
	EndDatetime   time.Time `json:"end_datetime"`
	Url           string    `json:"url"`
}

type IncidentSalesItem struct {
	Id          uint   `json:"id"`
	IncidentId  uint   `json:"incident_id"`
	SalesItemId uint   `json:"sales_item_id"`
	Count       uint   `json:"count"`
	Url         string `json:"url"`
}

type Prefectures struct {
	Id   uint   `json:"id"`
	Name string `json:"name"`
}

// Signin holds signin data
type Signin struct {
	CompanyCd  string
	StoreCd    string
	Password   string
	RememberMe bool
}

// CompanySignup holds company's all info for signup
type CompanySignupConfirm struct {
	From                         string `json:"from"`
	CompanyName                  string `json:"company_name"`
	RepresentativeFamilyName     string `json:"representative_family_name"`
	RepresentativeFirstName      string `json:"representative_first_name"`
	RepresentativeFamilyNameKana string `json:"representative_family_name_kana"`
	RepresentativeFirstNameKana  string `json:"representative_first_name_kana"`
	Zipcode                      string `json:"zipcode"`
	Prefecture                   string `json:"prefecture"`
	PrefectureId                 string `json:"prefecture_id"` // have to cast prefecture id from int to string
	City                         string `json:"city"`
	Street                       string `json:"street"`
	Building                     string `json:"building"`
	Tel                          string `json:"tel"`
	Mail                         string `json:"mail"`
	ManagerFamilyName            string `json:"manager_family_name"`
	ManagerFirstName             string `json:"manager_first_name"`
	ManagerFamilyNameKana        string `json:"manager_family_name_kana"`
	ManagerFirstNameKana         string `json:"manager_first_name_kana"`
	ManagerTel                   string `json:"manager_tel"`
	ManagerMail                  string `json:"manager_mail"`
	CardNo                       string `json:"card_no"`
	CardHolderFamilyNameKana     string `json:"card_holder_family_name_kana"`
	CardHolderFirstNameKana      string `json:"card_holder_first_name_kana"`
	CardMonth                    string `json:"card_month"`
	CardYear                     string `json:"card_year"`
	SecurityCd                   string `json:"security_cd"`
}

// CompanyBasicInfo holds company basic info data
type CompanyBasicInfo struct {
	CompanyName                  string
	RepresentativeFamilyName     string
	RepresentativeFirstName      string
	RepresentativeFamilyNameKana string
	RepresentativeFirstNameKana  string
	Zipcode                      string
	Prefecture                   string
	PrefectureId                 uint
	City                         string
	Street                       string
	Building                     string
	Tel                          string
	Mail                         string
}

// CompanyBasicInfoContd holds company basic info continued data
type CompanyBasicInfoContd struct {
	ManagerFamilyName     string
	ManagerFirstName      string
	ManagerFamilyNameKana string
	ManagerFirstNameKana  string
	ManagerTel            string
	ManagerMail           string
}

// CompanyPayment holds company payment data
type CompanyPayment struct {
	CardNo                   string
	CardHolderFamilyNameKana string
	CardHolderFirstNameKana  string
	CardMonth                string
	CardYear                 string
	SecurityCd               string
}

// StoreSignup holds store's all info for signup
type StoreSignupConfirm struct {
	From                  string `json:"from"`
	CompanyKey            string `json:"company_key"`
	CompanyCd             string `json:"company_cd"`
	StoreName             string `json:"store_name"`
	Zipcode               string `json:"zipcode"`
	Prefecture            string `json:"prefecture"`
	PrefectureId          string `json:"prefecture_id"` // have to cast prefecture id from int to string
	City                  string `json:"city"`
	Street                string `json:"street"`
	Building              string `json:"building"`
	ManagerFamilyName     string `json:"manager_family_name"`
	ManagerFirstName      string `json:"manager_first_name"`
	ManagerFamilyNameKana string `json:"manager_family_name_kana"`
	ManagerFirstNameKana  string `json:"manager_first_name_kana"`
	ManagerTel            string `json:"manager_tel"`
	ManagerMail           string `json:"manager_mail"`
	Password              string `json:"password"`
	PasswordConfirm       string `json:"password_confirm"`
}

// StoreBasicInfo holds store basic info data
type StoreBasicInfo struct {
	CompanyKey   string
	CompanyCd    string
	StoreName    string
	Zipcode      string
	Prefecture   string
	PrefectureId uint
	City         string
	Street       string
	Building     string
}

// StoreBasicInfoContd holds store basic info contd data
type StoreBasicInfoContd struct {
	ManagerFamilyName     string
	ManagerFirstName      string
	ManagerFamilyNameKana string
	ManagerFirstNameKana  string
	ManagerTel            string
	ManagerMail           string
}

// StorePassword holds store basic info contd data
type StorePassword struct {
	Password        string
	PasswordConfirm string
}

type StoreResetPassword struct {
	CompanyCd          string
	StoreCd            string
	StoreKey           string
	OldPassword        string
	NewPassword        string
	NewPasswordConfirm string
}
