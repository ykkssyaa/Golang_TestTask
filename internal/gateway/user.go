package gateway

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"testTask/internal/model"
	logger2 "testTask/pkg/logger"
)

type PostgresUserGateway interface {
	GetUsers(limit, offset int, filter model.UserFilter) ([]model.User, error)
	CreateUser(user model.User) (int, error)
	DeleteUser(id int) error
	UpdateUser(user model.User) error
}

type PostgresUserGatewayImpl struct {
	db     *sqlx.DB
	logger *logger2.Logger
}

func (u PostgresUserGatewayImpl) GetUsers(limit, offset int, filter model.UserFilter) ([]model.User, error) {

	u.logger.Info.Printf("Getting users from postgres. Limit: %d, Offset: %d", limit, offset)
	u.logger.Info.Println("Filter with params: ", filter)

	getUsersQuery, args := SelectQueryBuilder(limit, offset, filter)

	u.logger.Info.Println("Query: ", getUsersQuery)

	var users []model.User

	err := u.db.Select(&users, getUsersQuery, args...)
	if err != nil {
		return nil, err
	}
	return users, err

}

func SelectQueryBuilder(limit, offset int, filter model.UserFilter) (query string, args []interface{}) {

	getUsersQuery := fmt.Sprintf("SELECT * FROM %s WHERE 1=1", usersTable)

	i := 1

	if filter.Name != "" {
		getUsersQuery += fmt.Sprintf(" AND name=$%d ", i)
		i++
		args = append(args, filter.Name)
	}
	if filter.Surname != "" {
		getUsersQuery += fmt.Sprintf(" AND surname=$%d ", i)
		i++
		args = append(args, filter.Surname)
	}
	if filter.Patronymic != "" {

		getUsersQuery += fmt.Sprintf(" AND patronymic=$%d ", i)
		i++
		args = append(args, filter.Patronymic)
	}
	if filter.Age != 0 {

		getUsersQuery += fmt.Sprintf("AND age=$%d ", i)
		i++
		args = append(args, filter.Age)
	}
	if filter.Gender != "" {
		getUsersQuery += fmt.Sprintf("AND gender=$%d ", i)
		i++
		args = append(args, filter.Gender)
	}
	if filter.Country != "" {
		getUsersQuery += fmt.Sprintf("AND country=$%d ", i)
		i++
		args = append(args, filter.Country)
	}

	getUsersQuery += fmt.Sprintf(" LIMIT $%d OFFSET $%d", i, i+1)
	args = append(args, limit, offset)

	return getUsersQuery, args

}

func (u PostgresUserGatewayImpl) CreateUser(user model.User) (int, error) {

	u.logger.Info.Println("Saving user in postgres.")

	var id int

	createQuery := fmt.Sprintf("INSERT INTO %s "+
		"(name, surname, patronymic, age, country, gender) "+
		"VALUES ($1, $2, $3, $4, $5, $6)"+
		"RETURNING id", usersTable)

	u.logger.Info.Println("Query: ", createQuery)

	tx, err := u.db.DB.Begin()

	if err != nil {
		return 0, err
	}
	row := tx.QueryRow(createQuery, user.Name, user.Surname, user.Patronymic, user.Age, user.Country, user.Gender)
	if err := row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}

	if err = tx.Commit(); err != nil {
		return 0, err
	}

	return id, nil
}

func (u PostgresUserGatewayImpl) DeleteUser(id int) error {

	// TODO: check if exists

	u.logger.Info.Println("Deleting user from postgres.")

	deleteQuery := fmt.Sprintf("DELETE FROM %s WHERE id = $1", usersTable)
	u.logger.Info.Println("Query: ", deleteQuery)

	tx, err := u.db.DB.Begin()

	if err != nil {
		return err
	}
	if _, err = tx.Exec(deleteQuery, id); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil
}

func (u PostgresUserGatewayImpl) UpdateUser(user model.User) error {

	u.logger.Info.Println("Deleting user from postgres.")

	UpdateQuery, args := UpdateQueryBuilder(user)
	u.logger.Info.Println("Query: ", UpdateQuery)

	tx, err := u.db.DB.Begin()

	if err != nil {
		return err
	}
	if _, err = tx.Exec(UpdateQuery, args...); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}

	return nil

}

func UpdateQueryBuilder(user model.User) (query string, args []interface{}) {

	query = fmt.Sprintf("UPDATE %s SET ", usersTable)

	i := 1
	if user.Name != "" {
		query += fmt.Sprintf(" name = $%d", i)
		args = append(args, user.Name)
		i++
	}
	if user.Surname != "" {
		if i > 1 {
			query += ","
		}
		query += fmt.Sprintf(" surname = $%d", i)
		args = append(args, user.Surname)
		i++
	}
	if user.Patronymic != "" {
		if i > 1 {
			query += ","
		}
		query += fmt.Sprintf(" patronymic = $%d", i)
		args = append(args, user.Patronymic)
		i++
	}
	if user.Country != "" {
		if i > 1 {
			query += ","
		}
		query += fmt.Sprintf(" country = $%d", i)
		args = append(args, user.Country)
		i++
	}
	if user.Gender != "" {
		if i > 1 {
			query += ","
		}
		query += fmt.Sprintf(" gender = $%d", i)
		args = append(args, user.Gender)
		i++
	}
	if user.Age != 0 {
		if i > 1 {
			query += ","
		}
		query += fmt.Sprintf(" age = $%d", i)
		args = append(args, user.Age)
		i++
	}

	query += fmt.Sprintf(" WHERE id = $%d", i)
	args = append(args, user.Id)

	return
}
