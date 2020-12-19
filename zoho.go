package zoho

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

//--------------------------------------------------------
// structure holds Constants used to send as parameters for zoho API calls
//--------------------------------------------------------

type Params struct {
	OauthAuthorizationRequestSlug string
	OauthGenerateTokenRequestSlug string
	OauthRevokeTokenRequestSlug   string
	AthURL                        string

	Scope            string
	AccessType       string
	ResponseTypeCode string
	ClientID         string
	ClientSecret     string
	RedirectURI      string
	GrantType        string
	RefreshToken     string
	RgrantType       string
	Org              string
}

//-----------------------------------------------------
// struct to capture the OAuth token from zoho rest API
//-----------------------------------------------------
type OAuthToken struct {
	Accesstoken string `json:"access_token"`
	Tokentype   string `json:"token_type"`
	Expiresin   int    `json:"expires_in"`
}

//--------------------------------------------------------------
// struct to capture the Organizaiton details from zoho rest API
//--------------------------------------------------------------
type Organization struct {
	Data []struct {
		Id int `json:"id"`
	}
}

//--------------------------------------------------------------
// struct to capture the Department details from zoho rest API
//--------------------------------------------------------------
type Departments struct {
	Data []struct {
		Id          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
		IsEnabled   bool   `json:"isEnabled"`
	}
}

//--------------------------------------------------------------
// struct to capture the Agent details from zoho rest API
//--------------------------------------------------------------
type Agents struct {
	Data []struct {
		Id                      string   `json:"id"`
		LastName                string   `json:"lastName"`
		FirstName               string   `json:"firstName"`
		RoleId                  string   `json:"roleId"`
		EmailId                 string   `json:"emailId"`
		Name                    string   `json:"name"`
		Mobile                  string   `json:"mobile"`
		Status                  string   `json:"status"`
		STcode                  string   `json:"aboutInfo"`
		AssociatedDepartmentIds []string `json:"associatedDepartmentIds"`
	}
}

//--------------------------------------------------------------
// struct to capture the update ticket respones from zoho rest API
//--------------------------------------------------------------
type UpdateTicketResponse struct {
	TicketNumber string `json:"ticketNumber"`
	ModifiedTime string `json:"modifiedTime"`
	StatusType   string `json:"statusType"`
	DepartmentId string `json:"departmentId"`
	IsDeleted    bool   `json:"isDeleted"`
	Id           string `json:"id"`
	AssigneeId   string `json:"assigneeId"`
	Status       string `json:"status"`
	ErrorMsg     string `json:"message"`
}

type UpdateTicketDetails struct {
	TicketId     string
	DepartmentId string
	AssigneeId   string
}

//--------------------------------------------------------------
// struct to capture ticket details from zoho rest API
//--------------------------------------------------------------
type TicketDetailsResponse struct {
	Id           string `json:"id"`
	TicketNumber string `json:"ticketNumber"`
	ModifiedTime string `json:"modifiedTime"`
	ClosedTime   string `json:"closedTime"`
	DepartmentId string `json:"departmentId"`
	AssigneeId   string `json:"assigneeId"`
	IsDeleted    bool   `json:"isDeleted"`
	Status       string `json:"status"`
	ErrorMsg     string `json:"message"`
	Assignee     struct {
		FirstName string `json:"firstName"`
		LastName  string `json:"lastName"`
	}
	Cf struct {
		Cf_reason       string `json:"cf_reason"`
		Cf_reminderdate string `json:"cf_reminderdate"`
	}
}

//-------------------------------------------------
//function return Auth Token based on a refresh token
//-------------------------------------------------
func GetOAuthToken(params Params) (string, error) {

	//urla := "https://accounts.zoho.com/oauth/v2/token"
	urla := params.AthURL + params.OauthGenerateTokenRequestSlug

	parm := url.Values{}
	parm.Set("refresh_token", params.RefreshToken)
	parm.Add("client_id", params.ClientID)
	parm.Add("client_secret", params.ClientSecret)
	parm.Add("scope", params.Scope)
	parm.Add("redirect_uri", params.RedirectURI)
	parm.Add("grant_type", params.RgrantType)
	reqs, err := http.NewRequest("POST", urla, strings.NewReader(parm.Encode()))
	reqs.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	response, err := client.Do(reqs)
	//var jsondata string
	if err != nil {
		return "", fmt.Errorf("The HTTP request failed with error :", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		v_OAuthToken := OAuthToken{}
		err = json.Unmarshal((data), &v_OAuthToken)
		if err != nil {
			return "", fmt.Errorf("The Json request failed with error : ", err)
		} else {
			return v_OAuthToken.Accesstoken, nil
		}
	}

}

//-------------------------------------------------
//function return the organization value from Zoho
//-------------------------------------------------
func GetOrganization(authtoken string) (string, error) {

	v_Organization := Organization{}
	urla := "https://desk.zoho.com/api/v1/organizations"

	reqs, err := http.NewRequest("GET", urla, nil)
	reqs.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	reqs.Header.Add("Authorization", "Zoho-oauthtoken "+authtoken)
	client := &http.Client{}
	response, err := client.Do(reqs)

	if err != nil {
		return "", fmt.Errorf("The HTTP request failed with error :", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		err = json.Unmarshal([]byte(data), &v_Organization)
		if err != nil {
			return "", fmt.Errorf("The Json request failed with error : ", err)
		} else if len(v_Organization.Data) > 0 {
			return strconv.Itoa(v_Organization.Data[0].Id), nil
		} else {
			return "", fmt.Errorf("Unknown error.. please investigate")
		}

	}
}

//-------------------------------------------------
//function return the departments from Zoho
//-------------------------------------------------
func GetDepartments(authtoken string) (Departments, error) {

	v_Departments := Departments{}
	urla := "https://desk.zoho.com/api/v1/departments?limit=100"

	reqs, err := http.NewRequest("GET", urla, nil)
	reqs.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	reqs.Header.Add("Authorization", "Zoho-oauthtoken "+authtoken)
	client := &http.Client{}
	response, err := client.Do(reqs)

	if err != nil {
		return v_Departments, fmt.Errorf("The HTTP request failed with error :", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		if string(data) != "" {
			err = json.Unmarshal([]byte(data), &v_Departments)
			if err != nil {
				return v_Departments, fmt.Errorf("The Json request failed with error : ", err)
			}
		}

		return v_Departments, nil

	}
}

//-------------------------------------------------
//function return the departments from Zoho
//-------------------------------------------------
func GetAgents(authtoken string, from string, limit string, emailID string) (Agents, error) {

	v_Agents := Agents{}
	urla := "https://desk.zoho.com/api/v1/agents?"
	if from != "" {
		urla = urla + "&from=" + from
	}
	if limit != "" {
		urla = urla + "&limit=" + limit
	}
	if emailID != "" {
		urla = urla + "&fieldName=emailId&searchStr=" + emailID
	}
	//log.Println(urla)
	reqs, err := http.NewRequest("GET", urla, nil)
	reqs.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	reqs.Header.Add("Authorization", "Zoho-oauthtoken "+authtoken)
	//reqs.Header.Add("orgId", getOrg())
	client := &http.Client{}
	response, err := client.Do(reqs)

	if err != nil {
		return v_Agents, fmt.Errorf("The HTTP request failed with error :", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		//log.Println(string(data))
		if string(data) != "" {
			err = json.Unmarshal([]byte(data), &v_Agents)
			if err != nil {
				return v_Agents, fmt.Errorf("The Json request failed with error : ", err)
			}
		}
		return v_Agents, nil

	}
}

//-------------------------------------------------
//funciton return the org id already saved locally
//-------------------------------------------------
func GetOrg(params Params) string {
	return params.Org
}

//-------------------------------------------------
//function update ticket
//-------------------------------------------------
func UpdateTicketStatus(authtoken string, orgID string, ticketID string, status string) (UpdateTicketResponse, error) {
	v_Response := UpdateTicketResponse{}

	urla := "https://desk.zoho.com/api/v1/tickets/" + ticketID

	postBody, _ := json.Marshal(map[string]string{
		"status": status,
	})
	postJsonBody := bytes.NewBuffer(postBody)

	reqs, err := http.NewRequest("PATCH", urla, postJsonBody)
	reqs.Header.Add("Content-Type", "application/json")
	reqs.Header.Add("orgId", orgID)
	reqs.Header.Add("Authorization", "Zoho-oauthtoken "+authtoken)

	client := &http.Client{}
	response, err := client.Do(reqs)

	if err != nil {
		return v_Response, fmt.Errorf("The HTTP request failed with error :", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		//log.Println(string(data))
		if string(data) != "" {
			err = json.Unmarshal([]byte(data), &v_Response)
			if err != nil {
				return v_Response, fmt.Errorf("The Json request failed with error : ", err)
			}
		}
		return v_Response, nil

	}

}

//-------------------------------------------------
//function update ticket
//-------------------------------------------------
func UpdateTicketAgent(authtoken string, orgID string, ticket UpdateTicketDetails) (UpdateTicketResponse, error) {
	v_Response := UpdateTicketResponse{}

	urla := "https://desk.zoho.com/api/v1/tickets/" + ticket.TicketId

	postBody, _ := json.Marshal(map[string]string{
		"assigneeId":   ticket.AssigneeId,
		"departmentId": ticket.DepartmentId,
		"status":       "Open",
	})
	postJsonBody := bytes.NewBuffer(postBody)

	reqs, err := http.NewRequest("PATCH", urla, postJsonBody)
	reqs.Header.Add("Content-Type", "application/json")
	reqs.Header.Add("orgId", orgID)
	reqs.Header.Add("Authorization", "Zoho-oauthtoken "+authtoken)

	client := &http.Client{}
	response, err := client.Do(reqs)

	if err != nil {
		return v_Response, fmt.Errorf("The HTTP request failed with error :", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		//log.Println(string(data))
		if string(data) != "" {
			err = json.Unmarshal([]byte(data), &v_Response)
			if err != nil {
				return v_Response, fmt.Errorf("The Json request failed with error : ", err)
			}
		}
		return v_Response, nil

	}

}

//-------------------------------------------------
//function update ticket
//-------------------------------------------------
func UpdateTicketMove(authtoken string, orgID string, ticket UpdateTicketDetails) (UpdateTicketResponse, error) {
	v_Response := UpdateTicketResponse{}

	urla := "https://desk.zoho.com/api/v1/tickets/" + ticket.TicketId + "/move"

	postBody, _ := json.Marshal(map[string]string{
		"departmentId": ticket.DepartmentId,
	})
	postJsonBody := bytes.NewBuffer(postBody)

	reqs, err := http.NewRequest("POST", urla, postJsonBody)
	reqs.Header.Add("Content-Type", "application/json")
	reqs.Header.Add("orgId", orgID)
	reqs.Header.Add("Authorization", "Zoho-oauthtoken "+authtoken)

	client := &http.Client{}
	response, err := client.Do(reqs)
	log.Println(response)

	if err != nil {
		return v_Response, fmt.Errorf("The HTTP request failed with error :", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		//log.Println(string(data))
		if string(data) != "" {
			err = json.Unmarshal([]byte(data), &v_Response)
			if err != nil {
				return v_Response, fmt.Errorf("The Json request failed with error : ", err)
			}
		}
		return v_Response, nil

	}

}

//-------------------------------------------------
//function gets ticket details from zoho
//-------------------------------------------------
func GetTicketDetails(authtoken string, orgID string, ticketid string) (TicketDetailsResponse, error) {
	v_Response := TicketDetailsResponse{}

	urla := "https://desk.zoho.com/api/v1/tickets/" + ticketid + "?include=assignee,departments"

	reqs, err := http.NewRequest("GET", urla, nil)
	reqs.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	reqs.Header.Add("orgId", orgID)
	reqs.Header.Add("Authorization", "Zoho-oauthtoken "+authtoken)

	client := &http.Client{}
	response, err := client.Do(reqs)

	if err != nil {
		return v_Response, fmt.Errorf("The HTTP request failed with error :", err)
	} else {
		data, _ := ioutil.ReadAll(response.Body)
		//log.Println(string(data))
		if string(data) != "" {
			err = json.Unmarshal([]byte(data), &v_Response)
			if err != nil {
				return v_Response, fmt.Errorf("The Json request failed with error : ", err)
			}
		}
		return v_Response, nil

	}

}
