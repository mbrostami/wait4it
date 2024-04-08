package mysql

import (
    "context"
    "database/sql"
    "errors"
    "fmt"
    "wait4it/pkg/model"
    "github.com/go-sql-driver/mysql"
)

// BuildConnectionString creates a MySQL connection string
func (m *MySQLConnection) BuildConnectionString() string {
    // Standard MySQL connection string format
    return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", m.Username, m.Password, m.Host, m.Port, m.DatabaseName)
}

func (m *MySQLConnection) BuildContext(cx model.CheckContext) {
    m.Port = cx.Port
    m.Host = cx.Host
    m.Username = cx.Username
    m.Password = cx.Password
    m.DatabaseName = cx.DatabaseName
}

// Validate confirms that the host and username are provided and that the port is in the valid range
func (m *MySQLConnection) Validate() error {
    if m.Host == "" || m.Username == "" {
        return errors.New("host or username can't be empty")
    }

    if m.Port < 1 || m.Port > 65535 {
        return errors.New("invalid port range for MySQL")
    }

    return nil
}

func (m *MySQLConnection) Check(ctx context.Context) (bool, bool, error) {
    dsl := m.BuildConnectionString()

    db, err := sql.Open("mysql", dsl)
    if err != nil {
        return false, true, err
    }

    // Setting logger for MySQL driver
    err = mysql.SetLogger(mysql.Logger)
    if err != nil {
        return false, true, err
    }

    // Ping the database to check connectivity
    err = db.PingContext(ctx)
    if err != nil {
        return false, false, err
    }
    _ = db.Close()

    return true, true, nil
}
