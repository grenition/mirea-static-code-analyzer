// Client-side JavaScript for static code analyzer

document.addEventListener('DOMContentLoaded', function() {
    // Character counter for code snippet textarea
    const codeTextarea = document.getElementById('code');
    if (codeTextarea) {
        const maxLength = 50000;
        codeTextarea.addEventListener('input', function() {
            const remaining = maxLength - this.value.length;
            const counter = document.getElementById('char-counter') || createCounter();
            counter.textContent = `${this.value.length} / ${maxLength} characters`;
            if (remaining < 1000) {
                counter.classList.add('text-warning');
            } else {
                counter.classList.remove('text-warning');
            }
        });
    }
});

function createCounter() {
    const counter = document.createElement('small');
    counter.id = 'char-counter';
    counter.className = 'form-text text-muted';
    const textarea = document.getElementById('code');
    if (textarea && textarea.parentNode) {
        textarea.parentNode.appendChild(counter);
    }
    return counter;
}

