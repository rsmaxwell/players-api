package model

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rsmaxwell/players-api/internal/codeerror"
	"github.com/rsmaxwell/players-api/internal/debug"
	"golang.org/x/crypto/bcrypt"
)

var (
	functionUpdatePersonFieldsTx = debug.NewFunction(pkg, "UpdatePersonFieldsTx")
	functionUpdatePersonFields   = debug.NewFunction(pkg, "UpdatePersonFields")
)

// UpdatePerson method
func UpdatePersonFieldsTx(db *sql.DB, personID int, fields map[string]interface{}) error {
	f := functionUpdatePersonFieldsTx
	ctx := context.Background()

	// Begin a transaction
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		message := "Could not begin a new transaction"
		f.DumpError(err, message)
		return err
	}

	err = UpdatePersonFields(ctx, db, personID, fields)
	if err != nil {
		tx.Rollback()
		return err
	}

	err = tx.Commit()
	if err != nil {
		message := "Could not commit the transaction"
		f.DumpError(err, message)
	}

	return nil
}

func UpdatePersonFields(ctx context.Context, db *sql.DB, personID int, fields map[string]interface{}) error {
	f := functionUpdatePersonFields

	var p FullPerson
	p.ID = personID
	err := p.LoadPerson(ctx, db)
	if err != nil {
		message := fmt.Sprintf("could not load person: %d", personID)
		f.DebugVerbose(message)
		d := f.DumpError(err, message)
		d.AddObject("fields", fields)
		return codeerror.NewInternalServerError(message)
	}

	if val, ok := fields["firstname"]; ok {
		p.FirstName, ok = val.(string)
		if !ok {
			message := fmt.Sprintf("unexpected type for [%s]: %v", "firstName", val)
			f.DebugVerbose(message)
			f.DumpError(err, message)
			return codeerror.NewBadRequest(message)
		}
	}

	if val, ok := fields["lastname"]; ok {
		p.LastName, ok = val.(string)
		if !ok {
			message := fmt.Sprintf("unexpected type for [%s]: %v", "lastName", val)
			f.DebugVerbose(message)
			f.DumpError(err, message)
			return codeerror.NewBadRequest(message)
		}
	}

	if val, ok := fields["knownas"]; ok {
		p.Knownas, ok = val.(string)
		if !ok {
			message := fmt.Sprintf("unexpected type for [%s]: %v", "knownas", val)
			f.DebugVerbose(message)
			f.DumpError(err, message)
			return codeerror.NewBadRequest(message)
		}
	}

	if val, ok := fields["email"]; ok {
		p.Email, ok = val.(string)
		if !ok {
			message := fmt.Sprintf("unexpected type for [%s]: %v", "email", val)
			f.DebugVerbose(message)
			f.DumpError(err, message)
			return codeerror.NewBadRequest(message)
		}
	}

	if val, ok := fields["phone"]; ok {
		p.Phone, ok = val.(string)
		if !ok {
			message := fmt.Sprintf("unexpected type for [%s]: %v", "phone", val)
			f.DebugVerbose(message)
			f.DumpError(err, message)
			return codeerror.NewBadRequest(message)
		}
	}

	if val, ok := fields["password"]; ok {
		password, ok := val.(string)
		if !ok {
			message := fmt.Sprintf("unexpected type for [%s]: %v", "password", val)
			f.DebugVerbose(message)
			f.DumpError(err, message)
			return codeerror.NewBadRequest(message)
		}

		p.Hash, err = bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			message := "problem hashing password"
			f.DebugVerbose(message)
			f.DumpError(err, message)
			return codeerror.NewInternalServerError(message)
		}
	}

	if val, ok := fields["status"]; ok {
		p.Status, ok = val.(string)
		if !ok {
			message := fmt.Sprintf("unexpected type for [%s]: %v", "status", val)
			f.DebugVerbose(message)
			f.DumpError(err, message)
			return codeerror.NewInternalServerError(message)
		}
	}

	err = p.UpdatePerson(ctx, db)
	if err != nil {
		message := fmt.Sprintf("problem updating person: %d", personID)
		f.DebugVerbose(message)
		f.DumpError(err, message)
		return codeerror.NewInternalServerError(message)
	}

	return nil
}
