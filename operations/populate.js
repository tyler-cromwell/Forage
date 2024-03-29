print('======================================')
let database = connect('127.0.0.1:27017/forage')

let resultIngredientsDrop = database.ingredients.drop()
print('Ingredients Dropped:', resultIngredientsDrop)
let resultRecipesDrop = database.recipes.drop()
print('Recipes Dropped:', resultRecipesDrop)

// Production will include expiration date
let dateUpdated = new Date()
let ingredients = [
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
        name: 'Aioli'
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
        name: 'Aioli'
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
        name: 'Apples'
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
        name: 'Bacon'
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
        name: 'Bacon'
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
        name: 'Bagel'
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
        name: 'Bananas'
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
        name: 'Beef'
    },
    {
        amount: {
            value: 1,
            unit: 'pound'
        },
        attributes: {
            cooked: true,
            opened: true
        },
        lifespan: {
            refrigerator: {
                value: 2,
                unit: 'day'
            }
        },
        name: 'Beef'
    },
    {
        amount: {
            value: 1,
            unit: 'piece'
        },
        attributes: {
            cooked: false,
        },
        lifespan: {
            refrigerator: {
                value: 10,
                unit: 'day'
            }
        },
        name: 'Bell Peppers'
    },
    {
        amount: {
            value: 1,
            unit: 'piece'
        },
        attributes: {
            cooked: true,
        },
        lifespan: {
            refrigerator: {
                value: 5,
                unit: 'day'
            }
        },
        name: 'Bell Peppers'
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
        name: 'Broccoli'
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
        name: 'Brussel Sprouts'
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
        name: 'Bucatini'
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
        name: 'Bucatini'
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
        name: 'Butter'
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
        name: 'Carrots'
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
        name: 'Celery'
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
        name: 'Cheerios'
    },
    {
        amount: {
            unit: 'ounce',
            value: 1
        },
        attributes: {
            opened: false,
            type: 'American'
        },
        lifespan: {
            refrigerator: {
                unit: 'month',
                value: 5
            }
        },
        name: 'Cheese'
    },
    {
        amount: {
            unit: 'ounce',
            value: 1
        },
        attributes: {
            opened: false,
            type: 'Cream'
        },
        lifespan: {
            freezer: {
                unit: 'month',
                value: 2
            }
        },
        name: 'Cheese'
    },
    {
        amount: {
            unit: 'ounce',
            value: 1
        },
        attributes: {
            opened: true,
            type: 'Cream'
        },
        lifespan: {
            refrigerator: {
                value: 10,
                unit: 'day'
            }
        },
        name: 'Cheese'
    },
    {
        amount: {
            unit: 'ounce',
            value: 1
        },
        attributes: {
            opened: false,
            type: 'Mozzarella'
        },
        lifespan: {
            refrigerator: {
                value: 6,
                unit: 'week'
            }
        },
        name: 'Cheese'
    },
    {
        amount: {
            unit: 'ounce',
            value: 1
        },
        attributes: {
            opened: true,
            type: 'Mozzarella'
        },
        lifespan: {
            refrigerator: {
                value: 5,
                unit: 'day'
            }
        },
        name: 'Cheese'
    },
    {
        amount: {
            unit: 'ounce',
            value: 1
        },
        attributes: {
            type: 'Muenster'
        },
        lifespan: {
            refrigerator: {
                value: 3,
                unit: 'week'
            }
        },
        name: 'Cheese'
    },
    {
        amount: {
            unit: 'ounce',
            value: 1
        },
        attributes: {
            opened: false,
            type: 'Parmesan'
        },
        lifespan: {
            refrigerator: {
                unit: 'month',
                value: 3
            }
        },
        name: 'Cheese'
    },
    {
        amount: {
            unit: 'ounce',
            value: 1
        },
        attributes: {
            opened: true,
            type: 'Parmesan'
        },
        lifespan: {
            refrigerator: {
                unit: 'week',
                value: 6
            }
        },
        name: 'Cheese'
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
        name: 'Chicken'
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
        name: 'Chicken'
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
                value: 2,
                unit: 'month'
            },
            pantry: {
                value: 3,
                unit: 'day'
            },
            refrigerator: {
                value: 4,
                unit: 'day'
            }
        },
        name: 'Chicken Broth'
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
                value: 1,
                unit: 'year'
            },
            pantry: {
                value: 1,
                unit: 'year'
            },
            refrigerator: {
                value: 1,
                unit: 'year'
            }
        },
        name: 'Chicken Broth'
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
        name: 'Chili Flakes'
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
        name: 'Cucumber'
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
        name: 'Cumin'
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
        name: 'Dijon Mustard'
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
        name: 'Dijon Mustard'
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
        name: 'Eggs'
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
        name: 'Garlic'
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
        type: "ingredient"
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
        name: 'Grape Tomatoes'
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
        name: 'Greek Yogurt'
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
        name: 'Greek Yogurt'
    },
    {
        amount: {
            unit: 'head',
            value: 1
        },
        attributes: {
            type: 'Green Leaf'
        },
        lifespan: {
            refrigerator: {
                unit: 'day',
                value: 12
            }
        },
        name: 'Lettuce'
    },
    /*
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
        type: 'Meal'
    },
    */
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
        name: 'Hamburger Seasoning'
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
        name: 'Hot Sauce'
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
        name: 'Hot Sauce'
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
        name: 'Hot Sauce'
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
        name: 'Hot Sauce'
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
        name: 'Jam'
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
        name: 'Jam'
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
        name: 'Jam'
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
        name: 'Jam'
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
        name: 'Linguine'
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
        name: 'Linguine'
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
        name: 'Maple Syrup'
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
        name: 'Maple Syrup'
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
        name: 'Mayonnaise'
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
        name: 'Mayonnaise'
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
        name: 'Milk'
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
        name: 'Mirin'
    },
    {
        amount: {
            value: 1,
            unit: 'piece'
        },
        attributes: {
            type: 'Baby Bella'
        },
        lifespan: {
            pantry: {
                unit: 'week',
                value: 1
            }
        },
        name: 'Mushrooms'
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
        name: 'Napa Cabbage'
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
        name: 'Nori'
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
        name: 'Olive Oil'
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
        name: 'Olive Oil'
    },
    {
        amount: {
            unit: 'head',
            value: 1
        },
        attributes: {
            'cut': false,
            'type': 'Yellow'
        },
        lifespan: {
            freezer: {
                value: 8,
                unit: 'month'
            },
            refrigerator: {
                value: 6,
                unit: 'week'
            },
            pantry: {
                value: 4,
                unit: 'week'
            }
        },
        name: 'Onion'
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
        name: 'Onion'
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
        name: 'Panko Bread Crumbs'
    },
    {
        amount: {
            value: 1,
            unit: 'ounce'
        },
        attributes: {
            flavor: 'Smoked'
        },
        lifespan: {
            pantry: {
                unit: 'year',
                value: 3
            }
        },
        name: 'Paprika'
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
        name: 'Pasta'
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
        name: 'Peanut Butter'
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
        name: 'Peas'
    },
    {
        amount: {
            value: 1,
            unit: 'ounce'
        },
        attributes: {
            type: 'Black',
        },
        lifespan: {
            pantry: {
                value: 3,
                unit: 'year'
            }
        },
        name: 'Pepper'
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
        name: 'Pickles'
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
        name: 'Pork'
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
        name: 'Pork'
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
        name: 'Quinoa'
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
        name: 'Ranch Dressing'
    },
    {
        amount: {
            value: 1,
            unit: 'cup'
        },
        attributes: {
            cooked: false,
            lowHumidity: true,
            opened: false,
            type: 'White'
        },
        lifespan: {
            pantry: {
                unit: 'month',
                value: 6
            }
        },
        name: 'Rice'
    },
    {
        amount: {
            value: 1,
            unit: 'cup'
        },
        attributes: {
            cooked: true,
            lowHumidity: true,
            opened: false,
            type: 'White'
        },
        lifespan: {
            refrigerator: {
                unit: 'day',
                value: 3
            }
        },
        name: 'Rice'
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
        name: 'Rice Vinegar'
    },
    {
        amount: {
            value: 1,
            unit: 'ounce'
        },
        attributes: {
            type: 'Dried'
        },
        lifespan: {
            pantry: {
                unit: 'year',
                value: 1
            }
        },
        name: 'Rosemary'
    },
    {
        amount: {
            value: 1,
            unit: 'ounce'
        },
        lifespan: {
            pantry: {
                value: 0,
                unit: 'day'
            }
        },
        name: 'Salt'
    },
    {
        amount: {
            value: 1,
            unit: 'stalk'
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
        name: 'Scallions'
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
        name: 'Shallot'
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
        name: 'Soy Sauce'
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
        name: 'Soy Sauce'
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
        name: 'Sugar'
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
        name: 'Turkey'
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
        name: 'Turkey'
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
        name: 'Turkey'
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
        name: 'Unagi'
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
        name: 'Vegetable Oil'
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
        name: 'White Bread'
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
        name: 'Whole Wheat Bread'
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
        name: 'Worcestershire Sauce'
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
        name: 'Worcestershire Sauce'
    },
    {
        amount: {
            value: 1,
            unit: 'piece'
        },
        lifespan: {
            refrigerator: {
                value: 10,
                unit: 'day'
            }
        },
        name: 'Zucchini'
    }
]

// Prep ingredient documents for insertion
ingredients.forEach((document) => {
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
        dateUpdated.getFullYear(),
        dateUpdated.getMonth(),
        dateUpdated.getDate()+maxDays
    ).getTime()

    if (document.expirationDate === undefined && maxDays == 0) {
        expirationDate = 0
    }

    // Update document
    document.expirationDate = expirationDate
    document.haveStocked = false
    document.stockedDate = 0
    document.storeIn = maxEnv
    document.updated = dateUpdated.getTime()
})

// Insert the ingredients
let resultInsertIngredients = database.ingredients.insertMany(ingredients)
let ingredientIDs = resultInsertIngredients.insertedIds
print('Inserted', ingredients.length, 'of', ingredientIDs.length, 'ingredient documents')
if (ingredients.length != ingredientIDs.length) {
    // error
}

/*
idBellPeppers = database.ingredients.find({ name: 'Bell Peppers', "attributes.cooked": { $exists: true, $eq: false } }, { _id: 1 })
idChicken = database.ingredients.find({ name: 'Chicken', "attributes.cooked": { $exists: true, $eq: false } }, { _id: 1 })
idGarlic = database.ingredients.find({ name: 'Garlic' }, { _id: 1 })
idMushroomsBabyBella = database.ingredients.find({ name: 'Mushrooms', "attributes.type": { $exists: true, $eq: "Baby Bella" } }, { _id: 1 })
idOnionYellow = database.ingredients.find({ name: 'Onion', "attributes.type": { $exists: true, $eq: "Yellow" } }, { _id: 1 })
idZucchini = database.ingredients.find({ name: 'Zucchini' }, { _id: 1 })
*/

let recipes = [
//    { name: 'Bacon, Egg, and Cheese' },
    {
        name: 'Chicken & Vegetable Quinoa',
        ingredients: [
            "Bell Peppers",
            "Chicken",
            "Chicken Broth",
            "Garlic",
            "Mushrooms",
            "Olive Oil", // Extra Virgin
            "Onion",
            "Paprika", // Smoked
            "Pepper", // Black
            "Quinoa",
            "Rosemary", // Dried, Crushed
            "Salt",
            // ? Spinach Leaves
            "Zucchini"
        ]
    },
    {
        name: 'Chicken Fried Rice',
        ingredients: [
            "Carrots",
            "Chicken",
            "Eggs",
            "Garlic",
            "Peas",
            "Rice",
            "Scallions"
            // Soy Sauce
        ]
    },
//    { name: 'Gyoza' },
    {
        name: 'Hamburgers',
        ingredients: [
            "Beef",
            "Cheese", // American
            "Lettuce", // Green Leaf
            "Pepper", // Black
            "Salt",
        ]
    },
//    { name: 'Baked Gnocchi & Broccoli' },
//    { name: 'Blueberry Muffins' }
//    { name: 'Bucatini Carbonara (Modern)' },
//    { name: 'Chicken Saltimbocca Alla Romana' },
//    { name: 'Classic Caponata' },
//    { name: 'Italian Wedding Soup' },
//    { name: 'Meatballs' },
//    { name: 'One Pan Pasta' },
//    { name: 'Oyako Don' },
//    { name: 'Slow-Cooked Bolognese Sauce' },
//    { name: 'Strawberry Nutella Semifreddo' },
//    { name: 'Sushi' },
//    { name: 'Wings' },
]

// Prep recipe documents for insertion
recipes.forEach((document) => {
    document.isCookable = false
    document.updated = dateUpdated.getTime()
})

// Insert the recipes
let resultInsertRecipies = database.recipes.insertMany(recipes)
let recipeIDs = resultInsertRecipies.insertedIds
print('Inserted', recipes.length, 'of', recipeIDs.length, 'recipe documents')
if (recipes.length != recipeIDs.length) {
    // error
}

/*
let resultFind = database.ingredients.find()
print('Find:', resultFind)

while (resultFind.hasNext()) {
    print(tojson(resultFind.next()))
}
*/
print('======================================')