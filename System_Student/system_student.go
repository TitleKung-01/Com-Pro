package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"runtime"
	"sort"
	"strings"
)

type Student struct {
	Name          string            `json:"name"`
	Score         float64           `json:"score"`
	Grade         string            `json:"grade"`
	Money         float64           `json:"money"`
	PurchasedItems map[string]int   `json:"purchasedItems"` 
}

type Product struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

var shopItems = []Product{
	{"Notebook", 10.0},
	{"Pen", 2.0},
	{"Backpack", 50.0},
	{"Calculator", 30.0},
	{"USB Drive", 15.0},
}

var students map[string]*Student

func main() {
	for {
		students = loadStudentsFromFile()

		clearScreen()
		fmt.Println("===== เมนูระบบจัดการคะแนนศึกษา =====")
		fmt.Println("1. สร้างรายชื่อนักศึกษา")
		fmt.Println("2. เช็คคะแนน")
		fmt.Println("3. แก้ไขคะแนน")
		fmt.Println("4. เพิ่มคะแนน")
		fmt.Println("5. เช็ครายชื่อนักศึกษา")
		fmt.Println("6. ลบข้อมูลนักศึกษา")
		fmt.Println("7. ร้านค้า")
		fmt.Println("8. เพิ่มเงิน")
		fmt.Println("9. ออกจากโปรแกรม")
		fmt.Println("=====================================")

		var choice string
		fmt.Print("เลือกเมนู: ")
		fmt.Scanln(&choice)

		clearScreen()

		switch choice {
		case "1":
			addStudent()
		case "2":
			checkScore()
		case "3":
			editScore()
		case "4":
			addScore()
		case "5":
			listStudents()
		case "6":
			deleteStudent()
		case "7":
			purchaseItem()
		case "8":
			addMoney()
		case "9":
			fmt.Println("ออกจากโปรแกรม")
			return
		default:
			fmt.Println("กรุณาเลือกเมนูที่ถูกต้อง")
		}

		fmt.Println("\nกด Enter เพื่อกลับสู่เมนูหลัก")
		fmt.Scanln()
	}
}

func clearScreen() {
	switch runtime.GOOS {
	case "windows":
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	default:
		fmt.Print("\033[2J")
		fmt.Print("\033[H")
	}
}

func addStudent() {
	var id, name string
	var score, money float64

	fmt.Print("ใส่รหัสนักศึกษา: ")
	fmt.Scanln(&id)
	fmt.Print("ใส่ชื่อนักศึกษา: ")
	fmt.Scanln(&name)
	fmt.Print("ใส่คะแนน: ")
	fmt.Scanln(&score)
	fmt.Print("ใส่เงิน: ")
	fmt.Scanln(&money)

	if _, exists := students[id]; exists {
		fmt.Println("รหัสนักศึกษานี้มีอยู่ในระบบแล้ว")
		return
	}

	students[id] = &Student{Name: name, Score: score, Money: money, PurchasedItems: make(map[string]int)}
	fmt.Println("เพิ่มข้อมูลนักศึกษาเรียบร้อยแล้ว")
	saveStudentsToFile()
}

func checkScore() {
	var id string
	fmt.Print("ใส่รหัสนักศึกษา: ")
	fmt.Scanln(&id)

	student, exists := students[id]
	if !exists {
		fmt.Println("ไม่พบรหัสนักศึกษานี้ในระบบ")
		return
	}

	calculateGrade(student) // คำนวณเกรดใหม่ทุกครั้งที่ตรวจสอบคะแนน
	fmt.Printf("นักศึกษา %s มีคะแนน %.2f และเกรด %s\n", student.Name, student.Score, student.Grade)
}

func editScore() {
	var id string
	var newScore float64

	fmt.Print("ใส่รหัสนักศึกษา: ")
	fmt.Scanln(&id)

	student, exists := students[id]
	if !exists {
		fmt.Println("ไม่พบรหัสนักศึกษานี้ในระบบ")
		return
	}

	fmt.Printf("คะแนนปัจจุบันของ %s: %.2f\n", student.Name, student.Score)
	fmt.Print("ใส่คะแนนใหม่: ")
	fmt.Scanln(&newScore)

	student.Score = newScore
	calculateGrade(student) // คำนวณเกรดใหม่หลังจากแก้ไขคะแนน

	fmt.Println("แก้ไขคะแนนเรียบร้อยแล้ว")
	saveStudentsToFile() // บันทึกข้อมูลหลังการแก้ไข
}

func addScore() {
	var id string
	var additionalScore float64

	fmt.Print("ใส่รหัสนักศึกษา: ")
	fmt.Scanln(&id)

	student, exists := students[id]
	if !exists {
		fmt.Println("ไม่พบรหัสนักศึกษานี้ในระบบ")
		return
	}

	fmt.Printf("คะแนนปัจจุบันของ %s: %.2f\n", student.Name, student.Score)
	fmt.Print("ใส่คะแนนที่ต้องการเพิ่ม: ")
	fmt.Scanln(&additionalScore)

	student.Score += additionalScore
	if student.Score > 100 {
		student.Score = 100
		fmt.Println("คะแนนเกิน 100 จึงปรับให้เป็น 100")
	}
	calculateGrade(student) // คำนวณเกรดใหม่หลังจากเพิ่มคะแนน

	fmt.Printf("เพิ่มคะแนนเรียบร้อยแล้ว คะแนนใหม่ของ %s คือ %.2f\n", student.Name, student.Score)
	saveStudentsToFile() // บันทึกข้อมูลหลังการแก้ไข
}

func listStudents() {
	if len(students) == 0 {
		fmt.Println("ไม่มีนักศึกษาในระบบ")
		return
	}

	fmt.Println("===== รายชื่อนักศึกษาทั้งหมด =====")
	keys := make([]string, 0, len(students))
	for k := range students {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, id := range keys {
		student := students[id]
		calculateGrade(student)
		purchasedItems := getPurchasedItems(student)
		
		fmt.Printf("รหัส: %s\nชื่อ: %s\nคะแนน: %.2f\nเกรด: %s\nเงินคงเหลือ: %.2f บาท\nสินค้าที่ซื้อ: %s\n",
			id, student.Name, student.Score, student.Grade, student.Money, purchasedItems)
		fmt.Println(strings.Repeat("-", 40))
	}
}

func calculateGrade(student *Student) {
	switch {
	case student.Score >= 90:
		student.Grade = "S+"
	case student.Score >= 80:
		student.Grade = "A"
	case student.Score >= 70:
		student.Grade = "B"
	case student.Score >= 60:
		student.Grade = "C"
	case student.Score >= 50:
		student.Grade = "D"
	default:
		student.Grade = "F"
	}
}

func getPurchasedItems(student *Student) string {
	if len(student.PurchasedItems) == 0 {
		return "ไม่มีการซื้อสินค้า"
	}
	items := make([]string, 0, len(student.PurchasedItems))
	for itemName, quantity := range student.PurchasedItems {
		items = append(items, fmt.Sprintf("%s: %d ชิ้น", itemName, quantity))
	}
	return strings.Join(items, ", ")
}

func saveStudentsToFile() {
	studentsData, err := json.Marshal(students)
	if err != nil {
		fmt.Println("เกิดข้อผิดพลาดในการบันทึกข้อมูล:", err)
		return
	}

	err = os.WriteFile("students.json", studentsData, fs.FileMode(0644))
	if err != nil {
		fmt.Println("เกิดข้อผิดพลาดในการบันทึกข้อมูล:", err)
		return
	}

	fmt.Println("บันทึกข้อมูลนักศึกษาเรียบร้อยแล้ว")
}

func loadStudentsFromFile() map[string]*Student {
	studentsData, err := os.ReadFile("students.json")
	if err != nil {
		fmt.Println("ไม่พบไฟล์ข้อมูลนักศึกษา, กำลังสร้างไฟล์ใหม่...")
		return make(map[string]*Student)
	}

	var students map[string]*Student
	err = json.Unmarshal(studentsData, &students)
	if err != nil {
		fmt.Println("เกิดข้อผิดพลาดในการอ่านข้อมูล:", err)
		return make(map[string]*Student)
	}

	return students
}

func deleteStudent() {
	var id string
	fmt.Print("ใส่รหัสนักศึกษา: ")
	fmt.Scanln(&id)

	student, exists := students[id]
	if !exists {
		fmt.Println("ไม่พบรหัสนักศึกษานี้ในระบบ")
		return
	}

	delete(students, id)
	fmt.Printf("ลบข้อมูลนักศึกษา %s (%s) เรียบร้อยแล้ว\n", id, student.Name)
	saveStudentsToFile()
}

func displayShop() {
	fmt.Println("===== ร้านค้า =====")
	for i, item := range shopItems {
		fmt.Printf("%d. %s ราคา: %.2f บาท\n", i+1, item.Name, item.Price)
	}
	fmt.Println(strings.Repeat("-", 40))
}

func purchaseItem() {
    var id string
    fmt.Print("ใส่รหัสนักศึกษา: ")
    fmt.Scanln(&id)

    student, exists := students[id]
    if !exists {
        fmt.Println("ไม่พบรหัสนักศึกษานี้ในระบบ")
        return
    }

    displayShop()

    for {
        var choice int
        fmt.Print("เลือกสินค้าที่ต้องการซื้อ (ใส่หมายเลข): ")
        fmt.Scanln(&choice)

        if choice < 1 || choice > len(shopItems) {
            fmt.Println("หมายเลขสินค้าที่เลือกไม่ถูกต้อง")
            continue
        }

        selectedItem := shopItems[choice-1]

        fmt.Print("ใส่จำนวนที่ต้องการซื้อ: ")
        var quantity int
        fmt.Scanln(&quantity)

        totalCost := selectedItem.Price * float64(quantity)
        if totalCost > student.Money {
            fmt.Println("เงินไม่เพียงพอสำหรับการซื้อสินค้า")
            continue
        }

        student.Money -= totalCost
        student.PurchasedItems[selectedItem.Name] += quantity
        fmt.Printf("ซื้อสินค้า %s จำนวน %d ชิ้น รวมเป็นเงิน %.2f บาท\n", selectedItem.Name, quantity, totalCost)

        fmt.Print("ต้องการซื้อสินค้าอื่นอีกหรือไม่? (y/n): ")
        var more string
        fmt.Scanln(&more)

        if strings.ToLower(more) != "y" {
            break
        }
    }

    saveStudentsToFile()
}

func addMoney() {
	var id string
	var amount float64

	fmt.Print("ใส่รหัสนักศึกษา: ")
	fmt.Scanln(&id)

	student, exists := students[id]
	if !exists {
		fmt.Println("ไม่พบรหัสนักศึกษานี้ในระบบ")
		return
	}

	fmt.Print("ใส่จำนวนเงินที่ต้องการเพิ่ม: ")
	fmt.Scanln(&amount)

	student.Money += amount
	fmt.Printf("เพิ่มเงิน %.2f บาท ให้กับนักศึกษา %s เรียบร้อยแล้ว เงินคงเหลือ: %.2f บาท\n", amount, student.Name, student.Money)
	saveStudentsToFile()
}
