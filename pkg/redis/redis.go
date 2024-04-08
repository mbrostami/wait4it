package redis

import (
    "context"
    "errors"
    "fmt"
    "strconv"
    "wait4it/pkg/model"
    "github.com/go-redis/redis/v8"
)

const (
    Cluster    = "cluster"
    Standalone = "standalone"
)

// BuildConnectionString generates a connection string for Redis
func (m *RedisConnection) BuildConnectionString() string {
    // Standard Redis connection string format
    return fmt.Sprintf("%s:%d", m.Host, m.Port)
}

func (m *RedisConnection) BuildContext(cx model.CheckContext) {
    m.Host = cx.Host
    m.Port = cx.Port
    m.Password = cx.Password

    d, err := strconv.Atoi(cx.DatabaseName)
    if err != nil {
        d = 0 // Default to the first database if conversion fails
    }
    m.Database = d

    m.OperationMode = cx.DBConf.OperationMode
    if m.OperationMode != Cluster && m.OperationMode != Standalone {
        m.OperationMode = Standalone
    }
}

func (m *RedisConnection) Validate() error {
    if m.Host == "" {
        return errors.New("host cannot be empty")
    }

    if m.OperationMode != Cluster && m.OperationMode != Standalone {
        return errors.New("invalid operation mode")
    }

    if m.Port < 1 || m.Port > 65535 {
        return errors.New("invalid port range for Redis")
    }

    return nil
}

func (m *RedisConnection) Check(ctx context.Context) (bool, bool, error) {
    switch m.OperationMode {
    case Standalone:
        return m.checkStandAlone(ctx)
    case Cluster:
        return m.checkCluster(ctx)
    default:
        return false, false, nil
    }
}

func (m *RedisConnection) checkStandAlone(ctx context.Context) (bool, bool, error) {
    rdb := redis.NewClient(&redis.Options{
        Addr:     m.BuildConnectionString(),
        Password: m.Password, // no password set
        DB:       m.Database, // use default DB
    })

    _, err := rdb.Ping(ctx).Result()
    if err != nil {
        return false, false, err
    }

    _ = rdb.Close()

    return true, true, nil
}

func (m *RedisConnection) checkCluster(ctx context.Context) (bool, bool, error) {
    rdb := redis.NewClusterClient(&redis.ClusterOptions{
        Addrs:    []string{m.BuildConnectionString()}, // Cluster mode
        Password: m.Password,
    })
    defer rdb.Close()

    _, err := rdb.Ping(ctx).Result()
    if err != nil {
        return false, false, err
    }

    result, err := rdb.ClusterInfo(ctx).Result()
    if err != nil {
        return false, false, err
    }

    if result != "" && !strings.Contains(result, "cluster_state:ok") {
        return false, false, errors.New("cluster is not healthy")
    }

    return true, true, nil
}
