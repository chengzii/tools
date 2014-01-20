package mydb

import (
	"database/sql"
	"github.com/coopernurse/gorp"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"sync"
	"time"
)

var (
	dbmap *gorp.DbMap
	mutex sync.Mutex
)

func init() {

	mutex.Lock()
	defer mutex.Unlock()

	if dbmap != nil {
		return
	}

	db, err := sql.Open("mysql", "root:root@/chengzii")
	checkErr(err, "sql.Open failed")
	// construct a gorp DbMap
	dbmap = &gorp.DbMap{Db: db, Dialect: gorp.MySQLDialect{"InnoDB", "UTF8"}}

	// register the structs you wish to use with gorp
	// you can also use the shorter dbmap.AddTable() if you
	// don't want to override the table name
	//
	// SetKeys(true) means we have a auto increment primary key, which
	// will get automatically bound to your struct post-insert
	//
	dbmap.AddTable(Chz_food{}).SetKeys(true, "Id")
	dbmap.AddTable(Chz_play{}).SetKeys(true, "Id")
	err = dbmap.CreateTablesIfNotExists()
	checkErr(err, "Create tables failed")
}

type Chz_food struct {
	Id        int64  `db:"id"`
	Title     string `db:"title"`
	Summary   string `db:"summary"`
	Link      string `db:"link"`
	Pic       string `db:"pic"`
	Flag      string `db:"flag"`
	Date_time int64  `db:"date_time"`
	Save_time int64  `db:"save_time"`
}
type Chz_play struct {
	Id        int64  `db:"id"`
	Title     string `db:"title"`
	Summary   string `db:"summary"`
	Link      string `db:"link"`
	Pic       string `db:"pic"`
	Flag      string `db:"flag"`
	Date_time int64  `db:"date_time"`
	Save_time int64  `db:"save_time"`
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
func Insert_fd(title string, summary string, link string, pic string, tag string, date_time int64) error {
	save_time := time.Now().Unix()
	sum := &Chz_food{1, title, summary, link, pic, tag, date_time, save_time}
	err := dbmap.Insert(sum)
	//fmt.Println("insert id is : ", sum.Id)
	return err
}
func Update_fd(l *Chz_food) (int64, error) {
	count, err := dbmap.Update(l)
	return count, err
}
func Delete_fd(l *Chz_food) (int64, error) {
	count, err := dbmap.Delete(l)
	return count, err
}
func Deletebyid_fd(id int64) (count int64, err error) {
	l, err := Selectbyid_fd(id)
	if l == nil || err != nil {
		//fmt.Println(err)
		checkErr(err, "DELETE BY ID: "+strconv.Itoa(int(id))+" ERROR!")
		return
	}
	count, err = Delete_fd(l)
	return
}
func Selectbyid_fd(id int64) (inv *Chz_food, err error) {
	obj, err := dbmap.Get(Chz_food{}, id)
	if obj == nil {
		return
	}
	inv = obj.(*Chz_food)
	return
}
func Selectbynum_fd(num int64) (res []Chz_food, err error) {
	_, err = dbmap.Select(&res, "select * from chz_food order by date_time desc limit "+strconv.Itoa(int(num)))
	return
}
func SelectAll_fd() (res []Chz_food, err error) {
	_, err = dbmap.Select(&res, "select * from chz_food order by id desc")
	return
}
func Insert_pl(title string, summary string, link string, pic string, tag string, date_time int64) error {
	save_time := time.Now().Unix()
	sum := &Chz_play{1, title, summary, link, pic, tag, date_time, save_time}
	err := dbmap.Insert(sum)
	//fmt.Println("insert id is : ", sum.Id)
	return err
}
func Update_pl(l *Chz_play) (int64, error) {
	count, err := dbmap.Update(l)
	return count, err
}
func Delete_pl(l *Chz_play) (int64, error) {
	count, err := dbmap.Delete(l)
	return count, err
}
func Deletebyid_pl(id int64) (count int64, err error) {
	l, err := Selectbyid_pl(id)
	if l == nil || err != nil {
		//fmt.Println(err)
		checkErr(err, "DELETE BY ID: "+strconv.Itoa(int(id))+" ERROR!")
		return
	}
	count, err = Delete_pl(l)
	return
}
func Selectbyid_pl(id int64) (inv *Chz_play, err error) {
	obj, err := dbmap.Get(Chz_food{}, id)
	if obj == nil {
		return
	}
	inv = obj.(*Chz_play)
	return
}
func Selectbynum_pl(num int64) (res []Chz_play, err error) {
	_, err = dbmap.Select(&res, "select * from chz_play order by date_time desc limit "+strconv.Itoa(int(num)))
	return
}
func SelectAll_pl() (res []Chz_play, err error) {
	_, err = dbmap.Select(&res, "select * from chz_play order by id desc")
	return
}
