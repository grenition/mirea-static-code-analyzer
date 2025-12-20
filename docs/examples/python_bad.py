# Пример проблемного Python кода для тестирования Pylint
# Ожидаемые проблемы: неиспользуемые переменные, неправильное именование, неиспользуемые импорты

import os  # Неиспользуемый импорт
import sys  # Неиспользуемый импорт

def calculate_sum(numbers):
    unused = 10  # Неиспользуемая переменная
    result = 0
    
    for num in numbers:
        result = result + num
    
    temp = result * 2  # Неиспользуемая переменная
    
    # Плохое именование функции
    def BadFunction():
        pass
    
    # Использование print для отладки
    print("Result:", result)
    
    return result

def process_data(data):
    # Плохое именование переменной
    X = data[0]
    return X

def divide(a, b):
    # Возможное деление на ноль
    return a / b

