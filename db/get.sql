SELECT
    id,
    address,
    visit_date,
    notes,
    rating,
    price,
    created_at,
    updated_at
FROM apartments
WHERE
    id = ?
