**Malls list**
----
* **URL:**

    /malls/

* **Query Params:**

    shop [integer] - shop id, malls with such shop

    city [integer] - city id

    query [string] - text query to search

    ids [list] - list of specific mall ids

    sort [string] - which field to use for sorting

    limit [integer] - number of results

    offset [integer] - index of first result


**Mall details**
----
* **URL:**

    /malls/:id/

* **Query Params:**

    None


**Shops list**
----
* **URL:**

    /shops/

* **Query Params:**

    mall [integer] - mall id, shops in such mall

    category [integer] - category id

    query [string] - text query to search

    ids [list] - list of specific shop ids

    sort [string] - which field to use for sorting

    limit [integer] - number of results

    offset [integer] - index of first result


**Shop details**
----
* **URL:**

    /shops/:id/

* **Query Params:**

    None


**Categories list**
----
* **URL:**

    /categories/

* **Query Params:**

    shop [integer] - shop id, list of categories for this shop

    ids [list] - list of specific category ids

    sort [string] - which field to use for sorting


**Category details**
----
* **URL:**

    /category/:id/

* **Query Params:**

    None


**Cities list**
----
* **URL:**

    /cities/

* **Query Params:**

    ids [list] - list of specific city ids


**City details**
----
* **URL:**

    /cities/:id/

* **Query Params:**

    None

**Mall Object**
----
```javascript
{
  "id": 228,
  "name": "Some name",
  "site": "http://domain.com/",
  "phone": "+79250741413",
  "address": "ул. Перерва, 45, Москва, Россия, 10934",
  "location": {
    "lat": 22.33334,
    "lon": 33.35533
  },
  "logo": {
    "large": "https://storage.domain.com/path/to/large/logo.png",
    "small": "https://storage.domain.com/path/to/small/logo.png"
  },
  "subway_station": "Кантимировская",
  "day_and_night": false,
  "working_hours": [
    {
      "closing": {
        "day": 0,
        "time": "19:00:00"
      },
      "opening": {
        "day": 0,
        "time": "10:00:00"
      }
    },
    {
      "closing": {
        "day": 1,
        "time": "19:00:00"
      },
      "opening": {
        "day": 1,
        "time": "10:00:00"
      }
    },
    {
      "closing": {
        "day": 2,
        "time": "19:00:00"
      },
      "opening": {
        "day": 2,
        "time": "10:00:00"
      }
    },
    {
      "closing": {
        "day": 3,
        "time": "19:00:00"
      },
      "opening": {
        "day": 3,
        "time": "10:00:00"
      }
    },
    {
      "closing": {
        "day": 4,
        "time": "19:00:00"
      },
      "opening": {
        "day": 4,
        "time": "10:00:00"
      }
    },
    {
      "closing": {
        "day": 5,
        "time": "19:00:00"
      },
      "opening": {
        "day": 5,
        "time": "10:00:00"
      }
    },
    {
      "closing": {
        "day": 6,
        "time": "19:00:00"
      },
      "opening": {
        "day": 6,
        "time": "10:00:00"
      }
    }
  ]
}
```


**Shop Object**
----
```javascript
{
  "id": 228,
  "name": "Some name",
  "site": "http://domain.com/",
  "phone": "+79250741413",
  "score": 400,
  "logo": {
    "large": "https://storage.domain.com/path/to/large/logo.png",
    "small": "https://storage.domain.com/path/to/small/logo.png"
  }
}
```


**Category Object**
----
```javascript
{
  "id": 228,
  "name": "Some name",
  "logo": {
    "large": "https://storage.domain.com/path/to/large/logo.png",
    "small": "https://storage.domain.com/path/to/small/logo.png"
  }
}
```


**City Object**
----
```javascript
{
  "id": 228,
  "name": "Some name"
}
```