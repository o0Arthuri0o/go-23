# Примеры для тестирования лексического анализатора

## Пример 1: Простые выражения
```
a := 1;
b := 0;
result := a and b;
```

## Пример 2: Все операции
```
x := 1;
y := 0;
res1 := x or y;
res2 := x xor y;
res3 := x and y;
res4 := not x;
```

## Пример 3: Сложные выражения со скобками
```
a := 1;
b := 0;
c := 1;
result := (a or b) and (not c xor b);
complex := ((a and b) or (c xor not a));
```

## Пример 4: С комментариями
```
// Это однострочный комментарий
a := 1;
b := 0;

/* Это многострочный
   комментарий */
result := a and b;

output := (a or b) xor not a; // комментарий в конце строки
```

## Пример 5: Длинные идентификаторы
```
very_long_variable_name := 1;
another_var123 := 0;
_underscore_var := 1;
result_final := very_long_variable_name and another_var123;
```

## Пример 6: Ошибки (для проверки обработки)
```
a = 1;           // Ошибка: должно быть :=
b := 2;          // Ошибка: 2 не является допустимой константой
c := a & b;      // Ошибка: & недопустимый символ
d := a AND b;    // Ошибка: AND должно быть строчными буквами
```

## Пример 7: Все допустимые лексемы
```
// Идентификаторы
var1 := 1;
var_2 := 0;
_var3 := 1;

// Константы
zero := 0;
one := 1;

// Операции
or_result := var1 or zero;
xor_result := var1 xor one;
and_result := var1 and one;
not_result := not zero;

// Скобки
grouped := (var1 or var_2) and one;

// Разделители
statement1 := 1; statement2 := 0;
```

## Пример 8: Вложенные скобки
```
a := 1;
b := 0;
c := 1;
d := 0;
result := ((a or b) and (c xor d)) or (not (a and c));
```

## Пример 9: Минимальная программа
```
x := 1;
```

## Пример 10: Комбинация всего
```
// Программа вычисления логических функций
// Автор: Студент

/* Инициализация входных значений
   для логических операций */
input_a := 1;
input_b := 0;
input_c := 1;

// Базовые операции
or_op := input_a or input_b;      // Логическое ИЛИ
and_op := input_a and input_b;    // Логическое И
xor_op := input_a xor input_c;    // Исключающее ИЛИ
not_op := not input_b;             // Логическое НЕ

/* Сложные выражения */
expr1 := (input_a or input_b) and input_c;
expr2 := not (input_a and input_b) xor input_c;
expr3 := ((input_a or not input_b) and input_c) xor not input_a;

// Финальный результат
final_result := (expr1 and expr2) or (expr3 xor not input_c);
```
