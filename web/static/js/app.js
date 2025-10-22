// Litepost Web UI Application
class LitepostApp {
    constructor() {
        this.currentRequest = null;
        this.currentResponse = null;
        this.init();
    }

    init() {
        this.setupEventListeners();
        this.setupTabs();
        this.setupAuthFields();
        this.loadSampleData();
    }

    setupEventListeners() {
        // Send button
        document.getElementById('sendButton').addEventListener('click', () => {
            this.sendRequest();
        });

        // Tab switching - use event delegation for better reliability
        document.addEventListener('click', (e) => {
            if (e.target.classList.contains('tab')) {
                console.log('Tab clicked:', e.target.dataset.tab);
                this.switchTab(e.target.dataset.tab);
            }
        });

        // Add/remove param buttons
        document.querySelector('.add-param').addEventListener('click', () => {
            this.addParamRow();
        });

        document.querySelector('.add-header').addEventListener('click', () => {
            this.addHeaderRow();
        });

        // Body type change
        document.getElementById('bodyType').addEventListener('change', (e) => {
            this.updateBodyType(e.target.value);
        });

        // Auth type change
        document.getElementById('authType').addEventListener('change', (e) => {
            this.updateAuthType(e.target.value);
        });
    }

    setupTabs() {
        // Request tabs
        const requestTabs = document.querySelectorAll('.request-tabs .tab');
        requestTabs.forEach(tab => {
            tab.addEventListener('click', (e) => {
                requestTabs.forEach(t => t.classList.remove('active'));
                e.target.classList.add('active');
                
                const tabContent = document.getElementById(e.target.dataset.tab + 'Tab');
                if (tabContent) {
                    document.querySelectorAll('.request-content .tab-content').forEach(content => {
                        content.classList.remove('active');
                    });
                    tabContent.classList.add('active');
                }
            });
        });

        // Response tabs
        const responseTabs = document.querySelectorAll('.response-tabs .tab');
        responseTabs.forEach(tab => {
            tab.addEventListener('click', (e) => {
                responseTabs.forEach(t => t.classList.remove('active'));
                e.target.classList.add('active');
                
                const tabContent = document.getElementById(e.target.dataset.tab + 'Tab');
                if (tabContent) {
                    document.querySelectorAll('.response-content .tab-content').forEach(content => {
                        content.classList.remove('active');
                    });
                    tabContent.classList.add('active');
                }
            });
        });
    }

    setupAuthFields() {
        this.updateAuthType('none');
    }

    loadSampleData() {
        // Load sample URL
        document.getElementById('urlInput').value = 'https://httpbin.org/get';
        
        // Add sample headers
        this.addHeaderRow();
        const headerRows = document.querySelectorAll('.header-row');
        if (headerRows.length > 0) {
            const firstRow = headerRows[0];
            firstRow.querySelector('.header-key').value = 'User-Agent';
            firstRow.querySelector('.header-value').value = 'Litepost/1.0';
        }
    }

    switchTab(tabName) {
        console.log('Switching to tab:', tabName);
        
        // Handle response tabs
        if (tabName.startsWith('response-')) {
            console.log('Handling response tab');
            
            // Remove active class from all response tabs
            const responseTabs = document.querySelectorAll('.response-tabs .tab');
            responseTabs.forEach(tab => tab.classList.remove('active'));
            
            // Add active class to clicked tab
            const activeTab = document.querySelector(`[data-tab="${tabName}"]`);
            if (activeTab) {
                activeTab.classList.add('active');
                console.log('Active tab set:', activeTab);
            }
            
            // Hide all response content
            const allResponseContent = document.querySelectorAll('.response-content .tab-content');
            console.log('Found response content elements:', allResponseContent.length);
            allResponseContent.forEach(content => {
                content.classList.remove('active');
                content.style.display = 'none';
            });
            
            // Show selected content based on tabName
            let targetId;
            if (tabName === 'response-body') {
                targetId = 'responseBodyTab';
            } else if (tabName === 'response-headers') {
                targetId = 'responseHeadersTab';
            } else if (tabName === 'response-cookies') {
                targetId = 'responseCookiesTab';
            }
            
            const tabContent = document.getElementById(targetId);
            console.log('Target ID:', targetId, 'Tab content element:', tabContent);
            if (tabContent) {
                tabContent.classList.add('active');
                tabContent.style.display = 'block';
                console.log('Tab content activated and shown');
            }
        }
        // Handle request tabs
        else if (tabName.startsWith('request-')) {
            // Remove active class from all request tabs
            const requestTabs = document.querySelectorAll('.request-tabs .tab');
            requestTabs.forEach(tab => tab.classList.remove('active'));
            
            // Add active class to clicked tab
            const activeTab = document.querySelector(`[data-tab="${tabName}"]`);
            if (activeTab) {
                activeTab.classList.add('active');
            }
            
            // Hide all request content
            document.querySelectorAll('.request-content .tab-content').forEach(content => {
                content.classList.remove('active');
                content.style.display = 'none';
            });
            
            // Show selected content
            const tabContent = document.getElementById(tabName + 'Tab');
            if (tabContent) {
                tabContent.classList.add('active');
                tabContent.style.display = 'block';
            }
        }
    }

    addParamRow() {
        const paramList = document.getElementById('paramList');
        const paramRow = document.createElement('div');
        paramRow.className = 'param-row';
        paramRow.innerHTML = `
            <input type="text" placeholder="Key" class="param-key" />
            <input type="text" placeholder="Value" class="param-value" />
            <button class="remove-param">×</button>
        `;
        paramList.appendChild(paramRow);
        
        // Add remove functionality
        paramRow.querySelector('.remove-param').addEventListener('click', () => {
            paramRow.remove();
        });
    }

    addHeaderRow() {
        const headerList = document.getElementById('headerList');
        const headerRow = document.createElement('div');
        headerRow.className = 'header-row';
        headerRow.innerHTML = `
            <input type="text" placeholder="Header" class="header-key" />
            <input type="text" placeholder="Value" class="header-value" />
            <button class="remove-header">×</button>
        `;
        headerList.appendChild(headerRow);
        
        // Add remove functionality
        headerRow.querySelector('.remove-header').addEventListener('click', () => {
            headerRow.remove();
        });
    }

    updateBodyType(type) {
        const bodyContent = document.getElementById('bodyContent');
        
        switch (type) {
            case 'json':
                bodyContent.placeholder = 'Enter JSON body...\n{\n  "key": "value"\n}';
                break;
            case 'xml':
                bodyContent.placeholder = 'Enter XML body...\n<root>\n  <item>value</item>\n</root>';
                break;
            case 'form':
                bodyContent.placeholder = 'Enter form data...\nkey1=value1\nkey2=value2';
                break;
            case 'raw':
                bodyContent.placeholder = 'Enter raw body...';
                break;
            default:
                bodyContent.placeholder = 'Enter request body...';
        }
    }

    updateAuthType(type) {
        const authFields = document.getElementById('authFields');
        
        switch (type) {
            case 'basic':
                authFields.innerHTML = `
                    <div class="auth-field">
                        <label>Username:</label>
                        <input type="text" id="authUsername" placeholder="Enter username" />
                    </div>
                    <div class="auth-field">
                        <label>Password:</label>
                        <input type="password" id="authPassword" placeholder="Enter password" />
                    </div>
                `;
                break;
            case 'bearer':
                authFields.innerHTML = `
                    <div class="auth-field">
                        <label>Token:</label>
                        <input type="text" id="authToken" placeholder="Enter bearer token" />
                    </div>
                `;
                break;
            case 'apikey':
                authFields.innerHTML = `
                    <div class="auth-field">
                        <label>Key:</label>
                        <input type="text" id="authKey" placeholder="Enter API key name" />
                    </div>
                    <div class="auth-field">
                        <label>Value:</label>
                        <input type="text" id="authValue" placeholder="Enter API key value" />
                    </div>
                    <div class="auth-field">
                        <label>Header:</label>
                        <input type="text" id="authHeader" placeholder="X-API-Key" />
                    </div>
                `;
                break;
            default:
                authFields.innerHTML = '<p>No authentication required</p>';
        }
    }

    async sendRequest() {
        const sendButton = document.getElementById('sendButton');
        const originalText = sendButton.textContent;
        
        try {
            // Show loading state
            sendButton.textContent = 'Sending...';
            sendButton.disabled = true;
            document.body.classList.add('loading');

            // Build request object
            const request = this.buildRequest();
            
            // Create request in backend
            const response = await fetch('/api/requests', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(request)
            });

            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }

            const createdRequest = await response.json();
            
            // Execute the request
            const executeResponse = await fetch(`/api/requests/${createdRequest.id}/execute`, {
                method: 'POST'
            });

            if (!executeResponse.ok) {
                throw new Error(`HTTP error! status: ${executeResponse.status}`);
            }

            const result = await executeResponse.json();
            
            // Display response
            this.displayResponse(result);
            
        } catch (error) {
            console.error('Request failed:', error);
            this.displayError(error.message);
        } finally {
            // Reset loading state
            if (sendButton) {
                sendButton.textContent = originalText;
                sendButton.disabled = false;
            }
            document.body.classList.remove('loading');
        }
    }

    buildRequest() {
        const method = document.getElementById('methodSelect').value;
        const url = document.getElementById('urlInput').value;
        
        // Build headers
        const headers = {};
        document.querySelectorAll('.header-row').forEach(row => {
            const key = row.querySelector('.header-key').value;
            const value = row.querySelector('.header-value').value;
            if (key && value) {
                headers[key] = value;
            }
        });

        // Build query parameters
        const queryParams = {};
        document.querySelectorAll('.param-row').forEach(row => {
            const key = row.querySelector('.param-key').value;
            const value = row.querySelector('.param-value').value;
            if (key && value) {
                queryParams[key] = value;
            }
        });

        // Build body
        const bodyType = document.getElementById('bodyType').value;
        const bodyContent = document.getElementById('bodyContent').value;
        
        let body = null;
        if (bodyType !== 'none' && bodyContent) {
            body = {
                type: bodyType,
                content: bodyContent
            };
        }

        // Build auth
        const authType = document.getElementById('authType').value;
        let auth = null;
        if (authType !== 'none') {
            auth = {
                type: authType,
                config: this.getAuthConfig(authType)
            };
        }

        return {
            name: `Request to ${url}`,
            method: method,
            url: url,
            headers: headers,
            query_params: queryParams,
            body: body,
            auth: auth
        };
    }

    getAuthConfig(authType) {
        const config = {};
        
        switch (authType) {
            case 'basic':
                config.username = document.getElementById('authUsername')?.value || '';
                config.password = document.getElementById('authPassword')?.value || '';
                break;
            case 'bearer':
                config.token = document.getElementById('authToken')?.value || '';
                break;
            case 'apikey':
                config.key = document.getElementById('authKey')?.value || '';
                config.value = document.getElementById('authValue')?.value || '';
                config.header = document.getElementById('authHeader')?.value || 'X-API-Key';
                break;
        }
        
        return config;
    }

    displayResponse(response) {
        // Update status
        const statusCodeElement = document.getElementById('statusCode');
        const statusTextElement = document.getElementById('statusText');
        if (statusCodeElement) {
            statusCodeElement.textContent = response.status_code;
        }
        if (statusTextElement) {
            statusTextElement.textContent = this.getStatusText(response.status_code);
        }
        
        // Update response info
        const responseTimeElement = document.getElementById('responseTime');
        const responseSizeElement = document.getElementById('responseSize');
        if (responseTimeElement) {
            responseTimeElement.textContent = `${response.duration}ms`;
        }
        if (responseSizeElement) {
            responseSizeElement.textContent = `${response.size} bytes`;
        }
        
        // Update status code color
        if (statusCodeElement) {
            statusCodeElement.className = 'status-code';
            if (response.status_code >= 200 && response.status_code < 300) {
                statusCodeElement.style.backgroundColor = '#4CAF50';
            } else if (response.status_code >= 300 && response.status_code < 400) {
                statusCodeElement.style.backgroundColor = '#FF9800';
            } else if (response.status_code >= 400) {
                statusCodeElement.style.backgroundColor = '#F44336';
            }
        }
        
            // Update response body
            const responseBody = document.getElementById('responseBody');
            if (responseBody) {
                try {
                    // Try to format JSON
                    const jsonData = JSON.parse(response.body);
                    responseBody.textContent = JSON.stringify(jsonData, null, 2);
                    responseBody.style.whiteSpace = 'pre-wrap';
                    responseBody.style.fontFamily = 'Monaco, Menlo, Ubuntu Mono, monospace';
                } catch (e) {
                    // Not JSON, check if it's HTML and escape it
                    if (response.body.trim().startsWith('<')) {
                        // It's HTML, escape it for display
                        const escaped = this.escapeHtml(response.body);
                        responseBody.textContent = escaped;
                    } else {
                        // Regular text content
                        responseBody.textContent = response.body;
                    }
                    responseBody.style.whiteSpace = 'pre-wrap';
                    responseBody.style.fontFamily = 'Monaco, Menlo, Ubuntu Mono, monospace';
                }
            }
        
        // Update response headers
        this.displayResponseHeaders(response.headers);
        
        // Update response cookies
        this.displayResponseCookies(response.headers);
        
        // Show response area
        const responseArea = document.getElementById('responseArea');
        if (responseArea) {
            responseArea.style.display = 'block';
        }
    }

    displayResponseHeaders(headers) {
        console.log('Displaying response headers:', headers);
        const headersContainer = document.getElementById('responseHeaders');
        console.log('Headers container:', headersContainer);
        if (headersContainer) {
            headersContainer.innerHTML = '';
            
            Object.entries(headers).forEach(([key, value]) => {
                console.log('Adding header:', key, value);
                const headerRow = document.createElement('div');
                headerRow.className = 'header-row';
                headerRow.innerHTML = `
                    <span class="header-key">${key}</span>
                    <span class="header-value">${value}</span>
                `;
                headersContainer.appendChild(headerRow);
            });
            console.log('Headers container children:', headersContainer.children.length);
        }
    }

    displayResponseCookies(headers) {
        console.log('Displaying response cookies:', headers);
        const cookiesContainer = document.getElementById('responseCookies');
        console.log('Cookies container:', cookiesContainer);
        if (cookiesContainer) {
            cookiesContainer.innerHTML = '';
            
            // Look for Set-Cookie headers
            const setCookieHeaders = Object.entries(headers).filter(([key, value]) => 
                key.toLowerCase() === 'set-cookie'
            );
            
            console.log('Set-Cookie headers found:', setCookieHeaders.length);
            
            if (setCookieHeaders.length > 0) {
                setCookieHeaders.forEach(([key, value]) => {
                    console.log('Adding cookie:', key, value);
                    const cookieRow = document.createElement('div');
                    cookieRow.className = 'cookie-row';
                    cookieRow.innerHTML = `
                        <span class="cookie-name">${key}</span>
                        <span class="cookie-value">${value}</span>
                    `;
                    cookiesContainer.appendChild(cookieRow);
                });
            } else {
                // No cookies found
                console.log('No cookies found, adding no cookies message');
                const noCookiesRow = document.createElement('div');
                noCookiesRow.className = 'cookie-row';
                noCookiesRow.innerHTML = `
                    <span class="cookie-name">No Cookies</span>
                    <span class="cookie-value">No cookies were set in this response</span>
                `;
                cookiesContainer.appendChild(noCookiesRow);
            }
            console.log('Cookies container children:', cookiesContainer.children.length);
        }
    }

    displayError(error) {
        const statusCodeElement = document.getElementById('statusCode');
        const statusTextElement = document.getElementById('statusText');
        const responseBodyElement = document.getElementById('responseBody');
        
        if (statusCodeElement) {
            statusCodeElement.textContent = 'Error';
            statusCodeElement.style.backgroundColor = '#F44336';
        }
        if (statusTextElement) {
            statusTextElement.textContent = error;
        }
        if (responseBodyElement) {
            responseBodyElement.textContent = `Error: ${error}`;
            responseBodyElement.style.whiteSpace = 'pre-wrap';
            responseBodyElement.style.fontFamily = 'Monaco, Menlo, Ubuntu Mono, monospace';
        }
        
        // Show response area even for errors
        const responseArea = document.getElementById('responseArea');
        if (responseArea) {
            responseArea.style.display = 'block';
        }
    }

    escapeHtml(text) {
        const div = document.createElement('div');
        div.textContent = text;
        return div.innerHTML;
    }

    getStatusText(statusCode) {
        const statusTexts = {
            200: 'OK',
            201: 'Created',
            204: 'No Content',
            400: 'Bad Request',
            401: 'Unauthorized',
            403: 'Forbidden',
            404: 'Not Found',
            500: 'Internal Server Error',
            502: 'Bad Gateway',
            503: 'Service Unavailable'
        };
        return statusTexts[statusCode] || 'Unknown';
    }
}

// Initialize the application when the page loads
document.addEventListener('DOMContentLoaded', () => {
    new LitepostApp();
});
