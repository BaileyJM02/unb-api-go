package main

import(
    "fmt"
    "./v1"
)

var (
    token string = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMzk4MTk3MTEzNDk1NzQ4NjI2IiwiaWF0IjoxNTUyODU1NTQ4fQ.woAdiyZ9T7rkl-owuIpl205k1RPouwzanX6N8eMCYKs"
)

func main() {
    api := v1.New(token)
    
    data, err := api.Check()
    if err != nil {
        fmt.Print(err)
    } else {
        fmt.Printf("Ping: %s \nIs up: %v \n", data.Ping, data.Up)
    }
    
    user, err := api.UserBalance("411898639737421824", "398197113495748626")
    if err != nil {
        fmt.Print(err)
    } else {
        fmt.Printf("Ping: %v \nIs up: %v \n", user.Total, user.Infinite)
    }
}