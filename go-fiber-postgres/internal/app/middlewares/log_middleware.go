package middlewares

// import (
// 	"go-fiber-postgres/internal/app/exception"
// 	"time"

// 	"github.com/gofiber/fiber/v3"
// 	"github.com/sirupsen/logrus"
// )

// func LogMiddleware(cfg *config.AppConfig) fiber.Handler {
// 	return func(ctx *fiber.Ctx) error {
// 		start := time.Now()
// 		err := ctx.Next()
// 		latency := time.Since(start)

// 		// Default status code
// 		statusCode := ctx.Response().StatusCode()

// 		// Jika ada error, kita bedah error-nya untuk dapet status code asli
// 		if err != nil {
// 			switch e := err.(type) {
// 			case *exception.BaseException:
// 				statusCode = e.StatusCode
// 			case *fiber.Error:
// 				statusCode = e.Code
// 			default:
// 				// Jika status masih 200 padahal ada error yang ga dikenal
// 				if statusCode == 200 {
// 					statusCode = 500
// 				}
// 			}
// 		}

// 		fields := logrus.Fields{
// 			"request_id": ctx.Get(fiber.HeaderXRequestID),
// 			"user_id":    ctx.Locals("user_id"),
// 			"method":     ctx.Method(),
// 			"path":       ctx.Path(),
// 			"status":     statusCode,
// 			"latency":    latency.String(),
// 		}

// 		if statusCode >= 500 {
// 			if err != nil {
// 				fields["error"] = err.Error()
// 			}
// 			config.ErrorLog.WithFields(fields).Error("Internal Server Error")
// 		} else if statusCode >= 400 && statusCode < 500 {
// 			config.ErrorLog.WithFields(fields).Warn("Client Error")
// 		} else {
// 			config.AccessLog.WithFields(fields).Info("Request Success")
// 		}

// 		return err
// 	}
// }
