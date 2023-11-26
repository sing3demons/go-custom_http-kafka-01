package utils

import (
	"fmt"
	"os"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func ConvertToObjectID(id any) (primitive.ObjectID, error) {
	return primitive.ObjectIDFromHex(fmt.Sprintf("%s", id))
}

func Href(typeName string, id string) string {
	uri := os.Getenv("HOST_URL")
	return uri + fmt.Sprintf("/%s/%s", typeName, id)
}
