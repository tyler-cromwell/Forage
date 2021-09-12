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
                value: 1,
                unit: 'week'
            },
            name: 'Apples',
            refrigerate: true,
            updated: dateUpdated
        },
        {
            amount: {
                value: 4,
                unit: 'count'
            },
            lifespan: {
                value: 4,
                unit: 'day'
            },
            name: 'Bananas',
            shelved: true,
            updated: dateUpdated
        },
        {
            amount: {
                value: 12,
                unit: 'count'
            },
            lifespan: {
                value: 1,
                unit: 'month'
            },
            name: 'Eggs',
            refrigerate: true,
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'pound'
            },
            cooked: false,
            freeze: true,
            lifespan: {
                value: 3,
                unit: 'month'
            },
            name: 'Meat',
            sealed: true,
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'pound'
            },
            cooked: true,
            lifespan: {
                value: 2,
                unit: 'day'
            },
            name: 'Meat',
            refrigerate: true,
            sealed: true,
            updated: dateUpdated
        },
        {
            amount: {
                value: 1,
                unit: 'gallon'
            },
            lifespan: {
                value: 10,
                unit: 'day'
            },
            name: 'Milk',
            refrigerate: true,
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
