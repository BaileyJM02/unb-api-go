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
)

type userData struct {
    token string
    client *http.Client
}

type errorResponse struct {
    Error string `json:"error"`
    Message string `json:"message"`
}

type check struct {
    Ping time.Duration
    Up bool
}

type userBalance struct {
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

type userBalanceRaw struct {
    Rank int `json:"rank"`
    UserId string `json:"user_id"`
    Cash int `json:"cash"`
    Bank int `json:"bank"`
    Total interface{} `json:"total"`
}

func (u *userData) Request(protocol, url string) ([]byte, error) {
	req, err := http.NewRequest(protocol, "https://unbelievable.pizza/api/v1"+url, nil)
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
    data, err := u.Request("GET", "")
    elapsed := time.Since(start)
    if err != nil {
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

func (u *userData) UserBalance(guild, user string) (userBalance, error) {
    balUser := userBalance{}
    data, err := u.Request("GET", fmt.Sprintf("/guilds/%v/users/%v", guild, user))
    if err != nil {
        return userBalance{}, err
    }
    var objmap map[string]interface{}
    err = json.Unmarshal(data, &objmap)
    if err != nil {
        return userBalance{}, err
    }
    _, totalIsString := objmap["total"].(string)
    if totalIsString {
        switch x := objmap["total"]; x {
            case "Infinity":
                fmt.Print(x)
                objmap["total"] = 0
                objmap["infinite_total"] = true
            case "-Infinity":
                fmt.Print(x)
		        objmap["total"] = -0
                objmap["n-infinite_total"] = true
	        default:
	            fmt.Print(x)
		        objmap["total"] = 0
                objmap["infinite_total"] = true
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
		        objmap["cash"] = 0
                objmap["infinite_cash"] = true
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
		        objmap["bank"] = 0
                objmap["infinite_bank"] = true
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
        return userBalance{}, err
    }    
		
	return balUser, err
}