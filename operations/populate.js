print('======================================')
let database = connect('127.0.0.1:27017/forage')

let resultDrop = database.data.drop()
print('Dropped:', resultDrop)

// Production will include expiration date
let dateUpdated = new Date()
let resultInsertMany = database.data.insertMany(
    [
        {
            amount: {
                value: 5,
                unit: 'count'
            },
            lifespan: {
                refrigerator: {
                    value: 1,
                    unit: 'week'
                }
            },
            name: 'Apples',
            type: "Ingredient",
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'count'
            },
            lifespan: {
                refrigerator: {
                    value: 1,
                    unit: 'month'
                },
                shelf: {
                    value: 3,
                    unit: 'day'
                }
            },
            name: 'Bagel',
            type: "Ingredient",
            updated: dateUpdated
        },
        {
            amount: {
                value: 4,
                unit: 'count'
            },
            lifespan: {
                shelf: {
                    value: 4,
                    unit: 'day'
                }
            },
            name: 'Bananas',
            type: "Ingredient",
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'loaf'
            },
            lifespan: {
                shelf: {
                    value: 10,
                    unit: 'days'
                }
            },
            lowHumidity: true,
            name: 'Bread',
            sealed: true,
            type: "Ingredient",
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'stick'
            },
            lifespan: {
                refrigerator: {
                    value: 3,
                    unit: 'month'
                }
            },
            name: 'Butter',
            type: "Ingredient",
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'count'
            },
            lifespan: {
                refrigerator: {
                    value: 1,
                    unit: 'month'
                }
            },
            name: 'Carrots',
            type: "Ingredient",
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'count'
            },
            lifespan: {
                refrigerator: {
                    value: 2,
                    unit: 'week'
                }
            },
            name: 'Celery',
            type: "Ingredient",
            updated: dateUpdated
        },
        {
            amount: {
                value: 12,
                unit: 'count'
            },
            lifespan: {
                refrigerator: {
                    value: 1,
                    unit: 'month'
                }
            },
            name: 'Eggs',
            type: "Ingredient",
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'clove'
            },
            lifespan: {
                refrigerator: {
                    value: 2,
                    unit: 'month'
                }
            },
            name: 'Garlic',
            type: "Ingredient",
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'piece'
            },
            lifespan: {
                refrigerator: {
                    value: 2,
                    unit: 'day'
                }
            },
            name: 'Gyoza',
            type: "Meal",
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'jar'
            },
            lifespan: {
                shelf: {
                    value: 2,
                    unit: 'month'
                }
            },
            name: 'Jelly',
            type: "Ingredient",
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'head'
            },
            lifespan: {
                refrigerator: {
                    value: 1,
                    unit: 'week'
                }
            },
            name: 'Lettuce',
            type: "Ingredient",
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'pound'
            },
            cooked: false,
            lifespan: {
                freezer: {
                    value: 3,
                    unit: 'month'
                },
                refrigerator: {
                    value: 2,
                    unit: 'day'
                }
            },
            name: 'Meat',
            refreeze: false,
            sealed: true,
            type: "Ingredient",
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'pound'
            },
            cooked: true,
            lifespan: {
                refrigerator: {
                    value: 2,
                    unit: 'day'
                }
            },
            name: 'Meat',
            sealed: true,
            type: "Ingredient",
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'gallon'
            },
            lifespan: {
                refrigerator: {
                    value: 10,
                    unit: 'day'
                }
            },
            name: 'Milk',
            type: "Ingredient",
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'head'
            },
            lifespan: {
                refrigerator: {
                    value: 2,
                    unit: 'week'
                }
            },
            name: 'Onion',
            type: "Ingredient",
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'box'
            },
            cooked: false,
            lifespan: {
                shelf: {
                    value: 6,
                    unit: 'month'
                }
            },
            name: 'Pasta',
            sealed: true,
            type: "Ingredient",
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'pound'
            },
            cooked: true,
            lifespan: {
                refrigerator: {
                    value: 3,
                    unit: 'day'
                }
            },
            name: 'Pasta',
            sealed: true,
            type: "Ingredient",
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'jar'
            },
            lifespan: {
                shelf: {
                    value: 4,
                    unit: 'month'
                }
            },
            name: 'Peanut Butter',
            type: "Ingredient",
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'jar'
            },
            lifespan: {
                refrigerator: {
                    value: 3,
                    unit: 'week'
                }
            },
            name: 'Ranch Dressing',
            sealed: true,
            type: "Ingredient",
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'cup'
            },
            cooked: false,
            lifespan: {
                shelf: {
                    value: 6,
                    unit: 'month'
                }
            },
            lowHumidity: true,
            name: 'Rice',
            sealed: true,
            type: "Ingredient",
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'cup'
            },
            cooked: true,
            lifespan: {
                refrigerator: {
                    value: 3,
                    unit: 'day'
                }
            },
            lowHumidity: true,
            name: 'Rice',
            sealed: true,
            type: "Ingredient",
            updated: dateUpdated
        }
    ]
)
print('InsertMany:', tojson(resultInsertMany))
/*
let resultFind = database.data.find()
print('Find:', resultFind)

while (resultFind.hasNext()) {
    print(tojson(resultFind.next()))
}
*/
print('======================================')
