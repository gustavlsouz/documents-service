package persistence

import (
	"context"
	"database/sql"
	"log"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

type migrator struct {
	persistence Persistence
}

func newMigrator(persistence Persistence) *migrator {
	return &migrator{persistence: persistence}
}

const queryExistsTable = `
SELECT EXISTS (
   SELECT FROM information_schema.tables 
   WHERE  table_schema = 'public'
   AND    table_name   = 'migrations'
   );
`

const statmentCreateMigrations = `
create table public.migrations (
version int8,
updatedAt timestamp without time zone
);
`

const statmentInsertFirstMigrationsRow = `
insert into public.migrations
values (0 ,now());`

func (migrator *migrator) Migrate(migrationPath string) {
	row := migrator.persistence.Database().QueryRow(queryExistsTable)

	if err := row.Err(); err != nil {
		log.Panic(err)
	}

	var migrationsExists bool
	err := row.Scan(&migrationsExists)

	if err != nil {
		log.Panic(err)
	}

	if !migrationsExists {
		log.Println("migrations column does not exists")
		_, err := migrator.persistence.Database().Exec(statmentCreateMigrations)
		if err != nil {
			log.Panic(err)
		}

		_, err = migrator.persistence.Database().Exec(statmentInsertFirstMigrationsRow)
		if err != nil {
			log.Panic(err)
		}
	}

	row = migrator.persistence.Database().
		QueryRow(`select version from migrations`)

	if err := row.Err(); err != nil {
		log.Panic(err)
	}

	var version int64

	err = row.Scan(&version)

	if err != nil {
		log.Panic(err)
	}

	fileNames, err := listFiles(migrationPath)
	if err != nil {
		log.Panic("Error to load files in migrations path", err)
	}

	log.Println(fileNames)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)

	defer cancel()

	fileVersion := getVersionFromFileName(fileNames[len(fileNames)-1])

	if fileVersion == version {
		return
	}

	for _, fileName := range fileNames {
		log.Println("reading", fileName)
		filenameSections := strings.Split(fileName, ".")

		fileVersion, err := strconv.ParseInt(filenameSections[0], 10, 64)

		if err != nil {
			log.Panic(err)
		}

		if filenameSections[2] == "down" || fileVersion <= version {
			continue
		}
		filePath := path.Join("../deployments/migrations", fileName)
		log.Println("opening", fileName)
		fileBytes, err := os.ReadFile(filePath)

		if err != nil {
			log.Panic(err)
		}

		fileContent := string(fileBytes)

		tx, err := migrator.persistence.Database().
			BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
		if err != nil {
			log.Panic(err)
		}

		_, err = tx.Exec(fileContent)
		if err != nil {
			log.Panic(err)
		}

		_, err = tx.Exec(`update migrations set version = $1`, fileVersion)
		if err != nil {
			log.Panic(err)
		}

		tx.Commit()
	}

}

func listFiles(path string) ([]string, error) {
	files, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	fileNames := make([]string, 0, len(files))
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}

	sort.Strings(fileNames)

	return fileNames, nil
}

func getVersionFromFileName(fileName string) int64 {
	filenameSections := strings.Split(fileName, ".")

	fileVersion, err := strconv.ParseInt(filenameSections[0], 10, 64)

	if err != nil {
		log.Panic(err)
	}
	return fileVersion
}
