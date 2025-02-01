package internal

import (
	"context"
	"errors"
	"fmt"
	"github.com/solarwinds/swo-sdk-go/swov1"
	"github.com/solarwinds/swo-sdk-go/swov1/models/components"
	"github.com/solarwinds/swo-sdk-go/swov1/models/operations"
	"log"
	"os"
	"strconv"
	"strings"
)

var swo *swov1.Swo

func init() {
	tok := os.Getenv("SWO_API_TOKEN")
	if len(tok) == 0 {
		log.Fatal("SWO_API_TOKEN environment variable not set")
	}
	swo = swov1.New(swov1.WithSecurity(os.Getenv("SWO_API_TOKEN")))
}

func TagEntities(entityType string, filter string, tags []string) error {
	splitFilter := strings.Split(filter, ":")
	if len(splitFilter) != 2 {
		return errors.New("invalid tag filter")
	}
	tagsMap := make(map[string]string)
	for _, tag := range tags {
		splitTag := strings.Split(tag, ":")
		if len(splitTag) != 2 {
			return errors.New("invalid tag")
		}
		tagsMap[splitTag[0]] = splitTag[1]
	}

	res, err := swo.Entities.ListEntities(context.Background(), operations.ListEntitiesRequest{
		Type: entityType,
	})
	if err != nil {
		return err
	}

	if res.Object != nil {
		for {
			entities := res.Object.GetEntities()
			for _, entity := range entities {
				if matchEntity(splitFilter[0], splitFilter[1], entity) {
					err := tagEntity(tagsMap, entity)
					if err != nil {
						return err
					}
				}
			}

			res, err := res.Next()

			if err != nil {
				return err
			}

			if res == nil {
				break
			}
		}
	}

	return nil
}

func tagEntity(tagsMap map[string]string, entity components.Entity) error {
	_, err := swo.Entities.UpdateEntityByID(context.Background(),
		operations.UpdateEntityByIDRequest{
			ID: entity.GetID(),
			EntityUpdate: components.EntityUpdate{
				Tags: tagsMap,
			},
		})

	if err != nil {
		return err
	}

	return nil
}

func matchEntity(attribute string, value string, entity components.Entity) bool {
	if strings.HasPrefix(attribute, "tags.") {
		return entity.GetTags()[attribute[5:]] == value
	}

	// try to match base entity attribute
	switch strings.ToLower(attribute) {
	case "id":
		return entity.GetID() == value
	case "name":
		return *entity.GetName() == value
	case "displayname":
		return *entity.GetDisplayName() == value
	}

	// try to map extended entity attributes
	v, err := getAttribute(attribute, entity)
	if err != nil {
		log.Printf("error getting attribute %s: %v", attribute, err)
		return false
	}
	return v == value
}

func getAttribute(attribute string, entity components.Entity) (string, error) {
	for k, v := range entity.GetAttributes() {
		if strings.ToLower(k) == strings.ToLower(attribute) {
			switch typedV := v.(type) {
			case string:
				return typedV, nil
			case int:
				return strconv.Itoa(typedV), nil
			case float64:
				return strconv.FormatFloat(typedV, 'f', -1, 64), nil
			default:
				return "", fmt.Errorf("could not handle attribute (%s): %v of type %T\n", k, v, typedV)
			}
		}
	}
	return "", nil
}
