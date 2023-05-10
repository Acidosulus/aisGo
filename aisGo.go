package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"
	"time"

	"os"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/go-ini/ini"
	"github.com/xuri/excelize/v2"
)

var Connection *sql.DB

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

func (self *Routers) Report_of_Rejects_from_Paperprint_DocumentsHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Println("Report_of_Rejects_from_Paperprint_DocumentsHandler", self.Connection)
	rows, err := self.Connection.Query(FileToString("./SQL/Report_of_Rejects_from_Paperprint_Documents.sql"))
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

func (self *Routers) viewHandlerIndex(writer http.ResponseWriter, request *http.Request) {

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
	fmt.Println("./patterns/index.html")
}

func (self *Routers) PrepareAndReturnExcel() *excelize.File {

	rows, err := self.Connection.Query(FileToString("./SQL/Report_of_Rejects_from_Paperprint_Documents.sql"))
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

	file := excelize.NewFile()
	counter := 0
	counter++

	file.SetCellValue("Sheet1", fmt.Sprintf("A%v", counter), "Переход на безбумажные платёжные документы")
	counter++
	file.SetCellValue("Sheet1", fmt.Sprintf("A%v", counter), "ФИО")
	file.SetCellValue("Sheet1", fmt.Sprintf("B%v", counter), "Всего")
	file.SetCellValue("Sheet1", fmt.Sprintf("C%v", counter), "ЛК")
	file.SetCellValue("Sheet1", fmt.Sprintf("D%v", counter), "ЭДО")
	file.SetCellValue("Sheet1", fmt.Sprintf("E%v", counter), "Отказ от БМ")

	for _, rr := range Rows {
		counter++
		file.SetCellValue("Sheet1", fmt.Sprintf("A%v", counter), rr.Fio)
		file.SetCellValue("Sheet1", fmt.Sprintf("B%v", counter), rr.Total)
		file.SetCellValue("Sheet1", fmt.Sprintf("C%v", counter), rr.Lk)
		file.SetCellValue("Sheet1", fmt.Sprintf("D%v", counter), rr.Edo)
		file.SetCellValue("Sheet1", fmt.Sprintf("E%v", counter), rr.Reject)
	}

	style, err := file.NewStyle(
		&excelize.Style{
			Alignment: &excelize.Alignment{Horizontal: "left"},
			Font:      &excelize.Font{Bold: true, Color: "00000000"},
			Border: []excelize.Border{
				{Type: "left", Color: "00FF0000", Style: 1},
				{Type: "right", Color: "00FF0000", Style: 1},
				{Type: "top", Color: "00FF0000", Style: 1},
				{Type: "bottom", Color: "00FF0000", Style: 1},
			},
			Fill: excelize.Fill{
				Type:    "pattern",
				Color:   []string{"#DDDDDD"},
				Pattern: 1,
			},
		},
	)
	if err != nil {
		log.Fatal(err)
	}
	styleTitle, err := file.NewStyle(&excelize.Style{
		Fill: excelize.Fill{
			Type: "pattern",
			// navy blue
			Color:   []string{"#000000"},
			Pattern: 1,
		},
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF"},
		Alignment: &excelize.Alignment{Horizontal: "left"},
	})
	if err != nil {
		log.Fatal(err)
	}

	file.SetCellStyle("Sheet1", "A1", "e1", styleTitle)
	file.SetCellStyle("Sheet1", "A2", fmt.Sprintf("E%v", counter), style)
	file.SetColWidth("Sheet1", "A", "A", 45)
	file.SetColWidth("Sheet1", "B", "E", 15)
	return file
}

func (self *Routers) downloadExcel(writer http.ResponseWriter, r *http.Request) {
	// Get the Excel file with the user input data
	file := self.PrepareAndReturnExcel()
	// Set the headers necessary to get browsers to interpret the downloadable file

	writer.Header().Set("Content-Type", "application/octet-stream")
	writer.Header().Set("Content-Disposition", `attachment;filename="userInputData.xlsx`)
	writer.Header().Set("File-Name", "userInputData.xlsx")
	writer.Header().Set("Content-Transfer-Encoding", "binary")
	writer.Header().Set("Expires", "0")

	err := file.Write(writer)
	if err != nil {
		fmt.Println(err)
	}
}

type Agreement_Folder struct {
	Nc       string
	Folder   string
	District string
}

type Agreement_Folders struct {
	Last_Update_Time time.Time
	Items            []Agreement_Folder
}

func (self *Agreement_Folders) UpdateData() {

	rows, err := Connection.Query(FileToString("./SQL/Data_Agreement_Folders.sql"))
	if err != nil {
		log.Fatal("", err.Error())
	}
	fmt.Println(rows.Columns())
	return
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

	self.Last_Update_Time = time.Now()
}

type Routers struct {
	Connection *sql.DB
}

func (self *Routers) LoginToMSSQL() {
	Login = GetLoginOptions("settings.ini")
	connString := fmt.Sprintf("server=%s;user id=%s;password=%s;database=%s;encrypt=disable", Login.Server, Login.User, Login.Password, Login.Database)
	self.Connection, er = sql.Open("mssql", connString)
	fmt.Printf("%v =====>> %s\n", time.Now().Format("15:04:05.999"), self.Connection)
	if er != nil {
		log.Fatal("Open connection failed:", er.Error())
	}
	//defer self.Connection.Close()
}

var er error
var Router Routers
var aFolders Agreement_Folders

func main() {
	fmt.Println("=======================================================================================================================================================")
	Router.LoginToMSSQL()

	//aFolders.UpdateData()

	//os.Exit(0)

	fmt.Println("!!!SERVER READY!!!")
	http.HandleFunc("/", Router.viewHandlerIndex)
	http.HandleFunc("/reports/Report_of_Rejects_from_Paperprint_Documents", Router.Report_of_Rejects_from_Paperprint_DocumentsHandler)
	http.HandleFunc("/reports/Report_of_Rejects_from_Paperprint_Documents_Excel", Router.downloadExcel)
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))
	http.ListenAndServe("localhost:8080", nil)

}
