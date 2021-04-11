package routes

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/joyrex2001/kubedock/internal/container"
)

type ContainerExecRequest struct {
	Cmd []string `json:"Cmd"`
}

type ExecStartRequest struct {
	Detach bool `json:"Detach"`
	Tty    bool `json:"Tty"`
}

// POST "/containers/:id/exec"
func ContainerExec(c *gin.Context) {
	in := &ContainerExecRequest{}
	if err := json.NewDecoder(c.Request.Body).Decode(&in); err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}
	id := c.Param("id")
	ctainr, err := container.Load(id)
	if err != nil {
		Error(c, http.StatusNotFound, err)
		return
	}
	log.Printf("cmd = %v", in.Cmd)
	// TODO: implement prepare exec
	c.JSON(http.StatusCreated, gin.H{
		"Id": ctainr.ID,
	})
}

// POST "/exec/:id/start"
func ExecStart(c *gin.Context) {
	in := &ExecStartRequest{}
	if err := json.NewDecoder(c.Request.Body).Decode(&in); err != nil {
		Error(c, http.StatusInternalServerError, err)
		return
	}
	id := c.Param("id")
	// TODO: implement exec
	c.Writer.WriteHeader(http.StatusOK)
	if !in.Detach {
		r := c.Request
		w := c.Writer

		in, out, err := HijackConnection(w)
		if err != nil {
			Error(c, http.StatusInternalServerError, err)
			return
		}
		defer CloseStreams(in, out)

		if _, ok := r.Header["Upgrade"]; ok {
			fmt.Fprint(out, "HTTP/1.1 101 UPGRADED\r\nContent-Type: application/vnd.docker.raw-stream\r\nConnection: Upgrade\r\nUpgrade: tcp\r\n")
		} else {
			fmt.Fprint(out, "HTTP/1.1 200 OK\r\nContent-Type: application/vnd.docker.raw-stream\r\n")
		}

		// copy headers that were removed as part of hijack
		if err := w.Header().WriteSubset(out, nil); err != nil {
			Error(c, http.StatusInternalServerError, err)
			return
		}
		fmt.Fprint(out, "\r\n")

		log.Printf("attached mode for %s, return empty stdout/stderr", id)
		fmt.Fprintf(out, "") // nonohing, no stdout and no stderr result ;-)
	}
}

// GET "/exec/:id/json"
func ExecInfo(c *gin.Context) {
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"ID":       id,
		"Running":  false,
		"ExitCode": 0,
	})
}