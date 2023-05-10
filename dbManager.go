package main

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"os"
	"strconv"
	"time"
)

func createDb(dbName string) *sql.DB {
	log.Println("createDb called")
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println("Could not open the database: ", err)
		panic(err)
	} else {
		log.Println("Opened: ", os.Getenv("DATABASE_URL"))
	}

	log.Println("Database creation: ", dbName)
	_, err = db.Exec(`CREATE DATABASE ` + dbName)
	if err != nil {
		log.Println("Database ", dbName, " already exists")
	} else {
		log.Println("Database created: ", dbName)
	}

	return db
}

func createTable(db *sql.DB, tableName int) {
	log.Println(db.Stats())
	name := strconv.Itoa(tableName)
	log.Println("Table ", name, " creation")
	query := fmt.Sprintf("CREATE TABLE public.u%d ( date timestamp, mood varchar)", tableName)
	_, err := db.Exec(query)

	log.Println("Executing: ", query)
	if err != nil {
		log.Println("Table creation for " + name + "failed")
		log.Println("Error: " + err.Error())
	} else {
		log.Println("Table" + name + "created")
	}
	log.Println("createTable done;")
}

func addToDb(userId int, mood string) {
	log.Println("addToDb called")

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println("Could not open the database: ", err)
		panic(err)
	} else {
		log.Println("Opened: ", os.Getenv("DATABASE_URL"))
	}
	defer db.Close()

	query := fmt.Sprintf(`INSERT INTO public.u%d (date, mood)  VALUES ($1, $2)`, userId)

	log.Println("Executing: ", query)
	_, err = db.Exec(query, time.Now(), mood)
	if err != nil {
		fmt.Println("Could not add row dating " + time.Now().Format("2006-01-02 15:04:05"))
		panic(err)
	}
}

//func getDataFromDb(userId int) string {
//	message := ""
//	log.Println("getDataFromDb called")
//	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
//	if err != nil {
//		log.Println("Could not open the database: ", err)
//		message = NO_DATA
//		return message
//	}
//	log.Println("Opened: ", os.Getenv("DATABASE_URL"))
//	defer db.Close()
//	query := fmt.Sprintf(`SELECT * FROM u%d`, userId)
//	rows, _ := db.Query(query)
//	log.Println(rows)
//	for row = rows.Next(){
//		log.Println(row)
//	}
//	return message
//}
