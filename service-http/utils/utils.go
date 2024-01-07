package utils

import (
	"fmt"
	"net"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertToObjectID(id any) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(fmt.Sprintf("%s", id))
}

func Href(typeName string, id string) string {
	uri := os.Getenv("HOST_URL")
	return uri + fmt.Sprintf("/%s/%s", typeName, id)
}

func GetHeaders(ctx *gin.Context) map[string]any {
	// Request user agent
	userAgent := ctx.Request.UserAgent()
	platform := strings.Split(ctx.Request.Header.Get("sec-ch-ua"), ",")
	mobile := ctx.Request.Header.Get("sec-ch-ua-mobile")
	operatingSystem := ctx.Request.Header.Get("sec-ch-ua-platform")
	clientIP := ctx.ClientIP()
	reqId := ctx.Writer.Header().Get("X-Request-Id")
	if reqId == "" {
		reqId = uuid.NewString()
	}

	macIp := getMACAndIP()

	return map[string]any{
		"user_agent": userAgent,
		"Platform":   platform,
		"Mobile":     mobile,
		"OS":         operatingSystem,
		"client_ip":  clientIP,
		"request_id": reqId,
		"remote_ip":  ctx.Request.RemoteAddr,
		"mac_ip":     macIp,
	}
}

func getMACAndIP() MacIP {
	interfaces, _ := net.Interfaces()
	macAddr := MacIP{}
	for _, iface := range interfaces {

		if iface.Name != "" {
			macAddr.InterfaceName = iface.Name
		}

		if iface.HardwareAddr != nil {
			macAddr.HardwareAddr = iface.HardwareAddr.String()
		}

		var ips []string
		addrs, _ := iface.Addrs()

		for _, addr := range addrs {
			ips = append(ips, addr.String())
		}

		if len(ips) > 0 {
			macAddr.IPs = ips
		}
	}

	return macAddr
}

type MacIP struct {
	InterfaceName string   `json:"interface_name"`
	HardwareAddr  string   `json:"hardware_addr"`
	IPs           []string `json:"ips"`
}
