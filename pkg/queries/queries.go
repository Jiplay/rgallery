package queries

import (
	"context"
	"fmt"
	"time"

	"github.com/robbymilo/rgallery/pkg/database"
	"github.com/robbymilo/rgallery/pkg/types"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"
)

type Media = types.Media
type DatabaseMedia = types.DatabaseMedia
type Subject = types.Subject
type Subjects = types.Subjects
type Day = types.Day
type Days = types.Days
type FilterParams = types.FilterParams

var columns = database.Columns()

// GetTotalMediaItems returns the number of media items.
func GetTotalMediaItems(rating int, from, to, camera, lens string, c Conf) (int, error) {
	pool, err := sqlitex.NewPool(database.NewSqlConnectionString(c), sqlitex.PoolOptions{
		Flags:    sqlite.OpenReadOnly,
		PoolSize: 5,
	})
	if err != nil {
		return 0, fmt.Errorf("error opening sqlite db pool: %v", err)
	}
	defer pool.Close()

	conn, err := pool.Take(context.Background())
	if err != nil {
		return 0, fmt.Errorf("failed to take connection from pool: %w", err)
	}
	defer pool.Put(conn)

	query := `SELECT COUNT(*) FROM media WHERE rating >=? AND date >=? AND date <=? AND DATE != '0001-01-01T00:00:00.000Z'`

	stmt, err := conn.Prepare(query)
	if err != nil {
		return 0, fmt.Errorf("error preparing query: %v", err)
	}
	defer func() {
		if finalizeErr := stmt.Finalize(); finalizeErr != nil {
			err = fmt.Errorf("error finalizing statement: %v", finalizeErr)
		}
	}()

	stmt.BindInt64(1, int64(rating))
	stmt.BindText(2, from)
	stmt.BindText(3, to)

	var total int
	hasRow, err := stmt.Step()
	if err != nil {
		return 0, fmt.Errorf("error executing query: %v", err)
	}
	if hasRow {
		total = int(stmt.ColumnInt64(0))
	}

	return total, nil
}

// getTimeLocal takes a UTC time and offset and returns a local date
func getTimeLocal(dateString string, offset float64) (time.Time, error) {

	// Parse UTC timestamp
	t, err := time.Parse(time.RFC3339, dateString)
	if err != nil {
		return time.Time{}, err
	}

	offsetDuration := time.Duration(offset/60)*time.Hour + time.Duration(offset/60)*time.Minute

	// Add offset to UTC time
	localTime := t.Add(offsetDuration)

	return localTime, nil
}
