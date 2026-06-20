package utils

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

func SendSuccess(c *gin.Context, message string, data interface{}) {
    c.JSON(http.StatusOK, gin.H{
        "success": true,
        "message": message,
        "data":    data,
    })
}
func SendError(c *gin.Context, statusCode int, message string) {
    c.JSON(statusCode, gin.H{
        "success": false,
        "message": message,
    })
}