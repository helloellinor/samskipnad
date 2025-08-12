// Minimal HTMX-like functionality for Samskipnad
// This provides the core dynamic features we need

class HTMXLite {
    constructor() {
        this.init();
    }

    init() {
        // Load initial content
        this.loadBalanceOnPageLoad();
        this.loadFirstCategoryOnPageLoad();
        
        // Set up event listeners
        this.setupEventListeners();
    }

    setupEventListeners() {
        // Handle category tab clicks
        document.addEventListener('click', (e) => {
            if (e.target.classList.contains('category-tab')) {
                e.preventDefault();
                this.handleCategoryClick(e.target);
            }
        });

        // Handle purchase button clicks
        document.addEventListener('click', (e) => {
            if (e.target.classList.contains('purchase-btn')) {
                e.preventDefault();
                this.handlePurchaseClick(e.target);
            }
        });

        // Handle refresh balance clicks
        document.addEventListener('click', (e) => {
            if (e.target.closest('#balance-refresh') || e.target.textContent.includes('REFRESH BALANCE')) {
                e.preventDefault();
                this.loadBalance();
            }
        });
    }

    async handleCategoryClick(button) {
        // Get category ID
        const categoryId = this.extractCategoryId(button);
        if (!categoryId) return;

        // Update active tab
        document.querySelectorAll('.category-tab').forEach(tab => {
            tab.classList.remove('active');
        });
        button.classList.add('active');

        // Load category content
        await this.loadCategory(categoryId);
    }

    async handlePurchaseClick(button) {
        const categoryId = button.getAttribute('data-category');
        const packageIndex = button.getAttribute('data-package');
        
        if (!categoryId || packageIndex === null) {
            // Extract from hx-vals if present
            const hxVals = button.getAttribute('hx-vals');
            if (hxVals) {
                try {
                    const vals = JSON.parse(hxVals);
                    await this.makePurchase(vals.category_id, vals.package_index);
                } catch (e) {
                    console.error('Error parsing purchase data:', e);
                }
            }
            return;
        }

        await this.makePurchase(categoryId, packageIndex);
    }

    extractCategoryId(button) {
        // Extract category ID from button text or data attributes
        const text = button.textContent.trim();
        const categoryMap = {
            'Personlig Trening': 'personal_training',
            'PT 30 min': 'pt_30min',
            'Reformer / Apparater': 'reformer',
            'Spesialtimer': 'special'
        };
        return categoryMap[text];
    }

    async loadBalanceOnPageLoad() {
        const balanceContainer = document.getElementById('balance-display');
        if (balanceContainer) {
            await this.loadBalance();
        }
    }

    async loadFirstCategoryOnPageLoad() {
        const categoryContent = document.getElementById('category-content');
        if (categoryContent) {
            // Load first category (personal_training)
            await this.loadCategory('personal_training');
        }
    }

    async loadBalance() {
        const balanceContainer = document.getElementById('balance-display');
        if (!balanceContainer) return;

        try {
            const response = await fetch('/api/klippekort/balance', {
                credentials: 'same-origin'
            });
            if (response.ok) {
                const html = await response.text();
                balanceContainer.innerHTML = html;
            } else {
                balanceContainer.innerHTML = '<div class="text-center p-4 text-danger">Failed to load balance</div>';
            }
        } catch (error) {
            console.error('Error loading balance:', error);
            balanceContainer.innerHTML = '<div class="text-center p-4 text-danger">Error loading balance</div>';
        }
    }

    async loadCategory(categoryId) {
        const categoryContent = document.getElementById('category-content');
        if (!categoryContent) return;

        // Show loading state
        categoryContent.innerHTML = '<div class="text-center p-5"><div class="cyber-loader"></div><p class="mt-3 text-muted">Loading category...</p></div>';

        try {
            const response = await fetch(`/api/klippekort/category?category=${categoryId}`, {
                credentials: 'same-origin'
            });
            if (response.ok) {
                const html = await response.text();
                categoryContent.innerHTML = html;
                this.setupPurchaseButtons();
            } else {
                categoryContent.innerHTML = '<div class="text-center p-4 text-danger">Failed to load category</div>';
            }
        } catch (error) {
            console.error('Error loading category:', error);
            categoryContent.innerHTML = '<div class="text-center p-4 text-danger">Error loading category</div>';
        }
    }

    setupPurchaseButtons() {
        // Add click handlers to purchase buttons in the loaded content
        const purchaseButtons = document.querySelectorAll('.purchase-btn');
        purchaseButtons.forEach(button => {
            button.addEventListener('click', (e) => {
                e.preventDefault();
                this.handlePurchaseClick(button);
            });
        });
    }

    async makePurchase(categoryId, packageIndex) {
        const purchaseResult = document.getElementById('purchase-result');
        if (!purchaseResult) return;

        // Show loading state
        purchaseResult.innerHTML = '<div class="text-center p-4"><div class="cyber-loader"></div><p class="mt-3 text-muted">Processing purchase...</p></div>';

        try {
            const formData = new FormData();
            formData.append('category_id', categoryId);
            formData.append('package_index', packageIndex);

            const response = await fetch('/api/klippekort/purchase', {
                method: 'POST',
                body: formData,
                credentials: 'same-origin'
            });

            if (response.ok) {
                const html = await response.text();
                purchaseResult.innerHTML = html;
                
                // Refresh balance after purchase
                setTimeout(() => {
                    this.loadBalance();
                }, 1000);
            } else {
                purchaseResult.innerHTML = '<div class="error-message">Failed to process purchase. Please try again.</div>';
            }
        } catch (error) {
            console.error('Error making purchase:', error);
            purchaseResult.innerHTML = '<div class="error-message">Error processing purchase. Please try again.</div>';
        }
    }

    // Utility function to simulate HTMX's htmx.ajax
    static ajax(method, url, target) {
        fetch(url, { credentials: 'same-origin' })
            .then(response => response.text())
            .then(html => {
                const targetElement = document.querySelector(target);
                if (targetElement) {
                    targetElement.innerHTML = html;
                }
            })
            .catch(error => {
                console.error('HTMX Ajax error:', error);
            });
    }
}

// Initialize when DOM is loaded
document.addEventListener('DOMContentLoaded', () => {
    window.htmxLite = new HTMXLite();
    
    // Provide htmx compatibility
    window.htmx = {
        ajax: HTMXLite.ajax
    };
});

// Global function for package selection
function selectPackage(card, categoryId, packageIndex) {
    // Visual feedback
    document.querySelectorAll('.package-card').forEach(c => c.classList.remove('selected'));
    card.classList.add('selected');
}