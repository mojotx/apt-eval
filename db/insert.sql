INSERT INTO
    apartments (
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
    )
VALUES (
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        ?,
        CURRENT_TIMESTAMP,
        CURRENT_TIMESTAMP
    ) RETURNING id,
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
