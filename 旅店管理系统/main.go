package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// User 定义了用户结构体，Role 字段为 "admin" 或 "customer"。
// 对于顾客，CustomerType 表示会员或普通账号，Balance 表示账户余额。
type User struct {
	ID           int     `json:"id"`
	Username     string  `json:"username"`
	Password     string  `json:"password"`
	Role         string  `json:"role"`          // "admin" 或 "customer"
	CustomerType string  `json:"customer_type"` // "member" 或 "regular"，仅当 Role 为 "customer" 时有效
	Balance      float64 `json:"balance"`       // 仅当 Role 为 "customer" 时有效
}

// Room 定义了酒店房间结构体，Available 表示当前剩余可预订数量。
type Room struct {
	ID        int     `json:"id"`
	Type      string  `json:"type"`      // 房间类型，如单人间、双人间等
	Price     float64 `json:"price"`     // 房间价格
	Total     int     `json:"total"`     // 房间总数量
	Available int     `json:"available"` // 当前剩余数量
}

var users []User
var rooms []Room

const usersFile = "users.json"
const roomsFile = "rooms.json"

var reader = bufio.NewReader(os.Stdin)

func main() {
	// 加载用户和房间数据
	loadUsers()
	loadRooms()

	for {
		fmt.Println("================================")
		fmt.Println("欢迎使用酒店管理系统")
		fmt.Println("1. 登录")
		fmt.Println("2. 注册（仅限顾客）")
		fmt.Println("3. 退出")
		fmt.Print("请选择操作：")
		choice := readLine()
		switch choice {
		case "1":
			user := login()
			if user != nil {
				if user.Role == "admin" {
					adminMenu(user)
				} else if user.Role == "customer" {
					customerMenu(user)
				}
			}
		case "2":
			registerCustomer()
		case "3":
			fmt.Println("退出系统")
			return
		default:
			fmt.Println("无效的选项，请重试。")
		}
	}
}

// readLine 从标准输入读取一行数据并去掉末尾换行符
func readLine() string {
	input, _ := reader.ReadString('\n')
	return strings.TrimSpace(input)
}

// ------------------------- 数据持久化相关 ----------------------------

// 加载用户数据，如果文件不存在则初始化默认管理员账号
func loadUsers() {
	data, err := ioutil.ReadFile(usersFile)
	if err != nil {
		// 文件不存在，初始化默认管理员账号
		fmt.Println("未找到用户数据文件，初始化默认管理员账号。")
		users = []User{
			{
				ID:       1,
				Username: "admin",
				Password: "admin",
				Role:     "admin",
			},
		}
		saveUsers()
		return
	}
	err = json.Unmarshal(data, &users)
	if err != nil {
		fmt.Println("加载用户数据错误：", err)
		os.Exit(1)
	}
}

// 保存用户数据到文件
func saveUsers() {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		fmt.Println("保存用户数据错误：", err)
		return
	}
	err = ioutil.WriteFile(usersFile, data, 0644)
	if err != nil {
		fmt.Println("写入用户数据文件错误：", err)
	}
}

// 加载房间数据，如果文件不存在则初始化为空房间列表
func loadRooms() {
	data, err := ioutil.ReadFile(roomsFile)
	if err != nil {
		fmt.Println("未找到房间数据文件，初始化空房间列表。")
		rooms = []Room{}
		saveRooms()
		return
	}
	err = json.Unmarshal(data, &rooms)
	if err != nil {
		fmt.Println("加载房间数据错误：", err)
		os.Exit(1)
	}
}

// 保存房间数据到文件
func saveRooms() {
	data, err := json.MarshalIndent(rooms, "", "  ")
	if err != nil {
		fmt.Println("保存房间数据错误：", err)
		return
	}
	err = ioutil.WriteFile(roomsFile, data, 0644)
	if err != nil {
		fmt.Println("写入房间数据文件错误：", err)
	}
}

// ------------------------- 登录与注册 ----------------------------

// login 实现用户登录，输入用户名和密码后返回对应的用户指针（成功则返回，不成功返回 nil）
func login() *User {
	fmt.Print("请输入用户名：")
	username := readLine()
	fmt.Print("请输入密码：")
	password := readLine()

	for i := range users {
		if users[i].Username == username && users[i].Password == password {
			fmt.Println("登录成功！")
			return &users[i]
		}
	}
	fmt.Println("用户名或密码错误！")
	return nil
}

// registerCustomer 仅允许注册顾客账号（会员或普通），默认初始余额 1000 元
func registerCustomer() {
	fmt.Println("注册新顾客账号")
	fmt.Print("请输入用户名：")
	username := readLine()
	// 检查用户名是否已存在
	for _, user := range users {
		if user.Username == username {
			fmt.Println("用户名已存在！")
			return
		}
	}
	fmt.Print("请输入密码：")
	password := readLine()
	fmt.Print("请选择顾客类型（1. 会员账号 2. 普通账号）：")
	choice := readLine()
	var customerType string
	if choice == "1" {
		customerType = "member"
	} else {
		customerType = "regular"
	}
	newUser := User{
		ID:           getNextUserID(),
		Username:     username,
		Password:     password,
		Role:         "customer",
		CustomerType: customerType,
		Balance:      1000.0,
	}
	users = append(users, newUser)
	saveUsers()
	fmt.Println("注册成功！初始余额为 1000 元。")
}

// getNextUserID 获取下一个用户 ID（自动递增）
func getNextUserID() int {
	maxID := 0
	for _, user := range users {
		if user.ID > maxID {
			maxID = user.ID
		}
	}
	return maxID + 1
}

// getNextRoomID 获取下一个房间 ID（自动递增）
func getNextRoomID() int {
	maxID := 0
	for _, room := range rooms {
		if room.ID > maxID {
			maxID = room.ID
		}
	}
	return maxID + 1
}

// ------------------------- 管理员功能 ----------------------------

// adminMenu 为管理员提供用户管理和房间管理的菜单
func adminMenu(user *User) {
	for {
		fmt.Println("================================")
		fmt.Println("管理员菜单")
		fmt.Println("1. 用户管理")
		fmt.Println("2. 房间管理")
		fmt.Println("3. 退出")
		fmt.Print("请选择操作：")
		choice := readLine()
		switch choice {
		case "1":
			adminUserManagement()
		case "2":
			adminRoomManagement()
		case "3":
			fmt.Println("注销成功")
			return
		default:
			fmt.Println("无效的选项，请重试。")
		}
	}
}

// adminUserManagement 实现管理员对用户的增删改查操作
func adminUserManagement() {
	for {
		fmt.Println("--------- 用户管理 ---------")
		fmt.Println("1. 查看所有用户")
		fmt.Println("2. 添加用户")
		fmt.Println("3. 修改用户")
		fmt.Println("4. 删除用户")
		fmt.Println("5. 返回上一层")
		fmt.Print("请选择操作：")
		choice := readLine()
		switch choice {
		case "1":
			listUsers()
		case "2":
			addUser()
		case "3":
			updateUser()
		case "4":
			deleteUser()
		case "5":
			return
		default:
			fmt.Println("无效的选项，请重试。")
		}
	}
}

// listUsers 显示所有用户信息
func listUsers() {
	fmt.Println("----- 所有用户列表 -----")
	for _, user := range users {
		fmt.Printf("ID: %d, 用户名: %s, 角色: %s", user.ID, user.Username, user.Role)
		if user.Role == "customer" {
			fmt.Printf(", 类型: %s, 余额: %.2f", user.CustomerType, user.Balance)
		}
		fmt.Println()
	}
}

// addUser 由管理员添加新用户，可以添加管理员或顾客账号
func addUser() {
	fmt.Println("----- 添加新用户 -----")
	fmt.Print("请输入用户名：")
	username := readLine()
	// 检查用户名是否已存在
	for _, user := range users {
		if user.Username == username {
			fmt.Println("用户名已存在！")
			return
		}
	}
	fmt.Print("请输入密码：")
	password := readLine()
	fmt.Print("请选择角色（1. 管理员 2. 顾客）：")
	roleChoice := readLine()
	var role string
	var customerType string
	var balance float64
	if roleChoice == "1" {
		role = "admin"
	} else if roleChoice == "2" {
		role = "customer"
		fmt.Print("请选择顾客类型（1. 会员账号 2. 普通账号）：")
		ctChoice := readLine()
		if ctChoice == "1" {
			customerType = "member"
		} else {
			customerType = "regular"
		}
		balance = 1000.0 // 初始余额
	} else {
		fmt.Println("无效的角色选项")
		return
	}
	newUser := User{
		ID:           getNextUserID(),
		Username:     username,
		Password:     password,
		Role:         role,
		CustomerType: customerType,
		Balance:      balance,
	}
	users = append(users, newUser)
	saveUsers()
	fmt.Println("用户添加成功！")
}

// updateUser 修改指定用户的信息
func updateUser() {
	fmt.Print("请输入要修改的用户ID：")
	idStr := readLine()
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("无效的ID")
		return
	}
	var user *User
	for i := range users {
		if users[i].ID == id {
			user = &users[i]
			break
		}
	}
	if user == nil {
		fmt.Println("未找到该用户")
		return
	}
	fmt.Printf("当前用户名: %s\n", user.Username)
	fmt.Print("请输入新的用户名（直接回车保持不变）：")
	newUsername := readLine()
	if newUsername != "" {
		user.Username = newUsername
	}
	fmt.Print("请输入新的密码（直接回车保持不变）：")
	newPassword := readLine()
	if newPassword != "" {
		user.Password = newPassword
	}
	// 如果是顾客，则可修改顾客类型和余额
	if user.Role == "customer" {
		fmt.Printf("当前顾客类型: %s\n", user.CustomerType)
		fmt.Print("请选择新的顾客类型（1. 会员账号 2. 普通账号，回车保持不变）：")
		ctChoice := readLine()
		if ctChoice == "1" {
			user.CustomerType = "member"
		} else if ctChoice == "2" {
			user.CustomerType = "regular"
		}
		fmt.Printf("当前余额: %.2f\n", user.Balance)
		fmt.Print("请输入新的余额（回车保持不变）：")
		balanceStr := readLine()
		if balanceStr != "" {
			b, err := strconv.ParseFloat(balanceStr, 64)
			if err == nil {
				user.Balance = b
			} else {
				fmt.Println("无效的余额输入")
			}
		}
	}
	saveUsers()
	fmt.Println("用户信息更新成功")
}

// deleteUser 删除指定用户（管理员操作）
func deleteUser() {
	fmt.Print("请输入要删除的用户ID：")
	idStr := readLine()
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("无效的ID")
		return
	}
	index := -1
	for i, user := range users {
		if user.ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("未找到该用户")
		return
	}
	fmt.Print("确定要删除该用户吗？(y/n): ")
	confirm := readLine()
	if confirm != "y" && confirm != "Y" {
		return
	}
	users = append(users[:index], users[index+1:]...)
	saveUsers()
	fmt.Println("用户删除成功")
}

// adminRoomManagement 管理员对房间的增删改查操作
func adminRoomManagement() {
	for {
		fmt.Println("--------- 房间管理 ---------")
		fmt.Println("1. 查看所有房间")
		fmt.Println("2. 添加房间")
		fmt.Println("3. 修改房间")
		fmt.Println("4. 删除房间")
		fmt.Println("5. 返回上一层")
		fmt.Print("请选择操作：")
		choice := readLine()
		switch choice {
		case "1":
			listRooms()
		case "2":
			addRoom()
		case "3":
			updateRoom()
		case "4":
			deleteRoom()
		case "5":
			return
		default:
			fmt.Println("无效的选项，请重试。")
		}
	}
}

// listRooms 显示所有房间信息
func listRooms() {
	if len(rooms) == 0 {
		fmt.Println("当前无房间信息")
		return
	}
	fmt.Println("----- 房间列表 -----")
	for _, room := range rooms {
		fmt.Printf("ID: %d, 类型: %s, 价格: %.2f, 总数: %d, 剩余: %d\n",
			room.ID, room.Type, room.Price, room.Total, room.Available)
	}
}

// addRoom 添加新房间（仅管理员操作）
func addRoom() {
	fmt.Println("----- 添加新房间 -----")
	fmt.Print("请输入房间类型：")
	roomType := readLine()
	fmt.Print("请输入房间价格：")
	priceStr := readLine()
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		fmt.Println("无效的价格输入")
		return
	}
	fmt.Print("请输入房间总数：")
	totalStr := readLine()
	total, err := strconv.Atoi(totalStr)
	if err != nil {
		fmt.Println("无效的房间数量")
		return
	}
	newRoom := Room{
		ID:        getNextRoomID(),
		Type:      roomType,
		Price:     price,
		Total:     total,
		Available: total,
	}
	rooms = append(rooms, newRoom)
	saveRooms()
	fmt.Println("房间添加成功！")
}

// updateRoom 修改房间信息
func updateRoom() {
	fmt.Print("请输入要修改的房间ID：")
	idStr := readLine()
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("无效的ID")
		return
	}
	var room *Room
	for i := range rooms {
		if rooms[i].ID == id {
			room = &rooms[i]
			break
		}
	}
	if room == nil {
		fmt.Println("未找到该房间")
		return
	}
	fmt.Printf("当前房间类型: %s\n", room.Type)
	fmt.Print("请输入新的房间类型（回车保持不变）：")
	newType := readLine()
	if newType != "" {
		room.Type = newType
	}
	fmt.Printf("当前价格: %.2f\n", room.Price)
	fmt.Print("请输入新的价格（回车保持不变）：")
	priceStr := readLine()
	if priceStr != "" {
		price, err := strconv.ParseFloat(priceStr, 64)
		if err == nil {
			room.Price = price
		} else {
			fmt.Println("无效的价格输入")
		}
	}
	fmt.Printf("当前总数: %d\n", room.Total)
	fmt.Print("请输入新的总数（回车保持不变）：")
	totalStr := readLine()
	if totalStr != "" {
		total, err := strconv.Atoi(totalStr)
		if err == nil {
			// 调整剩余数量（假设已有预订时不允许负数）
			diff := total - room.Total
			room.Total = total
			room.Available += diff
			if room.Available < 0 {
				room.Available = 0
			}
		} else {
			fmt.Println("无效的数量输入")
		}
	}
	saveRooms()
	fmt.Println("房间信息更新成功")
}

// deleteRoom 删除房间（仅管理员操作）
func deleteRoom() {
	fmt.Print("请输入要删除的房间ID：")
	idStr := readLine()
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("无效的ID")
		return
	}
	index := -1
	for i, room := range rooms {
		if room.ID == id {
			index = i
			break
		}
	}
	if index == -1 {
		fmt.Println("未找到该房间")
		return
	}
	fmt.Print("确定要删除该房间吗？(y/n): ")
	confirm := readLine()
	if confirm != "y" && confirm != "Y" {
		return
	}
	rooms = append(rooms[:index], rooms[index+1:]...)
	saveRooms()
	fmt.Println("房间删除成功")
}

// ------------------------- 顾客功能 ----------------------------

// customerMenu 为顾客提供房间查询、预订及查看余额的菜单
func customerMenu(user *User) {
	for {
		fmt.Println("================================")
		fmt.Println("顾客菜单")
		fmt.Println("1. 查看房间信息")
		fmt.Println("2. 预订房间")
		fmt.Println("3. 查看余额")
		fmt.Println("4. 退出")
		fmt.Print("请选择操作：")
		choice := readLine()
		switch choice {
		case "1":
			listRooms()
		case "2":
			bookRoom(user)
		case "3":
			fmt.Printf("当前余额: %.2f\n", user.Balance)
		case "4":
			fmt.Println("注销成功")
			saveUsers() // 保存余额变动
			return
		default:
			fmt.Println("无效的选项，请重试。")
		}
	}
}

// bookRoom 实现顾客预订房间：检查房间剩余数量和余额，预订成功后扣款并更新房间状态
func bookRoom(customer *User) {
	if len(rooms) == 0 {
		fmt.Println("当前无可预订的房间")
		return
	}
	listRooms()
	fmt.Print("请输入要预订的房间ID：")
	idStr := readLine()
	id, err := strconv.Atoi(idStr)
	if err != nil {
		fmt.Println("无效的房间ID")
		return
	}
	var room *Room
	for i := range rooms {
		if rooms[i].ID == id {
			room = &rooms[i]
			break
		}
	}
	if room == nil {
		fmt.Println("未找到该房间")
		return
	}
	fmt.Printf("选择的房间: %s, 单价: %.2f, 剩余数量: %d\n", room.Type, room.Price, room.Available)
	fmt.Print("请输入预订数量：")
	quantityStr := readLine()
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil || quantity <= 0 {
		fmt.Println("无效的数量")
		return
	}
	if quantity > room.Available {
		fmt.Println("预订数量超过剩余房间数")
		return
	}
	totalCost := room.Price * float64(quantity)
	if customer.Balance < totalCost {
		fmt.Println("余额不足，无法预订")
		return
	}
	// 扣减余额并更新房间剩余数量
	customer.Balance -= totalCost
	room.Available -= quantity
	saveUsers()
	saveRooms()
	fmt.Printf("预订成功！共扣款 %.2f 元，剩余余额: %.2f\n", totalCost, customer.Balance)
}
