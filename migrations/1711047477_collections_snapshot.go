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
		jsonData := `[
			{
				"id": "ryjreaz08zdji3u",
				"created": "2024-03-19 18:27:45.669Z",
				"updated": "2024-03-21 18:13:02.783Z",
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
					},
					{
						"system": false,
						"id": "zb0eec7o",
						"name": "followups",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "w879h6ggcfuixkr",
							"cascadeDelete": true,
							"minSelect": null,
							"maxSelect": null,
							"displayFields": null
						}
					}
				],
				"indexes": [],
				"listRule": "      @request.auth.id != \"\" &&\n      @request.auth.images.id ?= id\n",
				"viewRule": "@request.auth.id != \"\" &&\n@request.auth.images.id ?= id\n",
				"createRule": "@request.auth.id != \"\"",
				"updateRule": "@request.auth.id != \"\" &&\n@request.auth.images.id ?= id",
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "w879h6ggcfuixkr",
				"created": "2024-03-19 18:35:25.779Z",
				"updated": "2024-03-21 18:43:25.436Z",
				"name": "followups",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "idkdigjk",
						"name": "user",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "7gbf23yk",
						"name": "text",
						"type": "text",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					}
				],
				"indexes": [],
				"listRule": "  @request.auth.id != \"\" &&\n@request.auth.images.followups.id ?= id",
				"viewRule": "  @request.auth.id != \"\" &&\n@request.auth.images.followups.id ?= id",
				"createRule": "  @request.auth.id != \"\"",
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "_pb_users_auth_",
				"created": "2024-03-21 18:02:26.605Z",
				"updated": "2024-03-21 18:02:26.614Z",
				"name": "users",
				"type": "auth",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "users_name",
						"name": "name",
						"type": "text",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"pattern": ""
						}
					},
					{
						"system": false,
						"id": "users_avatar",
						"name": "avatar",
						"type": "file",
						"required": false,
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
							"thumbs": null,
							"maxSelect": 1,
							"maxSize": 5242880,
							"protected": false
						}
					},
					{
						"system": false,
						"id": "lxqxcu5n",
						"name": "images",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "ryjreaz08zdji3u",
							"cascadeDelete": true,
							"minSelect": null,
							"maxSelect": null,
							"displayFields": null
						}
					}
				],
				"indexes": [],
				"listRule": "id = @request.auth.id",
				"viewRule": "id = @request.auth.id",
				"createRule": "",
				"updateRule": "id = @request.auth.id",
				"deleteRule": "id = @request.auth.id",
				"options": {
					"allowEmailAuth": true,
					"allowOAuth2Auth": true,
					"allowUsernameAuth": true,
					"exceptEmailDomains": null,
					"manageRule": null,
					"minPasswordLength": 8,
					"onlyEmailDomains": null,
					"onlyVerified": false,
					"requireEmail": false
				}
			}
		]`

		collections := []*models.Collection{}
		if err := json.Unmarshal([]byte(jsonData), &collections); err != nil {
			return err
		}

		return daos.New(db).ImportCollections(collections, true, nil)
	}, func(db dbx.Builder) error {
		return nil
	})
}
