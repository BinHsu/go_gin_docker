package main

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

var local_port int = 8080

func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		fmt.Printf("cannot open an SQLite memory database: %v\n", err)
		os.Exit(3)
	}
	defer db.Close()
	default_status_true := 1
	default_status_false := 0
	create_table_task := fmt.Sprintf("CREATE TABLE IF NOT EXISTS task(id INTEGER PRIMARY KEY AUTOINCREMENT, name TEXT NOT NULL, status BOOLEAN NOT NULL CHECK (status IN (%d, %d)))", default_status_false, default_status_true)
	_, err = db.Exec(create_table_task)
	if err != nil {
		fmt.Printf("cannot create schema: %v\n", err)
		os.Exit(3)
	}

	fmt.Println("hello gin")

	type Result struct {
		Id     int    `json:"id"`
		Name   string `json:"name" binding:"required"`
		Status int    `json:"status"`
	}

	type Result_p struct {
		Id     *int    `json:"id,omitempty"`
		Name   *string `json:"name,omitempty"`
		Status *int    `json:"status,omitempty"`
	}

	r := gin.Default()
	r.GET("/tasks", func(c *gin.Context) {
		var results []Result
		rows, err := db.Query("SELECT * FROM task")
		if err != nil {
			fmt.Printf("cannot select: %v\n", err)
			c.AbortWithStatus(500)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var id int
			var name string
			var status bool
			err = rows.Scan(&id, &name, &status)
			if err != nil {
				fmt.Printf("cannot fetch: %v\n", err)
				continue
			}
			//fmt.Printf("result: %d, %s, %v\n", id, name, status)
			var iStatus int
			if status {
				iStatus = default_status_true
			} else {
				iStatus = default_status_false
			}
			results = append(results, Result{id, name, iStatus})
		}

		if nil == results {
			c.AbortWithStatus(204)
			return
		}

		c.JSON(200, gin.H{
			"result": results,
		})
	})
	r.POST("/task", func(c *gin.Context) {
		var input Result
		err := c.BindJSON(&input)
		if err != nil {
			fmt.Printf("binding input failed: %v\n", err)
			c.AbortWithStatus(400)
			return
		}

		stmt, err := db.Prepare("INSERT INTO task(name, status) VALUES($1, $2) RETURNING id")
		if err != nil {
			fmt.Printf("cannot prepare stmt: %v\n", err)
			c.AbortWithStatus(500)
			return
		}
		defer stmt.Close()

		err = stmt.QueryRow(
			input.Name,
			0,
		).Scan(&input.Id)
		if err != nil {
			fmt.Printf("cannot run stmt: %v\n", err)
			c.AbortWithStatus(500)
			return
		}

		c.JSON(201, gin.H{
			"result": input,
		})
	})
	r.PUT("/task/:id", func(c *gin.Context) {
		id_string := c.Param("id")
		id, err := strconv.Atoi(id_string)
		if err != nil {
			// handle error
			fmt.Printf("get id failed: %v\n", err)
			c.AbortWithStatus(400)
			return
		}
		//default_name := "64c4c2f2-b32a-11ec-b909-0242ac120002"
		//default_status := -1
		input := Result_p{}
		err = c.BindJSON(&input)
		if err != nil {
			fmt.Printf("binding input failed: %v\n", err)
			c.AbortWithStatus(400)
			return
		}
		if input.Status != nil && *input.Status != default_status_true && *input.Status != default_status_false {
			fmt.Printf("invalid input status: %d\n", input.Status)
			c.AbortWithStatus(400)
			return
		}

		rows, err := db.Query("SELECT * FROM task WHERE id = " + id_string)
		if err != nil {
			fmt.Printf("cannot select: %v\n", err)
			c.AbortWithStatus(500)
			return
		}
		defer rows.Close()

		var results []Result
		for rows.Next() {
			var id int
			var name string
			var status bool
			err = rows.Scan(&id, &name, &status)
			if err != nil {
				fmt.Printf("cannot fetch: %v\n", err)
				continue
			}
			fmt.Printf("result: %d, %s, %v\n", id, name, status)
			var iStatus int
			if status {
				iStatus = default_status_true
			} else {
				iStatus = default_status_false
			}
			results = append(results, Result{id, name, iStatus})
		}

		if nil == results {
			c.AbortWithStatus(404)
			return
		}

		name_new := results[0].Name
		status_new := results[0].Status
		if input.Name != nil || input.Status != nil {
			if input.Name != nil {
				name_new = *input.Name
			}
			if input.Status != nil {
				status_new = *input.Status
			}
			cmd := fmt.Sprintf(`UPDATE task SET name = "%s", status = %d WHERE id = %d`, name_new, status_new, id)
			_, err = db.Exec(cmd)
			if err != nil {
				fmt.Printf("update failed: %v\n", err)
				c.AbortWithStatus(500)
				return
			}
		}
		c.JSON(200, gin.H{
			"result": input,
		})

	})
	r.DELETE("/task/:id", func(c *gin.Context) {
		id_string := c.Param("id")
		_, err := strconv.Atoi(id_string)
		if err != nil {
			// handle error
			fmt.Printf("get id failed: %v\n", err)
			c.AbortWithStatus(400)
			return
		}

		_, err = db.Exec("DELETE FROM task WHERE id = " + id_string)
		if err != nil {
			fmt.Printf("cannot delete: %v\n", err)
			c.AbortWithStatus(500)
			return
		}
	})
	path := fmt.Sprintf("0.0.0.0:%d", local_port)
	r.Run(path) // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
