package migrations

import (
	"encoding/json"

	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
	m "github.com/pocketbase/pocketbase/migrations"
	"github.com/pocketbase/pocketbase/models/schema"
)

func init() {
	m.Register(func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("ryjreaz08zdji3u")
		if err != nil {
			return err
		}

		// add
		new_followups := &schema.SchemaField{}
		if err := json.Unmarshal([]byte(`{
			"system": false,
			"id": "zb0eec7o",
			"name": "followups",
			"type": "relation",
			"required": true,
			"presentable": false,
			"unique": false,
			"options": {
				"collectionId": "w879h6ggcfuixkr",
				"cascadeDelete": true,
				"minSelect": null,
				"maxSelect": null,
				"displayFields": null
			}
		}`), new_followups); err != nil {
			return err
		}
		collection.Schema.AddField(new_followups)

		return dao.SaveCollection(collection)
	}, func(db dbx.Builder) error {
		dao := daos.New(db);

		collection, err := dao.FindCollectionByNameOrId("ryjreaz08zdji3u")
		if err != nil {
			return err
		}

		// remove
		collection.Schema.RemoveField("zb0eec7o")

		return dao.SaveCollection(collection)
	})
}
