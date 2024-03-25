Инструкция.

Перед запуском необходимо прописать переменные окружения:
1) Порт на котором будет работать сервер.


    export PORT=":8080"

3) Данные БД (логин, пароль, имя БД).

    export user="username" && export  password="password123" && export dbName="dbname"

Запуск. 
    
    go run main.go

Описание методов:
1) Сохранение новой команды в БД. Метод POST.

    localhost:8080/save
    
    body request:
    
    {"command_name" : "test"}

2) Получение списка команд. Метод GET.

    localhost:8080/getall

4) Получение одной команды + вызов. Метод POST.

    localhost:8080/get/id=n
   
    n - id команды
   
    optional body request:
   
    {"parameters": : ....}

6) Удаление команды из БД. Метод DELETE.

   localhost:8080/delete/id=n
   n - id команды

7) Остановка запущенной команды. Метод GET.

   localhost:8080/stop/id=n
   n - id команды
