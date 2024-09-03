package main

import (
	"fmt"
)

// Product เป็นโครงสร้างที่เก็บข้อมูลของสินค้า
type Product struct {
	ID    int
	Name  string
	Price float64
}

// Shop เป็นอินเตอร์เฟซที่กำหนดฟังก์ชันพื้นฐานของร้านค้า
type Shop interface {
	AddProduct(p Product)
	ListProducts()
	BuyProduct(id int)
}

// PetShop เป็นโครงสร้างที่มีรายการสินค้าและทำการ implement ฟังก์ชันของอินเตอร์เฟซ Shop
type PetShop struct {
	products []Product
}

// AddProduct เป็นฟังก์ชันที่ใช้เพิ่มสินค้าลงในร้านค้า
func (s *PetShop) AddProduct(p Product) {
	s.products = append(s.products, p) // เพิ่มสินค้าลงในร้านค้า
	fmt.Println("Product added:", p.Name)
}

// ListProducts แสดงรายการสินค้าทั้งหมดในร้านค้า
func (s *PetShop) ListProducts() {
	fmt.Println("Available Products:")
	for _, p := range s.products { // วนลูปเพื่อแสดงรายการสินค้าทั้งหมด
		fmt.Printf("ID: %d, Name: %s, Price: %.2f\n", p.ID, p.Name, p.Price)
	}
}

// BuyProduct เป็นฟังก์ชันที่ใช้ซื้อสินค้าจากร้านค้า
func (s *PetShop) BuyProduct(id int) { 
	for i, p := range s.products { // วนลูปเพื่อหาสินค้าที่ต้องการซื้อ
		if p.ID == id { // หาสินค้าที่ตรงกับ ID ที่ระบุ
			fmt.Printf("You bought: %s for %.2f\n", p.Name, p.Price) // แสดงสินค้าที่ซื้อ
			// ลบสินค้าที่ซื้อออกจากรายการ
			s.products = append(s.products[:i], s.products[i+1:]...) // ลบสินค้าที่ซื้อออกจากรายการ
			return
		}
	}
	fmt.Println("Product not found")
}

// ทำฟังก์ชันเพื่อทำงานกับอินเตอร์เฟซ Shop โดยเฉพาะ
func operateShop(shop Shop) {
	// เพิ่มสินค้าในร้านค้า
	shop.AddProduct(Product{ID: 3, Name: "Fish Food", Price: 2.99})
	shop.AddProduct(Product{ID: 4, Name: "Bird Cage", Price: 29.99})

	// แสดงรายการสินค้าในร้านค้า
	shop.ListProducts()

	// ซื้อสินค้า
	shop.BuyProduct(3)
	shop.BuyProduct(4)

	// แสดงรายการสินค้าหลังจากซื้อแล้ว
	shop.ListProducts()
}

func main() {
	// สร้างร้านค้าใหม่
	var myShop Shop = &PetShop{} // สร้างร้านค้าใหม่

	// ใช้งานร้านค้าผ่านอินเตอร์เฟซ
	operateShop(myShop) // ผลลัพธ์จะแสดงรายการสินค้าที่เพิ่ม และรายการสินค้าหลังจากซื้อ
}
