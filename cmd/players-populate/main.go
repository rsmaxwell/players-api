package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/rsmaxwell/players-api/internal/model"

	"github.com/rsmaxwell/players-api/internal/basic"
	"github.com/rsmaxwell/players-api/internal/config"
	"github.com/rsmaxwell/players-api/internal/debug"

	_ "github.com/jackc/pgx/stdlib"
)

// Person type
type PersonData struct {
	Data   model.Registration
	Status string
}

// Court type
type CourtData struct {
	Name string
}

var (
	pkg                = debug.NewPackage("main")
	functionMain       = debug.NewFunction(pkg, "main")
	functionMakePeople = debug.NewFunction(pkg, "makePeople")
	functionMakeCourts = debug.NewFunction(pkg, "makeCourts")
)

func init() {
	debug.InitDump("com.rsmaxwell.players", "players-createdb", "https://server.rsmaxwell.co.uk/archiva")
}

// http://go-database-sql.org/retrieving.html
func main() {
	f := functionMain
	ctx := context.Background()

	f.Infof("Players Populate: Version: %s", basic.Version())

	db, _, err := config.Setup()
	if err != nil {
		f.Errorf("Error setting up")
		os.Exit(1)
	}
	defer db.Close()

	_, err = makePeople(db)
	if err != nil {
		f.Errorf("Error making people")
		os.Exit(1)
	}

	/* courtIDs */
	_, err = makeCourts(ctx, db)
	if err != nil {
		f.Errorf("Error making courts")
		os.Exit(1)
	}

	count, err := model.CheckConistencyTx(db, true)
	if err != nil {
		f.Errorf("Error checking consistency")
		os.Exit(1)
	}

	fmt.Printf("Made %d database updates", count)
}

func makePeople(db *sql.DB) (map[int]int, error) {
	f := functionMakePeople

	peopleData := []PersonData{
		{Data: model.Registration{FirstName: "James", LastName: "Bond", Knownas: "007", Email: "007@mi6.gov.uk", Phone: "01632 960573", Password: "TopSecret123"}, Status: model.StatusPlayer},
		{Data: model.Registration{FirstName: "Alice", LastName: "Frombe", Knownas: "ali", Email: "ali@mikymouse.com", Phone: "01632 960372", Password: "ali1234567"}, Status: model.StatusPlayer},
		{Data: model.Registration{FirstName: "Tom", LastName: "Smith", Knownas: "tom", Email: "tom@hotmail.com", Phone: "01632 960512", Password: "tom12378909876"}, Status: model.StatusPlayer},
		{Data: model.Registration{FirstName: "Sandra", LastName: "Smythe", Knownas: "tom", Email: "sandra@hotmail.com", Phone: "01632 960966", Password: "sandra12334567"}, Status: model.StatusInactive},
		{Data: model.Registration{FirstName: "George", LastName: "Washington", Knownas: "george", Email: "george@hotmail.com", Phone: "01632 960278", Password: "george789"}, Status: model.StatusPlayer},
		{Data: model.Registration{FirstName: "Margret", LastName: "Tiffington", Knownas: "maggie", Email: "marg@hotmail.com", Phone: "01632 960165", Password: "magie876"}, Status: model.StatusPlayer},
		{Data: model.Registration{FirstName: "James", LastName: "Ernest", Knownas: "jamie", Email: "jamie@ntlworld.com", Phone: "01632 960757", Password: "jamie5293645284"}, Status: model.StatusInactive},
		{Data: model.Registration{FirstName: "Elizabeth", LastName: "Tudor", Knownas: "liz", Email: "liz@buck.palice.com", Phone: "01632 960252", Password: "liz1756453423"}, Status: model.StatusPlayer},
		{Data: model.Registration{FirstName: "Dick", LastName: "Whittington", Knownas: "dick", Email: "dick@ntlworld.com", Phone: "01746 352413", Password: "dick3296846734524"}, Status: model.StatusPlayer},
		{Data: model.Registration{FirstName: "Victoria", LastName: "Hempworth", Knownas: "vickie", Email: "vickie@waitrose.com", Phone: "0195 76863241", Password: "vickie846"}, Status: model.StatusPlayer},

		{Data: model.Registration{FirstName: "Shanika", LastName: "Pierre", Knownas: "pete", Email: "IcyGamer@gmail.com", Phone: "01632 960576", Password: "Top12345Secret"}, Status: model.StatusPlayer},
		{Data: model.Registration{FirstName: "Wanangwa", LastName: "Czajkowski", Knownas: "wan", Email: "torphy.dayana@dicki.com", Phone: "01632 960628", Password: "ali12387654"}, Status: model.StatusPlayer},
		{Data: model.Registration{FirstName: "Cormac", LastName: "Dwight", Knownas: "cor", Email: "adela.kunze@schmitt.com", Phone: "01632 960026", Password: "tom123frgthyj"}, Status: model.StatusSuspended},
		{Data: model.Registration{FirstName: "Ramóna", LastName: "Jonker", Knownas: "ram", Email: "ariel07@hotmail.com", Phone: "01632 960801", Password: "sandra123frr"}, Status: model.StatusSuspended},
		{Data: model.Registration{FirstName: "Quinctilius", LastName: "Jack", Knownas: "qui", Email: "kara.johnston@runte.com", Phone: "01632 960334", Password: "george789ed5"}, Status: model.StatusInactive},
		{Data: model.Registration{FirstName: "Radu", LastName: "Godfrey", Knownas: "rad", Email: "ella.vonrueden@kuhic.com", Phone: "01632 960450", Password: "magie87689ilom"}, Status: model.StatusSuspended},
		{Data: model.Registration{FirstName: "Aleksandrina", LastName: "Couture", Knownas: "ale", Email: "archibald.stark@hotmail.com", Phone: "01632 960928", Password: "jamie529re5gb"}, Status: model.StatusPlayer},
		{Data: model.Registration{FirstName: "Catrin", LastName: "Wooldridge", Knownas: "cat", Email: "sauer.luciano@hotmail.com", Phone: "01632 960126", Password: "liz14rdgujmbvr43"}, Status: model.StatusSuspended},
		{Data: model.Registration{FirstName: "Souleymane", LastName: "Walter", Knownas: "sou", Email: "damon.toy@swaniawski.com", Phone: "01632 960403", Password: "dick3287uyh5fredw"}, Status: model.StatusPlayer},
		{Data: model.Registration{FirstName: "Dorotėja", LastName: "Antúnez", Knownas: "dor", Email: "omante@marks.com", Phone: "01632 961252", Password: "vickie846y6"}, Status: model.StatusPlayer},
	}

	peopleIDs := make(map[int]int)
	for i, r := range peopleData {

		p, err := r.Data.ToPerson()
		if err != nil {
			message := "Could not register person"
			f.Errorf(message)
			f.DumpError(err, message)
			os.Exit(1)
		}

		p.Status = r.Status

		err = p.SavePersonTx(db)
		if err != nil {
			message := fmt.Sprintf("Could not save person: firstName: %s, lastname: %s, email: %s", p.FirstName, p.LastName, p.Email)
			f.Errorf(message)
			f.DumpError(err, message)
			os.Exit(1)
		}

		peopleIDs[i] = p.ID

		fmt.Printf("Added person:\n")
		fmt.Printf("    ID:        %d\n", p.ID)
		fmt.Printf("    FirstName: %s\n", p.FirstName)
		fmt.Printf("    LastName:  %s\n", p.LastName)
		fmt.Printf("    Knownas:   %s\n", p.Knownas)
		fmt.Printf("    Email:     %s\n", p.Email)
		fmt.Printf("    Password:  %s\n", r.Data.Password)
		fmt.Printf("    Hash:      %s\n", p.Hash)
		fmt.Printf("    Status:    %s\n", p.Status)
	}

	return peopleIDs, nil
}

func makeCourts(ctx context.Context, db *sql.DB) (map[int]int, error) {
	f := functionMakeCourts

	courtsData := []CourtData{
		{Name: "A"},
		{Name: "B"},
	}

	courtIDs := make(map[int]int)
	for i, c := range courtsData {

		court := model.Court{Name: c.Name}

		err := court.SaveCourt(ctx, db)
		if err != nil {
			message := fmt.Sprintf("Could not save court: Name: %s", court.Name)
			f.Errorf(message)
			f.DumpError(err, message)
			os.Exit(1)
		}

		courtIDs[i] = court.ID

		fmt.Printf("Added court:\n")
		fmt.Printf("    ID:    %d\n", court.ID)
		fmt.Printf("    Name:  %s\n", court.Name)
	}

	return courtIDs, nil
}
