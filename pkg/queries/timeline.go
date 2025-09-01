package queries

import (
	"context"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"

	_ "modernc.org/sqlite"
	"zombiezen.com/go/sqlite"
	"zombiezen.com/go/sqlite/sqlitex"

	"github.com/robbymilo/rgallery/pkg/database"
	"github.com/robbymilo/rgallery/pkg/types"
)

type RawMinimalMedia = types.RawMinimalMedia
type Segment = types.Segment
type SegmentMedia = types.SegmentMedia
type SegmentGroup = types.SegmentGroup
type ResponseSegment = types.ResponseSegment
type GearItem = types.GearItem
type GearItems = types.GearItems
type Conf = types.Conf
type PrevNext = types.PrevNext

// GetTimeline returns all media items in a timeline format.
func GetTimeline(params *FilterParams, c Conf) (ResponseSegment, int, error) {
	mediaItems, err := getTimelineItems(params, c)
	if err != nil {
		return nil, 0, fmt.Errorf("error getting timeline items: %v", err)
	}

	mediaByDay, err := groupMediaItemsByDay(params, mediaItems)
	if err != nil {
		return nil, 0, fmt.Errorf("error grouping mediaItems by day: %v", err)
	}

	mediaByMonth := make(map[string][]Segment)
	for i := range mediaByDay {
		t, err := time.Parse("2006-01-02", mediaByDay[i].SegmentId)
		if err != nil {
			return nil, 0, fmt.Errorf("error parsing time for mediaByDay: %v", err)
		}
		month := t.Format("2006-01")
		mediaByMonth[month] = append(mediaByMonth[month], mediaByDay[i])
	}

	var final ResponseSegment
	for k, v := range mediaByMonth {
		days := v
		total := 0
		for i := range v {
			total += len(v[i].Media)
		}

		// sort days
		sort.Slice(days, func(i, j int) bool {
			if params.Direction == "desc" {
				return days[i].SegmentId > days[j].SegmentId
			}
			return days[i].SegmentId < days[j].SegmentId
		})

		segmentGroup := &SegmentGroup{
			SectionId: &k,
			Total:     &total,
			Segments:  &days,
		}
		final = append(final, segmentGroup)
	}

	// sort months
	sort.Slice(final, func(i, j int) bool {
		if params.Direction == "desc" {
			return *final[i].SectionId > *final[j].SectionId
		}
		return *final[i].SectionId < *final[j].SectionId
	})

	return final, len(mediaItems), nil
}

// getTimelineItems returns all filtered media items from the db.
func getTimelineItems(params *FilterParams, c Conf) ([]*RawMinimalMedia, error) {
	pool, err := sqlitex.NewPool(database.NewSqlConnectionString(c), sqlitex.PoolOptions{
		Flags:    sqlite.OpenReadOnly,
		PoolSize: 1,
	})
	if err != nil {
		return nil, fmt.Errorf("error opening sqlite db pool: %v", err)
	}
	defer pool.Close()

	conn, err := pool.Take(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to take connection from pool: %w", err)
	}
	defer pool.Put(conn)

	table := "media"
	term, err := sanitizeSearchInput(params.Term)
	if err != nil {
		c.Logger.Error("error formatting term query", "error", err)
	}
	if term != "" {
		table = `SELECT * FROM images_virtual(?)`
	}

	camera := params.Camera
	if camera != "" {
		camera = `AND i.camera =? `
	}

	lens := params.Lens
	if lens != "" {
		var p string
		for k, v := range c.Aliases.Lenses {
			if k == params.Lens {
				for _, value := range c.Aliases.Lenses {
					if v == value {
						p = fmt.Sprintf("%s%s", p, "?,")
					}
				}
			}
			if v == params.Lens {
				p = fmt.Sprintf("%s%s", p, "?,")
			}
		}

		if p != "" {
			lens = fmt.Sprintf(`AND i.lens in (%s)`, strings.TrimSuffix(p, ","))
		} else {
			lens = `AND i.lens =? `
		}
	}

	mediatype := params.MediaType
	if mediatype != "" {
		mediatype = `AND i.mediatype =? `
	}

	software := params.Software
	if software != "" {
		software = `AND i.software =? `
	}

	f35 := ""
	focallength35 := params.FocalLength35
	if focallength35 != 0 {
		f35 = `AND i.focallength35 =? `
	}

	folder := ""
	if params.Folder != "" {
		folder = `AND folder =? `
	}

	tag_join := ""
	tag := ""
	if params.Subject != "" {
		tag_join = `JOIN images_tags i_t ON i.hash = i_t.image_id
									JOIN tags t ON i_t.tag_id = t.id`
		tag = `AND t.key =?`
	}

	query := fmt.Sprintf(
		`SELECT DISTINCT
			i.hash,
			i.width,
			i.height,
			i.date,
			i.color,
			i.mediatype,
			i.offset,
			i.modified
		FROM (%s) i
		%s
		WHERE i.rating >=?
		AND i.date >=?
		%s
		%s
		%s
		%s
		%s
		%s
		%s
		GROUP BY i.date
		ORDER BY %s %s`, table, tag_join, tag, folder, camera, lens, mediatype, software, f35, params.OrderBy, params.Direction)

	stmt, err := conn.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error preparing SELECT statement: %v", err)
	}

	paramIdx := 1
	if term != "" {
		stmt.BindText(paramIdx, term)
		paramIdx++
	}

	stmt.BindInt64(paramIdx, int64(params.Rating))
	paramIdx++
	stmt.BindText(paramIdx, params.From)
	paramIdx++

	if camera != "" {
		stmt.BindText(paramIdx, params.Camera)
		paramIdx++
	}
	if lens != "" {
		var exists bool
		for k, v := range c.Aliases.Lenses {
			if k == params.Lens || v == params.Lens {
				exists = true
			}
		}
		if exists {
			for k, v := range c.Aliases.Lenses {
				if k == params.Lens {
					for key, value := range c.Aliases.Lenses {
						if v == value {
							stmt.BindText(paramIdx, key)
							paramIdx++
						}
					}

				}
				if v == params.Lens {
					stmt.BindText(paramIdx, k)
					paramIdx++
				}
			}
		} else {
			stmt.BindText(paramIdx, params.Lens)
			paramIdx++
		}
	}
	if mediatype != "" {
		stmt.BindText(paramIdx, params.MediaType)
		paramIdx++
	}
	if software != "" {
		stmt.BindText(paramIdx, params.Software)
		paramIdx++
	}
	if focallength35 != 0 {
		stmt.BindFloat(paramIdx, focallength35)
		paramIdx++
	}
	if params.Folder != "" {
		stmt.BindText(paramIdx, params.Folder)
		paramIdx++
	}

	if params.Subject != "" {
		stmt.BindText(paramIdx, params.Subject)
		paramIdx++ //nolint:all
	}

	var result []*RawMinimalMedia

	for {
		hasRow, err := stmt.Step()
		if err != nil {
			return nil, fmt.Errorf("error stepping through result set: %v", err)
		}
		if !hasRow {
			break
		}

		media := &RawMinimalMedia{
			Hash:      uint32Ptr(stmt.ColumnInt64(0)),
			Width:     intPtr(stmt.ColumnInt64(1)),
			Height:    intPtr(stmt.ColumnInt64(2)),
			Date:      stringPtr(stmt.ColumnText(3)),
			Color:     stringPtr(stmt.ColumnText(4)),
			MediaType: stringPtr(stmt.ColumnText(5)),
			Offset:    float64Ptr(stmt.ColumnFloat(6)),
			Modified:  stringPtr(stmt.ColumnText(7)),
		}
		result = append(result, media)
	}

	err = stmt.Finalize()
	if err != nil {
		return nil, fmt.Errorf("error finalizing statement: %v", err)
	}

	return result, nil
}

// groupMediaItemsByDay returns all filtered media items grouped by YYYY-MM-DD.
func groupMediaItemsByDay(params *FilterParams, mediaItems []*RawMinimalMedia) (map[string]Segment, error) {
	mediaByDay := make(map[string]Segment)
	for i := range mediaItems {
		var t time.Time
		var err error
		if params.OrderBy == "date" {
			if *mediaItems[i].Offset != 0 {
				t, err = getTimeLocal(*mediaItems[i].Date, *mediaItems[i].Offset)
				if err != nil {
					return nil, fmt.Errorf("error parsing local time: %v", err)
				}
			} else {
				t, err = time.Parse("2006-01-02T15:04:05.000Z", *mediaItems[i].Date)
				if err != nil {
					return nil, fmt.Errorf("error parsing raw time: %v", err)
				}
			}
		} else if params.OrderBy == "modified" {
			t, err = time.Parse("2006-01-02T15:04:05.000Z", *mediaItems[i].Modified)
			if err != nil {
				return nil, fmt.Errorf("error parsing modified time: %v", err)
			}
		}

		if t.Year() >= 1900 {
			day := t.Format("2006-01-02")
			imgs := mediaByDay[day].Media

			// Create a slice with base metadata
			metaSlice := make([]interface{}, 3)
			metaSlice[0] = *mediaItems[i].Width
			metaSlice[1] = *mediaItems[i].Height
			metaSlice[2] = *mediaItems[i].Hash

			// Add color if present
			if *mediaItems[i].Color != "" {
				metaSlice = append(metaSlice, *mediaItems[i].Color)
			}

			// Add video marker if needed
			if *mediaItems[i].MediaType == "video" {
				metaSlice = append(metaSlice, "v")
			}

			imgs = append(imgs, metaSlice)

			mediaByDay[day] = Segment{
				SegmentId: day,
				Media:     imgs,
			}
		}
	}

	return mediaByDay, nil
}

// sanitizeSearchInput sanitizes user input for search queries
// Allows alphanumeric characters, spaces, and Unicode letters (including diacritics)
func sanitizeSearchInput(input string) (string, error) {
	// Use a whitelist approach that allows:
	// - Alphanumeric (A-Z, a-z, 0-9)
	// - Spaces
	// - Unicode letters (including characters like š, é, ñ, etc.)
	re, err := regexp.Compile(`[^\p{L}\p{N} ]+`)
	if err != nil {
		return "", err
	}
	// Replace all non-matching characters with an empty string
	cleaned := re.ReplaceAllString(input, "")
	return cleaned, nil
}

// Helper functions for pointer conversion
func intPtr(v int64) *int           { val := int(v); return &val }
func uint32Ptr(v int64) *uint32     { val := uint32(v); return &val }
func stringPtr(v string) *string    { return &v }
func float64Ptr(v float64) *float64 { return &v }
