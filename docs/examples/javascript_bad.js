// Пример проблемного JavaScript кода для тестирования ESLint
// Ожидаемые проблемы: неиспользуемые переменные, использование ==, console.log

function calculateSum(numbers) {
    var unused = 10; // Неиспользуемая переменная
    let result = 0;
    
    for (let i = 0; i < numbers.length; i++) {
        result += numbers[i];
    }
    
    // Неиспользуемая переменная
    const temp = result * 2;
    
    // Использование == вместо ===
    if (result == 0) {
        return null;
    }
    
    // console.log в коде
    console.log("Result:", result);
    
    return result;
}

function processData(data) {
    // Использование var (устаревшее)
    for (var i = 0; i < data.length; i++) {
        console.log(data[i]);
    }
    
    // Доступ к i вне цикла (проблема с var)
    console.log("Last index:", i);
    
    return data;
}

