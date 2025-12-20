// Пример проблемного TypeScript кода для тестирования ESLint
// Ожидаемые проблемы: использование any, неиспользуемые переменные, использование ==

function processUser(user: any) {
    // Использование any
    const name = user.name;
    const age = user.age;
    
    // Неиспользуемая переменная
    const unused = 10;
    
    // Использование == вместо ===
    if (age == 0) {
        console.log("Zero age");
    }
    
    // Неправильное использование типов
    return name + age; // Конкатенация строки и числа
}

interface User {
    name: string;
    age: number;
}

function getUser(): User {
    // Возврат объекта без всех полей
    return {
        name: "John"
    };
}
