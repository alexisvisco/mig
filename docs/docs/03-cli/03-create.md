# Create

Create will add a new migration file to the migrations folder.

There is two kinds of migrations files: 

- classic: they have an Up, Down, Date and Name function.
- change: they have a Change, Date and Name function.

To create a new migration file, run the following command:

```sh
mig create <name>
```

It will create a new migration file `<version>_<name>.go` in the `migrations` folder.

Example: 

```go
package migrations

import (
    "github.com/alexisvisco/mig/pkg/schema/pg"
    "github.com/alexisvisco/mig/pkg/schema"
    "time"
)

type Migration20240502155033SchemaVersion struct {}

func (m Migration20240502155033SchemaVersion) Change(s *pg.Schema) {
    s.CreateTable("public.mig_schema_versions", func(s *pg.PostgresTableDef) {
        s.String("id")
    }, schema.TableOptions{ IfNotExists: true })
}

func (m Migration20240502155033SchemaVersion) Name() string {
    return "schema_version"
}

func (m Migration20240502155033SchemaVersion) Date() time.Time {
    t, _  := time.Parse(time.RFC3339, "2024-05-02T17:50:33+02:00")
    return t
}
```

## Flags

- `--dump` will dump the schema of your database with `pg_dump` (make sure you have it).
- `--skip` will insert the current version of the schema into the `mig_schema_versions` table without running the migration (because the schema already exists).
- `--dump-schema` to specify the schema to dump. (default is `public`)
- `--pg-dump-path` to specify the path to the `pg_dump` command.
- `--type` to specify the type of migration to create, possible values are [classic, change] (default is `change`)

## Difference between classic and change migration

In the classic migration, you have to implement the Up and Down function. The Up function is used to apply the migration, and the Down function is used to rollback the migration.

While in the change migration, you have to implement the Change function. The Change function is used to apply the migration. 
The rollback is done by applying the inverse of the Change function. 

The inverse of the Change function is automatically generated by all mig methods. 

For example, if you have this method in the Change function: 
```go
s.AddColumn("name", schema.StringType)
```

The inverse of this method is automatically generated by mig and will be like this: 
```go
s.DropColumn("name")
```

If you have a custom SQL to execute, you can call the `Reversible` method. 

To understand why you need to use the `Reversible` method, you need to know that mig needs to know how to rollback the migration in a Change function.

Suppose you have this migration:

```go
func (m Migration20240502155033SchemaVersion) Change(s *pg.Schema) {
    s.CreateTable("public.mig_schema_versions", func(s *pg.PostgresTableDef) {
        s.String("id")
    })

    s.Exec("INSERT INTO public.mig_schema_versions (id) VALUES ('1')")
}
```

If you run this migration and then rollback it, mig will not know how to rollback the `INSERT INTO` statement.
So it will execute it, the problem that the table will be dropped and the `INSERT INTO` statement will fail.

To avoid this problem, you need to use the `Reversible` method.


```go
func (m Migration20240502155033SchemaVersion) Change(s *pg.Schema) {
    s.CreateTable("public.mig_schema_versions", func(s *pg.PostgresTableDef) {
        s.String("id")
    })

    s.Reversible(schema.Directions{
        Up: func() {
            s.Exec("INSERT INTO public.mig_schema_versions (id) VALUES ('1')")
        },
    
        Down: func() { },
    })
}
```




