package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rsmaxwell/players-api/internal/debug"
)

var (
	functionCheckConistencyTx = debug.NewFunction(pkg, "CheckConistencyTx")
	functionCheckConistency   = debug.NewFunction(pkg, "CheckConistency")
)

func CheckConistencyTx(db *sql.DB, fix bool) (int, error) {
	f := functionCheckConistencyTx
	ctx := context.Background()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		message := "Could not begin a new transaction"
		f.DumpError(err, message)
		return 0, err
	}

	count, err := CheckConistency(ctx, db, fix)
	if err != nil {
		tx.Rollback()
		return 0, err
	}

	err = tx.Commit()
	if err != nil {
		message := "Could not commit the transaction"
		f.DumpError(err, message)
	}

	return count, nil
}

func CheckConistency(ctx context.Context, db *sql.DB, fix bool) (int, error) {
	f := functionCheckConistency

	list, err := ListPeople(ctx, db, "")
	if err != nil {
		message := "Could not list people"
		f.Errorf(message)
		f.DumpError(err, message)
		return 0, err
	}

	count := 0

	for _, person := range list {

		waiters, err := ListWaitersForPerson(ctx, db, person.ID)
		if err != nil {
			message := fmt.Sprintf("Could not list waiters for person: [%d: %s]", person.ID, person.Knownas)
			f.Errorf(message)
			f.DumpError(err, message)
			return 0, err
		}

		players, err := ListPlayersForPerson(ctx, db, person.ID)
		if err != nil {
			message := fmt.Sprintf("Could not list players for person: [%d: %s]", person.ID, person.Knownas)
			f.Errorf(message)
			f.DumpError(err, message)
			return 0, err
		}

		if person.Status == StatusPlayer {
			if len(waiters) < 1 {
				if len(players) < 1 {

					count++
					f.DebugError(fmt.Sprintf("Inconsistant data: person [%d: %s] is a player but has no waiter or player records", person.ID, person.Knownas))

					if fix {
						err := AddWaiter(ctx, db, person.ID)
						if err != nil {
							message := fmt.Sprintf("Could not add waiter: [%d: %s]", person.ID, person.Knownas)
							f.Errorf(message)
							f.DumpError(err, message)
							return 0, err
						}
					}
				} else if len(players) > 1 {

					count++
					f.DebugError(fmt.Sprintf("Inconsistant data: person [%d: %s] is a player and has %d player records", person.ID, person.Knownas, len(players)))

					if fix {
						err = RemovePlayer(ctx, db, person.ID)
						if err != nil {
							message := fmt.Sprintf("Could not remove player: id: [%d: %s]", person.ID, person.Knownas)
							f.Errorf(message)
							f.DumpError(err, message)
							return 0, err
						}

						err := AddWaiter(ctx, db, person.ID)
						if err != nil {
							message := fmt.Sprintf("Could not add waiter: id: [%d: %s]", person.ID, person.Knownas)
							f.Errorf(message)
							f.DumpError(err, message)
							return 0, err
						}
					}
				} else {
					// NOP
				}
			} else if len(waiters) > 1 {

				count++
				f.DebugError(fmt.Sprintf("Inconsistant data: person [%d: %s] is a player but has %d waiter records", person.ID, person.Knownas, len(waiters)))

				if fix {
					err = RemoveWaiter(ctx, db, person.ID)
					if err != nil {
						message := fmt.Sprintf("Could not remove waiter: [%d: %s]", person.ID, person.Knownas)
						f.Errorf(message)
						f.DumpError(err, message)
						return 0, err
					}

					err = RemovePlayer(ctx, db, person.ID)
					if err != nil {
						message := fmt.Sprintf("Could not remove player: [%d: %s]", person.ID, person.Knownas)
						f.Errorf(message)
						f.DumpError(err, message)
						return 0, err
					}

					err := AddWaiter(ctx, db, person.ID)
					if err != nil {
						message := fmt.Sprintf("Could not add waiter: [%d: %s]", person.ID, person.Knownas)
						f.Errorf(message)
						f.DumpError(err, message)
						return 0, err
					}
				}
			} else if len(players) < 1 {
				// NOP
			} else {

				count++
				f.DebugError(fmt.Sprintf("Inconsistant data: person [%d: %s] is a player but has 1 waiter record and %d player records", person.ID, person.Knownas, len(players)))

				if fix {
					err = RemovePlayer(ctx, db, person.ID)
					if err != nil {
						message := fmt.Sprintf("Could not remove player: [%d: %s]", person.ID, person.Knownas)
						f.Errorf(message)
						f.DumpError(err, message)
						return 0, err
					}
				}
			}
		} else {

			if len(waiters) > 0 {

				count++
				f.DebugError(fmt.Sprintf("Inconsistant data: person [%d: %s] is not a player but has %d waiter records", person.ID, person.Knownas, len(waiters)))

				if fix {
					err = RemoveWaiter(ctx, db, person.ID)
					if err != nil {
						message := fmt.Sprintf("Could not remove waiter: [%d: %s]", person.ID, person.Knownas)
						f.Errorf(message)
						f.DumpError(err, message)
						return 0, err
					}
				}
			}

			if len(players) > 0 {

				count++
				f.DebugError(fmt.Sprintf("Inconsistant data: person [%d: %s] is not a player but has %d player records", person.ID, person.Knownas, len(players)))

				if fix {
					err = RemovePlayer(ctx, db, person.ID)
					if err != nil {
						message := fmt.Sprintf("Could not remove player: [%d: %s]", person.ID, person.Knownas)
						f.Errorf(message)
						f.DumpError(err, message)
						return 0, err
					}
				}
			}
		}
	}

	return count, nil
}
