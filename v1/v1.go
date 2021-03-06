package v1

import(
	"io/ioutil"
	"net/http"
	"time"
	"errors"
	"fmt"
	"encoding/json"
	"strconv"
	"strings"
	"bytes"
)

type userData struct {
    token string
    client *http.Client
}

type errorResponse struct {
    Error string `json:"error"`
    Message string `json:"message"`
}

type timeoutResponse struct {
    Message string `json:"message"`
    RetryAfter time.Duration `json:"retry_after"`
}

type check struct {
    Ping time.Duration
    Up bool
}

type userObj struct {
    Rank int `json:"rank"`
    UserId string `json:"user_id"`
    Cash int `json:"cash"`
    CashInfinite bool `json:"infinite_cash"`
    CashNinfinite bool `json:"n-infinite_cash"`
    Bank int `json:"bank"`
    BankInfinite bool `json:"infinite_bank"`
    BankNinfinite bool `json:"n-infinite_bank"`
    Total int `json:"total"`
    Infinite bool `json:"infinite_total"`
    Ninfinite bool `json:"n-infinite_total"`
}

type userObjwReason struct {
    UserId string `json:"user_id"`
    Cash int `json:"cash"`
    CashInfinite bool `json:"infinite_cash"`
    CashNinfinite bool `json:"n-infinite_cash"`
    Bank int `json:"bank"`
    BankInfinite bool `json:"infinite_bank"`
    BankNinfinite bool `json:"n-infinite_bank"`
    Reason string `json:"reason"`
}

type userObjRaw struct {
    Rank interface{} `json:"rank"`
    UserId interface{} `json:"user_id"`
    Cash interface{} `json:"cash"`
    Bank interface{} `json:"bank"`
    Total interface{} `json:"total"`
}

type userObjPut struct {
    Cash interface{} `json:"cash,omitempty"`
    Bank interface{} `json:"bank,omitempty"`
    Reason interface{} `json:"reason,omitempty"`
}

func (u *userData) Request(protocol, url string, payload []byte) ([]byte, error) {
    b := bytes.NewBuffer(payload)
	req, err := http.NewRequest(protocol, "https://unbelievable.pizza/api/v1"+url, b)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", u.token)
	resp, err := u.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	respo, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode == 429 {
	    err := timeoutResponse{}
	    srsly := json.Unmarshal(respo, &err)
	    if srsly != nil {
	        // This is a srsly bad error -_-
	        panic(err)
	    }
	    return respo, errors.New(fmt.Sprintf("%v Retry after: %s", err.Message, err.RetryAfter))
	}
	// Bit hacky, test if the response contains the error body
	if strings.Contains(string(respo), "error") {
	    err := errorResponse{}
	    srsly := json.Unmarshal(respo, &err)
	    if srsly != nil {
	        // This is a srsly bad error -_-
	        panic(err)
	    }
	    return respo, errors.New(fmt.Sprintf("%v (%v)", err.Error, err.Message))
	}
	return respo, err
}

func fixTypesToStruct(data []byte) (userObj, error) {
    balUser := userObj{}
    var objmap map[string]interface{}
    err := json.Unmarshal(data, &objmap)
    if err != nil {
        return userObj{}, err
    }
    _, totalIsString := objmap["total"].(string)
    if totalIsString {
        switch x := objmap["total"]; x {
            case "Infinity":
                objmap["total"] = 0
                objmap["infinite_total"] = true
            case "-Infinity":
		        objmap["total"] = -0
                objmap["n-infinite_total"] = true
	        default:
	            objmap["total"], _ = strconv.ParseInt(objmap["total"].(string), 0, 64)
	    }
        
    }
    _, cashIsString := objmap["cash"].(string)
    if cashIsString {
        switch x := objmap["cash"]; x {
            case "Infinity":
                objmap["cash"] = 0
                objmap["infinite_cash"] = true
            case "-Infinity":
		        objmap["cash"] = -0
                objmap["n-infinite_cash"] = true
	        default:
		        objmap["cash"], _ = strconv.ParseInt(objmap["cash"].(string), 0, 64)
	    }
    }
    _, bankIsString := objmap["bank"].(string)
    if bankIsString {
        switch x := objmap["bank"]; x {
            case "Infinity":
                objmap["bank"] = 0
                objmap["infinite_bank"] = true
            case "-Infinity":
		        objmap["bank"] = -0
                objmap["n-infinite_bank"] = true
	        default:
		        objmap["bank"], _ = strconv.ParseInt(objmap["bank"].(string), 0, 64)
	    }
    }
    _, rankIsString := objmap["rank"].(string)
    if rankIsString {
        objmap["rank"], _ = strconv.ParseInt(objmap["rank"].(string), 0, 64)
    }
    
    b, err := json.Marshal(objmap)
    if err != nil {
        panic(err)
    }
    
    err = json.Unmarshal([]byte(b), &balUser)
    if err != nil {
        return userObj{}, err
    }
    
    return balUser, err
}

func New(token string) userData {
    client := &http.Client{}
    u := userData{token, client}
    return u
}

func Custom(token string, client *http.Client) userData {
    u := userData{token, client}
    return u
}


func (u *userData) Check() (check, error) {
    start := time.Now()
    data, err := u.Request("GET", "", nil)
    elapsed := time.Since(start)
    if err != nil {
        // because we never know how long.
        if strings.Contains(string(data), `{"message":"You are being rate limited.","retry_after"`) {
            return check{time.Since(time.Now()), true}, err
        }
        switch x := string(data); x {
            case `{"error":"404: Not found"}`:
                return check{elapsed, true}, nil
                
            case `{"error":"401: Unauthorized"}`:
		        return check{elapsed, true}, errors.New("401 Unauthorized (Check your token)")
		        
	        default:
		        return check{time.Since(time.Now()), false}, err
	    }
    }	
	return check{time.Since(time.Now()), false}, errors.New("Cannot Connect to API url.")
}

func (u *userData) GetBalance(guild, user string) (userObj, error) {
    data, err := u.Request("GET", fmt.Sprintf("/guilds/%v/users/%v", guild, user), nil)
    if err != nil {
        return userObj{}, err
    }
    userBal, err := fixTypesToStruct(data)
    if err != nil {
        return userObj{}, err
    }
	return userBal, err
}

func (u *userData) SetBalance(guild, user string, cash, bank, reason interface{}) (userObj, error) {
    var payloadTypes = make(map[string]interface{})
    switch x := cash; x.(type) {
        case string:
            if cash == "Infinity" {
                payloadTypes["Cash"] = "Infinity"
            } else if cash == "-Infinity" {
                payloadTypes["Cash"] = "-Infinity"
            }
        case int:
            payloadTypes["Cash"] = cash
    }
    switch x := bank; x.(type) {
        case string:
            if bank == "Infinity" {
                payloadTypes["Bank"] = "Infinity"
            } else if bank == "-Infinity" {
                payloadTypes["Bank"] = "-Infinity"
            }
        case int:
            payloadTypes["Bank"] = bank
    }
    switch x := reason; x.(type) {
        case string:
            payloadTypes["Reason"] = reason
        case nil:
            payloadTypes["Reason"] = "No reason provided."
    }
    value, err := json.Marshal(payloadTypes)
    if err != nil {
        return userObj{}, err
    }
    data, err := u.Request("PUT", fmt.Sprintf("/guilds/%v/users/%v", guild, user), value)
    if err != nil {
        return userObj{}, err
    }
    userBal, err := fixTypesToStruct(data)
    if err != nil {
        return userObj{}, err
    }
	return userBal, err
}

func (u *userData) UpdateBalance(guild, user string, cash, bank int, reason interface{}) (userObj, error) {
    var payloadTypes = make(map[string]interface{})
    payloadTypes["Cash"] = cash
    payloadTypes["Bank"] = bank
    switch x := reason; x.(type) {
        case string:
            payloadTypes["Reason"] = reason
        case nil:
            payloadTypes["Reason"] = "No reason provided."
    }
    value, err := json.Marshal(payloadTypes)
    if err != nil {
        return userObj{}, err
    }
    data, err := u.Request("PATCH", fmt.Sprintf("/guilds/%v/users/%v", guild, user), value)
    if err != nil {
        return userObj{}, err
    }
    userBal, err := fixTypesToStruct(data)
    if err != nil {
        return userObj{}, err
    }
	return userBal, err
}

func (u *userData) Leaderboard(guild string) ([]userObj, error) {
    var leaderboardRaw []userObjRaw
    var leaderboard []userObj
    
    data, err := u.Request("GET", fmt.Sprintf("/guilds/%v/users", guild), nil)
    if err != nil {
        return []userObj{}, err
    }
    
    if err := json.Unmarshal(data, &leaderboardRaw)
    err != nil {
        return []userObj{}, err
    }
    for _, v := range leaderboardRaw {
        value := fmt.Sprintf(`{"rank":"%v","user_id":"%v","cash":"%v","bank":"%v","total":"%v"}`,v.Rank,v.UserId,v.Cash,v.Bank,v.Total)
        user, err := fixTypesToStruct([]byte(value))
        if err != nil {
            return []userObj{}, err
        }
        leaderboard = append(leaderboard, user)
    }

    if err != nil {
        return []userObj{}, err
    }
	return leaderboard, err
}