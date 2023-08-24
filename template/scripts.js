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

            // Add download link for the converted result file
            if (formData.get('format') !== 'pdf') {
                const downloadLink = document.createElement('a');
                downloadLink.href = `/download?file=${formData.get('file')}.${formData.get('format')}`;
                downloadLink.textContent = 'Download Result';
                resultDiv.appendChild(downloadLink);
            }
        } catch (error) {
            resultDiv.innerHTML = '<p>An error occurred during conversion. Please try again later.</p>';
            console.error(error);
        }
    });
});
