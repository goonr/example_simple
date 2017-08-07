// Package models includes the functions on the model Comment.
package models

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
)

// set flags to output more detailed log
func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

type Comment struct {
	Id        int64     `json:"id,omitempty" db:"id" valid:"-"`
	Commenter string    `json:"commenter,omitempty" db:"commenter" valid:"required"`
	Body      string    `json:"body,omitempty" db:"body" valid:"required,length(20|4294967295)"`
	ArticleId int64     `json:"article_id,omitempty" db:"article_id" valid:"-"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at" valid:"-"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at" valid:"-"`
	Article   Article   `json:"article,omitempty" db:"article" valid:"-"`
}

// FindComment find a single comment by an id
func FindComment(id int64) (*Comment, error) {
	if id == 0 {
		return nil, errors.New("Invalid id: it can't be zero")
	}
	_comment := Comment{}
	err := db.Get(&_comment, db.Rebind(`SELECT * FROM comments WHERE id = ? LIMIT 1`), id)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_comment, nil
}

// FirstComment find the first one comment by id ASC order
func FirstComment() (*Comment, error) {
	_comment := Comment{}
	err := db.Get(&_comment, db.Rebind(`SELECT * FROM comments ORDER BY id ASC LIMIT 1`))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_comment, nil
}

// FirstComments find the first N comments by id ASC order
func FirstComments(n uint32) ([]Comment, error) {
	_comments := []Comment{}
	sql := fmt.Sprintf("SELECT * FROM comments ORDER BY id ASC LIMIT %v", n)
	err := db.Select(&_comments, db.Rebind(sql))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _comments, nil
}

// LastComment find the last one comment by id DESC order
func LastComment() (*Comment, error) {
	_comment := Comment{}
	err := db.Get(&_comment, db.Rebind(`SELECT * FROM comments ORDER BY id DESC LIMIT 1`))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_comment, nil
}

// LastComments find the last N comments by id DESC order
func LastComments(n uint32) ([]Comment, error) {
	_comments := []Comment{}
	sql := fmt.Sprintf("SELECT * FROM comments ORDER BY id DESC LIMIT %v", n)
	err := db.Select(&_comments, db.Rebind(sql))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _comments, nil
}

// FindComments find one or more comments by one or more ids
func FindComments(ids ...int64) ([]Comment, error) {
	if len(ids) == 0 {
		msg := "At least one or more ids needed"
		log.Println(msg)
		return nil, errors.New(msg)
	}
	_comments := []Comment{}
	idsHolder := strings.Repeat(",?", len(ids)-1)
	sql := db.Rebind(fmt.Sprintf(`SELECT * FROM comments WHERE id IN (?%s)`, idsHolder))
	idsT := []interface{}{}
	for _, id := range ids {
		idsT = append(idsT, interface{}(id))
	}
	err := db.Select(&_comments, sql, idsT...)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _comments, nil
}

// FindCommentBy find a single comment by a field name and a value
func FindCommentBy(field string, val interface{}) (*Comment, error) {
	_comment := Comment{}
	sqlFmt := `SELECT * FROM comments WHERE %s = ? LIMIT 1`
	sqlStr := fmt.Sprintf(sqlFmt, field)
	err := db.Get(&_comment, db.Rebind(sqlStr), val)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_comment, nil
}

// FindCommentsBy find all comments by a field name and a value
func FindCommentsBy(field string, val interface{}) (_comments []Comment, err error) {
	sqlFmt := `SELECT * FROM comments WHERE %s = ?`
	sqlStr := fmt.Sprintf(sqlFmt, field)
	err = db.Select(&_comments, db.Rebind(sqlStr), val)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _comments, nil
}

// AllComments get all the Comment records
func AllComments() (comments []Comment, err error) {
	err = db.Select(&comments, "SELECT * FROM comments")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return comments, nil
}

// CommentCount get the count of all the Comment records
func CommentCount() (c int64, err error) {
	err = db.Get(&c, "SELECT count(*) FROM comments")
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return c, nil
}

// CommentCountWhere get the count of all the Comment records with a where clause
func CommentCountWhere(where string, args ...interface{}) (c int64, err error) {
	sql := "SELECT count(*) FROM comments"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := db.Preparex(db.Rebind(sql))
	if err != nil {
		log.Println(err)
		return 0, err
	}
	err = stmt.Get(&c, args...)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return c, nil
}

// CommentIncludesWhere get the Comment associated models records, it's just the eager_load function
func CommentIncludesWhere(assocs []string, sql string, args ...interface{}) (_comments []Comment, err error) {
	_comments, err = FindCommentsWhere(sql, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(assocs) == 0 {
		log.Println("No associated fields ard specified")
		return _comments, err
	}
	if len(_comments) <= 0 {
		return nil, errors.New("No results available")
	}
	ids := make([]interface{}, len(_comments))
	for _, v := range _comments {
		ids = append(ids, interface{}(v.Id))
	}
	return _comments, nil
}

// CommentIds get all the Ids of Comment records
func CommentIds() (ids []int64, err error) {
	err = db.Select(&ids, "SELECT id FROM comments")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ids, nil
}

// CommentIdsWhere get all the Ids of Comment records by where restriction
func CommentIdsWhere(where string, args ...interface{}) ([]int64, error) {
	ids, err := CommentIntCol("id", where, args...)
	return ids, err
}

// CommentIntCol get some int64 typed column of Comment by where restriction
func CommentIntCol(col, where string, args ...interface{}) (intColRecs []int64, err error) {
	sql := "SELECT " + col + " FROM comments"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := db.Preparex(db.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&intColRecs, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return intColRecs, nil
}

// CommentStrCol get some string typed column of Comment by where restriction
func CommentStrCol(col, where string, args ...interface{}) (strColRecs []string, err error) {
	sql := "SELECT " + col + " FROM comments"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := db.Preparex(db.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&strColRecs, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return strColRecs, nil
}

// FindCommentsWhere query use a partial SQL clause that usually following after WHERE
// with placeholders, eg: FindUsersWhere("first_name = ? AND age > ?", "John", 18)
// will return those records in the table "users" whose first_name is "John" and age elder than 18
func FindCommentsWhere(where string, args ...interface{}) (comments []Comment, err error) {
	sql := "SELECT * FROM comments"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := db.Preparex(db.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&comments, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return comments, nil
}

// FindCommentBySql query use a complete SQL clause
// with placeholders, eg: FindUserBySql("SELECT * FROM users WHERE first_name = ? AND age > ? ORDER BY DESC LIMIT 1", "John", 18)
// will return only One record in the table "users" whose first_name is "John" and age elder than 18
func FindCommentBySql(sql string, args ...interface{}) (*Comment, error) {
	stmt, err := db.Preparex(db.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	_comment := &Comment{}
	err = stmt.Get(_comment, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return _comment, nil
}

// FindCommentsBySql query use a complete SQL clause
// with placeholders, eg: FindUsersBySql("SELECT * FROM users WHERE first_name = ? AND age > ?", "John", 18)
// will return those records in the table "users" whose first_name is "John" and age elder than 18
func FindCommentsBySql(sql string, args ...interface{}) (comments []Comment, err error) {
	stmt, err := db.Preparex(db.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&comments, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return comments, nil
}

// CreateComment use a named params to create a single Comment record.
// A named params is key-value map like map[string]interface{}{"first_name": "John", "age": 23} .
func CreateComment(am map[string]interface{}) (int64, error) {
	if len(am) == 0 {
		return 0, fmt.Errorf("Zero key in the attributes map!")
	}
	t := time.Now()
	for _, v := range []string{"created_at", "updated_at"} {
		if am[v] == nil {
			am[v] = t
		}
	}
	keys := make([]string, len(am))
	i := 0
	for k := range am {
		keys[i] = k
		i++
	}
	sqlFmt := `INSERT INTO comments (%s) VALUES (%s)`
	sqlStr := fmt.Sprintf(sqlFmt, strings.Join(keys, ","), ":"+strings.Join(keys, ",:"))
	result, err := db.NamedExec(sqlStr, am)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return lastId, nil
}

// Create is a method for Comment to create a record
func (_comment *Comment) Create() (int64, error) {
	ok, err := govalidator.ValidateStruct(_comment)
	if !ok {
		errMsg := "Validate Comment struct error: Unknown error"
		if err != nil {
			errMsg = "Validate Comment struct error: " + err.Error()
		}
		log.Println(errMsg)
		return 0, errors.New(errMsg)
	}
	t := time.Now()
	_comment.CreatedAt = t
	_comment.UpdatedAt = t
	sql := `INSERT INTO comments (commenter,body,article_id,created_at,updated_at) VALUES (:commenter,:body,:article_id,:created_at,:updated_at)`
	result, err := db.NamedExec(sql, _comment)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	lastId, err := result.LastInsertId()
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return lastId, nil
}

// CreateArticle is a method for a Comment object to create an associated Article record
func (_comment *Comment) CreateArticle(am map[string]interface{}) error {
	am["comment_id"] = _comment.Id
	_, err := CreateArticle(am)
	return err
}

// Destroy is method used for a Comment object to be destroyed.
func (_comment *Comment) Destroy() error {
	if _comment.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := DestroyComment(_comment.Id)
	return err
}

// DestroyComment will destroy a Comment record specified by id parameter.
func DestroyComment(id int64) error {
	stmt, err := db.Preparex(db.Rebind(`DELETE FROM comments WHERE id = ?`))
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

// DestroyComments will destroy Comment records those specified by the ids parameters.
func DestroyComments(ids ...int64) (int64, error) {
	if len(ids) == 0 {
		msg := "At least one or more ids needed"
		log.Println(msg)
		return 0, errors.New(msg)
	}
	idsHolder := strings.Repeat(",?", len(ids)-1)
	sql := fmt.Sprintf(`DELETE FROM comments WHERE id IN (?%s)`, idsHolder)
	idsT := []interface{}{}
	for _, id := range ids {
		idsT = append(idsT, interface{}(id))
	}
	stmt, err := db.Preparex(db.Rebind(sql))
	result, err := stmt.Exec(idsT...)
	if err != nil {
		return 0, err
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// DestroyCommentsWhere delete records by a where clause
// like: DestroyCommentsWhere("name = ?", "John")
// And this func will not call the association dependent action
func DestroyCommentsWhere(where string, args ...interface{}) (int64, error) {
	sql := `DELETE FROM comments WHERE `
	if len(where) > 0 {
		sql = sql + where
	} else {
		return 0, errors.New("No WHERE conditions provided")
	}
	stmt, err := db.Preparex(db.Rebind(sql))
	result, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return cnt, nil
}

// Save method is used for a Comment object to update an existed record mainly.
// If no id provided a new record will be created. A UPSERT action will be implemented further.
func (_comment *Comment) Save() error {
	ok, err := govalidator.ValidateStruct(_comment)
	if !ok {
		errMsg := "Validate Comment struct error: Unknown error"
		if err != nil {
			errMsg = "Validate Comment struct error: " + err.Error()
		}
		log.Println(errMsg)
		return errors.New(errMsg)
	}
	if _comment.Id == 0 {
		_, err = _comment.Create()
		return err
	}
	_comment.UpdatedAt = time.Now()
	sqlFmt := `UPDATE comments SET %s WHERE id = %v`
	sqlStr := fmt.Sprintf(sqlFmt, "commenter = :commenter, body = :body, article_id = :article_id, updated_at = :updated_at", _comment.Id)
	_, err = db.NamedExec(sqlStr, _comment)
	return err
}

// UpdateComment is used to update a record with a id and map[string]interface{} typed key-value parameters.
func UpdateComment(id int64, am map[string]interface{}) error {
	if len(am) == 0 {
		return errors.New("Zero key in the attributes map!")
	}
	am["updated_at"] = time.Now()
	keys := make([]string, len(am))
	i := 0
	for k := range am {
		keys[i] = k
		i++
	}
	sqlFmt := `UPDATE comments SET %s WHERE id = %v`
	setKeysArr := []string{}
	for _, v := range keys {
		s := fmt.Sprintf(" %s = :%s", v, v)
		setKeysArr = append(setKeysArr, s)
	}
	sqlStr := fmt.Sprintf(sqlFmt, strings.Join(setKeysArr, ", "), id)
	_, err := db.NamedExec(sqlStr, am)
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

// Update is a method used to update a Comment record with the map[string]interface{} typed key-value parameters.
func (_comment *Comment) Update(am map[string]interface{}) error {
	if _comment.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateComment(_comment.Id, am)
	return err
}

func (_comment *Comment) UpdateAttributes(am map[string]interface{}) error {
	if _comment.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateComment(_comment.Id, am)
	return err
}

// UpdateColumns method is supposed to be used to update Comment records as corresponding update_columns in Rails
func (_comment *Comment) UpdateColumns(am map[string]interface{}) error {
	if _comment.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateComment(_comment.Id, am)
	return err
}

// UpdateCommentsBySql is used to update Comment records by a SQL clause
// that use '?' binding syntax
func UpdateCommentsBySql(sql string, args ...interface{}) (int64, error) {
	if sql == "" {
		return 0, errors.New("A blank SQL clause")
	}
	sql = strings.Replace(strings.ToLower(sql), "set", "set updated_at = ?, ", 1)
	args = append([]interface{}{time.Now()}, args...)
	stmt, err := db.Preparex(db.Rebind(sql))
	result, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	cnt, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return cnt, nil
}
