package resources

import (
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
	"time"

	"gorm.io/gorm" // gorm.DB

	"gorm.io/driver/postgres"
)

type User struct {
	UserId   string
	Username string
}

type Resources struct {
	// rooms and memberships postgresql database
	ROOMSDB *gorm.DB
	// client user mongodb database
	USERDB *mongo.Client
	// connection timeout
	Timeout time.Duration
	// user from Token
	User User
}

type RoomModel struct {
	// general
	Id        string `json:"id"`
	Title     string `json:"title"`
	Image     string `json:"image"`
	Biography string `json:"biography"`
	Creator   string `json:"creator"`
	CreatorId string `json:"creator_id"`
	// types is mutual, chat_group and channel
	Type string `json:"type"`
	// status
	IsReported bool `json:"is_reported"`
	IsActive   bool `json:"is_active"`
	IsPublic   bool `json:"is_public"`
	// date
	CreatedAt  int    `json:"created_at"`
	UpdatedAt  int    `json:"updated_at"`
	DeletedAt  int    `json:"deleted_at"`
	DeletedFor string `json:"deleted_for"`
}

type RoomMembershipModel struct {
	Id       string `json:"id"`
	RoomId   string `json:"room_id"`
	JoinBy   string `json:"join_by"`
	JoinById string `json:"join_by_id"`
	Username string `json:"username"`
	UserId   string `json:"user_id"`
	IsAdmin  bool   `json:"is_admin"`
	// mute related's
	MuteUntil int `json:"mute_until"`
	// date
	JoinAt  int `json:"join_at"`
	LeaveAt int `json:"leave_at"`
}

func GetROOMSDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", os.Getenv("POSTGRESQLROOMSHOST"), os.Getenv("POSTGRESQLROOMSPORT"), os.Getenv("POSTGRESQLROOMSUSER"), os.Getenv("POSTGRESQLROOMSPASSWORD"), os.Getenv("POSTGRESQLROOMSDBNAME"), os.Getenv("POSTGRESQLROOMSSSLMODE"))
}

func SetConnectionPool(db *gorm.DB) error {
	// Create the connection pool

	sqlDB, getDbErr := db.DB()
	if getDbErr != nil {
		return getDbErr
	}

	sqlDB.SetConnMaxIdleTime(time.Minute)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDB.SetMaxIdleConns(100)

	// SetMaxOpenConns sets the maximum number of open connections to the database.
	sqlDB.SetMaxOpenConns(5000)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDB.SetConnMaxLifetime(time.Hour * -1)
	return nil
}

func CreateDBConnection(stringDSM string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.New(postgres.Config{
		DSN:                  stringDSM,
		PreferSimpleProtocol: true, // disables implicit prepared statement usage
	}), &gorm.Config{})

	if err != nil {
		return nil, err
	}

	poolErr := SetConnectionPool(db)
	if poolErr != nil {
		return nil, poolErr
	}

	return db, nil
}

func GetROOMSDatabaseConnection() (*gorm.DB, error) {
	db, connErr := CreateDBConnection(GetROOMSDSN())
	if connErr != nil {
		return nil, connErr
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}
	if err := sqlDb.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

func getMongoURI() string {
	// mongodb://username:password@localhost:27017/?authSource=admin
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/?authSource=admin", os.Getenv("MONGODBUSERSUSER"), os.Getenv("MONGODBUSERSPASSWORD"), os.Getenv("MONGODBUSERSHOST"), os.Getenv("MONGODBUSERSPORT"))
	return uri
}

func GetUSERMongodbConnection() (*mongo.Client, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(getMongoURI()))
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return client, nil
}

func ConstructResource() (Resources, error) {
	r := Resources{
		ROOMSDB: nil,
		USERDB:  nil,
		Timeout: DefaultTimeout,
	}

	// TODO : uncomment in init
	sessInitRoom, err := GetROOMSDatabaseConnection()
	if err != nil {
		return r, err
	}
	var room RoomModel
	sessInitRoom.Table("rooms").AutoMigrate(room)

	sessInitMembership, err := GetROOMSDatabaseConnection()
	if err != nil {
		return r, err
	}
	var membership RoomMembershipModel
	sessInitMembership.Table("membership").AutoMigrate(membership)

	sess, err := GetROOMSDatabaseConnection()
	if err != nil {
		return r, err
	}
	r.ROOMSDB = sess

	cli, err := GetUSERMongodbConnection()
	if err != nil {
		return r, err
	}
	r.USERDB = cli

	return r, nil
}
