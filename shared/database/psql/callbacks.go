package psql

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/esmailemami/chess/shared/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func loadCallbacks(dbConn *gorm.DB) {
	_ = dbConn.Callback().
		Create().
		Before("gorm:create").
		Register("uuid_setter", func(scope *gorm.DB) {
			modelRef := scope.Statement.ReflectValue

			switch modelRef.Kind() {
			case reflect.Struct:
				{
					idVal := modelRef.FieldByName("ID")
					id := fmt.Sprintf("%v", idVal)

					if isNiIDString(id) && idVal.Type() == reflect.TypeOf(uuid.UUID{}) {
						uid := uuid.New()
						idVal.Set(reflect.ValueOf(uid))
					}
				}
			case reflect.Slice:
				{
					// handle setting of ids from db layer
				}
			}
		})

	_ = dbConn.Callback().
		Create().
		Before("gorm:create").
		Register("before_create_logger", beforeCallback(beforeCreateCallback))
	_ = dbConn.Callback().
		Update().
		Before("gorm:update").
		Register("before_update_logger", beforeCallback(beforeUpdateCallback))
	_ = dbConn.Callback().
		Delete().
		Before("gorm:delete").
		Register("before_delete_logger", beforeCallback(beforeDeleteCallback))
}

func isNiIDString(id string) bool {
	return id == "" || id == "<nil>" || id == "0" || uuid.Nil.String() == id

	// id = strings.ReplaceAll(id, "-", "")
	// id = strings.ReplaceAll(id, "0", "")
	// return id == ""
}

func extractIDFromWhereClause(w clause.Where) (ids []string) {
	for _, c := range w.Exprs {
		switch t := c.(type) {
		case clause.Eq:
			{
				colName := ""
				switch col := t.Column.(type) {
				case string:
					colName = col
				case clause.Column:
					colName = col.Name
				}

				if colName == "id" || colName == clause.PrimaryKey {
					ids = append(ids, fmt.Sprintf("%v", t.Value))
				}
			}
		case clause.IN:
			{
				colName := ""
				switch col := t.Column.(type) {
				case string:
					colName = col
				case clause.Column:
					colName = col.Name
				}

				if colName != "id" && colName != clause.PrimaryKey {
					continue
				}
				for _, i := range t.Values {
					ids = append(ids, fmt.Sprintf("%v", i))
				}
			}
		case clause.Expr:
			{

			}
		default:
			{
				return nil
			}
		}
	}

	return ids
}

func beforeCallback(callback func(scope *gorm.DB, v reflect.Value, user *models.User) error) func(*gorm.DB) {
	return func(scope *gorm.DB) {
		if scope.Statement.Error != nil {
			return
		}

		modelRef := scope.Statement.ReflectValue

		var structsToAudit []reflect.Value

		switch modelRef.Kind() {
		case reflect.Struct:
			{
				structsToAudit = append(structsToAudit, modelRef)
			}
		case reflect.Slice:
			{
				for i := 0; i < modelRef.Len(); i++ {
					v := modelRef.Index(i)
					if v.Kind() == reflect.Ptr {
						v = v.Elem()
					}

					if v.Kind() == reflect.Struct {
						structsToAudit = append(structsToAudit, v)
					}
				}
			}
		}

		user, ok := scope.Statement.Context.Value("user").(*models.User)
		if !ok {
			return
		}

		for _, st := range structsToAudit {
			_ = callback(scope, st, user)
		}
	}
}

func beforeCreateCallback(scope *gorm.DB, st reflect.Value, user *models.User) error {
	v := st.FieldByName("CreatedByID")
	if !v.IsValid() {
		return nil
	}

	if !v.CanSet() {
		err := errors.New("CreatedByID can't be set")
		_ = scope.Statement.AddError(err)
		return err
	}

	v.Set(reflect.ValueOf(&user.ID))
	return nil
}

func beforeUpdateCallback(scope *gorm.DB, st reflect.Value, user *models.User) error {
	v := st.FieldByName("UpdatedByID")
	if !v.IsValid() {
		return nil
	}
	if !v.CanSet() {
		err := errors.New("UpdatedByID can't be set")
		_ = scope.Statement.AddError(err)
		return err
	}

	v.Set(reflect.ValueOf(&user.ID))
	return nil
}

func beforeDeleteCallback(scope *gorm.DB, st reflect.Value, user *models.User) error {
	if scope.Statement.Unscoped {
		return nil
	}

	var id uuid.UUID
	curTime := scope.Statement.DB.NowFunc()
	// handle soft delete deleted_by_id
	//region soft delete
	modelVal := reflect.ValueOf(scope.Statement.Model)

	if modelVal.Kind() == reflect.Ptr {

		idField := reflect.Indirect(modelVal).FieldByName("ID")
		if idField.IsValid() {
			idPtr := idField.Interface().(*uuid.UUID)
			if idPtr == nil {
				id = uuid.Nil
			}
		}

		val := reflect.Indirect(modelVal).FieldByName("DeletedByID")
		if !val.IsValid() {
			return nil
		}
		if !val.CanSet() {
			err := errors.New("DeletedByID can't be set")
			_ = scope.Statement.AddError(err)
			return err
		}
	} else {
		idField := reflect.Indirect(modelVal).FieldByName("ID")
		if idField.IsValid() {
			idPtr := idField.Interface().(*uuid.UUID)
			if idPtr == nil {
				id = uuid.Nil
			}
		}

		val := modelVal.FieldByName("DeletedByID")
		if !val.IsValid() {
			return nil
		}
		if !val.CanSet() {
			err := errors.New("DeletedByID can't be set")
			_ = scope.Statement.AddError(err)
			return err
		}
	}

	// create update statement
	scope.Statement.AddClause(clause.Update{})
	scope.Statement.AddClause(clause.Set{
		{Column: clause.Column{Name: "deleted_at"}, Value: curTime},
		{Column: clause.Column{Name: "deleted_by_id"}, Value: &user.ID},
	})

	scope.Statement.SetColumn("deleted_at", curTime)
	scope.Statement.SetColumn("deleted_by_id", &user.ID)

	if id.String() != "00000000-0000-0000-0000-000000000000" {
		scope.Statement.AddClause(clause.Where{Exprs: []clause.Expression{
			clause.Eq{Column: clause.PrimaryColumn, Value: id},
			clause.Eq{Column: "deleted_at", Value: nil},
		}})
	}
	scope.Statement.Build(
		clause.Update{}.Name(),
		clause.Set{}.Name(),
		clause.Where{}.Name(),
	)

	//endregion
	return nil
}
