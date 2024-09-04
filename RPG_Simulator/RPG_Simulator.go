package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

// User struct เก็บข้อมูลของผู้เล่น
type User struct {
	Username string
	Password string
	Job      string
}

// Job interface สำหรับอาชีพ
type Job interface {
	Describe() string
	Attack() int
}

// Warrior struct สำหรับอาชีพ Warrior
type Warrior struct{}

func (w Warrior) Describe() string {
	return "You are a brave Warrior, strong and fearless."
}

func (w Warrior) Attack() int {
	return rand.Intn(10) + 10 // Warrior โจมตีแรง
}

// Mage struct สำหรับอาชีพ Mage
type Mage struct{}

func (m Mage) Describe() string {
	return "You are a wise Mage, mastering the elements."
}

func (m Mage) Attack() int {
	return rand.Intn(15) + 5 // Mage โจมตีเวทมนตร์
}

// Thief struct สำหรับอาชีพ Thief
type Thief struct{}

func (t Thief) Describe() string {
	return "You are a stealthy Thief, quick and cunning."
}

func (t Thief) Attack() int {
	return rand.Intn(8) + 8 // Thief โจมตีรวดเร็ว
}

// Monster struct สำหรับมอนสเตอร์
type Monster struct {
	Name   string
	Health int
	Attack int
}

// checkDuplicateUser ฟังก์ชันตรวจสอบว่ามีผู้ใช้ซ้ำหรือไม่
func checkDuplicateUser(username string, users []User) bool {
	for _, user := range users {
		if user.Username == username {
			return true
		}
	}
	return false
}

// chooseJob ฟังก์ชันให้ผู้เล่นเลือกอาชีพ
func chooseJob() Job {
	var choice int
	fmt.Println()
	printTitle("Choose Your Job")
	fmt.Println("1. Warrior - Brave and fearless")
	fmt.Println("2. Mage - Wise and mastering the elements")
	fmt.Println("3. Thief - Stealthy and cunning")
	fmt.Print("Enter the number of your choice: ")
	fmt.Scan(&choice)
	fmt.Println()

	switch choice {
	case 1:
		return Warrior{}
	case 2:
		return Mage{}
	case 3:
		return Thief{}
	default:
		fmt.Println("Invalid choice, defaulting to Warrior.")
		return Warrior{}
	}
}

// printTitle ฟังก์ชันสำหรับพิมพ์หัวข้อ
func printTitle(title string) {
	border := strings.Repeat("=", len(title)+4)
	fmt.Println(border)
	fmt.Printf("= %s =\n", title)
	fmt.Println(border)
}

// saveUsersToFile บันทึกข้อมูลผู้ใช้ลงไฟล์ JSON
func saveUsersToFile(users []User, filename string) error {
	data, err := json.MarshalIndent(users, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// loadUsersFromFile โหลดข้อมูลผู้ใช้จากไฟล์ JSON
func loadUsersFromFile(filename string) ([]User, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		if os.IsNotExist(err) {
			return []User{}, nil // ถ้าไฟล์ไม่พบ ให้เริ่มต้นเป็น list ว่าง
		}
		return nil, err
	}
	var users []User
	err = json.Unmarshal(data, &users)
	if err != nil {
		return nil, err
	}
	return users, nil
}

// createMonsters ฟังก์ชันสำหรับสร้างรายการมอนสเตอร์
func createMonsters() []Monster {
	return []Monster{
		{Name: "Goblin", Health: rand.Intn(20) + 30, Attack: rand.Intn(5) + 5},
		{Name: "Orc", Health: rand.Intn(30) + 40, Attack: rand.Intn(10) + 10},
		{Name: "Dragon", Health: rand.Intn(50) + 80, Attack: rand.Intn(15) + 20},
		{Name: "Troll", Health: rand.Intn(40) + 50, Attack: rand.Intn(10) + 15},
	}
}

// battleMonster ฟังก์ชันสำหรับการต่อสู้กับมอนสเตอร์
func battleMonster(job Job, monster Monster) {
	fmt.Println("A wild", monster.Name, "appears!")
	for monster.Health > 0 {
		// Player attack
		playerDamage := job.Attack()
		monster.Health -= playerDamage
		fmt.Printf("You attack the %s for %d damage.\n", monster.Name, playerDamage)

		if monster.Health <= 0 {
			fmt.Println("You defeated the", monster.Name, "!")
			break
		}

		// Monster attack
		monsterDamage := monster.Attack
		fmt.Printf("The %s attacks you for %d damage.\n", monster.Name, monsterDamage)

		fmt.Println("Monster Health:", monster.Health)
		fmt.Println(strings.Repeat("-", 20))
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	const filename = "users.json"
	users, err := loadUsersFromFile(filename)
	if err != nil {
		fmt.Println("Error loading users:", err)
		return
	}

	// ล็อกอินและเช็คผู้เล่นซ้ำ
	var username, password string
	fmt.Println()
	printTitle("Welcome to RPG Simulator")
	fmt.Print("Enter username: ")
	fmt.Scan(&username)

	if checkDuplicateUser(username, users) {
		fmt.Println("Username already exists! Please try again.")
		return
	}

	fmt.Print("Enter password: ")
	fmt.Scan(&password)

	// เลือกอาชีพ
	job := chooseJob()

	// สร้างผู้ใช้ใหม่และเพิ่มเข้าไปในระบบ
	newUser := User{
		Username: username,
		Password: password,
		Job:      job.Describe(),
	}

	users = append(users, newUser)

	// บันทึกผู้ใช้ลงไฟล์ JSON
	err = saveUsersToFile(users, filename)
	if err != nil {
		fmt.Println("Error saving users:", err)
		return
	}

	fmt.Println()
	printTitle("Character Creation Complete")
	fmt.Printf("Welcome %s!\n", newUser.Username)
	fmt.Println(newUser.Job)
	fmt.Println(strings.Repeat("=", 30))

	// สร้างมอนสเตอร์หลายตัว
	monsters := createMonsters()

	// สุ่มเลือกมอนสเตอร์เพื่อต่อสู้
	monster := monsters[rand.Intn(len(monsters))]
	battleMonster(job, monster)
}
