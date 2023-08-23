// scripts.js
const form = document.getElementById('upload-form');
const resultDiv = document.getElementById('conversion-result');

form.addEventListener('submit', async (event) => {
    event.preventDefault();

    const formData = new FormData(form);
    const response = await fetch('/upload', {
        method: 'POST',
        body: formData,
    });

    const resultText = await response.text();
    resultDiv.innerHTML = `<h2>Conversion Result:</h2><pre>${resultText}</pre>`;
});
