package httpx

type User struct {
	Id              int      `json:"id"`
	Role            int      `json:"role"`
	Username        string   `json:"username"`
	Email           string   `json:"email"`
	Portrait        string   `json:"portrait"`
	Language        string   `json:"language"`
	Timezone        string   `json:"timezone"`
	WebTitle        string   `json:"webTitle"`
	ClientType      string   `json:"clientType"`
	Item            string   `json:"item"`
	Version         string   `json:"version"`
	CreatedAt       int      `json:"createdAt"`
	IsEmailConfirm  int      `json:"isEmailConfirm"`
	IsMechanismSub  bool     `json:"isMechanismSub"`  //是否机构子账号
	MechanismId     int      `json:"mechanismId"`     //机构ID 用于标识机构登录机构账号操作
	MechanismDomain []string `json:"mechanismDomain"` //所属机构domain
	Organization    int      `json:"organization"`    //所属机构ID
}
