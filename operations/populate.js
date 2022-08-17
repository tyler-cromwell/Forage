print('======================================')
let database = connect('127.0.0.1:27017/forage')

let resultDrop = database.data.drop()
print('Dropped:', resultDrop)

// Production will include expiration date
let dateUpdated = new Date().getTime()
let documents = [
    {
        amount: {
            unit: 'ounces',
            value: 10.25
        },
        attributes: {
            family: 'sauce',
            flavor: 'Roasted Garlic',
            opened: false
        },
        lifespan: {
            pantry: {
                unit: 'month',
                value: 4
            }
        },
        name: 'Aioli',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'ounces',
            value: 10.25
        },
        attributes: {
            family: 'sauce',
            flavor: 'Roasted Garlic',
            opened: true
        },
        lifespan: {
            refrigerator: {
                unit: 'day',
                value: 4
            }
        },
        name: 'Aioli',
        type: 'ingredient',
        updated: dateUpdated
    },
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
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'piece'
        },
        attributes: {
            opened: false
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
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'pound'
        },
        attributes: {
            opened: true
        },
        lifespan: {
            refrigerator: {
                value: 3,
                unit: 'day'
            }
        },
        name: 'Bacon',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'count'
        },
        lifespan: {
            freezer: {
                value: 1,
                unit: 'year'
            },
            pantry: {
                value: 2,
                unit: 'day'
            },
            refrigerator: {
                value: 1,
                unit: 'month'
            },
        },
        name: 'Bagel',
        type: 'ingredient',
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
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'pound'
        },
        attributes: {
            cooked: false,
            opened: false,
            refreeze: false
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
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'pound'
        },
        attributes: {
            cooked: true,
            opened: false
        },
        lifespan: {
            refrigerator: {
                value: 2,
                unit: 'day'
            }
        },
        name: 'Beef',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'ounce'
        },
        lifespan: {
            pantry: {
                value: 3,
                unit: 'year'
            }
        },
        name: 'Black Pepper',
        type: 'ingredient',
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
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'ounce'
        },
        attributes: {
            wrapped: true
        },
        lifespan: {
            freezer: {
                value: 3,
                unit: 'month'
            },
            refrigerator: {
                value: 5,
                unit: 'day'
            }
        },
        name: 'Brussel Sprouts',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'pound',
            value: 1
        },
        attributes: {
            cooked: false,
            opened: false,
        },
        lifespan: {
            pantry: {
                unit: 'year',
                value: 3
            }
        },
        name: 'Bucatini',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'pound',
            value: 1
        },
        attributes: {
            cooked: false,
            opened: true,
        },
        lifespan: {
            pantry: {
                unit: 'year',
                value: 1
            }
        },
        name: 'Bucatini',
        type: 'ingredient',
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
        type: 'ingredient',
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
        type: 'ingredient',
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
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'box'
        },
        lifespan: {
            pantry: {
                value: 10,
                unit: 'month'
            }
        },
        name: 'Cheerios',
        type: 'ingredient',
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
            opened: false
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
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'pound'
        },
        attributes: {
            cooked: true,
            opened: false
        },
        lifespan: {
            refrigerator: {
                value: 2,
                unit: 'day'
            }
        },
        name: 'Chicken',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'ounce'
        },
        lifespan: {
            pantry: {
                value: 9,
                unit: 'month'
            }
        },
        name: 'Chili Flakes',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'ounce',
            value: 1
        },
        attributes: {
            opened: false
        },
        lifespan: {
            freezer: {
                unit: 'month',
                value: 2
            }
        },
        name: 'Cream Cheese',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'ounce',
            value: 1
        },
        attributes: {
            opened: true
        },
        lifespan: {
            refrigerator: {
                value: 10,
                unit: 'day'
            }
        },
        name: 'Cream Cheese',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'ounce'
        },
        lifespan: {
            refrigerator: {
                value: 10,
                unit: 'day'
            }
        },
        name: 'Cucumber',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'ounce'
        },
        lifespan: {
            pantry: {
                value: 3,
                unit: 'year'
            }
        },
        name: 'Cumin',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'ounce'
        },
        attributes: {
            opened: false
        },
        lifespan: {
            pantry: {
                value: 2,
                unit: 'year'
            }
        },
        name: 'Dijon Mustard',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'ounce'
        },
        attributes: {
            opened: true
        },
        lifespan: {
            refrigerator: {
                value: 1,
                unit: 'year'
            }
        },
        name: 'Dijon Mustard',
        type: 'ingredient',
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
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'head'
        },
        lifespan: {
            pantry: {
                unit: 'month',
                value: 3
            }
        },
        name: 'Garlic',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: "pound",
            value: 1.0
        },
        lifespan: {
            freezer: {
                unit: "month",
                value: 3
            },
            refrigerator: {
                unit: "month",
                value: 1
            }
        },
        name: "Gnocchi",
        type: "ingredient",
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'ounce',
            value: 10
        },
        lifespan: {
            refrigerator: {
                unit: 'day',
                value: 5
            }
        },
        name: 'Grape Tomatoes',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'cup'
        },
        attributes: {
            opened: false
        },
        lifespan: {
            freezer: {
                value: 2,
                unit: 'month'
            },
            refrigerator: {
                value: 2,
                unit: 'week'
            }
        },
        name: 'Greek Yogurt',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'cup'
        },
        attributes: {
            opened: true
        },
        lifespan: {
            freezer: {
                value: 1,
                unit: 'month'
            },
            refrigerator: {
                value: 1,
                unit: 'week'
            }
        },
        name: 'Greek Yogurt',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'head',
            value: 1
        },
        lifespan: {
            refrigerator: {
                unit: 'day',
                value: 12
            }
        },
        name: 'Green Leaf Lettuce',
        type: 'ingredient',
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
            unit: 'ounce',
            value: 7.75
        },
        lifespan: {
            pantry: {
                unit: 'year',
                value: 2
            }
        },
        name: 'Hamburger Seasoning',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 5,
            unit: 'fluid ounce'
        },
        attributes: {
            brand: 'The Heatonist',
            family: 'sauce',
            flavor: 'Los Calientes',
            opened: false
        },
        lifespan: {
            pantry: {
                unit: 'year',
                value: 6
            }
        },
        name: 'Hot Sauce',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 5,
            unit: 'fluid ounce'
        },
        attributes: {
            brand: 'The Heatonist',
            family: 'sauce',
            flavor: 'Los Calientes',
            opened: true
        },
        lifespan: {
            refrigerator: {
                unit: 'year',
                value: 2
            }
        },
        name: 'Hot Sauce',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'fluid ounce',
            value: 5
        },
        attributes: {
            brand: 'The Spicy Shark',
            family: 'sauce',
            flavor: 'Original Habanero',
            opened: false
        },
        lifespan: {
            pantry: {
                unit: 'year',
                value: 6
            }
        },
        name: 'Hot Sauce',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'fluid ounce',
            value: 5
        },
        attributes: {
            brand: 'The Spicy Shark',
            family: 'sauce',
            flavor: 'Original Habanero',
            opened: true
        },
        lifespan: {
            refrigerator: {
                unit: 'year',
                value: 2
            }
        },
        name: 'Hot Sauce',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 12,
            unit: 'ounce'
        },
        attributes: {
            brand: "Smucker's",
            flavor: 'Grape',
            opened: false
        },
        lifespan: {
            refrigerator: {
                unit: 'year',
                value: 2
            },
            pantry: {
                unit: 'year',
                value: 2
            }
        },
        name: 'Jam',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 12,
            unit: 'ounce'
        },
        attributes: {
            brand: "Smucker's",
            flavor: 'Grape',
            opened: true
        },
        lifespan: {
            refrigerator: {
                value: 8,
                unit: 'month'
            }
        },
        name: 'Jam',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 12,
            unit: 'ounce'
        },
        attributes: {
            brand: "Smucker's",
            flavor: 'Raspberry',
            opened: false
        },
        lifespan: {
            refrigerator: {
                unit: 'month',
                value: 20
            },
            pantry: {
                unit: 'month',
                value: 20
            }
        },
        name: 'Jam',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 12,
            unit: 'ounce'
        },
        attributes: {
            brand: "Smucker's",
            flavor: 'Raspberry',
            opened: true
        },
        lifespan: {
            refrigerator: {
                value: 8,
                unit: 'month'
            }
        },
        name: 'Jam',
        type: 'ingredient',
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
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'pound',
            value: 1
        },
        attributes: {
            cooked: false,
            opened: false,
        },
        lifespan: {
            pantry: {
                unit: 'year',
                value: 3
            }
        },
        name: 'Linguine',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'pound',
            value: 1
        },
        attributes: {
            cooked: false,
            opened: true,
        },
        lifespan: {
            pantry: {
                unit: 'year',
                value: 1
            }
        },
        name: 'Linguine',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'cup',
            value: 1
        },
        attributes: {
            genuine: true,
            opened: false
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
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'cup',
            value: 1
        },
        attributes: {
            genuine: true,
            opened: true
        },
        lifespan: {
            refrigerator: {
                value: 10,
                unit: 'month'
            }
        },
        name: 'Maple Syrup',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'cup',
            value: 1
        },
        attributes: {
            opened: false,
            store: true
        },
        lifespan: {
            pantry: {
                value: 3,
                unit: 'month'
            }
        },
        name: 'Mayonnaise',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'cup',
            value: 1
        },
        attributes: {
            opened: true,
            store: true
        },
        lifespan: {
            refrigerator: {
                value: 6,
                unit: 'weeks'
            }
        },
        name: 'Mayonnaise',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'gallon',
            value: 1
        },
        lifespan: {
            refrigerator: {
                value: 12,
                unit: 'day'
            }
        },
        name: 'Milk',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'fluid ounce',
            value: 1
        },
        lifespan: {
            pantry: {
                value: 3,
                unit: 'month'
            }
        },
        name: 'Mirin',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'pound',
            value: 1
        },
        attributes: {
            form: 'Loaf'
        },
        lifespan: {
            refrigerator: {
                unit: 'week',
                value: 3
            }
        },
        name: 'Mozzarella Cheese',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'ounce',
            value: 1
        },
        attributes: {
            opened: false
        },
        lifespan: {
            refrigerator: {
                value: 6,
                unit: 'week'
            }
        },
        name: 'Mozzarella Cheese',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'ounce',
            value: 1
        },
        attributes: {
            opened: true
        },
        lifespan: {
            refrigerator: {
                value: 5,
                unit: 'day'
            }
        },
        name: 'Mozzarella Cheese',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'ounce',
            value: 1
        },
        lifespan: {
            refrigerator: {
                value: 3,
                unit: 'week'
            }
        },
        name: 'Muenster Cheese',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'cup',
            value: 1
        },
        lifespan: {
            refrigerator: {
                value: 1,
                unit: 'week'
            }
        },
        name: 'Napa Cabbage',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'ounce',
            value: 1
        },
        lifespan: {
            pantry: {
                value: 18,
                unit: 'day'
            }
        },
        name: 'Nori',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'fluid ounces',
            value: 1
        },
        attributes: {
            type: 'Extra Virgin',
            opened: false
        },
        lifespan: {
            pantry: {
                unit: 'month',
                value: 20
            }
        },
        name: 'Olive Oil',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'fluid ounces',
            value: 1
        },
        attributes: {
            type: 'Extra Virgin',
            opened: true
        },
        lifespan: {
            pantry: {
                unit: 'month',
                value: 6
            }
        },
        name: 'Olive Oil',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'head',
            value: 1
        },
        attributes: {
            "cut": false,
            "type": "Red"
        },
        lifespan: {
            freezer: {
                value: 8,
                unit: 'month'
            },
            pantry: {
                value: 2,
                unit: 'month'
            }
        },
        name: 'Onion',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'box',
            value: 1
        },
        lifespan: {
            pantry: {
                value: 6,
                unit: 'month'
            }
        },
        name: 'Panko Bread Crumbs',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'pound',
            value: 1
        },
        attributes: {
            cooked: true,
            opened: true
        },
        lifespan: {
            refrigerator: {
                value: 3,
                unit: 'day'
            }
        },
        name: 'Pasta',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'ounce',
            value: 1
        },
        attributes: {
            opened: false
        },
        lifespan: {
            refrigerator: {
                unit: 'month',
                value: 3
            }
        },
        name: 'Parmesan Cheese',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'ounce',
            value: 1
        },
        attributes: {
            opened: true
        },
        lifespan: {
            refrigerator: {
                unit: 'week',
                value: 6
            }
        },
        name: 'Parmesan Cheese',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'jar',
            value: 1
        },
        lifespan: {
            pantry: {
                value: 4,
                unit: 'month'
            }
        },
        name: 'Peanut Butter',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            unit: 'jar',
            value: 1
        },
        lifespan: {
            pantry: {
                value: 18,
                unit: 'month'
            }
        },
        name: 'Pickles',
        type: 'ingredient',
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
            opened: false
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
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'pound'
        },
        attributes: {
            cooked: true,
            opened: false
        },
        lifespan: {
            refrigerator: {
                value: 2,
                unit: 'day'
            }
        },
        name: 'Pork',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'jar'
        },
        attributes: {
            opened: false
        },
        lifespan: {
            refrigerator: {
                value: 3,
                unit: 'week'
            }
        },
        name: 'Ranch Dressing',
        type: 'ingredient',
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
            opened: false
        },
        lifespan: {
            pantry: {
                value: 6,
                unit: 'month'
            }
        },
        name: 'Rice',
        type: 'ingredient',
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
            opened: false
        },
        lifespan: {
            refrigerator: {
                value: 3,
                unit: 'day'
            }
        },
        name: 'Rice',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'cup'
        },
        lifespan: {
            pantry: {
                value: 2,
                unit: 'year'
            }
        },
        name: 'Rice Vinegar',
        type: 'ingredient',
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
        name: 'Salt',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'group'
        },
        attributes: {
            lowHumidity: true,
            opened: false
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
        type: 'ingredient',
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
                unit: 'month'
            },
            pantry: {
                value: 6,
                unit: 'week'
            }
        },
        name: 'Shallot',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 30,
            unit: 'fluid ounce'
        },
        attributes: {
            opened: false
        },
        lifespan: {
            pantry: {
                value: 3,
                unit: 'year'
            }
        },
        name: 'Soy Sauce',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'fluid ounce'
        },
        attributes: {
            opened: true
        },
        lifespan: {
            refrigerator: {
                value: 1,
                unit: 'month'
            }
        },
        name: 'Soy Sauce',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'ounce'
        },
        lifespan: {
            pantry: {
                value: 2,
                unit: 'year'
            }
        },
        name: 'Sugar',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'pound'
        },
        attributes: {
            type: 'Cold cut'
        },
        lifespan: {
            refrigerator: {
                value: 2,
                unit: 'day'
            }
        },
        name: 'Turkey',
        type: 'ingredient',
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
            opened: false
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
        name: 'Turkey',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'pound'
        },
        attributes: {
            cooked: true,
            opened: false
        },
        lifespan: {
            refrigerator: {
                value: 2,
                unit: 'day'
            }
        },
        name: 'Turkey',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'ounce'
        },
        lifespan: {
            freezer: {
                value: 15,
                unit: 'day'
            },
            refrigerator: {
                value: 3,
                unit: 'day'
            }
        },
        name: 'Unagi',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 48,
            unit: 'fluid ounce'
        },
        lifespan: {
            pantry: {
                value: 2,
                unit: 'year'
            }
        },
        name: 'Vegetable Oil',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'loaf'
        },
        attributes: {
            lowHumidity: true
        },
        lifespan: {
            pantry: {
                value: 14,
                unit: 'days'
            }
        },
        name: 'White Bread',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'loaf'
        },
        attributes: {
            lowHumidity: true
        },
        lifespan: {
            pantry: {
                value: 14,
                unit: 'days'
            }
        },
        name: 'Whole Wheat Bread',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'cup'
        },
        attributes: {
            opened: false
        },
        lifespan: {
            pantry: {
                value: 0,
                unit: 'year'
            }
        },
        name: 'Worcestershire Sauce',
        type: 'ingredient',
        updated: dateUpdated
    },
    {
        amount: {
            value: 1,
            unit: 'cup'
        },
        attributes: {
            opened: true
        },
        lifespan: {
            pantry: {
                value: 3,
                unit: 'year'
            }
        },
        name: 'Worcestershire Sauce',
        type: 'ingredient',
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
    let updatedDate = new Date(document.updated)
    let expirationDate = new Date(
        updatedDate.getFullYear(),
        updatedDate.getMonth(),
        updatedDate.getDate()+maxDays
    ).getTime()

    if (document.expirationDate === undefined && maxDays == 0) {
        expirationDate = 0
    }

    // Update document
    document.expirationDate = expirationDate
    document.haveStocked = false
    document.stockedDate = 0
    document.storeIn = maxEnv
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
