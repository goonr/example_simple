// Package models is generated by go-on-rails, a Rails generator gem:
// https://rubygems.org/gems/go-on-rails
// Or on Github: https://github.com/goonr/go-on-rails
// By B1nj0y <idegorepl@gmail.com>
package models

import (
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
)

type Article struct {
	Id        int64     `json:"id,omitempty" db:"id" valid:"-"`
	Title     string    `json:"title,omitempty" db:"title" valid:"required,length(10|30)"`
	Text      string    `json:"text,omitempty" db:"text" valid:"required,length(20|4294967295)"`
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at" valid:"-"`
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at" valid:"-"`
	Comments  []Comment `json:"comments,omitempty" db:"comments" valid:"-"`
}

// FindArticle find a single article by an id
func FindArticle(id int64) (*Article, error) {
	if id == 0 {
		return nil, errors.New("Invalid id: it can't be zero")
	}
	_article := Article{}
	err := db.Get(&_article, db.Rebind(`SELECT * FROM articles WHERE id = ? LIMIT 1`), id)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_article, nil
}

// FirstArticle find the first one article by id ASC order
func FirstArticle() (*Article, error) {
	_article := Article{}
	err := db.Get(&_article, db.Rebind(`SELECT * FROM articles ORDER BY id ASC LIMIT 1`))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_article, nil
}

// FirstArticles find the first N articles by id ASC order
func FirstArticles(n uint32) ([]Article, error) {
	_articles := []Article{}
	sql := fmt.Sprintf("SELECT * FROM articles ORDER BY id ASC LIMIT %v", n)
	err := db.Select(&_articles, db.Rebind(sql))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _articles, nil
}

// LastArticle find the last one article by id DESC order
func LastArticle() (*Article, error) {
	_article := Article{}
	err := db.Get(&_article, db.Rebind(`SELECT * FROM articles ORDER BY id DESC LIMIT 1`))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_article, nil
}

// LastArticles find the last N articles by id DESC order
func LastArticles(n uint32) ([]Article, error) {
	_articles := []Article{}
	sql := fmt.Sprintf("SELECT * FROM articles ORDER BY id DESC LIMIT %v", n)
	err := db.Select(&_articles, db.Rebind(sql))
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _articles, nil
}

// FindArticles find one or more articles by one or more ids
func FindArticles(ids ...int64) ([]Article, error) {
	if len(ids) == 0 {
		msg := "At least one or more ids needed"
		log.Println(msg)
		return nil, errors.New(msg)
	}
	_articles := []Article{}
	idsHolder := strings.Repeat(",?", len(ids)-1)
	sql := db.Rebind(fmt.Sprintf(`SELECT * FROM articles WHERE id IN (?%s)`, idsHolder))
	idsT := []interface{}{}
	for _, id := range ids {
		idsT = append(idsT, interface{}(id))
	}
	err := db.Select(&_articles, sql, idsT...)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _articles, nil
}

// FindArticleBy find a single article by a field name and a value
func FindArticleBy(field string, val interface{}) (*Article, error) {
	_article := Article{}
	sqlFmt := `SELECT * FROM articles WHERE %s = ? LIMIT 1`
	sqlStr := fmt.Sprintf(sqlFmt, field)
	err := db.Get(&_article, db.Rebind(sqlStr), val)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return &_article, nil
}

// FindArticlesBy find all articles by a field name and a value
func FindArticlesBy(field string, val interface{}) (_articles []Article, err error) {
	sqlFmt := `SELECT * FROM articles WHERE %s = ?`
	sqlStr := fmt.Sprintf(sqlFmt, field)
	err = db.Select(&_articles, db.Rebind(sqlStr), val)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	return _articles, nil
}

// AllArticles get all the Article records
func AllArticles() (articles []Article, err error) {
	err = db.Select(&articles, "SELECT * FROM articles")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return articles, nil
}

// ArticleCount get the count of all the Article records
func ArticleCount() (c int64, err error) {
	err = db.Get(&c, "SELECT count(*) FROM articles")
	if err != nil {
		log.Println(err)
		return 0, err
	}
	return c, nil
}

// ArticleCountWhere get the count of all the Article records with a where clause
func ArticleCountWhere(where string, args ...interface{}) (c int64, err error) {
	sql := "SELECT count(*) FROM articles"
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

// ArticleIncludesWhere get the Article associated models records, it's just the eager_load function
func ArticleIncludesWhere(assocs []string, sql string, args ...interface{}) (_articles []Article, err error) {
	_articles, err = FindArticlesWhere(sql, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	if len(assocs) == 0 {
		log.Println("No associated fields ard specified")
		return _articles, err
	}
	if len(_articles) <= 0 {
		return nil, errors.New("No results available")
	}
	ids := make([]interface{}, len(_articles))
	for _, v := range _articles {
		ids = append(ids, interface{}(v.Id))
	}
	idsHolder := strings.Repeat(",?", len(ids)-1)
	for _, assoc := range assocs {
		switch assoc {
		case "comments":
			where := fmt.Sprintf("article_id IN (?%s)", idsHolder)
			_comments, err := FindCommentsWhere(where, ids...)
			if err != nil {
				log.Printf("Error when query associated objects: %v\n", assoc)
				continue
			}
			for _, vv := range _comments {
				for i, vvv := range _articles {
					if vv.ArticleId == vvv.Id {
						vvv.Comments = append(vvv.Comments, vv)
					}
					_articles[i].Comments = vvv.Comments
				}
			}
		}
	}
	return _articles, nil
}

// ArticleIds get all the Ids of Article records
func ArticleIds() (ids []int64, err error) {
	err = db.Select(&ids, "SELECT id FROM articles")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ids, nil
}

// ArticleIdsWhere get all the Ids of Article records by where restriction
func ArticleIdsWhere(where string, args ...interface{}) ([]int64, error) {
	ids, err := ArticleIntCol("id", where, args...)
	return ids, err
}

// ArticleIntCol get some int64 typed column of Article by where restriction
func ArticleIntCol(col, where string, args ...interface{}) (intColRecs []int64, err error) {
	sql := "SELECT " + col + " FROM articles"
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

// ArticleStrCol get some string typed column of Article by where restriction
func ArticleStrCol(col, where string, args ...interface{}) (strColRecs []string, err error) {
	sql := "SELECT " + col + " FROM articles"
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

// FindArticlesWhere query use a partial SQL clause that usually following after WHERE
// with placeholders, eg: FindUsersWhere("first_name = ? AND age > ?", "John", 18)
// will return those records in the table "users" whose first_name is "John" and age elder than 18
func FindArticlesWhere(where string, args ...interface{}) (articles []Article, err error) {
	sql := "SELECT * FROM articles"
	if len(where) > 0 {
		sql = sql + " WHERE " + where
	}
	stmt, err := db.Preparex(db.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&articles, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return articles, nil
}

// FindArticlesBySQL query use a complete SQL clause
// with placeholders, eg: FindUsersBySQL("SELECT * FROM users WHERE first_name = ? AND age > ?", "John", 18)
// will return those records in the table "users" whose first_name is "John" and age elder than 18
func FindArticlesBySQL(sql string, args ...interface{}) (articles []Article, err error) {
	stmt, err := db.Preparex(db.Rebind(sql))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	err = stmt.Select(&articles, args...)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return articles, nil
}

func CreateArticle(am map[string]interface{}) (int64, error) {
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
	sqlFmt := `INSERT INTO articles (%s) VALUES (%s)`
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

func (_article *Article) Create() (int64, error) {
	ok, err := govalidator.ValidateStruct(_article)
	if !ok {
		errMsg := "Validate Article struct error: Unknown error"
		if err != nil {
			errMsg = "Validate Article struct error: " + err.Error()
		}
		log.Println(errMsg)
		return 0, errors.New(errMsg)
	}
	t := time.Now()
	_article.CreatedAt = t
	_article.UpdatedAt = t
	sql := `INSERT INTO articles (title,text,created_at,updated_at) VALUES (:title,:text,:created_at,:updated_at)`
	result, err := db.NamedExec(sql, _article)
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

// CommentsCreate used to create the associated objects
func (_article *Article) CommentsCreate(am map[string]interface{}) error {
	am["article_id"] = _article.Id
	_, err := CreateComment(am)
	return err
}

// GetComments used to get associated objects
func (_article *Article) GetComments() error {
	_comments, err := ArticleGetComments(_article.Id)
	if err == nil {
		_article.Comments = _comments
	}
	return err
}

// ArticleGetComments a helper fuction used to get associated objects for ArticleIncludesWhere()
func ArticleGetComments(id int64) ([]Comment, error) {
	_comments, err := FindCommentsBy("article_id", id)
	return _comments, err
}

func (_article *Article) Destroy() error {
	if _article.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := DestroyArticle(_article.Id)
	return err
}

func DestroyArticle(id int64) error {
	// Destroy association objects at first
	// Not care if exec properly temporarily
	destroyArticleAssociations(id)
	stmt, err := db.Preparex(db.Rebind(`DELETE FROM articles WHERE id = ?`))
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func DestroyArticles(ids ...int64) (int64, error) {
	if len(ids) == 0 {
		msg := "At least one or more ids needed"
		log.Println(msg)
		return 0, errors.New(msg)
	}
	// Destroy association objects at first
	// Not care if exec properly temporarily
	destroyArticleAssociations(ids...)
	idsHolder := strings.Repeat(",?", len(ids)-1)
	sql := fmt.Sprintf(`DELETE FROM articles WHERE id IN (?%s)`, idsHolder)
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

// DestroyArticlesWhere delete records by a where clause
// like: DestroyArticlesWhere("name = ?", "John")
// And this func will not call the association dependent action
func DestroyArticlesWhere(where string, args ...interface{}) (int64, error) {
	sql := `DELETE FROM articles WHERE `
	if len(where) > 0 {
		sql = sql + where
	} else {
		return 0, errors.New("No WHERE conditions provided")
	}
	ids, x_err := ArticleIdsWhere(where, args...)
	if x_err != nil {
		log.Printf("Delete associated objects error: %v\n", x_err)
	} else {
		destroyArticleAssociations(ids...)
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

// the func not return err temporarily
func destroyArticleAssociations(ids ...int64) {
	idsHolder := ""
	if len(ids) > 1 {
		idsHolder = strings.Repeat(",?", len(ids)-1)
	}
	idsT := []interface{}{}
	for _, id := range ids {
		idsT = append(idsT, interface{}(id))
	}
	var err error
	// make sure no declared-and-not-used exception
	_, _, _ = idsHolder, idsT, err
	where := fmt.Sprintf("article_id IN (?%s)", idsHolder)
	_, err = DestroyCommentsWhere(where, idsT...)
	if err != nil {
		log.Printf("Destroy associated object %s error: %v\n", "Comments", err)
	}
}

func (_article *Article) Save() error {
	ok, err := govalidator.ValidateStruct(_article)
	if !ok {
		errMsg := "Validate Article struct error: Unknown error"
		if err != nil {
			errMsg = "Validate Article struct error: " + err.Error()
		}
		log.Println(errMsg)
		return errors.New(errMsg)
	}
	if _article.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	_article.UpdatedAt = time.Now()
	sqlFmt := `UPDATE articles SET %s WHERE id = %v`
	sqlStr := fmt.Sprintf(sqlFmt, "title = :title, text = :text, updated_at = :updated_at", _article.Id)
	_, err = db.NamedExec(sqlStr, _article)
	return err
}

func UpdateArticle(id int64, am map[string]interface{}) error {
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
	sqlFmt := `UPDATE articles SET %s WHERE id = %v`
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

func (_article *Article) Update(am map[string]interface{}) error {
	if _article.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateArticle(_article.Id, am)
	return err
}

func (_article *Article) UpdateAttributes(am map[string]interface{}) error {
	if _article.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateArticle(_article.Id, am)
	return err
}

func (_article *Article) UpdateColumns(am map[string]interface{}) error {
	if _article.Id == 0 {
		return errors.New("Invalid Id field: it can't be a zero value")
	}
	err := UpdateArticle(_article.Id, am)
	return err
}

func UpdateArticlesBySQL(sql string, args ...interface{}) (int64, error) {
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
