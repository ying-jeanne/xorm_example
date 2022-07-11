package main

import (
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

type Team struct {
	ID        int       `xorm:"'id' pk autoincr"`
	Name      string    `xorm:"name"`
	OrgID     int       `xorm:"org_id"`
	CreatedAt time.Time `xorm:"'created'"`
	UpdatedAt time.Time `xorm:"'updated'"`
	Email     string
}

type Team2 struct {
	ID   int    `xorm:"'id' pk autoincr"`
	Name string `xorm:"name"`
	// OrgID     int  when we are trying to not put tag, xorm is actually error out for get
	OrgID     int       `xorm:"org_id"`
	CreatedAt time.Time `xorm:"'created'"`
	UpdatedAt time.Time `xorm:"'updated'"`
	Email     string
}

func (t Team2) TableName() string {
	return "team"
}

func insertTeam(e *xorm.Engine, team1 Team) error {
	// for insert, xorm is actually checking that all the filled field has a corresponding column name, if not, error out
	// not set field would be filled directly with default value when creation
	var err error
	if _, err = e.Insert(&team1); err != nil {
		log.Fatal(err)
	}
	return err
}

func getTeam(e *xorm.Engine, name string) Team {
	teams := []Team{}
	err := e.Find(&teams)
	if err != nil {
		log.Fatal(err)
	}
	var team1 Team
	_, err = e.Where("name=?", name).Get(&team1)
	if err != nil {
		log.Fatal(err)
	}
	return team1
}

func deleteTeam(e *xorm.Engine, name string) {
	_, err := e.Exec("DELETE FROM team WHERE name=?", name)
	if err != nil {
		log.Fatal(err)
	}
}

func updateTeam(e *xorm.Engine, team Team) (int64, error) {
	affected, err := e.ID(team.ID).Update(team)
	return affected, err
}

func main() {
	engine, err := xorm.NewEngine("sqlite3", "grafana.db")
	if err != nil {
		log.Fatal(err)
	}
	engine.SetTableMapper(names.GonicMapper{})
	team1 := Team{Name: "myname4", OrgID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	err = insertTeam(engine, team1)
	if err != nil {
		log.Fatal(err)
	}

	team3 := getTeam(engine, team1.Name)
	fmt.Printf("the team3 is %v \n", team3)

	// here it is very confusing, since xorm omit the default value OrgID 0 here, we would have to
	// force the update by calling .AllCols().Update or .Cols("org_id").Update
	// but if we put .AllCols() we have to put all the fields that are mandatory so it is also not convinient
	team2 := Team{ID: team3.ID, OrgID: 0, Name: "princess"}
	_, err = updateTeam(engine, team2)
	if err != nil {
		log.Fatal(err)
	}

	team4 := getTeam(engine, team2.Name)
	fmt.Printf("the team4 is %v after update \n", team4)
	deleteTeam(engine, team2.Name)
}
