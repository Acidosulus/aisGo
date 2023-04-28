package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-ini/ini"
	"github.com/xuri/excelize/v2"
)

type Agreement struct {
	row_id int64  `sql:"ROW_ID"`
	number string `sql:"Номер"`
}

var (
	agr []Agreement
)

type TLogin struct {
	Server   string
	User     string
	Password string
	Port     string
	Database string
}

type TRow struct {
	Fio                    string
	Total, Lk, Edo, Reject int
}

type TRows []TRow

var Rows TRows

// Загружает в структуру данные для авторизации MSSQL
func GetLoginOptions(filepath string) TLogin {
	var result TLogin
	cfg, err := ini.Load(filepath)
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}
	result.Server = cfg.Section("login").Key("SERVER").String()
	fmt.Println("SERVER:", result.Server)
	result.User = cfg.Section("login").Key("UID").String()
	fmt.Println("USER:", result.User)
	result.Password = cfg.Section("login").Key("PWD").String()
	fmt.Println("PASSWORD:", result.Password)
	result.Database = cfg.Section("login").Key("DATABASE").String()
	fmt.Println("DATABASE:", result.Database)

	return result
}

var Login TLogin

func FileToString(filepath string) string {
	str, err := os.ReadFile(filepath)
	if err != nil {
		fmt.Print(err)
	}
	return string(str) // convert content to a 'string'
}

func Report_of_Rejects_from_Paperprint_DocumentsHandler(writer http.ResponseWriter, request *http.Request) {

	Login = GetLoginOptions("settings.ini")
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;encrypt=disable", Login.Server, Login.User, Login.Password, Login.Database)

	conn, err := sql.Open("mssql", connString)
	fmt.Printf("=====>> %t", conn)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer conn.Close()
	rows, err := conn.Query(FileToString("./SQL/Report_of_Rejects_from_Paperprint_Documents.sql"))
	if err != nil {
		log.Fatal("", err.Error())
	}
	fmt.Println(rows.Columns())
	var fio string
	var total, lk, edo, reject int
	var row TRow
	Rows = Rows[:0]
	for rows.Next() {
		fio, total, lk, edo, reject, row.Fio, row.Total, row.Lk, row.Edo, row.Reject = "", 0, 0, 0, 0, "", 0, 0, 0, 0
		err = rows.Scan(&fio, &total, &lk, &edo, &reject)
		if err != nil {
			fmt.Println(err)
		}
		row.Fio, row.Total, row.Lk, row.Edo, row.Reject = fio, total, lk, edo, reject
		Rows = append(Rows, row)
	}

	fmt.Println(Rows)

	html, err := template.ParseFiles("./patterns/reports/Report_of_Rejects_from_Paperprint_Documents.html")
	if err != nil {
		fmt.Println(err)
		http.Error(writer, "Internal Server Error\n"+err.Error(), 500)
	}

	err = html.Execute(writer, Rows)
	if err != nil {
		fmt.Println(err)
		http.Error(writer, "Internal Server Error\n"+err.Error(), 500)
	}

}

func viewHandlerIndex(writer http.ResponseWriter, request *http.Request) {

	html, err := template.ParseFiles("./patterns/index.html")
	if err != nil {
		fmt.Println(err)
		http.Error(writer, "Internal Server Error\n"+err.Error(), 500)
	}

	err = html.Execute(writer, nil)
	if err != nil {
		fmt.Println(err)
		http.Error(writer, "Internal Server Error\n"+err.Error(), 500)
	}

}

func PrepareAndReturnExcel() *excelize.File {

	Login = GetLoginOptions("settings.ini")
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;encrypt=disable", Login.Server, Login.User, Login.Password, Login.Database)

	conn, err := sql.Open("mssql", connString)
	fmt.Printf("=====>> %t", conn)
	if err != nil {
		log.Fatal("Open connection failed:", err.Error())
	}
	defer conn.Close()
	rows, err := conn.Query(FileToString("./SQL/Report_of_Rejects_from_Paperprint_Documents.sql"))
	if err != nil {
		log.Fatal("", err.Error())
	}
	fmt.Println(rows.Columns())
	var fio string
	var total, lk, edo, reject int
	var row TRow
	Rows = Rows[:0]
	for rows.Next() {
		fio, total, lk, edo, reject, row.Fio, row.Total, row.Lk, row.Edo, row.Reject = "", 0, 0, 0, 0, "", 0, 0, 0, 0
		err = rows.Scan(&fio, &total, &lk, &edo, &reject)
		if err != nil {
			fmt.Println(err)
		}
		row.Fio, row.Total, row.Lk, row.Edo, row.Reject = fio, total, lk, edo, reject
		Rows = append(Rows, row)
	}

	f := excelize.NewFile()
	counter := 0
	counter++
	f.SetCellValue("Sheet1", fmt.Sprintf("A%v", counter), "Переход на безбумажные платёжные документы")
	counter++
	f.SetCellValue("Sheet1", fmt.Sprintf("A%v", counter), "ФИО")
	f.SetCellValue("Sheet1", fmt.Sprintf("B%v", counter), "Всего")
	f.SetCellValue("Sheet1", fmt.Sprintf("C%v", counter), "ЛК")
	f.SetCellValue("Sheet1", fmt.Sprintf("D%v", counter), "ЭДО")
	f.SetCellValue("Sheet1", fmt.Sprintf("E%v", counter), "Отказ от БМ")

	for _, rr := range Rows {
		counter++
		f.SetCellValue("Sheet1", fmt.Sprintf("A%v", counter), rr.Fio)
		f.SetCellValue("Sheet1", fmt.Sprintf("B%v", counter), rr.Total)
		f.SetCellValue("Sheet1", fmt.Sprintf("C%v", counter), rr.Lk)
		f.SetCellValue("Sheet1", fmt.Sprintf("D%v", counter), rr.Edo)
		f.SetCellValue("Sheet1", fmt.Sprintf("E%v", counter), rr.Reject)
	}
	return f
}

func downloadExcel(w http.ResponseWriter, r *http.Request) {
	// Get the Excel file with the user input data
	file := PrepareAndReturnExcel()
	// Set the headers necessary to get browsers to interpret the downloadable file
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment;filename="userInputData.xlsx`)
	w.Header().Set("File-Name", "userInputData.xlsx")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Expires", "0")
	err := file.Write(w)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	fmt.Println("READY!!!")
	http.HandleFunc("/", viewHandlerIndex)

	http.HandleFunc("/reports/Report_of_Rejects_from_Paperprint_Documents", Report_of_Rejects_from_Paperprint_DocumentsHandler)
	http.HandleFunc("/reports/Report_of_Rejects_from_Paperprint_Documents_Excel", downloadExcel)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.ListenAndServe("localhost:8080", nil)

}
