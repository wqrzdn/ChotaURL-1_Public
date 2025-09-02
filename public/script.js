
document.addEventListener('DOMContentLoaded', function() {
    const form = document.getElementById('shorten-form');
    const urlInput = document.getElementById('url-input');
    const customNameInput = document.getElementById('custom-name');
    const submitButton = form.querySelector('button[type="submit"]');

   
    customNameInput.addEventListener('input', function() {
        const value = this.value;
        if (value && (value.length < 3 || value.length > 20)) {
            this.setCustomValidity('Custom name must be 3-20 characters');
        } else if (value && !/^[a-zA-Z0-9_-]+$/.test(value)) {
            this.setCustomValidity('Only letters, numbers, hyphens, and underscores allowed');
        } else {
            this.setCustomValidity('');
        }
    });

   
    form.addEventListener('submit', async function(e) {
        e.preventDefault();

        const url = urlInput.value.trim();
        const customName = customNameInput.value.trim();
        
     
        let formattedUrl = url;
        if (!url.startsWith('http://') && !url.startsWith('https://')) {
            formattedUrl = 'https://' + url;
        }

        
        submitButton.innerHTML = '<i class="fas fa-spinner fa-spin"></i> Shortening...';
        submitButton.disabled = true;

        try {
            const requestBody = { url: formattedUrl };
            if (customName) {
                requestBody.custom = customName;
            }

           
            const isVercel = window.location.hostname.includes('vercel.app');
            const apiEndpoint = isVercel ? '/api/shorten' : '/shorten';

            const response = await fetch(apiEndpoint, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(requestBody)
            });

            const data = await response.json();

            if (response.ok && data.short_url) {
                showResult(data.short_url, formattedUrl);
                form.reset();
            } else {
                showError(data.error || 'Failed to shorten URL');
            }
        } catch (error) {
            showError('Network error: ' + error.message);
        } finally {
            
            submitButton.innerHTML = '<i class="fas fa-magic"></i> Shorten URL';
            submitButton.disabled = false;
        }
    });

    function showResult(shortUrl, originalUrl) {
        
        const existingResult = document.querySelector('.result-container');
        if (existingResult) {
            existingResult.remove();
        }

      
        const resultContainer = document.createElement('div');
        resultContainer.className = 'result-container';
        resultContainer.innerHTML = `
            <div class="result">
                <h3>URL Shortened Successfully!</h3>
                <div class="url-pair">
                    <div class="url-row">
                        <label>Original:</label>
                        <span class="original-url">${originalUrl}</span>
                    </div>
                    <div class="url-row">
                        <label>Shortened:</label>
                        <div class="short-url-container">
                            <a href="${shortUrl}" target="_blank" class="short-url">${shortUrl}</a>
                            <button class="copy-btn" data-url="${shortUrl}">
                                <i class="fas fa-copy"></i> Copy
                            </button>
                        </div>
                    </div>
                </div>
            </div>
        `;

       
        form.parentNode.insertBefore(resultContainer, form.nextSibling);

      
        const copyBtn = resultContainer.querySelector('.copy-btn');
        copyBtn.addEventListener('click', function() {
            navigator.clipboard.writeText(this.dataset.url).then(() => {
                this.innerHTML = '<i class="fas fa-check"></i> Copied!';
                setTimeout(() => {
                    this.innerHTML = '<i class="fas fa-copy"></i> Copy';
                }, 2000);
            });
        });
    }

    function showError(message) {
       
        const existingError = document.querySelector('.error-container');
        if (existingError) {
            existingError.remove();
        }

       
        const errorContainer = document.createElement('div');
        errorContainer.className = 'error-container';
        errorContainer.innerHTML = `
            <div class="error">
                <h3>❌ Error</h3>
                <p>${message}</p>
            </div>
        `;

        
        form.parentNode.insertBefore(errorContainer, form.nextSibling);

       
        setTimeout(() => {
            if (errorContainer.parentNode) {
                errorContainer.remove();
            }
        }, 5000);
    }
});
