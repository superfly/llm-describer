package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		jsonData := `{
			"id": "ryjreaz08zdji3u",
			"created": "2024-03-19 18:27:45.669Z",
			"updated": "2024-03-19 18:27:45.669Z",
			"name": "images",
			"type": "base",
			"system": false,
			"schema": [
				{
					"system": false,
					"id": "bbtbjjdu",
					"name": "file",
					"type": "file",
					"required": true,
					"presentable": false,
					"unique": false,
					"options": {
						"mimeTypes": [
							"image/jpeg",
							"image/png",
							"image/svg+xml",
							"image/gif",
							"image/webp"
						],
						"thumbs": [],
						"maxSelect": 1,
						"maxSize": 20000000,
						"protected": true
					}
				}
			],
			"indexes": [],
			"listRule": null,
			"viewRule": null,
			"createRule": null,
			"updateRule": null,
			"deleteRule": null,
			"options": {}
		}`

		collection := &models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collection); err != nil {
			return err
		}

		return daos.New(db).SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("ryjreaz08zdji3u")
		if err != nil {
			return err
		}

		return dao.DeleteCollection(collection)
	})
}
