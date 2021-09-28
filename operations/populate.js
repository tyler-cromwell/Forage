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
            type: 'Ingredient',
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'count'
            },
            lifespan: {
                pantry: {
                    value: 3,
                    unit: 'day'
                },
                refrigerator: {
                    value: 1,
                    unit: 'month'
                }
            },
            name: 'Bagel',
            type: 'Ingredient',
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'piece'
            },
            lifespan: {
                freezer: {
                    value: 8,
                    unit: 'month'
                },
                refrigerator: {
                    value: 2,
                    unit: 'week'
                }
            },
            name: 'Bacon',
            sealed: true,
            type: 'Ingredient',
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'piece'
            },
            lifespan: {
                pantry: {
                    value: 3,
                    unit: 'day'
                }
            },
            name: 'Bacon',
            sealed: false,
            type: 'Ingredient',
            updated: dateUpdated
        },
        {
            amount: {
                value: 4,
                unit: 'count'
            },
            lifespan: {
                pantry: {
                    value: 4,
                    unit: 'day'
                }
            },
            name: 'Bananas',
            type: 'Ingredient',
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
            name: 'Beef',
            refreeze: false,
            sealed: true,
            type: 'Ingredient',
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
            name: 'Beef',
            sealed: true,
            type: 'Ingredient',
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'oz'
            },
            lifespan: {
                pantry: {
                    value: 3,
                    unit: 'year'
                }
            },
            name: 'Black Pepper',
            type: 'Ingredient',
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'loaf'
            },
            lifespan: {
                pantry: {
                    value: 10,
                    unit: 'days'
                }
            },
            lowHumidity: true,
            name: 'Bread',
            sealed: true,
            type: 'Ingredient',
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
            type: 'Ingredient',
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
            type: 'Ingredient',
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
            type: 'Ingredient',
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'box'
            },
            lifespan: {
                refrigerator: {
                    value: 10,
                    unit: 'month'
                }
            },
            name: 'Cheerios',
            type: 'Ingredient',
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
            name: 'Chicken',
            refreeze: false,
            sealed: true,
            type: 'Ingredient',
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
            name: 'Chicken',
            sealed: true,
            type: 'Ingredient',
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
            type: 'Ingredient',
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
            type: 'Ingredient',
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
            type: 'Meal',
            updated: dateUpdated
        },
        {
            amount: {
                value: 8,
                unit: 'fl oz'
            },
            lifespan: {
                refrigerator: {
                    value: 6,
                    unit: 'year'
                }
            },
            name: 'Hot Sauce',
            sealed: true,
            type: 'Ingredient',
            updated: dateUpdated
        },
        {
            amount: {
                value: 8,
                unit: 'fl oz'
            },
            lifespan: {
                refrigerator: {
                    value: 3,
                    unit: 'year'
                }
            },
            name: 'Hot Sauce',
            sealed: false,
            type: 'Ingredient',
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'jar'
            },
            lifespan: {
                pantry: {
                    value: 2,
                    unit: 'month'
                }
            },
            name: 'Jelly',
            type: 'Ingredient',
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'count'
            },
            lifespan: {
                pantry: {
                    value: 0,
                    unit: 'day'
                }
            },
            name: 'Kosher Salt',
            type: 'Ingredient',
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
            type: 'Ingredient',
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'cup'
            },
            genuine: true,
            lifespan: {
                freezer: {
                    value: 0,
                    unit: 'day'
                },
                pantry: {
                    value: 1,
                    unit: 'year'
                }
            },
            name: 'Maple Syrup',
            sealed: true,
            type: 'Ingredient',
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'cup'
            },
            genuine: true,
            lifespan: {
                refrigerator: {
                    value: 10,
                    unit: 'month'
                }
            },
            name: 'Maple Syrup',
            sealed: false,
            type: 'Ingredient',
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'cup'
            },
            lifespan: {
                pantry: {
                    value: 3,
                    unit: 'month'
                }
            },
            name: 'Mayonnaise',
            sealed: true,
            store: true,
            type: 'Ingredient',
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'cup'
            },
            lifespan: {
                refrigerator: {
                    value: 6,
                    unit: 'weeks'
                }
            },
            name: 'Mayonnaise',
            sealed: false,
            store: true,
            type: 'Ingredient',
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
            type: 'Ingredient',
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'cup'
            },
            lifespan: {
                refrigerator: {
                    value: 1,
                    unit: 'week'
                }
            },
            name: 'Napa Cabbage',
            type: 'Ingredient',
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
            type: 'Ingredient',
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'box'
            },
            cooked: false,
            lifespan: {
                pantry: {
                    value: 6,
                    unit: 'month'
                }
            },
            name: 'Pasta',
            sealed: true,
            type: 'Ingredient',
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
            type: 'Ingredient',
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'jar'
            },
            lifespan: {
                pantry: {
                    value: 4,
                    unit: 'month'
                }
            },
            name: 'Peanut Butter',
            type: 'Ingredient',
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
            name: 'Pork',
            refreeze: false,
            sealed: true,
            type: 'Ingredient',
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
            name: 'Pork',
            sealed: true,
            type: 'Ingredient',
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
            type: 'Ingredient',
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'cup'
            },
            cooked: false,
            lifespan: {
                pantry: {
                    value: 6,
                    unit: 'month'
                }
            },
            lowHumidity: true,
            name: 'Rice',
            sealed: true,
            type: 'Ingredient',
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
            type: 'Ingredient',
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'group'
            },
            lifespan: {
                pantry: {
                    value: 3,
                    unit: 'day'
                },
                refrigerator: {
                    comment: 'Store in a ziplock bag with paper towels',
                    value: 10,
                    unit: 'day'
                }
            },
            lowHumidity: true,
            name: 'Scallions',
            sealed: true,
            type: 'Ingredient',
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
