// Global variables
let currentApartmentId = null;
let apartmentData = [];

// DOM elements
const apartmentsListEl = document.getElementById('apartmentsList');
const loadingIndicatorEl = document.getElementById('loadingIndicator');
const emptyStateEl = document.getElementById('emptyState');

// Modals
const apartmentModal = new bootstrap.Modal(document.getElementById('apartmentModal'));
const detailsModal = new bootstrap.Modal(document.getElementById('detailsModal'));
const deleteModal = new bootstrap.Modal(document.getElementById('deleteModal'));

// Event listeners
document.addEventListener('DOMContentLoaded', () => {
    loadApartments();
    setupEventListeners();
});

function setupEventListeners() {
    // New apartment button
    document.getElementById('newApartmentBtn').addEventListener('click', () => {
        resetForm();
        document.getElementById('modalTitle').textContent = 'Add Apartment';
        apartmentModal.show();
    });

    // Star rating
    document.querySelectorAll('.rating-stars .star').forEach(star => {
        star.addEventListener('click', (e) => {
            const value = parseInt(e.target.getAttribute('data-value'));
            setRating(value);
        });

        star.addEventListener('mouseenter', (e) => {
            const value = parseInt(e.target.getAttribute('data-value'));
            highlightStars(value);
        });
    });

    document.getElementById('ratingStars').addEventListener('mouseleave', () => {
        const currentRating = parseInt(document.getElementById('rating').value) || 0;
        highlightStars(currentRating);
    });

    // Save apartment button
    document.getElementById('saveApartment').addEventListener('click', saveApartment);

    // Edit button in details modal
    document.getElementById('editBtn').addEventListener('click', () => {
        detailsModal.hide();
        const apartment = apartmentData.find(a => a.id === currentApartmentId);
        if (apartment) {
            populateForm(apartment);
            document.getElementById('modalTitle').textContent = 'Edit Apartment';
            apartmentModal.show();
        }
    });

    // Delete button in details modal
    document.getElementById('deleteBtn').addEventListener('click', () => {
        detailsModal.hide();
        deleteModal.show();
    });

    // Confirm delete button
    document.getElementById('confirmDelete').addEventListener('click', () => {
        deleteApartment(currentApartmentId);
    });
}

// Load apartments from API
async function loadApartments() {
    loadingIndicatorEl.classList.remove('d-none');
    emptyStateEl.classList.add('d-none');

    try {
        const response = await fetch('/api/apartments');
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        apartmentData = await response.json();

        renderApartmentsList();

        if (apartmentData.length === 0) {
            emptyStateEl.classList.remove('d-none');
        }
    } catch (error) {
        console.error('Error loading apartments:', error);
        showAlert('Failed to load apartments. Please try again.', 'danger');
    } finally {
        loadingIndicatorEl.classList.add('d-none');
    }
}

// Render apartments list
function renderApartmentsList() {
    // Clear the list except for loading and empty state
    Array.from(apartmentsList.children).forEach(child => {
        if (child !== loadingIndicatorEl && child !== emptyStateEl) {
            child.remove();
        }
    });

    apartmentData.forEach(apartment => {
        const col = document.createElement('div');
        col.className = 'col-md-4';

        const visitDate = apartment.visit_date ? new Date(apartment.visit_date).toLocaleDateString() : 'Not visited';

        col.innerHTML = `
            <div class="card apartment-card">
                <div class="card-body">
                    <h5 class="card-title">${escapeHtml(apartment.address)}</h5>
                    <h6 class="card-subtitle mb-2 text-muted">$${apartment.price.toFixed(2)} | Floor: ${apartment.floor || 1}</h6>
                    <div class="mb-2">
                        ${renderStarRating(apartment.rating)}
                    </div>
                    <div class="mb-2 small">
                        <span class="badge ${apartment.is_gated ? 'bg-success' : 'bg-light text-dark border'}">Gated</span>
                        <span class="badge ${apartment.has_garage ? 'bg-success' : 'bg-light text-dark border'}">Garage</span>
                        <span class="badge ${apartment.has_laundry ? 'bg-success' : 'bg-light text-dark border'}">Laundry</span>
                    </div>
                    <p class="card-text">${apartment.notes ? escapeHtml(apartment.notes.substring(0, 100)) + (apartment.notes.length > 100 ? '...' : '') : 'No notes'}</p>
                    <div class="text-muted small mb-2">Visited: ${visitDate}</div>
                    <button class="btn btn-sm btn-outline-primary view-details" data-id="${apartment.id}">View Details</button>
                </div>
            </div>
        `;

        apartmentsList.appendChild(col);

        // Add event listener to the view details button
        const detailsBtn = col.querySelector('.view-details');
        detailsBtn.addEventListener('click', () => showApartmentDetails(apartment.id));
    });
}

// Show apartment details
function showApartmentDetails(id) {
    const apartment = apartmentData.find(a => a.id === id);
    if (!apartment) return;

    currentApartmentId = id;

    const visitDate = apartment.visit_date
        ? new Date(apartment.visit_date).toLocaleString()
        : 'Not visited';

    const detailsContent = document.getElementById('detailsContent');
    detailsContent.innerHTML = `
        <div class="apartment-details">
            <h3>${escapeHtml(apartment.address)}</h3>
            <div class="mb-3">${renderStarRating(apartment.rating)}</div>
            <div class="row mb-3">
                <div class="col-6">
                    <strong>Price:</strong> $${apartment.price.toFixed(2)}
                </div>
                <div class="col-6">
                    <strong>Visit Date:</strong> ${visitDate}
                </div>
            </div>
            <div class="row mb-3">
                <div class="col-6">
                    <strong>Floor:</strong> ${apartment.floor || 1}
                </div>
                <div class="col-6">
                    <strong>Features:</strong>
                    <ul class="list-unstyled">
                        <li>
                            <i class="bi ${apartment.is_gated ? 'bi-check-circle-fill text-success' : 'bi-x-circle text-muted'}"></i>
                            Gated Community
                        </li>
                        <li>
                            <i class="bi ${apartment.has_garage ? 'bi-check-circle-fill text-success' : 'bi-x-circle text-muted'}"></i>
                            Garage
                        </li>
                        <li>
                            <i class="bi ${apartment.has_laundry ? 'bi-check-circle-fill text-success' : 'bi-x-circle text-muted'}"></i>
                            In-unit Laundry
                        </li>
                    </ul>
                </div>
            </div>
            <div class="mb-3">
                <strong>Notes:</strong>
                <p>${apartment.notes ? escapeHtml(apartment.notes) : 'No notes'}</p>
            </div>
            <div class="text-muted small">
                <div>Created: ${new Date(apartment.created_at).toLocaleString()}</div>
                <div>Updated: ${new Date(apartment.updated_at).toLocaleString()}</div>
            </div>
        </div>
    `;

    detailsModal.show();
}

// Save apartment (create or update)
async function saveApartment() {
    // Validate the form
    const addressInput = document.getElementById('address');
    if (!addressInput.value.trim()) {
        showAlert('Address is required', 'warning');
        return;
    }

    // Collect form data
    let visitDateValue = document.getElementById('visitDate').value;

    // Format the date in RFC3339 format if provided
    if (visitDateValue) {
        // The datetime-local input returns YYYY-MM-DDThh:mm
        // We need to append seconds and timezone: YYYY-MM-DDThh:mm:00Z
        visitDateValue = visitDateValue + ':00Z';
    }

    const apartmentData = {
        address: addressInput.value.trim(),
        visit_date: visitDateValue || null,
        price: parseFloat(document.getElementById('price').value) || 0,
        rating: parseInt(document.getElementById('rating').value) || 0,
        notes: document.getElementById('notes').value.trim(),
        floor: parseInt(document.getElementById('floor').value) || 1,
        is_gated: document.getElementById('isGated').checked,
        has_garage: document.getElementById('hasGarage').checked,
        has_laundry: document.getElementById('hasLaundry').checked
    };

    try {
        let url = '/api/apartments';
        let method = 'POST';

        // If editing an existing apartment
        if (currentApartmentId) {
            url += `/${currentApartmentId}`;
            method = 'PUT';
        }

        const response = await fetch(url, {
            method: method,
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(apartmentData)
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        // Close the modal and reload apartments
        apartmentModal.hide();
        await loadApartments();

        // Show success message
        showAlert(
            `Apartment ${currentApartmentId ? 'updated' : 'added'} successfully!`,
            'success'
        );

    } catch (error) {
        console.error('Error saving apartment:', error);
        showAlert(`Failed to ${currentApartmentId ? 'update' : 'add'} apartment. Please try again.`, 'danger');
    }
}

// Delete an apartment
async function deleteApartment(id) {
    try {
        const response = await fetch(`/api/apartments/${id}`, {
            method: 'DELETE'
        });

        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }

        // Close the modal and reload apartments
        deleteModal.hide();
        await loadApartments();

        // Show success message
        showAlert('Apartment deleted successfully!', 'success');

    } catch (error) {
        console.error('Error deleting apartment:', error);
        showAlert('Failed to delete apartment. Please try again.', 'danger');
    }
}

// Helper functions
function resetForm() {
    document.getElementById('apartmentId').value = '';
    document.getElementById('address').value = '';
    document.getElementById('visitDate').value = '';
    document.getElementById('price').value = '';
    document.getElementById('notes').value = '';
    document.getElementById('floor').value = '1';
    document.getElementById('isGated').checked = false;
    document.getElementById('hasGarage').checked = false;
    document.getElementById('hasLaundry').checked = false;
    setRating(0);
    currentApartmentId = null;
}

function populateForm(apartment) {
    document.getElementById('apartmentId').value = apartment.id;
    document.getElementById('address').value = apartment.address;

    // Format date for datetime-local input
    if (apartment.visit_date) {
        const date = new Date(apartment.visit_date);
        const formattedDate = date.toISOString().slice(0, 16); // Format as YYYY-MM-DDTHH:MM
        document.getElementById('visitDate').value = formattedDate;
    } else {
        document.getElementById('visitDate').value = '';
    }

    document.getElementById('price').value = apartment.price;
    document.getElementById('notes').value = apartment.notes || '';
    document.getElementById('floor').value = apartment.floor || 1;
    document.getElementById('isGated').checked = apartment.is_gated || false;
    document.getElementById('hasGarage').checked = apartment.has_garage || false;
    document.getElementById('hasLaundry').checked = apartment.has_laundry || false;
    setRating(apartment.rating);

    currentApartmentId = apartment.id;
}

function setRating(rating) {
    document.getElementById('rating').value = rating;
    highlightStars(rating);
}

function highlightStars(count) {
    const stars = document.querySelectorAll('.rating-stars .star');
    stars.forEach((star, index) => {
        if (index < count) {
            star.classList.add('active');
        } else {
            star.classList.remove('active');
        }
    });
}

function renderStarRating(rating) {
    let stars = '';
    for (let i = 1; i <= 5; i++) {
        stars += `<span class="star ${i <= rating ? 'active' : ''}"">★</span>`;
    }
    return stars;
}

function showAlert(message, type) {
    const alertContainer = document.getElementById('alertContainer');
    const alert = document.createElement('div');
    alert.className = `alert alert-${type} alert-dismissible fade show`;
    alert.innerHTML = `
        ${message}
        <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
    `;

    alertContainer.appendChild(alert);

    // Auto-dismiss after 5 seconds
    setTimeout(() => {
        alert.classList.remove('show');
        setTimeout(() => alert.remove(), 150);
    }, 5000);
}

function escapeHtml(unsafe) {
    return unsafe
         .replace(/&/g, "&amp;")
         .replace(/</g, "&lt;")
         .replace(/>/g, "&gt;")
         .replace(/"/g, "&quot;")
         .replace(/'/g, "&#039;");
}
