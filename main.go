package main

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/olekukonko/tablewriter"
)

type Menu struct {
	kd_menu   string
	nama_menu string
	kategori  string
	harga     int
}

func connect() (*sql.DB, error) {
	return sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/restoran")
}

func clearScreen() {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func getAllMenus() {
	db, err := connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM menu")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var menus []Menu
	for rows.Next() {
		var menu Menu
		err := rows.Scan(&menu.kd_menu, &menu.nama_menu, &menu.kategori, &menu.harga)
		if err != nil {
			panic(err)
		}
		menus = append(menus, menu)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Kode Menu", "Nama Menu", "Kategori", "Harga"})

	if menus == nil {
		table.Append([]string{"", "", "", ""})
	} else {
		for _, menu := range menus {
			table.Append([]string{menu.kd_menu, menu.nama_menu, menu.kategori, fmt.Sprintf("%d", menu.harga)})
		}
	}

	table.Render()
}

func getMenuByCategory(kategori string) {
	db, err := connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM menu WHERE kategori = ?", kategori)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var menus []Menu
	for rows.Next() {
		var menu Menu
		err := rows.Scan(&menu.kd_menu, &menu.nama_menu, &menu.kategori, &menu.harga)
		if err != nil {
			panic(err)
		}
		menus = append(menus, menu)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Kode Menu", "Nama Menu", "Kategori", "Harga"})

	if menus == nil {
		table.Append([]string{"", "", "", ""})
	} else {
		for _, menu := range menus {
			table.Append([]string{menu.kd_menu, menu.nama_menu, menu.kategori, fmt.Sprintf("%d", menu.harga)})
		}
	}

	table.Render()
}

func getMenuByID(id string) {
	db, err := connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	rows, err := db.Query("SELECT * FROM menu WHERE kd_menu = ?", id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	var menu Menu
	for rows.Next() {
		err := rows.Scan(&menu.kd_menu, &menu.nama_menu, &menu.kategori, &menu.harga)
		if err != nil {
			panic(err)
		}
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Kode Menu", "Nama Menu", "Kategori", "Harga"})

	if menu.kd_menu == "" {
		table.Append([]string{"", "", "", ""})
	} else {
		table.Append([]string{menu.kd_menu, menu.nama_menu, menu.kategori, fmt.Sprintf("%d", menu.harga)})
	}

	table.Render()
}

func addMenu(id string, nama string, kategori string, harga int) {
	db, err := connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO menu (kd_menu, nama_menu, kategori, harga) VALUES (?, ?, ?, ?)", id, nama, kategori, harga)
	if err != nil {
		panic(err)
	}

	fmt.Println("Menu berhasil ditambahkan")
}

func updateMenu(id string, nama string, kategori string, harga int) {
	db, err := connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM menu WHERE kd_menu = ?)", id).Scan(&exists)
	if err != nil {
		panic(err)
	}

	if exists {
		_, err = db.Exec("UPDATE menu SET nama_menu = ?, kategori = ?, harga = ? WHERE kd_menu = ?", nama, kategori, harga, id)
		if err != nil {
			panic(err)
		}
		fmt.Println("Menu berhasil diupdate")
	} else {
		fmt.Println("Kode menu tidak ditemukan")
	}
}

func deleteMenu(id string) {
	db, err := connect()
	if err != nil {
		panic(err)
	}
	defer db.Close()

	var exists bool
	err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM menu WHERE kd_menu = ?)", id).Scan(&exists)
	if err != nil {
		panic(err)
	}

	if exists {
		_, err = db.Exec("DELETE FROM menu WHERE kd_menu = ?", id)
		if err != nil {
			panic(err)
		}
		fmt.Println("Menu berhasil dihapus")
	} else {
		fmt.Println("Kode menu tidak ditemukan")
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var choice int
	var id string
	var nama string
	var kategori string
	var harga int

	for {
		clearScreen()
		fmt.Println("1. Lihat Menu")
		fmt.Println("2. Lihat Menu Berdasarkan Kategori")
		fmt.Println("3. Lihat Menu Berdasarkan Kode Menu")
		fmt.Println("4. Tambah Menu")
		fmt.Println("5. Update Menu")
		fmt.Println("6. Hapus Menu")
		fmt.Println("7. Keluar")
		fmt.Print("Pilihan: ")
		fmt.Scanln(&choice)

		switch choice {
		case 1:
			getAllMenus()
		case 2:
			fmt.Print("Cari kategori: ")
			fmt.Scanln(&kategori)
			getMenuByCategory(kategori)
		case 3:
			fmt.Print("Cari Kode Menu: ")
			fmt.Scanln(&id)
			getMenuByID(id)
		case 4:
			fmt.Print("Kode Menu: ")
			fmt.Scanln(&id)
			fmt.Print("Nama Menu: ")
			nama, _ = reader.ReadString('\n')
			nama = strings.TrimSpace(nama)
			fmt.Print("Kategori: ")
			fmt.Scanln(&kategori)
			fmt.Print("Harga: ")
			fmt.Scanln(&harga)
			addMenu(id, nama, kategori, harga)
		case 5:
			fmt.Print("Kode Menu: ")
			fmt.Scanln(&id)
			fmt.Print("Nama Menu: ")
			nama, _ = reader.ReadString('\n')
			nama = strings.TrimSpace(nama)
			fmt.Print("Kategori: ")
			fmt.Scanln(&kategori)
			fmt.Print("Harga: ")
			fmt.Scanln(&harga)
			updateMenu(id, nama, kategori, harga)
		case 6:
			fmt.Print("Kode Menu: ")
			fmt.Scanln(&id)
			deleteMenu(id)
		case 7:
			for i := 3; i >= 0; i-- {
				fmt.Printf("\rProgram berhenti pada %d...", i)
				time.Sleep(1 * time.Second)
			}
			fmt.Println()
			clearScreen()
			os.Exit(0)
		default:
			fmt.Println("Pilihan tidak valid")
		}
		fmt.Print("Tekan enter untuk melanjutkan...")
		fmt.Scanln()
	}
}
