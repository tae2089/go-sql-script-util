package util

import (
	"database/sql"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

var (
	// regex for removing line comments
	lineCommentRegex = regexp.MustCompile(`--.*`)
	// regex for removing multi-line comments
	multiBlockCommentRegex = regexp.MustCompile(`\/\*\s*(.*?)\s*\*\/`)
)

func ExecuteSqlDir(db *sql.DB, dirPath string) error {
	files, err := os.ReadDir(dirPath)
	if err != nil {
		return err
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		filePath := dirPath + "/" + file.Name()
		if err := ExecuteSQLFile(db, filePath); err != nil {
			return err
		}
	}
	return nil
}

func ExecuteSqlFiles(db *sql.DB, filePaths ...string) error {
	for _, filePath := range filePaths {
		if err := ExecuteSQLFile(db, filePath); err != nil {
			return err
		}
	}
	return nil
}

func ExecuteSQLFile(db *sql.DB, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	return executeCommandsFromScanner(db, file)
}

func executeCommandsFromScanner(db *sql.DB, f io.Reader) error {
	data, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	originalQuery := removeComments(string(data))
	splitsQuerys := strings.Split(originalQuery, ";")
	for idx, query := range splitsQuerys {
		if query == "" {
			continue
		}
		err = executeSQLCommand(db, query)
		if err != nil {
			log.Printf("Error executing %d of SQL commands: %v", idx, err)
		}
	}
	return nil
}

func executeSQLCommand(db *sql.DB, query string) error {
	// execute the query
	_, err := db.Exec(strings.TrimSpace(query))
	return err
}

func removeComments(line string) (newline string) {
	// remove comments
	newline = lineCommentRegex.ReplaceAllString(line, "")
	newline = multiBlockCommentRegex.ReplaceAllString(newline, "")
	return compactSQLStatement(newline)
}

func compactSQLStatement(sql string) string {
	// checking line for remove multiple spaces
	spaceRegex := `\s+`
	re := regexp.MustCompile(spaceRegex)
	// replace multiple spaces with a single space
	compactedSQL := re.ReplaceAllString(sql, " ")
	return compactedSQL
}

// deprecated
func checkHasMultiLineComment(line string, inMultiLineComment *bool) bool {
	skipLine := false
	// Check for the start of a multi-line comment
	if strings.Contains(line, "/*") {
		*inMultiLineComment = true
	}

	// If we're in a multi-line comment, check for the end of it
	if *inMultiLineComment {
		if strings.Contains(line, "*/") {
			*inMultiLineComment = false
		}
		skipLine = true
	}
	return skipLine
}

// deprecated
// Check if the line ends with a semicolon indicating the end of a command
func checkHasSemicolon(line string) bool {
	return strings.HasSuffix(strings.TrimSpace(line), ";")
}

// deprecated
func writeScriptToBuilder(scriptBuilder *strings.Builder, line string) {
	// Add line to current script if it's not a comment
	if line != "" {
		newSpace := " "
		if strings.HasSuffix(line, " ") {
			newSpace = ""
		}
		scriptBuilder.WriteString(line + newSpace)
	}
}
