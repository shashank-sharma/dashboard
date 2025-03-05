/// <reference path="../pb_typings.js" />

migrate((db) => {
    const collection = new Collection({
        name: 'food_log',
        type: 'base',
        system: false,
        schema: [
            {
                system: false,
                id: 'fsrqwdve',
                name: 'user',
                type: 'relation',
                required: true,
                presentable: false,
                unique: false,
                options: {
                    collectionId: '_pb_users_auth_',
                    cascadeDelete: true,
                    minSelect: null,
                    maxSelect: 1,
                    displayFields: []
                }
            },
            {
                system: false,
                id: 'nhoqwgdw',
                name: 'name',
                type: 'text',
                required: true,
                presentable: true,
                unique: false,
                options: {
                    min: 1,
                    max: 200,
                    pattern: ''
                }
            },
            {
                system: false,
                id: 'sdqvwetr',
                name: 'image',
                type: 'file',
                required: true,
                presentable: false,
                unique: false,
                options: {
                    maxSelect: 1,
                    maxSize: 5242880,
                    mimeTypes: [
                        'image/jpeg',
                        'image/png',
                        'image/gif',
                        'image/webp'
                    ],
                    thumbs: ['100x100', '300x300', '600x600']
                }
            },
            {
                system: false,
                id: 'yiwertus',
                name: 'tag',
                type: 'select',
                required: true,
                presentable: false,
                unique: false,
                options: {
                    maxSelect: 1,
                    values: [
                        'breakfast',
                        'lunch',
                        'dinner',
                        'snack',
                        'dessert',
                        'drink'
                    ]
                }
            },
            {
                system: false,
                id: 'tzkoqpwe',
                name: 'date',
                type: 'date',
                required: true,
                presentable: false,
                unique: false,
                options: {
                    min: "",
                    max: ""
                }
            }
        ],
        indexes: [
            [
                'created',
                'DESC'
            ],
            [
                'user',
                'created',
                'DESC'
            ],
            [
                'tag',
                'created',
                'DESC'
            ],
            [
                'date',
                'DESC'
            ]
        ],
        listRule: "user.id = @request.auth.id",
        viewRule: "user.id = @request.auth.id",
        createRule: "user.id = @request.auth.id",
        updateRule: "user.id = @request.auth.id",
        deleteRule: "user.id = @request.auth.id",
    });

    return Dao(db).saveCollection(collection);
}, (db) => {
    const dao = new Dao(db);
    const collection = dao.findCollectionByNameOrId("food_log");

    return dao.deleteCollection(collection);
}); 