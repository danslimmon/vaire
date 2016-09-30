package main

import (
	"math/rand"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/contrib/ginrus"
	"github.com/gin-gonic/gin"
)

var (
	ReqIdChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	ReqIdLen   = 16
)

// Returns a logrus entry with fields based on the gin Context.
//
// This adds the `reqId` field containing the request ID populated by middlewareReqId().
func logger(c *gin.Context) *log.Entry {
	reqIdIface, exists := c.Get("reqId")
	if !exists {
		return log.WithFields(log.Fields{})
	}
	reqId, ok := reqIdIface.(string)
	if !ok {
		log.Error("reqId is not string")
		return log.WithFields(log.Fields{})
	}
	return log.WithFields(log.Fields{
		"reqId": reqId,
	})
}

// Middleware that assigns a random request ID to each request.
func middlewareReqId(c *gin.Context) {
	var reqIdBytes []byte
	for i := 0; i < ReqIdLen; i++ {
		reqIdBytes = append(reqIdBytes, ReqIdChars[rand.Intn(len(ReqIdChars))])
	}
	reqId := string(reqIdBytes)
	c.Set("reqId", reqId)
	c.Header("Vaire-ReqId", reqId)
	c.Next()
}

// Middleware that checks authentication/authorization
func middlewareCheckAuth(c *gin.Context) {
	authResult := checkAuth(c)
	if !authResult.Authorized {
		logger(c).WithFields(log.Fields{
			"error": authResult.Error,
		}).Info("Auth/auth check failed")
		c.JSON(403, gin.H{
			"error": authResult.Error,
		})
		c.Abort()
		return
	}
	c.Next()
}

func main() {
	// Seed the PRNG
	rand.Seed(time.Now().UnixNano())

	// Read config
	err := LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Set up router
	r := gin.New()
	r.Use(middlewareReqId)
	r.Use(middlewareCheckAuth)
	r.Use(ginrus.Ginrus(log.StandardLogger(), time.RFC3339, true))
	r.POST("/api/v1/queues/:queueName", handle_POST_Queues_QueueName)

	// Start listening
	log.WithFields(log.Fields{
		"listen_addr": Config.Listen,
	}).Info("Listening for connections")
	r.Run(Config.Listen)
}
