package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/profile"
	"xorm.io/xorm"
	"xorm.io/xorm/names"
)

const numrecords = 2000

type Team struct {
	ID        int       `xorm:"'id' pk autoincr"`
	Name      string    `xorm:"name"`
	OrgID     int       `xorm:"org_id"`
	CreatedAt time.Time `xorm:"'created'"`
	UpdatedAt time.Time `xorm:"'updated'"`
	Email     string
}

func insertTeam(e *xorm.Engine, team1 Team) error {
	// for insert, xorm is actually checking that all the filled field has a corresponding column name, if not, error out
	// not set field would be filled directly with default value when creation
	_, err := e.Insert(&team1)
	return err
}

func getTeam(e *xorm.Engine, name string) (bool, Team, error) {
	var team1 Team
	found, err := e.Where("name=?", name).Get(&team1)
	return found, team1, err
}

func deleteTeam(e *xorm.Engine, name string) error {
	_, err := e.Exec("DELETE FROM team WHERE name=? Returning id", name)
	return err
}

func updateTeam(e *xorm.Engine, team Team) error {
	_, err := e.ID(team.ID).Update(team)
	return err
}

func main() {
	// set engine of sqlite3 here
	engine, err := xorm.NewEngine("sqlite3", "grafana.db")
	if err != nil {
		log.Fatal(err)
	}
	engine.SetTableMapper(names.GonicMapper{})

	// The test senarios:
	// 1. Create 1000 teams with go struct
	// 2. Verify after each creation so 1000 get with go struct casting
	// 3. Update 1000 teams with only name field
	// 4. Verify after update so 1000 get with go struct casting
	// 5. Delete 1000 teams with name field, return the id deleted
	// 6. Get 1000 teams which is nolonger existing

	// start timer
	start := time.Now()

	// CPU profiling by default related doc here: https://flaviocopes.com/golang-profiling/
	// defer profile.Start().Stop()

	// Memory profiling
	defer profile.Start(profile.MemProfile).Stop()

	for i := 0; i < numrecords; i++ {
		num := strconv.Itoa(i)
		team1 := Team{Name: "mynamee" + num, OrgID: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}
		err = insertTeam(engine, team1)
		// fmt.Printf("the inserted team %+v\n", team1)
		if err != nil {
			log.Fatal(err)
		}

		_, team3, err := getTeam(engine, team1.Name)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("the get team %+v\n", team3)

		// here it is very confusing, since xorm omit the default value OrgID 0 here, we would have to
		// force the update by calling .AllCols().Update or .Cols("org_id").Update
		// but if we put .AllCols() we have to put all the fields that are mandatory so it is also not convinient
		team2 := Team{ID: team3.ID, OrgID: 0, Name: "princess" + num}
		err = updateTeam(engine, team2)
		if err != nil {
			log.Fatal(err)
		}

		_, team4, err := getTeam(engine, team2.Name)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("the updated team %+v\n", team4)

		err = deleteTeam(engine, team4.Name)
		if err != nil {
			log.Fatal(err)
		}

		_, _, err = getTeam(engine, team4.Name)
		if err != nil {
			log.Fatal(err)
		}
		// fmt.Printf("the deleted team %+v\n", f)
	}

	fmt.Println("the result is: ", time.Since(start))
}
