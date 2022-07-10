package user

import (
	"database/sql"
	"github.com/huandu/go-sqlbuilder"
	"github.com/nicjohnson145/mixer-service/pkg/settings"
)

func getPublicUsers(db *sql.DB) ([]string, error) {
	sb := sqlbuilder.NewSelectBuilder()
	sb.Select("u.username")
	sb.From("user u")
	sb.JoinWithOption(sqlbuilder.LeftOuterJoin, "user_setting us", "u.username = us.username")
	sb.Where(
		sb.Or(
			sb.IsNull("us.username"),
			sb.And(
				sb.Equal("us.key", settings.PublicProfile),
				sb.Equal("us.value", "true"),
			),
		),
	)

	sql, args := sb.Build()
	rows, err := db.Query(sql, args...)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	users := []string{}
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return []string{}, err
		}

		users = append(users, username)
	}

	return users, nil
}
