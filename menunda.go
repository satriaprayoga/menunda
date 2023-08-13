package menunda

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

const version = "1.0.0"

type Menunda struct {
	AppName  string
	Debug    bool
	Version  string
	ErrorLog *log.Logger
	InfoLog  *log.Logger
	DB       Database
	RootPath string
	Routes   *echo.Echo
	config   config
}

func (m *Menunda) New(rootPath string) error {
	folderNames := []string{"controllers", "migrations", "models", "data", "public", "tmp", "logs", "middleware"}
	err := m.Init(rootPath, folderNames)
	if err != nil {
		return err
	}

	err = m.checkDotEnv(rootPath)
	if err != nil {
		return err
	}

	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	infoLog, errorLog := m.setupLogs()

	if os.Getenv("DB_TYPE") != "" {
		db, err := m.OpenDB(os.Getenv("DB_TYPE"), m.buildConnString())
		if err != nil {
			errorLog.Println(err)
			os.Exit(1)
		}
		m.DB = Database{
			DbType: os.Getenv("DB_TYPE"),
			Pool:   db,
		}
	}

	m.InfoLog = infoLog
	m.ErrorLog = errorLog

	m.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	m.Version = version
	m.RootPath = rootPath
	m.Routes = m.routes()

	m.config = config{
		port: os.Getenv("PORT"),
	}

	return err
}

func (m *Menunda) Init(root string, folderNames []string) error {
	for _, path := range folderNames {
		err := m.CreateDirIfNotExist(root + "/" + path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Menunda) ListenAndServe() {
	serv := &http.Server{
		Addr:         fmt.Sprintf(":%s", os.Getenv("PORT")),
		ErrorLog:     m.ErrorLog,
		Handler:      m.Routes,
		IdleTimeout:  30 * time.Second,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	defer m.DB.Pool.Close()

	m.InfoLog.Printf("Listening on port %s", os.Getenv("PORT"))
	err := serv.ListenAndServe()
	m.ErrorLog.Fatal(err)
}

func (m *Menunda) checkDotEnv(path string) error {
	err := m.CreateFileIfNotExists(fmt.Sprintf("%s/.env", path))
	if err != nil {
		return err
	}
	return nil
}

func (m *Menunda) setupLogs() (*log.Logger, *log.Logger) {
	var (
		infoLog  *log.Logger
		errorLog *log.Logger
	)

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return infoLog, errorLog
}

func (m *Menunda) buildConnString() string {
	var conn string

	switch os.Getenv("DB_TYPE") {
	case "postgres", "postgresql":
		conn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5",
			os.Getenv("DB_HOST"),
			os.Getenv("DB_PORT"),
			os.Getenv("DB_USER"),
			os.Getenv("DB_NAME"),
			os.Getenv("DB_SSL_MODE"))
		if os.Getenv("DB_PASSWORD") != "" {
			conn = fmt.Sprintf("%s password=%s", conn, os.Getenv("DB_PASSWORD"))
		}
	default:
	}
	return conn
}
