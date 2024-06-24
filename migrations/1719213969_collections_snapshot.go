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
				"id": "_pb_users_auth_",
				"created": "2023-12-06 09:09:21.423Z",
				"updated": "2023-12-06 09:12:40.240Z",
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
					}
				],
				"indexes": [],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {
					"allowEmailAuth": true,
					"allowOAuth2Auth": false,
					"allowUsernameAuth": true,
					"exceptEmailDomains": null,
					"manageRule": null,
					"minPasswordLength": 8,
					"onlyEmailDomains": null,
					"onlyVerified": false,
					"requireEmail": false
				}
			},
			{
				"id": "t9sk3rywb3ap8vt",
				"created": "2023-12-06 09:29:38.045Z",
				"updated": "2023-12-08 13:37:52.587Z",
				"name": "habits",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "djgqz8o0",
						"name": "user",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "avmxjods",
						"name": "name",
						"type": "text",
						"required": true,
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
						"id": "cuzws3g1",
						"name": "type",
						"type": "text",
						"required": true,
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
						"id": "edrekxyy",
						"name": "status",
						"type": "text",
						"required": true,
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
						"id": "ccjvzri6",
						"name": "priority",
						"type": "number",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "qwb4x3gg",
						"name": "streak",
						"type": "number",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"noDecimal": false
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id = user",
				"viewRule": "@request.auth.id = user",
				"createRule": "@request.auth.id = user",
				"updateRule": "@request.auth.id = user && \n(@request.data.user:isset = false || @request.auth.id = @request.data.user)",
				"deleteRule": "@request.auth.id = user",
				"options": {}
			},
			{
				"id": "frl5rotrwv3eesj",
				"created": "2023-12-06 10:05:15.517Z",
				"updated": "2023-12-08 13:37:38.078Z",
				"name": "habit_progress",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "ru4ezxgg",
						"name": "habit",
						"type": "relation",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "t9sk3rywb3ap8vt",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "rmto1bva",
						"name": "user",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "3gllrfeg",
						"name": "date",
						"type": "date",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "ygy2luu5",
						"name": "status",
						"type": "select",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"completed",
								"pending",
								"skipped"
							]
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id = user",
				"viewRule": "@request.auth.id = user",
				"createRule": "@request.auth.id = user",
				"updateRule": "@request.auth.id = user && \n(@request.data.user:isset = false || @request.auth.id = @request.data.user)",
				"deleteRule": "@request.auth.id = user",
				"options": {}
			},
			{
				"id": "w7w8875kihaao80",
				"created": "2023-12-06 10:54:54.166Z",
				"updated": "2023-12-08 13:38:10.901Z",
				"name": "life_balance",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "pgfnlwox",
						"name": "user",
						"type": "relation",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "4xuu1iui",
						"name": "date",
						"type": "date",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "nkgw7pna",
						"name": "relationship",
						"type": "number",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": 1,
							"max": 10,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "b4dtrsoz",
						"name": "health",
						"type": "number",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": 1,
							"max": 10,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "rpotfkho",
						"name": "career",
						"type": "number",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": 1,
							"max": 10,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "szl0uc5f",
						"name": "growth",
						"type": "number",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": 1,
							"max": 10,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "n14uzvlu",
						"name": "life",
						"type": "number",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": 1,
							"max": 10,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "udc4d4vc",
						"name": "social",
						"type": "number",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": 1,
							"max": 10,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "opd7l8op",
						"name": "hobby",
						"type": "number",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": 1,
							"max": 10,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "erzuwcfq",
						"name": "finance",
						"type": "number",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": 1,
							"max": 10,
							"noDecimal": false
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id = user",
				"viewRule": "@request.auth.id = user",
				"createRule": "@request.auth.id = user",
				"updateRule": "@request.auth.id = user && \n(@request.data.user:isset = false || @request.auth.id = @request.data.user)",
				"deleteRule": "@request.auth.id = user",
				"options": {}
			},
			{
				"id": "oduxdvd88nnn3i0",
				"created": "2023-12-06 11:28:34.292Z",
				"updated": "2023-12-08 13:29:52.362Z",
				"name": "daily_log",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "tdggxldt",
						"name": "user",
						"type": "relation",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "g5mw43pj",
						"name": "summary",
						"type": "editor",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"convertUrls": false
						}
					},
					{
						"system": false,
						"id": "2mfocdla",
						"name": "score",
						"type": "number",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": 1,
							"max": 5,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "fezgfzju",
						"name": "date",
						"type": "date",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "cj8vulaj",
						"name": "bath",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "xd5g1m9p",
						"name": "feeling",
						"type": "text",
						"required": false,
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
				"listRule": "@request.auth.id = user",
				"viewRule": "@request.auth.id = user",
				"createRule": "@request.auth.id = user",
				"updateRule": "@request.auth.id = user && \n(@request.data.user:isset = false || @request.auth.id = @request.data.user)",
				"deleteRule": "@request.auth.id = user",
				"options": {}
			},
			{
				"id": "kcxn4npza3pe5w7",
				"created": "2023-12-08 08:09:59.869Z",
				"updated": "2023-12-08 13:36:51.882Z",
				"name": "food_log",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "y1zey4md",
						"name": "user",
						"type": "relation",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "8z7aqaur",
						"name": "name",
						"type": "text",
						"required": true,
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
						"id": "nc3hahpp",
						"name": "field",
						"type": "file",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"mimeTypes": [],
							"thumbs": [],
							"maxSelect": 1,
							"maxSize": 5242880,
							"protected": true
						}
					},
					{
						"system": false,
						"id": "e7vpmfz0",
						"name": "tag",
						"type": "select",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"breakfast",
								"lunch",
								"dinner",
								"snack",
								"extra"
							]
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id = user",
				"viewRule": "@request.auth.id = user",
				"createRule": "@request.auth.id = user",
				"updateRule": "@request.auth.id = user && \n(@request.data.user:isset = false || @request.auth.id = @request.data.user)",
				"deleteRule": "@request.auth.id = user",
				"options": {}
			},
			{
				"id": "am2sdcobu2q3nfw",
				"created": "2023-12-08 08:59:58.207Z",
				"updated": "2023-12-08 13:38:25.567Z",
				"name": "tasks",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "gqjjsfih",
						"name": "user",
						"type": "relation",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "eefatx88",
						"name": "title",
						"type": "text",
						"required": true,
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
						"id": "n8zxg8gu",
						"name": "description",
						"type": "editor",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"convertUrls": false
						}
					},
					{
						"system": false,
						"id": "7mxger5a",
						"name": "due",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "xagxgazf",
						"name": "image",
						"type": "file",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"mimeTypes": [],
							"thumbs": [],
							"maxSelect": 1,
							"maxSize": 5242880,
							"protected": false
						}
					},
					{
						"system": false,
						"id": "rziqit0r",
						"name": "category",
						"type": "select",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"focus",
								"goals",
								"fitin",
								"backburner"
							]
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id = user",
				"viewRule": "@request.auth.id = user",
				"createRule": "@request.auth.id = user",
				"updateRule": "@request.auth.id = user && \n(@request.data.user:isset = false || @request.auth.id = @request.data.user)",
				"deleteRule": "@request.auth.id = user",
				"options": {}
			},
			{
				"id": "z60l06ij6ugtt49",
				"created": "2023-12-08 09:19:10.959Z",
				"updated": "2023-12-11 08:11:24.244Z",
				"name": "track_items",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "amx7xc9k",
						"name": "user",
						"type": "relation",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "phy93evm",
						"name": "track_id",
						"type": "number",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "l2r1wraq",
						"name": "source",
						"type": "text",
						"required": true,
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
						"id": "j5ddjtbl",
						"name": "app",
						"type": "text",
						"required": true,
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
						"id": "alhtvrwk",
						"name": "task_name",
						"type": "text",
						"required": true,
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
						"id": "qsh511ki",
						"name": "title",
						"type": "text",
						"required": true,
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
						"id": "bzuaj2cq",
						"name": "begin_date",
						"type": "date",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "dhdk3194",
						"name": "end_date",
						"type": "date",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					}
				],
				"indexes": [
					"CREATE UNIQUE INDEX ` + "`" + `idx_0t9KXmL` + "`" + ` ON ` + "`" + `track_items` + "`" + ` (\n  ` + "`" + `begin_date` + "`" + `,\n  ` + "`" + `end_date` + "`" + `,\n  ` + "`" + `task_name` + "`" + `\n)"
				],
				"listRule": "@request.auth.id = user",
				"viewRule": "@request.auth.id = user",
				"createRule": "@request.auth.id = user",
				"updateRule": "@request.auth.id = user && \n(@request.data.user:isset = false || @request.auth.id = @request.data.user)",
				"deleteRule": "@request.auth.id = user",
				"options": {}
			},
			{
				"id": "9dpns8heuo4oukt",
				"created": "2023-12-08 09:29:07.002Z",
				"updated": "2023-12-12 11:35:55.160Z",
				"name": "track_upload",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "eo2l6bdp",
						"name": "user",
						"type": "relation",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": []
						}
					},
					{
						"system": false,
						"id": "ytzn6wgn",
						"name": "source",
						"type": "text",
						"required": true,
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
						"id": "ezv92del",
						"name": "file",
						"type": "file",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"mimeTypes": [],
							"thumbs": [],
							"maxSelect": 1,
							"maxSize": 50000000,
							"protected": true
						}
					},
					{
						"system": false,
						"id": "dcrgc9da",
						"name": "synced",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "qk9rbp7x",
						"name": "status",
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
						"id": "py0fx9dz",
						"name": "total_record",
						"type": "number",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "rl05mai2",
						"name": "duplicate_record",
						"type": "number",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"noDecimal": false
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id = user",
				"viewRule": "@request.auth.id = user",
				"createRule": "@request.auth.id = user",
				"updateRule": "@request.auth.id = user && \n(@request.data.user:isset = false || @request.auth.id = @request.data.user)",
				"deleteRule": "@request.auth.id = user",
				"options": {}
			},
			{
				"id": "c4gkvayo41sc93a",
				"created": "2023-12-19 11:04:13.291Z",
				"updated": "2023-12-19 11:24:37.315Z",
				"name": "dev_tokens",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "mem01g3k",
						"name": "user",
						"type": "relation",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "zbzihqll",
						"name": "token",
						"type": "text",
						"required": true,
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
						"id": "idqfi4qs",
						"name": "is_active",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					}
				],
				"indexes": [
					"CREATE UNIQUE INDEX ` + "`" + `idx_myq1gAs` + "`" + ` ON ` + "`" + `dev_tokens` + "`" + ` (\n  ` + "`" + `user` + "`" + `,\n  ` + "`" + `token` + "`" + `\n)"
				],
				"listRule": "@request.auth.id = user",
				"viewRule": "@request.auth.id = user",
				"createRule": "@request.auth.id = user",
				"updateRule": "@request.auth.id = user && \n(@request.data.user:isset = false || @request.auth.id = @request.data.user)",
				"deleteRule": "@request.auth.id = user",
				"options": {}
			},
			{
				"id": "650ed4q6e8zapgc",
				"created": "2023-12-19 11:18:12.435Z",
				"updated": "2023-12-19 11:23:26.715Z",
				"name": "devices",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "lksu8kd6",
						"name": "user",
						"type": "relation",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "kpqvs55v",
						"name": "name",
						"type": "text",
						"required": true,
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
						"id": "l10bwrov",
						"name": "hostname",
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
						"id": "jzpfffio",
						"name": "os",
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
						"id": "rbejzqoo",
						"name": "arch",
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
						"id": "uhdas8ja",
						"name": "is_online",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "unduf27t",
						"name": "is_active",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					}
				],
				"indexes": [
					"CREATE UNIQUE INDEX ` + "`" + `idx_g0z7UR3` + "`" + ` ON ` + "`" + `devices` + "`" + ` (\n  ` + "`" + `user` + "`" + `,\n  ` + "`" + `name` + "`" + `\n)"
				],
				"listRule": "@request.auth.id = user",
				"viewRule": "@request.auth.id = user",
				"createRule": "@request.auth.id = user",
				"updateRule": "@request.auth.id = user && \n(@request.data.user:isset = false || @request.auth.id = @request.data.user)",
				"deleteRule": "@request.auth.id = user",
				"options": {}
			},
			{
				"id": "q2557ntixlaryir",
				"created": "2023-12-27 13:36:25.720Z",
				"updated": "2023-12-27 13:38:22.921Z",
				"name": "music",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "1maiwjxe",
						"name": "user",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "qqxl9s8b",
						"name": "title",
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
						"id": "04ikf8pa",
						"name": "file",
						"type": "file",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"mimeTypes": [],
							"thumbs": [],
							"maxSelect": 1,
							"maxSize": 5242880,
							"protected": false
						}
					},
					{
						"system": false,
						"id": "uu1oknzn",
						"name": "category",
						"type": "text",
						"required": false,
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
				"listRule": "@request.auth.id = user",
				"viewRule": "@request.auth.id = user",
				"createRule": "@request.auth.id = user",
				"updateRule": "@request.auth.id = user && \n(@request.data.user:isset = false || @request.auth.id = @request.data.user)",
				"deleteRule": "@request.auth.id = user",
				"options": {}
			},
			{
				"id": "xkmnl0t8i3afn7q",
				"created": "2023-12-27 13:49:13.902Z",
				"updated": "2023-12-27 13:50:23.740Z",
				"name": "bookshelf",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "8osdkh1z",
						"name": "user",
						"type": "relation",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "cvb0e8de",
						"name": "title",
						"type": "text",
						"required": true,
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
						"id": "sxujtbu8",
						"name": "link",
						"type": "url",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"exceptDomains": null,
							"onlyDomains": null
						}
					},
					{
						"system": false,
						"id": "3tfqf0nl",
						"name": "category",
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
						"id": "xyaqwvok",
						"name": "status",
						"type": "select",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"maxSelect": 1,
							"values": [
								"pending",
								"toread",
								"reading",
								"done"
							]
						}
					},
					{
						"system": false,
						"id": "mhugpgou",
						"name": "image",
						"type": "file",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"mimeTypes": [],
							"thumbs": [],
							"maxSelect": 1,
							"maxSize": 5242880,
							"protected": false
						}
					},
					{
						"system": false,
						"id": "yxdxithl",
						"name": "rating",
						"type": "number",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"noDecimal": false
						}
					},
					{
						"system": false,
						"id": "ogqjr5xs",
						"name": "review",
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
						"id": "bexilw1s",
						"name": "progress",
						"type": "number",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": null,
							"max": null,
							"noDecimal": false
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id = user",
				"viewRule": "@request.auth.id = user",
				"createRule": "@request.auth.id = user",
				"updateRule": "@request.auth.id = user && \n(@request.data.user:isset = false || @request.auth.id = @request.data.user)",
				"deleteRule": "@request.auth.id = user",
				"options": {}
			},
			{
				"id": "wez4w5u1ntn52c1",
				"created": "2024-01-03 21:24:43.592Z",
				"updated": "2024-02-15 11:37:10.924Z",
				"name": "calendar_tokens",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "ygx5v94n",
						"name": "user",
						"type": "relation",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "8d61ck9b",
						"name": "account",
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
						"id": "6ulnkxna",
						"name": "access_token",
						"type": "text",
						"required": true,
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
						"id": "0rwjgaf2",
						"name": "token_type",
						"type": "text",
						"required": true,
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
						"id": "akzkixl2",
						"name": "refresh_token",
						"type": "text",
						"required": true,
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
						"id": "lyhudzsa",
						"name": "expiry",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					}
				],
				"indexes": [],
				"listRule": "@request.auth.id = user",
				"viewRule": "@request.auth.id = user",
				"createRule": "@request.auth.id = user",
				"updateRule": "@request.auth.id = user && \n(@request.data.user:isset = false || @request.auth.id = @request.data.user)",
				"deleteRule": "@request.auth.id = user",
				"options": {}
			},
			{
				"id": "jwn6wr5k0kzlglc",
				"created": "2024-02-15 11:44:42.647Z",
				"updated": "2024-02-15 11:44:42.647Z",
				"name": "calendar_sync",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "p0xcn1wi",
						"name": "user",
						"type": "relation",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "lo27ux9m",
						"name": "token",
						"type": "relation",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "wez4w5u1ntn52c1",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "tpv5mdoe",
						"name": "name",
						"type": "text",
						"required": true,
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
						"id": "qfppknnf",
						"name": "type",
						"type": "text",
						"required": true,
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
						"id": "tp62qixu",
						"name": "sync_token",
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
						"id": "fv7r21ub",
						"name": "is_active",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					},
					{
						"system": false,
						"id": "fx2dsbbs",
						"name": "last_synced",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
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
			},
			{
				"id": "c48cnpzhq9fci8h",
				"created": "2024-02-15 13:36:08.008Z",
				"updated": "2024-06-14 11:47:58.330Z",
				"name": "calendar_events",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "i66yhxn5",
						"name": "calendar_id",
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
						"id": "qbh7t792",
						"name": "calendar_uid",
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
						"id": "sn7lzm48",
						"name": "user",
						"type": "relation",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "t7zh9tvo",
						"name": "calendar",
						"type": "relation",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "jwn6wr5k0kzlglc",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "msfht7ky",
						"name": "etag",
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
						"id": "dsfdeiz1",
						"name": "summary",
						"type": "text",
						"required": true,
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
						"id": "12etv9qi",
						"name": "description",
						"type": "text",
						"required": true,
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
						"id": "fwdyeh1e",
						"name": "event_type",
						"type": "text",
						"required": true,
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
						"id": "h7ambt52",
						"name": "start",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "beeuyogs",
						"name": "end",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "8drf9mbk",
						"name": "creator",
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
						"id": "mkgyxzix",
						"name": "creator_email",
						"type": "email",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"exceptDomains": null,
							"onlyDomains": null
						}
					},
					{
						"system": false,
						"id": "ubvpbhhz",
						"name": "organizer",
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
						"id": "4tzc0iyr",
						"name": "organizer_email",
						"type": "email",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"exceptDomains": null,
							"onlyDomains": null
						}
					},
					{
						"system": false,
						"id": "3ce0gjqz",
						"name": "kind",
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
						"id": "xzllxdx6",
						"name": "location",
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
						"id": "cahec8az",
						"name": "status",
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
						"id": "6ywijhx7",
						"name": "event_created",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "aywutjkr",
						"name": "event_updated",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
						}
					},
					{
						"system": false,
						"id": "vbzksp65",
						"name": "is_day_event",
						"type": "bool",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {}
					}
				],
				"indexes": [
					"CREATE UNIQUE INDEX ` + "`" + `idx_5AFbbtM` + "`" + ` ON ` + "`" + `calendar_events` + "`" + ` (\n  ` + "`" + `user` + "`" + `,\n  ` + "`" + `calendar_id` + "`" + `,\n  ` + "`" + `calendar_uid` + "`" + `\n)"
				],
				"listRule": null,
				"viewRule": null,
				"createRule": null,
				"updateRule": null,
				"deleteRule": null,
				"options": {}
			},
			{
				"id": "bu8wqeaacmnvshf",
				"created": "2024-06-19 08:31:50.998Z",
				"updated": "2024-06-20 09:52:27.889Z",
				"name": "fold_tokens",
				"type": "base",
				"system": false,
				"schema": [
					{
						"system": false,
						"id": "2btpwsuo",
						"name": "user",
						"type": "relation",
						"required": true,
						"presentable": false,
						"unique": false,
						"options": {
							"collectionId": "_pb_users_auth_",
							"cascadeDelete": false,
							"minSelect": null,
							"maxSelect": 1,
							"displayFields": null
						}
					},
					{
						"system": false,
						"id": "t5uhjgd5",
						"name": "phone",
						"type": "text",
						"required": true,
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
						"id": "ipqxwh2u",
						"name": "uuid",
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
						"id": "6jokzu3k",
						"name": "first_name",
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
						"id": "s5noawjm",
						"name": "last_name",
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
						"id": "aiwmnczk",
						"name": "email",
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
						"id": "nawdxn6s",
						"name": "access_token",
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
						"id": "ia81f0b2",
						"name": "refresh_token",
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
						"id": "t5ppu6vc",
						"name": "user_agent",
						"type": "text",
						"required": true,
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
						"id": "ijkew4en",
						"name": "device_hash",
						"type": "text",
						"required": true,
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
						"id": "splhifeg",
						"name": "device_location",
						"type": "text",
						"required": true,
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
						"id": "fvkoqya0",
						"name": "device_name",
						"type": "text",
						"required": true,
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
						"id": "n89us3ku",
						"name": "device_type",
						"type": "text",
						"required": true,
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
						"id": "jffxsfxu",
						"name": "expires_at",
						"type": "date",
						"required": false,
						"presentable": false,
						"unique": false,
						"options": {
							"min": "",
							"max": ""
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
