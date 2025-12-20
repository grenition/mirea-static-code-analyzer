// Пример проблемного C++ кода для тестирования cppcheck
// Ожидаемые проблемы: утечки памяти, использование после освобождения, выход за границы массива

#include <iostream>

void memoryLeak() {
    int* ptr = new int[100];
    // Утечка памяти - нет delete[]
    return;
}

void useAfterFree() {
    int* ptr = new int;
    delete ptr;
    *ptr = 10; // Использование после освобождения
}

void nullPointer() {
    int* ptr = nullptr;
    *ptr = 5; // Разыменование нулевого указателя
}

void arrayBounds() {
    int arr[10];
    arr[15] = 100; // Выход за границы массива
}

void uninitialized() {
    int value;
    if (value > 0) { // Использование неинициализированной переменной
        std::cout << "Positive";
    }
}

int divide(int a, int b) {
    return a / b; // Возможное деление на ноль
}

