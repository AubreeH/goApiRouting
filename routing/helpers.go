package routing

import "github.com/gin-gonic/gin"

func (_ BaseApi) NoMiddleware(_ *gin.Context) bool { return true }
