package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

type RobloxUserApiResponse struct {
	IsBanned bool
}

type ItemDataResponse struct {
	Items map[int][]interface{}
}

type UsernameIDJsonResponse struct {
	ID int
	Username string
}

type InventoryJsonResponse struct {
	Data []struct {
		Name string
		InstanceID int
	}
}

var itemCount int
var userValue int
var userRAP int

func getUserName() (string, int) {
	var userID int
	var usernameJson *UsernameIDJsonResponse
	fmt.Println("Enter userID: ")
	_, _ = fmt.Scanln(&userID)
	resp, err := http.Get("https://api.roblox.com/users/" + fmt.Sprintf("%d",userID)) // NEEDS NEW API, WILL BE DEPRECATED SOON
	if err != nil {
		println(err)
	}
	err1 := json.NewDecoder(resp.Body).Decode(&usernameJson)
	if err != nil {
		fmt.Println(err1)
	}
	defer resp.Body.Close()
	username := usernameJson.Username
	return username, userID
}

func getUserID() (string, int) {
	var username string
	var usernameJson *UsernameIDJsonResponse
	fmt.Println("Enter username: ")
	_, _ = fmt.Scanln(&username)
	resp, err := http.Get("https://api.roblox.com/users/get-by-username?username=" + username) // NEEDS NEW API, WILL BE DEPRECATED SOON
	if err != nil {
		println(err)
	}
	err1 := json.NewDecoder(resp.Body).Decode(&usernameJson)
	if err1 != nil{
		println(err1)
	}
	defer resp.Body.Close()
	username = usernameJson.Username
	userID := usernameJson.ID
	return username, userID
}

func checkBan(userID int) bool {
	var responseJson *RobloxUserApiResponse
	resp9, err9 := http.Get("https://users.roblox.com/v1/users/" + fmt.Sprintf("%d",userID))
	if err9 != nil {
		fmt.Println(err9)
	}
	err10 := json.NewDecoder(resp9.Body).Decode(&responseJson)
	if err10 != nil {
		fmt.Print(err10)
	}
	defer resp9.Body.Close()
	return responseJson.IsBanned
}

func getItemData() ([]int, map[int][]interface{}) {
	var itemIds []int
	var data *ItemDataResponse
	resp, err := http.Get("https://www.rolimons.com/itemapi/itemdetails")
	if err != nil {
		fmt.Println("Error fetching Rolimon's item data:", err)
	}
	body, err3 := ioutil.ReadAll(resp.Body)
	if err3 != nil{
		fmt.Println(err3)
	}
	err4 := json.Unmarshal(body,&data)
	if err4 != nil{
		fmt.Println(err4)
	}
	defer resp.Body.Close()
	for id := range data.Items {
		itemIds = append(itemIds, id)
	}
	return itemIds, data.Items
}

func scannerHandler(username string, userID int) {
	var wg sync.WaitGroup
	itemCount = 0
	userValue = 0
	userRAP = 0
	limitedList, itemDataArray := getItemData()
	isBanned := checkBan(userID)
	wg.Add(len(limitedList))
	fmt.Printf("Currently scanning %s's (%d) inventory\n-------------------------------------------\n",username,userID)
	start := time.Now()
	if isBanned {
		for i := 0; i < len(limitedList); i++ {
			go termScan(limitedList[i],userID,&wg,itemDataArray)
		}
	} else {
		for i := 0; i < len(limitedList); i++ {
			go nonTermScan(limitedList[i],userID, &wg, itemDataArray)
		}
	}
	wg.Wait()
	fmt.Println(time.Since(start), "to scan", itemCount, "items.")
	fmt.Println(username, "has",userValue,"value and",userRAP,"RAP.")
}

func nonTermScan(limitedList int, userID int, wg *sync.WaitGroup, itemDataArray map[int][]interface{}) {
	var limitedJson *InventoryJsonResponse
	defer wg.Done()
	resp, err := http.Get("https://inventory.roblox.com/v1/users/" + fmt.Sprintf("%d",userID) + "/items/Asset/" + fmt.Sprintf("%d",limitedList))
	if err != nil {
		println(err)
	}
	err1 := json.NewDecoder(resp.Body).Decode(&limitedJson)
	if err1 != nil {
		println(err1)
	}
	defer resp.Body.Close()
	if len(limitedJson.Data) == 1 {
		fmt.Printf("%s (%d)\n", limitedJson.Data[0].Name, limitedJson.Data[0].InstanceID)
		itemCount++
		userValue += int(itemDataArray[limitedList][4].(float64))
		userRAP += int(itemDataArray[limitedList][2].(float64))
	} else if len(limitedJson.Data) > 1 {
		fmt.Printf("%dx %s ", len(limitedJson.Data), limitedJson.Data[0].Name)
		fmt.Print("(")
		for i, v := range limitedJson.Data{
			fmt.Printf("%d", v.InstanceID)
			if i != len(limitedJson.Data) - 1 {
				fmt.Print(",")
			}
		}
		fmt.Print(")\n")
		itemCount += len(limitedJson.Data)
		userValue += int(itemDataArray[limitedList][4].(float64)) * len(limitedJson.Data)
		userRAP += int(itemDataArray[limitedList][2].(float64)) * len(limitedJson.Data)
	}
}

func termScan(limitedList int, userID int, wg *sync.WaitGroup, itemDataArray map[int][]interface{}) {
	fmt.Println("To be completed.")
	defer wg.Done()
	return
	/*var limitedJson *InventoryJsonResponse
	defer wg.Done()
	resp, err := http.Get("https://api.roblox.com/ownership/hasasset?userId=" + fmt.Sprintf("%d",userID) + "&assetId=" + fmt.Sprintf("%d",limitedList))
	if err != nil {
		println(err)
	}
	err1 := json.NewDecoder(resp.Body).Decode(&limitedJson)
	if err1 != nil {
		fmt.Println(err1)
	}
	defer resp.Body.Close()
	if resp.Body
	fmt.Printf("%s (%d)\n", limitedJson.Data[0]) // pointer
	fmt.Println("")
	itemCount++
	userValue += int(itemDataArray[limitedList][4].(float64)) * len(limitedJson.Data)
	userRAP += int(itemDataArray[limitedList][2].(float64)) * len(limitedJson.Data)*/
}

func main() {
	var searchChoice int
	var killSwitch = true
	for killSwitch {
		fmt.Println("--[Offline#0909's Private Inventory Scanner]--\n1) Scan by username\n2) Scan by userID\n3) Quit")
		_, _ = fmt.Scanln(&searchChoice)
		switch searchChoice {
		case 1:
			username, userID := getUserID()
			if username == "" {
				fmt.Println("Invalid username.")
			} else {
				scannerHandler(username, userID)
			}
			break
		case 2:
			username, userID := getUserName()
			if username == "" {
				fmt.Println("Invalid userID.")
			} else {
				scannerHandler(username, userID)
			}
			break
		case 3:
			killSwitch = false
		default:
			fmt.Println("Invalid selection.")
		}
	}
}