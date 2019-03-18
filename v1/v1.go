package v1

import(
	"io/ioutil"
	"net/http"
	"time"
	"errors"
	"fmt"
	"encoding/json"
)

type userData struct {
    token string
    client *http.Client
}

type check struct {
    Ping time.Duration
    Up bool
}

type userBalance struct {
    Rank string `json:"rank"`
    UserId string `json:"user_id"`
    Cash int `json:"cash"`
    Bank int `json:"bank"`
    Total int `json:"total"`
    Infinite bool `json:"infinite"`
    Ninfinite bool `json:"Ninfinte"`
}

type userBalanceRaw struct {
    Rank string `json:"rank"`
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
    if err != nil {
        return check{time.Since(time.Now()), false}, err
    }
    elapsed := time.Since(start)
    if string(data) == `{"error":"404: Not found"}` {
	    return check{elapsed, true}, nil
	}
	if string(data) == `{"error":"401: Unauthorized"}` {
	    return check{elapsed, true}, errors.New("Unauthorized - Check your token")
	}
		
	return check{time.Since(time.Now()), false}, errors.New("Cannot Connect to API url.")
}

func (u *userData) UserBalance(guild, user string) (userBalance, error) {
    balUser := userBalance{}
    data, err := u.Request("GET", fmt.Sprintf("/guilds/%v/users/%v", guild, user))
    if err != nil {
        return userBalance{"","",0,0,0,false,false}, err
    }
    var objmap map[string]interface{}
    err = json.Unmarshal(data, &objmap)
    if err != nil {
        return userBalance{"","",0,0,0,false,false}, err
    }
    _, ok := objmap["total"].(string)
    if ok {
        objmap["total"] = 0
        objmap["Infinite"] = true
    }
    
    b, err := json.Marshal(objmap)
    if err != nil {
        panic(err)
    }
    err = json.Unmarshal([]byte(b), &balUser)
    if err != nil {
        return userBalance{"","",0,0,0,false,false}, err
    }    
    if string(data) == `{"error":"404: Not found"}` {
	    return userBalance{"","",0,0,0,false,false}, nil
	}
	if string(data) == `{"error":"401: Unauthorized"}` {
	    return userBalance{"","",0,0,0,false,false}, err
	}
		
	return userBalance{"","",0,0,0,false,false}, err
}