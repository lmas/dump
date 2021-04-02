package internal

import (
	"time"

	"github.com/jmoiron/sqlx"
	//_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

const sqlSchema = `CREATE TABLE IF NOT EXISTS accounts (
	username VARCHAR(50) PRIMARY KEY,
	hash VARCHAR(60),
	email VARCHAR(50),
	created TIMESTAMP,
	kicked TIMESTAMP,
	banned TIMESTAMP,
	lastlogin TIMESTAMP
);

CREATE TABLE IF NOT EXISTS characters (
	id VARCHAR(26) UNIQUE NOT NULL,
	username VARCHAR(50),
	firstname VARCHAR(30),
	lastname VARCHAR(30),
	ismale BOOL,
	credits INTEGER
);
`

var sqlQueries = map[string]string{
	"createaccount": `INSERT INTO accounts VALUES(:username, :hash, :email, :created, :kicked, :banned, :lastlogin);`,
	"updateaccount": `UPDATE accounts SET hash = :hash, email = :email, lastlogin = :lastlogin WHERE username = :username;`,
	"delaccount":    `DELETE FROM accounts WHERE username = :username;`,
	"getaccount":    `SELECT * FROM accounts WHERE username = :username;`,

	"listcharacters":  `SELECT * FROM characters WHERE username = :username ORDER BY firstname ASC;`,
	"countcharacters": `SELECT count(*) FROM characters WHERE username = :username;`,
	"createcharacter": `INSERT INTO characters (id, username, firstname, lastname, ismale, credits) VALUES(:id, :username, :firstname, :lastname, :ismale, :credits);`,
	"getcharacter":    `SELECT * FROM characters WHERE id = :id AND username = :username LIMIT 1;`,
	"delcharacter":    `DELETE FROM characters WHERE id = :id AND username = :username;`,
}

type DB struct {
	*sqlx.DB
	prepared map[string]*sqlx.NamedStmt
}

func OpenDB(dsn string) (*DB, error) {
	db, err := sqlx.Connect("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	d := &DB{db, make(map[string]*sqlx.NamedStmt)}
	_, err = d.Exec(sqlSchema)
	if err != nil {
		return nil, err
	}
	for name, sql := range sqlQueries {
		stmt, err := d.PrepareNamed(sql)
		if err != nil {
			return nil, err
		}
		d.prepared[name] = stmt
	}
	return d, nil
}

////////////////////////////////////////////////////////////////////////////////

type Account struct {
	Username  string    `db:"username"`
	Hash      string    `db:"hash"`
	Email     string    `db:"email"`
	Created   time.Time `db:"created"`
	Kicked    time.Time `db:"kicked"`
	Banned    time.Time `db:"banned"`
	LastLogin time.Time `db:"lastlogin"`
}

func (a Account) IsZero() bool {
	return a.Username == ""
}

func (d *DB) CreateAccount(a Account) error {
	_, err := d.prepared["createaccount"].Exec(a)
	return err
}

func (d *DB) UpdateAccount(a Account) error {
	_, err := d.prepared["updateaccount"].Exec(a)
	return err
}

func (d *DB) DelAccount(username string) error {
	_, err := d.prepared["delaccount"].Exec(Account{Username: username})
	return err
}

func (d *DB) GetAccount(username string) (Account, error) {
	var a Account
	err := d.prepared["getaccount"].Get(&a, Account{Username: username})
	if err != nil {
		return Account{}, err
	}
	return a, nil
}

////////////////////////////////////////////////////////////////////////////////

type Character struct {
	ID        string    `db:"id"`
	Username  string    `db:"username"`
	FirstName string    `db:"firstname"`
	LastName  string    `db:"lastname"`
	Birthday  time.Time `db:"birtday"` // created at
	IsMale    bool      `db:"ismale"`
	Credits   int       `db:"credits"`
}

func (c Character) String() string {
	return c.Username + "(" + c.FullName() + ")"
}

func (c Character) FullName() string {
	return c.FirstName + " " + c.LastName
}

func (d *DB) ListCharacters(username string) ([]Character, error) {
	var chars []Character
	err := d.prepared["listcharacters"].Select(&chars, Account{Username: username})
	if err != nil {
		return nil, err
	}
	return chars, nil
}

func (d *DB) CountCharacters(username string) (int, error) {
	var i int
	err := d.prepared["countcharacters"].Get(&i, Account{Username: username})
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (d *DB) CreateCharacter(c Character) error {
	id, err := genULID(time.Now())
	if err != nil {
		return err
	}
	c.ID = id
	_, err = d.prepared["createcharacter"].Exec(c)
	return err
}

type accountChar struct {
	ID       string
	Username string
}

func (d *DB) GetCharacter(id, username string) (Character, error) {
	var c Character
	err := d.prepared["getcharacter"].Get(&c, accountChar{ID: id, Username: username})
	if err != nil {
		return Character{}, err
	}
	return c, nil
}

func (d *DB) DelCharacter(id, username string) error {
	_, err := d.prepared["delcharacter"].Exec(accountChar{ID: id, Username: username})
	return err
}
