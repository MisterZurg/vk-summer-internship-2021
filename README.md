# vk-summer-internship-2021
Test assignment for an intern for the Core Infrastructure team ( Summer 2021 )

# Description
Создайте программу-калькулятор, работающую с римскими цифрами.

## Язык: 
C, C++ или Go.

Без внешних зависимостей, можно использовать только стандартную библиотеку выбранного языка.

• Генераторы лексеров и парсеров использовать нельзя.

## Сборка программы:

C: ```gcc -std=c99 -Werror -Wall -Wextra -Wpedantic -o calc *.c```

C++: ```gcc -std=c++17 -Werror -Wall -Wextra -Wpedantic -o calc *.cpp```

Go: ```go build -o calc```

Исходный код должен находиться в приватном репозитории на GitHub или GitLab.

## Ввод-вывод:

* Выражения для вычисления читаются из stdin, по одному выражению на строку.
* Строки разделены \n.

* Ответы пишутся в stdout, по одному ответу на строку.
* Строки разделены \n.
* Ничего, кроме ответа, в stdout не пишется.

* Программа завершает работу после обработки всех данных из stdin.

## Формат входных данных:

* Числа, состоящие из римских цифр: ```I, V, X, L, C, D, M, Z.```
* ```Z``` обозначает ```0```.
* Операции: ```+, -, *, /.```
* Унарный минус ```(-)``` перед числом обозначает, что оно отрицательное.
* Скобки ```(, )```.
* Между числами, операциями и скобками допустимы пробелы.

## Формат ответа:
* В случае успеха — число, состоящее из римских цифр.
* В случае ошибки — строка, начинающаяся с ```error:```.

Вычисления производятся с использованием 64-битных знаковых чисел.

## Бонусом будет:
* Описание BNF-грамматики в файле grammar.txt.
* Ответ ошибкой при определении целочисленного переполнения.


### [Разбор тестового задания команды Core Infrastructure](https://vk.com/@vkteam-razbor-testovogo-zadaniya-komandy-core-infrastructure)
