# go-sql-gen

Generates functions for scanning sql.Rows into structs:

```golang
type MyStruct struct {
    A, B, C int
    D       float64
    E       bool `go-sql-gen:"-"` // this field will be ignored
    f       string // if flag -i provided, this field would be ignored, because it is private field
}

// go-sql-gen generates following functions:
// * GetMyStructContext(ctx context.Context, query string, args ...interface{}) (MyStruct, error)
// * GetMyStructListContext(ctx context.Context, query string, args ...interface{}) ([]MyStruct, error)

func GetMyStructContext(ctx context.Context, db *sql.DB, query string, args ...interface{}) (result MyStruct, err error) {
    err = db.QueryRowContext(ctx, query, args).Scan(&result.A, &result.B, &result.C, &result.D)
    return result, err
}

func GetMyStructListContext(ctx context.Context, db *sql.DB, query string, args ...interface{}) (result []MyStruct, err error) {
    rows, err := db.QueryContext(ctx, query, args)
    if err != nil {
        return nil, err
    }
    defer rows.Close()  
    for rows.Next() {
        var item MyStruct
        if err := rows.Scan(&item.A, &item.B, &item.C, &item.D); err != nil {
            return nil, err
        }
        result = append(result, item)
    }
    return result, err
}

```

### Flag usage
 + **-i** - ignore private struct fields
 + **-m** - if you want generate code only for specific struct name(regexp)
