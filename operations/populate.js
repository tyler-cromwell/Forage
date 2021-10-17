print('======================================')
let database = connect('127.0.0.1:27017/forage')

let resultDrop = database.data.drop()
print('Dropped:', resultDrop)

// Production will include expiration date
let dateUpdated = new Date()
let documents = [
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
                value: 2,
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
        attributes: {
            sealed: true
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
        type: 'Ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'piece'
        },
        attributes: {
            sealed: false
        },
        lifespan: {
            pantry: {
                value: 3,
                unit: 'day'
            }
        },
        name: 'Bacon',
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
        attributes: {
            cooked: false,
            refreeze: false,
            sealed: true
        },
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
        type: 'Ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'pound'
        },
        attributes: {
            cooked: true,
            sealed: true
        },
        lifespan: {
            refrigerator: {
                value: 2,
                unit: 'day'
            }
        },
        name: 'Beef',
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
        attributes: {
            lowHumidity: true,
            sealed: true
        },
        lifespan: {
            pantry: {
                value: 10,
                unit: 'days'
            }
        },
        name: 'Bread',
        type: 'Ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'head'
        },
        attributes: {
           wrapped: true
        },
        lifespan: {
            refrigerator: {
                value: 1,
                unit: 'week'
            }
        },
        name: 'Broccoli',
        type: 'Ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'oz'
        },
        attributes: {
            wrapped: true
        },
        lifespan: {
            freezer: {
                value: 1,
                unit: 'year'
            },
            refrigerator: {
                value: 5,
                unit: 'day'
            }
        },
        name: 'Brussel Sprouts',
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
        attributes: {
            cooked: false,
            refreeze: false,
            sealed: true
        },
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
        type: 'Ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'pound'
        },
        attributes: {
            cooked: true,
            sealed: true
        },
        lifespan: {
            refrigerator: {
                value: 2,
                unit: 'day'
            }
        },
        name: 'Chicken',
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
        attributes: {
            sealed: true
        },
        lifespan: {
            refrigerator: {
                value: 6,
                unit: 'year'
            }
        },
        name: 'Hot Sauce',
        type: 'Ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 8,
            unit: 'fl oz'
        },
        attributes: {
            sealed: false
        },
        lifespan: {
            refrigerator: {
                value: 3,
                unit: 'year'
            }
        },
        name: 'Hot Sauce',
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
        attributes: {
            genuine: true,
            sealed: true
        },
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
        type: 'Ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'cup'
        },
        attributes: {
            genuine: true,
            sealed: false
        },
        lifespan: {
            refrigerator: {
                value: 10,
                unit: 'month'
            }
        },
        name: 'Maple Syrup',
        type: 'Ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'cup'
        },
        attributes: {
            sealed: true,
            store: true
        },
        lifespan: {
            pantry: {
                value: 3,
                unit: 'month'
            }
        },
        name: 'Mayonnaise',
        type: 'Ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'cup'
        },
        attributes: {
            sealed: false,
            store: true
        },
        lifespan: {
            refrigerator: {
                value: 6,
                unit: 'weeks'
            }
        },
        name: 'Mayonnaise',
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
        attributes: {
            cooked: false,
            sealed: true
        },
        lifespan: {
            pantry: {
                value: 6,
                unit: 'month'
            }
        },
        name: 'Pasta',
        type: 'Ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'pound'
        },
        attributes: {
            cooked: true,
            sealed: true
        },
        lifespan: {
            refrigerator: {
                value: 3,
                unit: 'day'
            }
        },
        name: 'Pasta',
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
            unit: 'jar'
        },
        attributes: {
            sealed: true
        },
        lifespan: {
            pantry: {
                value: 3,
                unit: 'year'
            },
            refrigerator: {
                value: 3,
                unit: 'year'
            }
        },
        name: 'Pickles',
        type: 'Ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'jar'
        },
        attributes: {
            sealed: false
        },
        lifespan: {
            pantry: {
                value: 2,
                unit: 'year'
            }
        },
        name: 'Pickles',
        type: 'Ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'pound'
        },
        attributes: {
            cooked: false,
            refreeze: false,
            sealed: true
        },
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
        type: 'Ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'pound'
        },
        attributes: {
            cooked: true,
            sealed: true
        },
        lifespan: {
            refrigerator: {
                value: 2,
                unit: 'day'
            }
        },
        name: 'Pork',
        type: 'Ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'jar'
        },
        attributes: {
            sealed: true
        },
        lifespan: {
            refrigerator: {
                value: 3,
                unit: 'week'
            }
        },
        name: 'Ranch Dressing',
        type: 'Ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'cup'
        },
        attributes: {
            cooked: false,
            lowHumidity: true,
            sealed: true
        },
        lifespan: {
            pantry: {
                value: 6,
                unit: 'month'
            }
        },
        name: 'Rice',
        type: 'Ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'cup'
        },
        attributes: {
            cooked: true,
            lowHumidity: true,
            sealed: true
        },
        lifespan: {
            refrigerator: {
                value: 3,
                unit: 'day'
            }
        },
        name: 'Rice',
        type: 'Ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'group'
        },
        attributes: {
            lowHumidity: true,
            sealed: true
        },
        lifespan: {
            pantry: {
                value: 3,
                unit: 'day'
            },
            refrigerator: {
                comment: 'Store in a ziplock bag with paper towels',
                value: 9,
                unit: 'day'
            }
        },
        name: 'Scallions',
        type: 'Ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'oz'
        },
        lifespan: {
            refrigerator: {
                value: 2,
                unit: 'month'
            },
            pantry: {
                value: 6,
                unit: 'week'
            }
        },
        name: 'Shallot',
        type: 'Ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'pound'
        },
        comment: 'Cold cut',
        lifespan: {
            refrigerator: {
                value: 2,
                unit: 'day'
            }
        },
        name: 'Turkey',
        type: 'Ingredient',
        updated: dateUpdated
    }
]

documents.forEach((document) => {
    let lifespan = document.lifespan
    let maxDays = 1
    let maxEnv = ""

    // For each storage environment
    Object.keys(lifespan).forEach((env) => {
        let {value, unit} = lifespan[env]
        let days = value

        // Convert all units to days
        if (unit === 'year') {
            days *= 365
        } else if (unit === 'month') {
            days *= 30
        } else if (unit === 'week') {
            days *= 7
        }

        // Determine the maximum storage time
        if (days === 0) {
            maxDays = days
            maxEnv = env
        } else if (maxDays >= 1 && maxDays < days) {
            maxDays = days
            maxEnv = env
        }
    })

    // Calculate max expiration date
    // Since this is before insertMany, updated is also created date
    let expirationDate = new Date(
        document.updated.getFullYear(),
        document.updated.getMonth(),
        document.updated.getDate()+maxDays
    )

    if (maxDays == 0) {
        expirationDate = new Date(8640000000000000)
    }

    // Update document
    document.storeIn = maxEnv
    document.expirationDate = expirationDate
})

let resultInsertMany = database.data.insertMany(documents)
print('Inserted', documents.length, 'of', resultInsertMany.insertedIds.length, 'documents')
/*
let resultFind = database.data.find()
print('Find:', resultFind)

while (resultFind.hasNext()) {
    print(tojson(resultFind.next()))
}
*/
print('======================================')
