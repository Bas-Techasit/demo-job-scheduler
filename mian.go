package main

import (
	"fmt"
	"gocron/repository"
	"gocron/services"
	"strings"
	"time"

	"github.com/go-co-op/gocron"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/spf13/viper"
)

func main() {
	initConfig()
	initTimeZone()
	dbConfig_1 := DBConfig{
		Driver:   viper.GetString("db1.driver"),
		Username: viper.GetString("db1.username"),
		Password: viper.GetString("db1.password"),
		Host:     viper.GetString("db1.host"),
		Port:     viper.GetString("db1.port"),
		Database: viper.GetString("db1.database"),
	}

	dbConfig_2 := DBConfig{
		Driver:   viper.GetString("db2.driver"),
		Username: viper.GetString("db2.username"),
		Password: viper.GetString("db2.password"),
		Host:     viper.GetString("db2.host"),
		Port:     viper.GetString("db2.port"),
		Database: viper.GetString("db2.database"),
	}

	db_1 := initDB(dbConfig_1)
	db_2 := initDB(dbConfig_2)

	jobRepoDB1 := repository.NewJobRepository(db_1)
	jobRepoDB2 := repository.NewJobRepository(db_2)

	jobServiceForDB1 := services.NewJobService(jobRepoDB1)
	jobServiceForDB2 := services.NewJobService(jobRepoDB2)

	loc, _ := time.LoadLocation("Asia/Bangkok")
	cron := gocron.NewScheduler(loc)
	cron.SingletonModeAll()

	jobServiceForDB1.ScheduleAllJob(cron, "check DB1")
	jobServiceForDB2.ScheduleAllJob(cron, "check DB2")

	cron.StartAsync()
	fmt.Println(cron.GetAllTags())
	fmt.Println(cron.Len())

	time.Sleep(2 * time.Minute)

}

func initTimeZone() {
	ict, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		panic(err)
	}
	time.Local = ict
}

func initDB(config DBConfig) *sqlx.DB {

	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	db, err := sqlx.Open(config.Driver, dsn+"?parseTime=true")
	if err != nil {
		panic(err)
	}
	return db
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

type DBConfig struct {
	Driver   string
	Username string
	Password string
	Host     string
	Port     string
	Database string
}
