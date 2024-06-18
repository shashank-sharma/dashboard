package routes

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v5"
)

func streamMP3(path string) (io.Reader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	return f, nil
}

func AudioStreamMP3(c echo.Context) error {
	path := c.QueryParam("path")
	if path == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "path parameter is required"})
	}

	reader, err := streamMP3(path)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("error streaming mp3: %v", err)})
	}

	// Create a buffered reader with a buffer size of 4096 bytes.
	bufferedReader := bufio.NewReaderSize(reader, 4096)

	c.Response().Header().Set("Content-Type", "audio/mpeg")
	_, err = io.CopyBuffer(c.Response(), bufferedReader, nil)
	if err != nil {
		return err
	}
	return nil
}
