pg docker:

docker run --name avito_intern -e POSTGRES_PASSWORD=docker -p 54322:5432 -d postgres

Реализовано:

1) Метод начисления средств на баланс.
    - Принимает id пользователя / id кошелька / email и сколько средств зачислить.
    - Возвращает message с email; balance.
2) Метод получения баланса по id пользователя.
    - Принимает id пользователя.
    - Возвращает message с email; balance.
3) Метод перевод средств от пользователя к пользователю.
    - Принимает id пользователя, от которого переводятся средства; id пользователя, которому переводятся средства;
      количество средст.
    - Возвращает message.
4) Метод резервирования средств с основного баланса на отдельном счете - реализован как метод оформления заказа на покупку услуги.
    - Принимает id пользователя, название услуги, сумму платежа.
    - Возвращает message и order_id.
5) Метод признания выручки - реализован как <br>
5.1) метод подтверждения заказа на покупку услуги <br>
5.2) отклонения заказа (и возвращение средств на основной баланс кошелька пользователя - разрезервирование средств пользователя).
    - Принимает user_id, order_id.
    - Возвращает message.

Вопросы:

1) По умолчанию сервис не содержит в себе никаких данных о балансах (пустая табличка в БД). Данные о балансе появляются
   при первом зачислении денег.
    - При создании пользователя сервиса у него нет кошелька, кошелек создается при первом зачислении денег на счет -
      создается кошелек с 0 балансом и следующим запросом баланс пополняется на указанную сумму.
    - Можно оптимизировать - делать все сразу в 1 запросе, но пока оставил такую реализацию.
2) Во всех методах, которые идентифицируют пользователей оставил только ид пользователя, кроме метода начисления
   средств, с целью уменьшения количества кода, но реализация показал :)
3) В идеале для метода "3) перевод средст между пользователями" нужна авторизация у того, кто переводит деньги, но в
   задании
   не было указано, что нужна авторизация, поэтому принял, что приходят запросы из другой системы, где пользователь уже
   авторизован и мы выполняем запросы от его имени.
4) Есть следующая формулировка "Для этого у нас есть специальный сервис управления услугами, который перед применением
   услуги резервирует деньги на отдельном счете и потом списывает в доход компании."
    - Сделал так, что при оформлении заказа на покупку услуги, средства переносятся с баланса на "удерживаемые" средства
      кошелька пользователя до того момента, пока не будет закрыт заказ (успех - списываются и переводятся продавцу, отказ -
      возвращаются пользователю).
5) Для доп.задания с отчетом решил, что не стоит делать отдельной таблицы и достаточно таблицы с транзакциями и
   фильтрацией по комментариям транзакций. Все транзакции с комментарием "income from..." - это получение прибыли с продажи услуг.
