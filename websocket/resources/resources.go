package resources

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"gorm.io/gorm" // gorm.DB

	"github.com/go-redis/redis/v8" // Redis driver
	"github.com/golang-jwt/jwt"
	"gorm.io/driver/postgres"
)

type User struct {
	ConnectionId string
	UserId       string
	Username     string
	AccessToken  string
}

type Resource struct {
	// messages postgresql database
	MESSAGEDB *gorm.DB
	// rooms and memberships postgresql database
	ROOMSDB *gorm.DB
	// connections redis database
	CONNDB *redis.Client
	// connection timeout
	Timeout time.Duration
	// user from Token
	User User
	// room Id
	RoomID string
}

type MessageModel struct {
	// general fields
	Id       string `json:"id"`
	GUID     string `json:"guid"`
	Username string `json:"username"`
	UserId   string `json:"user_id"`
	RoomId   string `json:"room_id"`
	// types can be text, link, image, voice, video, pdf and other files
	Type string `json:"type"`
	// text rebated's
	TextContent string `json:"text_content"`
	// link url
	LinkUrl string `json:"link_url,omitempty"`
	// date
	CreatedAt  int    `json:"created_at"`
	UpdatedAt  int    `json:"updated_at,omitempty"`
	DeletedAt  int    `json:"deleted_at,omitempty"`
	DeletedFor string `json:"deleted_for,omitempty"`
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

func GetMESSAGESDSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s", os.Getenv("POSTGRESQLMESSAGESHOST"), os.Getenv("POSTGRESQLMESSAGESPORT"), os.Getenv("POSTGRESQLMESSAGESUSER"), os.Getenv("POSTGRESQLMESSAGESPASSWORD"), os.Getenv("POSTGRESQLMESSAGESDBNAME"), os.Getenv("POSTGRESQLMESSAGESSSLMODE"))
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

// func GetDatabaseConnection() (*gorm.DB, error) {
// 	db, connErr := CreateDBConnection(GetDSN())
// 	if connErr != nil {
// 		return nil, connErr
// 	}

// 	sqlDb, err := db.DB()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if err := sqlDb.Ping(); err != nil {
// 		return nil, err
// 	}
// 	return db, nil
// }

func GetMESSAGESDatabaseConnection() (*gorm.DB, error) {
	db, connErr := CreateDBConnection(GetMESSAGESDSN())
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

// TODO : change the configs
func NewConnectionRedisSession() *redis.Client {
	dbInt, err := strconv.Atoi(os.Getenv("REDISCONNECTIONDB"))
	if err != nil {
		log.Fatalln("NewConnectionRedisSession/strconv.Atoi/err: ", err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDISCONNECTIONHOST") + ":" + os.Getenv("REDISCONNECTIONPORT"),
		Password: os.Getenv("REDISCONNECTIONPASSWORD"),
		DB:       dbInt,
	})

	return rdb
}

func SetupNotification() {
}

func ConstructResource() (Resource, error) {
	r := Resource{
		MESSAGEDB: nil,
		ROOMSDB:   nil,
		CONNDB:    nil,
		Timeout:   DefaultTimeout,
		User:      User{},
		RoomID:    "",
	}

	// TODO : uncomment in init
	sessInitMessages, err := GetMESSAGESDatabaseConnection()
	if err != nil {
		return r, err
	}
	var message MessageModel
	sessInitMessages.Table("messages").AutoMigrate(message)

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

	sess, err := GetMESSAGESDatabaseConnection()
	if err != nil {
		return r, err
	}
	r.MESSAGEDB = sess

	sess, err = GetROOMSDatabaseConnection()
	if err != nil {
		return r, err
	}
	r.ROOMSDB = sess

	r.CONNDB = NewConnectionRedisSession()

	return r, nil
}

func ConstructUserResource(tokenClaims jwt.MapClaims, res Resource, roomId string) (Resource, error) {
	user := User{
		UserId:      tokenClaims["user_id"].(string),
		Username:    tokenClaims["username"].(string),
		AccessToken: tokenClaims["accessToken"].(string),
	}

	r := Resource{
		MESSAGEDB: res.MESSAGEDB,
		ROOMSDB:   res.ROOMSDB,
		CONNDB:    res.CONNDB,
		Timeout:   DefaultTimeout,
		User:      user,
		RoomID:    roomId,
	}

	return r, nil
}
