**Общие замечания**
----
- Все запры GET.

- limit, offset - эти параметры отвечают за пагинацию, если указаны значит запрос подразумевает пагинацию.

- sort - как сортировать выборку. Везде значение по умолчанию: "id", то есть по айди по возрастанию.
У каждого значения есть обратное: "name" и "-name" по возрастанию и по убыванию соответственно.

- Успешный ответ выглядит следующим образом:
http status = 200
```json
{
  "data": {...}
}
```
на верхнем уровне поле "data", внутри данные соответствующие запросу.

- Ошибочный ответ выглядит следующим образом:
http status != 200
```json
{
  "error": {
    "code": "SOME_ERROR_CONSTANT",
    "status": 400,
    "details": "some text explanation"
  }
}
```
error.code это код ошибки на который нужно смотреть.
error.status дублирует http status.
Любой запрос точно может вернуть код "INCORRECT_REQUEST_DATA" это означает что данные пришли не в том типе, формате, диапазоне и т.д.

- Если запрос подразумевает пагинацию, то данные в ответе ( то что в поле "data") будут выглядеть так:
```json
{
  "total_count": 1000,
  "count": 20,
  "results": [...]
}
```
total_count общее количество элементов в выборке, без учета пагинации.
count размер выборки, равер длине массива results, может быть меньше параметра limit
results список объектов

- //details напротив какого либо поля в описании структуры объекта означает,
что данное поле будет только когда вы запрашиваете этот объект в единственном экземпляре.

- //null напротив какого либо поля означает что это поле может быть null у каких либо объектов.

- Если напротив параметра стоит filter. значит это основной фильтрующий параметр и в запросе не может участвовать два фильтровых параметра,
будет выбран только одни.

- Параметр city является дополнительным фильтрующим параметров и присутствует почти во всех запросах,
если указан то будет выборка только ко конкретному городу, если нет то по всей базе.


**Malls list**
----
Возвращает список тц по различным фильтрам.

* **URL:**

    /malls/

* **Query Params:**

    **Optional:**

    shop [integer] filter - дай тц в которых есть этот магаз

    subway_station [integer] filter - дай тц которые находятся на данной станции метро

    query [string] filter - дай тц по имени

    ids [list] filter - дай тц по этим айдишкам. С этим параметром не будет работать пагинация и сортировка

    city [integer] - city id

    sort [string] - возможные значения: "name", "shops_count", "id"

    limit [integer]

    offset [integer]

* **Error Responses:**

    400, "INCORRECT_REQUEST_DATA"

    404, "SHOP_NOT_FOUND"

    404, "SUBWAY_STATION_NOT_FOUND"

    404, "CITY_NOT_FOUND"


**Mall details**
----
Возращает детальную информацию о конкретном тц.

* **URL:**

    /malls/:id/

* **Query Params:**

    None

* **Error Responses:**

    400, "INCORRECT_REQUEST_DATA"

    404, "MALL_NOT_FOUND"

**Shops list**
----
Возращает список магазов по различным фильтрам.

* **URL:**

    /shops/

* **Query Params:**


    **Optional:**

    mall [integer] filter - дай магазы в этом тц

    category [integer] filter - дай магазы у которых есть такая категория

    query [string] filter - дай магазы по имени

    ids [list] filter - дай магазы по айдишкам, пагинация и сортировка не работают

    city [integer] - city id

    sort [string] - возможные значения: "name", "id", "score", "malls_count"

    limit [integer]

    offset [integer]

* **Error Responses:**

    400, "INCORRECT_REQUEST_DATA"

    404, "MALL_NOT_FOUND"

    404, "CATEGORY_NOT_FOUND"

    404, "CITY_NOT_FOUND"


**Shop details**
----
Возвращает подробную информацию о данном магазе.

* **URL:**

    /shops/:id/

* **Query Params:**

    **Optional:**

    location_lat [float] - x координата юзера, используются для определения ближайшего тц

    location_lon [float] - y координата юзера

    city [integer] - city id

* **Error Responses:**

    400, "INCORRECT_REQUEST_DATA"

    404, "SHOP_NOT_FOUND"

    404, "CITY_NOT_FOUND"


**Categories list**
----
Возращает список категорий.

* **URL:**

    /categories/

* **Query Params:**

    **Optional:**

    shop [integer] filter - дай категории данного магаза

    ids [list] filter - дай категории с такими айдишками, сортировка не работает

    city [integer] - city id

    sort [string] - возможные значения: "id", "name", "shops_count"

* **Error Responses:**

    400, "INCORRECT_REQUEST_DATA"

    404, "SHOP_NOT_FOUND"

    404, "CITY_NOT_FOUND"


**Category details**
----
Возращает инфу у конкретной категории, на данный момент там все таже инфа что и в списке.

* **URL:**

    /category/:id/

* **Query Params:**

    **Optional:**

    city [integer] - city id

* **Error Responses:**

    400, "INCORRECT_REQUEST_DATA"

    404, "CATEGORY_NOT_FOUND"

    404, "CITY_NOT_FOUND"


**Cities list**
----
Возвращает список городов.

* **URL:**

    /cities/

* **Query Params:**

    query [string] - text query to search

    sort [string] - возможные значения: "id", "name"

* **Error Responses:**

    400, "INCORRECT_REQUEST_DATA"


**Current mall**
----
Пытается определить тц в котором сейчас находится пользователь по местоположению. Возращает подробную инфу о тц.

* **URL:**

    /current_mall/

* **Query Params:**

* **Required:**

    location_lat [float] - x координата юзера

    location_lon [float] - y координата юзера

* **Error Responses:**

    400, "INCORRECT_REQUEST_DATA"

    404, "MALL_NOT_FOUND"


**Current city**
----
Пытается определить город пользователя по координатам.

* **URL:**

    /current_city/

* **Query Params:**

* **Required:**

    location_lat [float] - x координата юзера

    location_lon [float] - y координата юзера

* **Error Responses:**

    400, "INCORRECT_REQUEST_DATA"

    404, "CITY_NOT_FOUND"


**Shops In Malls**
----
Устанавлиет какие из указанных магазинов есть в указанных тц.

* **URL:**

    /shops_in_malls/

* **Query Params:**

* **Required:**

    shops [list] - список магазинов которые надо соотнести с тц

    malls [list] - список тц у которых надо опеределить магазины

* **Success Responses:**

```json
[
  {
    "mall": 1,
    "shops": [
      2,
      3,
      4,
    ]
  },
  {
    "mall": 2,
    "shops": [
      2,
      5,
      4
    ]
  },
  {
    "mall": 4,
    "shops": []
  }
]
```

* **Error Responses:**

    400, "INCORRECT_REQUEST_DATA"


**Search**
----
Ищет тц по указанному списку магазов

* **URL:**

    /shops_in_malls/

* **Query Params:**

    shops [list] - список магазинов по которым надо искать

    location_lat [float] - x координата юзера

    location_lon [float] - y координата юзера

    city [integer] - city id

    sort [string] - результаты будут отсортированы по кол-ву совпавжих магазинов,
                    это вторичная сортировка, для тц у которых кол-во совпавших равно
                    возможные значения: "mall_id", "mall_name", "mall_shops_count", "distance"
                    сортировка по "distance" возможно только если запрос содержал координты юзера.

    limit [integer]

    offset [integer]

* **Success Responses:**

```json
[
  {
    "mall": {
      "id": 228,
      "name": "Some name",
      "location": {
        "lat": 22.33334,
        "lon": 33.35533
      },
      "logo": {
        "large": "https://storage.domain.com/path/to/large/logo.png",
        "small": "https://storage.domain.com/path/to/small/logo.png"
      },
      "shops_count": 44,
    },
    "shops": [
      2,
      3,
      4,
    ],
    "distance": 1000.222, //null
  },
  {
    "mall": {
      "id": 322,
      "name": "Some another name",
      "location": {
        "lat": 22.33334,
        "lon": 33.35533
      },
      "logo": {
        "large": "https://storage.domain.com/path/to/large/logo.png",
        "small": "https://storage.domain.com/path/to/small/logo.png"
      },
      "shops_count": 44,
    },
    "shops": [
      2,
      3,
    ],
    "distance": 500.222, //null
  }
]
```

* **Error Responses:**

    400, "INCORRECT_REQUEST_DATA"

    404, "CITY_NOT_FOUND"

**Mall Object**
----
```json
{
  "id": 228,
  "name": "Some name",
  "site": "http://domain.com/", //details
  "phone": "+79250741413", //details
  "address": "ул. Перерва, 45, Москва, Россия, 10934", //details
  "location": {
    "lat": 22.33334,
    "lon": 33.35533
  },
  "logo": {
    "large": "https://storage.domain.com/path/to/large/logo.png",
    "small": "https://storage.domain.com/path/to/small/logo.png"
  },
  "subway_station": { //details
    "id": 228,
    "name": "Кантимировская"
  },
  "day_and_night": false, //details
  "shops_count": 44,
  "working_hours": [ //details
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
```json
{
  "id": 228,
  "name": "Some name",
  "site": "http://domain.com/", //detials
  "phone": "+79250741413", //details
  "nearest_mall": { //details, //null
    "id": 228,
    "name": "Some name",
    "location": {
      "lat": 22.33334,
      "lon": 33.35533
    },
    "logo": {
      "large": "https://storage.domain.com/path/to/large/logo.png",
      "small": "https://storage.domain.com/path/to/small/logo.png"
    },
    "shops_count": 44,
  },
  "score": 400,
  "malls_count": 23,
  "logo": {
    "large": "https://storage.domain.com/path/to/large/logo.png",
    "small": "https://storage.domain.com/path/to/small/logo.png"
  }
}
```


**Category Object**
----
```json
{
  "id": 228,
  "name": "Some name",
  "shops_count": 2222,
  "logo": {
    "large": "https://storage.domain.com/path/to/large/logo.png",
    "small": "https://storage.domain.com/path/to/small/logo.png"
  }
}
```


**City Object**
----
```json
{
  "id": 228,
  "name": "Some name"
}
```