document.addEventListener('DOMContentLoaded', () => {
    const form = document.getElementById('upload-form');
    const resultDiv = document.getElementById('conversion-result');

    form.addEventListener('submit', async (event) => {
        event.preventDefault();

        const formData = new FormData(form);
        try {
            const response = await fetch('/upload', {
                method: 'POST',
                body: formData,
            });

            if (!response.ok) {
                throw new Error('Conversion request failed');
            }

            const resultText = await response.text();
            resultDiv.innerHTML = `<h2>Conversion Result:</h2><pre>${resultText}</pre>`;
        } catch (error) {
            resultDiv.innerHTML = '<p>An error occurred during conversion. Please try again later.</p>';
            console.error(error);
        }
    });
});
