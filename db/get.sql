SELECT
    id,
    address,
    visit_date,
    notes,
    rating,
    price,
    floor,
    is_gated,
    has_garage,
    has_laundry,
    created_at,
    updated_at
FROM apartments
WHERE
    id = ?
